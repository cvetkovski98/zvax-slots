package service

import (
	"context"
	"fmt"
	"log"

	slots "github.com/cvetkovski98/zvax-slots/internal"
	"github.com/cvetkovski98/zvax-slots/internal/crypto"
	"github.com/cvetkovski98/zvax-slots/internal/dto"
	"github.com/cvetkovski98/zvax-slots/internal/mapper"
	"github.com/pkg/errors"
)

type SlotServiceImpl struct {
	r slots.Repository
}

func (s *SlotServiceImpl) GetSlotsAtLocationBetween(ctx context.Context, page *dto.SlotListRequest) (*dto.Page[dto.Slot], error) {
	slots, err := s.r.FindAllWithDateTimeBetween(ctx, page.StartDate, page.EndDate)
	if err != nil {
		return nil, err
	}
	items := make([]*dto.Slot, len(slots))
	for i, slot := range slots {
		items[i] = &dto.Slot{
			SlotID:    slot.SlotID,
			Location:  slot.Location,
			DateTime:  slot.DateTime,
			Available: slot.Available,
		}
	}
	return &dto.Page[dto.Slot]{
		Items: items,
	}, nil
}

func (s *SlotServiceImpl) CreateSlot(ctx context.Context, slot *dto.CreateSlotRequest) (*dto.Slot, error) {
	slotModel := mapper.CreateSlotDtoToModel(slot)
	createdSlot, err := s.r.InsertOne(ctx, slotModel)
	if err != nil {
		return nil, err
	}
	return mapper.SlotModelToDto(createdSlot), nil
}

func (s *SlotServiceImpl) CreateReservation(ctx context.Context, slotID string) (*dto.Reservation, error) {
	reservation, err := s.r.ReserveOneByKey(ctx, slotID)
	if err != nil {
		return nil, err
	}
	return mapper.ReservationModelToDto(reservation), nil
}

func (s *SlotServiceImpl) ConfirmReservation(ctx context.Context, reservationID string) (string, error) {
	confirmed, err := s.r.ConfirmOneByReservationID(ctx, reservationID)
	log.Println(confirmed)
	if err != nil {
		return "", errors.Wrap(err, fmt.Sprintf(
			"error confirming reservation with id=%s",
			confirmed.ReservationID,
		))
	}
	d := mapper.ReservationModelToDto(confirmed)
	return crypto.SignReservation(d)
}

func NewSlotServiceImpl(slotRepository slots.Repository) slots.Service {
	return &SlotServiceImpl{
		r: slotRepository,
	}
}
