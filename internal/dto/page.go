package dto

import (
	"time"
)

type DTO interface {
	Slot
}

type Page[T DTO] struct {
	Items []*T
}

type SlotListRequest struct {
	StartDate time.Time
	EndDate   time.Time
	Location  string
}
