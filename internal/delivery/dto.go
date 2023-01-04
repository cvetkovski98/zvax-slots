package delivery

import (
	"time"

	"github.com/cvetkovski98/zvax-common/gen/pbslot"
	"github.com/cvetkovski98/zvax-slots/internal/dto"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func SlotListRequestToDto(req *pbslot.SlotListRequest) (*dto.SlotListRequest, error) {
	from, err := time.Parse("2006-01-02", req.StartDate)
	if err != nil {
		return nil, ErrInvalidDateFormat
	}

	to, err := time.Parse("2006-01-02", req.EndDate)
	if err != nil {
		return nil, ErrInvalidDateFormat
	}

	if to.Before(from) {
		return nil, ErrInvalidDateRange
	}

	return &dto.SlotListRequest{
		Location:  req.Location,
		StartDate: from,
		EndDate:   to,
	}, nil
}

func SlotDtoToResponse(slot *dto.Slot) *pbslot.SlotResponse {
	dateStr := slot.DateTime.Format(time.RFC3339)
	return &pbslot.SlotResponse{
		SlotId:   slot.SlotID,
		Location: slot.Location,
		DateTime: dateStr,
	}
}

func SlotPageToResponse(page *dto.Page[dto.Slot]) *pbslot.SlotListResponse {
	items := make([]*pbslot.SlotResponse, len(page.Items))
	for i, slot := range page.Items {
		items[i] = SlotDtoToResponse(slot)
	}

	return &pbslot.SlotListResponse{
		Items: items,
	}
}

func ReservationDtoToResponse(reservation *dto.Reservation) *pbslot.SlotReservationResponse {
	return &pbslot.SlotReservationResponse{
		ReservationId: reservation.ReservationID,
		ValidUntil:    timestamppb.New(reservation.ValidUntil),
	}
}
