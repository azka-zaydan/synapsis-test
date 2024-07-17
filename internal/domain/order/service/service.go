package service

import (
	"github.com/azka-zaydan/synapsis-test/configs"
	"github.com/azka-zaydan/synapsis-test/infras"
	"github.com/azka-zaydan/synapsis-test/internal/domain/order/repository"
)

type OrderService interface{}

type OrderServiceImpl struct {
	Repo   repository.OrderRepository
	Redis  *infras.Redis
	config *configs.Config
}

func ProvideOrderServiceImpl(repo repository.OrderRepository, redis *infras.Redis, config *configs.Config) *OrderServiceImpl {
	return &OrderServiceImpl{
		Redis:  redis,
		Repo:   repo,
		config: config,
	}
}
