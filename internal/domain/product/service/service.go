package service

import (
	"context"

	"github.com/azka-zaydan/synapsis-test/internal/domain/product/model"
	"github.com/azka-zaydan/synapsis-test/internal/domain/product/model/dto"
	"github.com/azka-zaydan/synapsis-test/internal/domain/product/repository"
	"github.com/rs/zerolog/log"
)

type ProductService interface {
	GetProductByFilter(ctx context.Context, filter model.Filter) (res dto.ProductFilterResponse, err error)
	CreateProduct(ctx context.Context, req dto.ProductCreateRequest) (res dto.ProductResponse, err error)
}

type ProductServiceImpl struct {
	Repo repository.ProductRepository
}

func ProvideProductServiceImpl(repo repository.ProductRepository) *ProductServiceImpl {
	return &ProductServiceImpl{
		Repo: repo,
	}
}

func (s *ProductServiceImpl) GetProductByFilter(ctx context.Context, filter model.Filter) (res dto.ProductFilterResponse, err error) {
	err = dto.TransformToDBField(&filter)
	if err != nil {
		log.Error().Err(err).Msg("[GetProductByFilter] Failed TransformToDBField")
		return
	}

	data, totalData, err := s.Repo.GetProductByFilter(ctx, &filter)

	if err != nil {
		log.Error().Err(err).Msg("[GetProductByFilter] Failed GetProductByFilter")
		return
	}

	return dto.NewProductFilterResponse(data, filter.Page, filter.PageSize, totalData), nil
}

func (s *ProductServiceImpl) CreateProduct(ctx context.Context, req dto.ProductCreateRequest) (res dto.ProductResponse, err error) {
	prod, err := req.ToModel()
	if err != nil {
		log.Error().Err(err).Msg("[CreateProduct] Failed creating model")
		return
	}
	err = s.Repo.CreateProduct(ctx, &prod)
	if err != nil {
		log.Error().Err(err).Msg("[CreateProduct] Failed CreateProduct")
		return
	}
	return dto.NewProductResponse(prod), nil
}
