package repository

import (
	"context"
	"fmt"
	"strconv"
	"time"

	slots "github.com/cvetkovski98/zvax-slots/internal"
	"github.com/cvetkovski98/zvax-slots/internal/model"
	"github.com/go-redis/redis/v9"
	"github.com/pkg/errors"
)

const slotsOrderedSetKey = "slots"
const confirmationTimeout = time.Minute * 5

type RedisSlotRepository struct {
	rdb *redis.Client
}

// FindOneByKey returns a slot by key
func (r *RedisSlotRepository) FindOneByKey(ctx context.Context, key string) (*model.Slot, error) {
	result, err := r.rdb.HGetAll(ctx, key).Result()
	if err != nil && err != redis.Nil {
		return nil, errors.Wrap(err, fmt.Sprintf("error getting slot with key=%s", key))
	}
	return model.NewSlotFromMap(result)
}

func (repository *RedisSlotRepository) InsertOne(ctx context.Context, slot *model.Slot) (*model.Slot, error) {
	slot.SlotID = model.NewSlotRedisId(slot.Location, slot.DateTime)
	_, err := repository.rdb.Pipelined(ctx, func(pipe redis.Pipeliner) error {
		err := pipe.HSet(ctx, slot.SlotID, slot.ToMap()).Err()
		if err != nil {
			return errors.Wrap(err, fmt.Sprintf("error setting key=%s", slot.SlotID))
		}
		err = pipe.ZAdd(ctx, slotsOrderedSetKey, redis.Z{
			Score:  float64(slot.DateTime.Unix()),
			Member: slot.SlotID,
		}).Err()
		if err != nil {
			return errors.Wrap(err, fmt.Sprintf("error adding slot with key=%s to ordered set", slot.SlotID))
		}
		return nil
	})
	if err != nil && err != redis.Nil {
		return nil, errors.Wrap(err, fmt.Sprintf("error inserting slot with key=%s", slot.SlotID))
	}
	return slot, nil
}

func (repository *RedisSlotRepository) FindAllWithDateTimeBetween(ctx context.Context, from time.Time, to time.Time) ([]*model.Slot, error) {
	const rangeErrTmpl = "error getting slots in time range from=%s to=%s"

	slotIDs, err := repository.rdb.ZRangeByScore(ctx, slotsOrderedSetKey, &redis.ZRangeBy{
		Min:    strconv.FormatInt(from.Unix(), 10),
		Max:    strconv.FormatInt(to.Unix(), 10),
		Offset: 0,
		Count:  -1,
	}).Result()
	if err != nil {
		return nil, errors.Wrap(err, fmt.Sprintf(rangeErrTmpl, from, to))
	}
	commands, err := repository.rdb.Pipelined(ctx, func(pipe redis.Pipeliner) error {
		for _, slotID := range slotIDs {
			err := pipe.HGetAll(ctx, slotID).Err()
			if err != nil {
				return errors.Wrap(err, fmt.Sprintf("error getting slot with id=%s", slotID))
			}
		}
		return nil
	})
	if err != nil {
		return nil, errors.Wrap(err, fmt.Sprintf(rangeErrTmpl, from, to))
	}
	slots := make([]*model.Slot, 0, len(commands))
	for i, command := range commands {
		sMap, err := command.(*redis.MapStringStringCmd).Result()
		if err != nil && err != redis.Nil {
			return nil, errors.Wrap(err, fmt.Sprintf("error processing command result for slot=%s", slotIDs[i]))
		}
		if slot, err := model.NewSlotFromMap(sMap); err != nil {
			return nil, errors.Wrap(err, fmt.Sprintf(
				"error creating slot from map for slot=%s",
				slotIDs[i],
			))
		} else if slot.Available {
			slots = append(slots, slot)
		}
	}
	return slots, nil
}

func (r *RedisSlotRepository) ReserveOneByKey(ctx context.Context, key string) (*model.Reservation, error) {
	reservationID := model.NewReservationRedisId(key)
	slot, err := r.FindOneByKey(ctx, key)
	if err != nil {
		return nil, errors.Wrap(err, fmt.Sprintf("error finding slot with key=%s", key))
	}
	if !slot.Available {
		return nil, errors.New("slot is not available")
	}

	tfx := func(tx *redis.Tx) error {
		// update the slot to unavailable
		err = tx.HSet(ctx, slot.SlotID, "available", "false").Err()
		if err != nil {
			return errors.Wrap(err, fmt.Sprintf(
				"error setting key=%s to value=%v",
				slot.SlotID,
				slot,
			))
		}
		// create a reservation for the slot
		err = tx.Set(ctx, reservationID, slot.SlotID, confirmationTimeout).Err()
		if err != nil {
			return errors.Wrap(err, fmt.Sprintf("error setting key=%s", reservationID))
		}
		return nil
	}

	err = r.rdb.Watch(ctx, tfx, reservationID, key)
	if err != nil {
		return nil, errors.Wrap(err, fmt.Sprintf(
			"error reserving slot with key=%s",
			key,
		))
	}
	return &model.Reservation{
		ReservationID: reservationID,
		SlotID:        key,
		ValidUntil:    time.Now().Add(confirmationTimeout),
	}, nil
}

func (repository *RedisSlotRepository) ConfirmOneByReservationID(ctx context.Context, reservationID string) (*model.Reservation, error) {
	slotID := model.NewSlotRedisIDFromReservationID(reservationID)

	tfx := func(tx *redis.Tx) error {
		exists, err := tx.Exists(ctx, reservationID).Result()
		if err != nil {
			return errors.Wrap(err, fmt.Sprintf(
				"error checking if key=%s exists",
				reservationID,
			))
		}
		if exists != 0 {
			// we try to book the slot if available
			slot, err := repository.FindOneByKey(ctx, slotID)
			if err != nil {
				return errors.Wrap(err, fmt.Sprintf("error finding slot with key=%s", slotID))
			}
			if !slot.Available {
				return errors.New("slot is not available")
			}
			err = tx.HSet(ctx, slot.SlotID, "available", "false").Err()
			if err != nil {
				return errors.Wrap(err, "error setting slot to unavailable")
			}
			return nil
		}

		// try to expire the reservation
		err = tx.Expire(ctx, reservationID, 0).Err()
		if err != nil {
			return errors.Wrap(err, fmt.Sprintf(
				"error expiring reservation with id=%s",
				reservationID,
			))
		}
		return nil
	}

	err := repository.rdb.Watch(ctx, tfx, reservationID, slotID)
	if err != nil {
		return nil, errors.Wrap(err, fmt.Sprintf(
			"error confirming reservation with id=%s",
			reservationID,
		))
	}
	return &model.Reservation{
		ReservationID: reservationID,
		SlotID:        slotID,
		ValidUntil:    time.Now().Add(confirmationTimeout),
	}, nil
}

func NewRedisSlotRepository(rdb *redis.Client) slots.Repository {
	return &RedisSlotRepository{rdb: rdb}
}
