package service

import (
	"context"

	"github.com/azka-zaydan/synapsis-test/configs"
	"github.com/azka-zaydan/synapsis-test/infras"
	"github.com/azka-zaydan/synapsis-test/internal/domain/payment/model/dto"
	"github.com/azka-zaydan/synapsis-test/internal/domain/payment/repository"
	"github.com/rs/zerolog/log"
)

type PaymentService interface{}

type PaymentServiceImpl struct {
	Repo   repository.PaymentRepository
	Redis  *infras.Redis
	config *configs.Config
}

func ProvidePaymentServiceImpl(repo repository.PaymentRepository, redis *infras.Redis, config *configs.Config) *PaymentServiceImpl {
	return &PaymentServiceImpl{
		Redis:  redis,
		Repo:   repo,
		config: config,
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

	err = s.Repo.UpdatePayment(ctx, &mod)
	if err != nil {
		log.Error().Err(err).Msg("[Pay] Failed UpdatePayment")
		return
	}

	return dto.NewPaymentResponse(mod), nil
}
