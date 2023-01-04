package mapper

import (
	"time"

	"github.com/cvetkovski98/zvax-common/gen/pbslot"
	"github.com/cvetkovski98/zvax-slots/internal/model"
)

func NewSlotResponseFromSlot(slot model.Slot) *pbslot.SlotResponse {
	return &pbslot.SlotResponse{
		SlotId:   model.NewSlotRedisId(slot.Location, slot.DateTime),
		DateTime: slot.DateTime.Format(time.RFC3339),
	}
}
