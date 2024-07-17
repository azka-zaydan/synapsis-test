package service

import (
	"context"

	cartDto "github.com/azka-zaydan/synapsis-test/internal/domain/cart/model/dto"
	cartSvc "github.com/azka-zaydan/synapsis-test/internal/domain/cart/service"
	"github.com/azka-zaydan/synapsis-test/internal/domain/user/model/dto"
	"github.com/azka-zaydan/synapsis-test/internal/domain/user/repository"
	"github.com/rs/zerolog/log"
)

type UserService interface {
	CreateUser(ctx context.Context, req dto.CreateUserRequest) (res dto.UserResponse, err error)
}

type UserServiceImpl struct {
	repo    repository.UserRepo
	cartSvc cartSvc.CartService
}

func ProvideUserServiceImpl(repo repository.UserRepository, cartSvc cartSvc.CartService) *UserServiceImpl {
	return &UserServiceImpl{
		repo:    repo,
		cartSvc: cartSvc,
	}
}

func (s *UserServiceImpl) CreateUser(ctx context.Context, req dto.CreateUserRequest) (res dto.UserResponse, err error) {
	user := req.ToModel()

	err = s.repo.CreateUser(ctx, &user)
	if err != nil {
		log.Error().Err(err).Msg("[CreateUser] Failed CreateUser")
		return
	}

	_, err = s.cartSvc.CreateCart(ctx, cartDto.CreateCartRequest{
		UserID: user.ID.String(),
	})
	if err != nil {
		log.Error().Err(err).Msg("[CreateUser] Failed CreateCart")
		return
	}

	return dto.NewUserResponse(user), nil
}
