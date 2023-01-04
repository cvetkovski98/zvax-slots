package dto

import "time"

type Reservation struct {
	ReservationID string
	SlotID        string
	ValidUntil    time.Time
}
