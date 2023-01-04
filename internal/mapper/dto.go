package mapper

import (
	"github.com/cvetkovski98/zvax-slots/internal/dto"
	"github.com/cvetkovski98/zvax-slots/internal/model"
)

func SlotModelToDto(model *model.Slot) *dto.Slot {
	return &dto.Slot{
		SlotID:    model.SlotID,
		Location:  model.Location,
		DateTime:  model.DateTime,
		Available: model.Available,
	}
}

func ReservationModelToDto(model *model.Reservation) *dto.Reservation {
	return &dto.Reservation{
		SlotID:        model.SlotID,
		ReservationID: model.ReservationID,
		ValidUntil:    model.ValidUntil,
	}
}
