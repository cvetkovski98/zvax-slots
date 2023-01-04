package delivery

import (
	"context"
	"log"

	"github.com/cvetkovski98/zvax-common/gen/pbslot"
	slots "github.com/cvetkovski98/zvax-slots/internal"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type SlotGrpcServerImpl struct {
	pbslot.UnimplementedSlotGrpcServer

	s slots.Service
}

func (s SlotGrpcServerImpl) GetSlotList(ctx context.Context, request *pbslot.SlotListRequest) (*pbslot.SlotListResponse, error) {
	pageRequest, err := SlotListRequestToDto(request)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "invalid request: %v", err)
	}
	payload, err := s.s.GetSlotsAtLocationBetween(ctx, pageRequest)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to get slots: %v", err)
	}
	return SlotPageToResponse(payload), nil
}

func (s SlotGrpcServerImpl) CreateSlotReservation(ctx context.Context, request *pbslot.SlotReservationRequest) (*pbslot.SlotReservationResponse, error) {
	reservation, err := s.s.CreateReservation(ctx, request.SlotId)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to create slot: %v", err)
	}
	return ReservationDtoToResponse(reservation), nil
}

func (s SlotGrpcServerImpl) ConfirmSlotReservation(ctx context.Context, request *pbslot.SlotConfirmationRequest) (*pbslot.SlotConfirmationResponse, error) {
	token, err := s.s.ConfirmReservation(ctx, request.ReservationId)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to confirm slot: %v", err)
	}
	log.Println("token: ", token)
	return &pbslot.SlotConfirmationResponse{
		SlotConfirmationToken: token,
	}, nil
}

func NewSlotGrpcServerImpl(s slots.Service) pbslot.SlotGrpcServer {
	return &SlotGrpcServerImpl{
		s: s,
	}
}
