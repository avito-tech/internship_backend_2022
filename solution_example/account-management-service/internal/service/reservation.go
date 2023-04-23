package service

import (
	"account-management-service/internal/entity"
	"account-management-service/internal/repo"
	"context"
)

type ReservationService struct {
	reservationRepo repo.Reservation
}

func NewReservationService(reservationRepo repo.Reservation) *ReservationService {
	return &ReservationService{reservationRepo: reservationRepo}
}

func (s *ReservationService) CreateReservation(ctx context.Context, input ReservationCreateInput) (int, error) {
	reservation := entity.Reservation{
		AccountId: input.AccountId,
		ProductId: input.ProductId,
		OrderId:   input.OrderId,
		Amount:    input.Amount,
	}

	id, err := s.reservationRepo.CreateReservation(ctx, reservation)
	if err != nil {
		return 0, ErrCannotCreateReservation
	}

	return id, nil
}

func (s *ReservationService) RefundReservationByOrderId(ctx context.Context, orderId int) error {
	return s.reservationRepo.RefundReservationByOrderId(ctx, orderId)
}

func (s *ReservationService) RevenueReservationByOrderId(ctx context.Context, orderId int) error {
	return s.reservationRepo.RevenueReservationByOrderId(ctx, orderId)
}
