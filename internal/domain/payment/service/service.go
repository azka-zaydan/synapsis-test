package service

import (
	"context"
	"time"

	"github.com/azka-zaydan/synapsis-test/configs"
	"github.com/azka-zaydan/synapsis-test/infras"
	orderModel "github.com/azka-zaydan/synapsis-test/internal/domain/order/model"
	orderRepo "github.com/azka-zaydan/synapsis-test/internal/domain/order/repository"
	"github.com/azka-zaydan/synapsis-test/internal/domain/payment/model/dto"

	"github.com/azka-zaydan/synapsis-test/internal/domain/payment/repository"
	"github.com/guregu/null"

	"github.com/rs/zerolog/log"
)

type PaymentService interface {
	Pay(ctx context.Context, req dto.PayRequest) (res dto.PaymentResponse, err error)
}

type PaymentServiceImpl struct {
	Repo      repository.PaymentRepository
	Redis     *infras.Redis
	config    *configs.Config
	OrderRepo orderRepo.OrderRepository
}

func ProvidePaymentServiceImpl(repo repository.PaymentRepository, redis *infras.Redis, config *configs.Config, orderRepo orderRepo.OrderRepository) *PaymentServiceImpl {
	return &PaymentServiceImpl{
		Redis:     redis,
		Repo:      repo,
		config:    config,
		OrderRepo: orderRepo,
	}
}

func (s *PaymentServiceImpl) CreatePayment(ctx context.Context, req dto.CreatePaymentRequest) (res dto.PaymentResponse, err error) {
	mod, err := req.ToModel()
	if err != nil {
		log.Error().Err(err).Msg("[CreatePayment] Failed creating model")
		return
	}
	err = s.Repo.CreatePayment(ctx, &mod)
	if err != nil {
		log.Error().Err(err).Msg("[CreatePayment] Failed creating payment")
		return
	}
	return dto.NewPaymentResponse(mod), nil
}

func (s *PaymentServiceImpl) Pay(ctx context.Context, req dto.PayRequest) (res dto.PaymentResponse, err error) {
	mod, err := s.Repo.GetPaymentByOrderID(ctx, req.OrderID)
	if err != nil {
		log.Error().Err(err).Msg("[Pay] Failed GetPaymentByOrderID")
		return
	}

	mod.Pay()

	order, err := s.OrderRepo.GetOrderByID(ctx, mod.OrderID.String())
	if err != nil {
		log.Error().Err(err).Msg("[Pay] Failed GetOrderByID")
		return
	}
	order.Status = int(orderModel.OrderPaidStatus)
	order.PaymentAt = null.TimeFrom(time.Now())

	err = s.Repo.UpdatePayment(ctx, &mod)
	if err != nil {
		log.Error().Err(err).Msg("[Pay] Failed UpdatePayment")
		return
	}

	err = s.OrderRepo.UpdateOrder(ctx, &order)
	if err != nil {
		log.Error().Err(err).Msg("[Pay] Failed UpdateOrder")
		return
	}

	return dto.NewPaymentResponse(mod), nil
}
