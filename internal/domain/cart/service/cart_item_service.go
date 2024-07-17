package service

import (
	"context"

	"github.com/azka-zaydan/synapsis-test/internal/domain/cart/model/dto"
	"github.com/rs/zerolog/log"
)

func (s *CartServiceImpl) CreateCartItem(ctx context.Context, req dto.CartItemCreateRequest) (res dto.CartItemResponse, err error) {
	item, err := req.ToModel()
	if err != nil {
		log.Error().Err(err).Msg("[CreateCartItem] Failed Creating Model")
		return
	}
	err = s.Repo.CreateCartItem(ctx, &item)
	if err != nil {
		log.Error().Err(err).Msg("[CreateCartItem] Failed Create Item")
		return
	}
	return dto.NewCartItemResponse(item), nil
}
