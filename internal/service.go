package slots

import (
	"context"

	"github.com/cvetkovski98/zvax-slots/internal/dto"
)

type Service interface {
	GetSlotsAtLocationBetween(ctx context.Context, page *dto.SlotListRequest) (*dto.Page[dto.Slot], error)
	CreateSlot(ctx context.Context, slot *dto.CreateSlotRequest) (*dto.Slot, error)
	CreateReservation(ctx context.Context, slotID string) (*dto.Reservation, error)
	ConfirmReservation(ctx context.Context, reservationID string) (string, error)
}
