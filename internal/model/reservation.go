package model

import (
	"fmt"
	"time"

	"github.com/pkg/errors"
)

type Reservation struct {
	ReservationID string    `json:"reservationId" redis:"reservationId"`
	SlotID        string    `json:"slotId" redis:"slotId"`
	ValidUntil    time.Time `json:"validUntil" redis:"validUntil"`
}

func (reservation *Reservation) ToMap() map[string]string {
	return map[string]string{
		"reservationId": reservation.ReservationID,
		"slotId":        reservation.SlotID,
		"validUntil":    reservation.ValidUntil.Format(time.RFC3339),
	}
}

func NewReservationFromMap(rMap map[string]string) (*Reservation, error) {
	reservationID, ok := rMap["reservationId"]
	if !ok {
		return nil, errors.New("reservationId is not in hash")
	}
	slotID, ok := rMap["slotId"]
	if !ok {
		return nil, errors.New("slotId is not in hash")
	}
	validUntil, ok := rMap["validUntil"]
	if !ok {
		return nil, errors.New("valid_until is not in hash")
	}
	dateTime, err := time.Parse(time.RFC3339, validUntil)
	if err != nil {
		return nil, errors.Wrap(err, "failed to parse valid_until")
	}
	return &Reservation{
		ReservationID: reservationID,
		SlotID:        slotID,
		ValidUntil:    dateTime,
	}, err
}

func NewReservationRedisId(slotId string) string {
	return fmt.Sprintf("reservation:%s", slotId)
}
