package mapper

import (
	"github.com/cvetkovski98/zvax-slots/internal/dto"
	"github.com/cvetkovski98/zvax-slots/internal/model"
)

func CreateSlotDtoToModel(dto *dto.CreateSlotRequest) *model.Slot {
	return &model.Slot{
		Location:  dto.Location,
		DateTime:  dto.DateTime,
		Available: dto.Available,
	}
}

func ReservationDtoToModel(dto *dto.Reservation) *model.Reservation {
	return &model.Reservation{
		SlotID:        dto.SlotID,
		ReservationID: dto.ReservationID,
		ValidUntil:    dto.ValidUntil,
	}
}
