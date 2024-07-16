package service

import (
	"context"
	"fmt"

	"github.com/azka-zaydan/synapsis-test/configs"
	"github.com/azka-zaydan/synapsis-test/infras"
	"github.com/azka-zaydan/synapsis-test/internal/domain/auth/model/dto"
	"github.com/azka-zaydan/synapsis-test/internal/domain/user/model"
	userRepo "github.com/azka-zaydan/synapsis-test/internal/domain/user/repository"
	"github.com/azka-zaydan/synapsis-test/shared/failure"
	"github.com/azka-zaydan/synapsis-test/shared/hash"
	"github.com/azka-zaydan/synapsis-test/shared/jwt"
	"github.com/gofiber/fiber/v2"
	"github.com/redis/go-redis/v9"
	"github.com/rs/zerolog/log"
)

type AuthService interface {
	Register(ctx context.Context, req dto.RegisterDto) (res dto.JWTResponse, err error)
	Login(ctx context.Context, req dto.RegisterDto) (res dto.JWTResponse, err error)
}

type AuthServiceImpl struct {
	Redis      *infras.Redis
	Config     *configs.Config
	UserRepo   userRepo.UserRepository
	JwtService *jwt.JwtService
}

func ProvideAuthServiceImpl(cfg *configs.Config, redis *infras.Redis, userRepo userRepo.UserRepository) *AuthServiceImpl {
	return &AuthServiceImpl{
		Config:     cfg,
		Redis:      redis,
		UserRepo:   userRepo,
		JwtService: jwt.NewJwtService(cfg),
	}
}

func (s *AuthServiceImpl) Register(ctx context.Context, req dto.RegisterDto) (res dto.JWTResponse, err error) {

	registered, err := s.isUserRegistered(ctx, req.Username)
	if err != nil {
		log.Error().Err(err).Msg("[Register] Failed isUserRegistered")
		return
	}

	if registered {
		err = failure.Conflict("register", "user", "already registered")
		log.Error().Err(err).Msg("[Register] User Already Registered")
		return
	}

	hashedPass, err := hash.HashPassword(req.Password)
	if err != nil {
		log.Error().Err(err).Msg("[Register] Failed Hash Password")
		return
	}
	user := model.NewUser(req.Username, hashedPass)

	err = s.UserRepo.CreateUser(ctx, &user)
	if err != nil {
		log.Error().Err(err).Msg("[Register] Failed CreateUser")
		return
	}
	token, err := s.JwtService.GenerateJWT(user.Username, user.ID.String())
	if err != nil {
		log.Error().Err(err).Msg("[Register] Failed Generate Token")
		return
	}
	s.Redis.Client.Set(ctx, fmt.Sprintf("token:{%s}", req.Username), token, s.Config.Cache.Token.ExpiresIn)
	return dto.NewJWTResponse(token), nil
}

func (s *AuthServiceImpl) isUserRegistered(ctx context.Context, username string) (registered bool, err error) {
	user, err := s.UserRepo.FindByUsername(ctx, username)
	if err != nil {
		if failure.GetCode(err) != fiber.StatusNotFound {
			log.Error().Err(err).Msg("[isUserRegistered] Failed FindByUsername")
			return
		}
	}
	if user.Username != "" {
		return true, nil
	}

	return false, nil
}

func (s *AuthServiceImpl) Login(ctx context.Context, req dto.RegisterDto) (res dto.JWTResponse, err error) {

	registered, err := s.isUserRegistered(ctx, req.Username)
	if err != nil {
		log.Error().Err(err).Msg("[Login] Failed isUserRegistered")
		return
	}
	if !registered {
		err = failure.NotFound("user")
		log.Error().Err(err).Msg("[Login] User Has Not Registered")
		return
	}

	token, err := s.Redis.Client.Get(ctx, fmt.Sprintf("token:{%s}", req.Username)).Result()
	if err != nil {
		if err != redis.Nil {
			log.Error().Err(err).Msg("[Login] Failed Getting From Redis")
			return
		}
	}

	if token != "" {
		log.Info().Msg("[Login] Using From Cache")
		return dto.NewJWTResponse(token), nil
	}

	user, err := s.UserRepo.FindByUsername(ctx, req.Username)
	if err != nil {
		log.Error().Err(err).Msg("[Login] Failed FindByUsername")
		return
	}
	valid := hash.CheckPasswordHash(req.Password, user.Password)
	if !valid {
		err = failure.BadRequestFromString("invalid password")
		log.Error().Err(err).Msg("[Login] Invalid Password")
		return
	}
	token, err = s.JwtService.GenerateJWT(user.Username, user.ID.String())
	if err != nil {
		log.Error().Err(err).Msg("[Login] Failed Generate Token")
		return
	}
	s.Redis.Client.Set(ctx, fmt.Sprintf("token:{%s}", req.Username), token, s.Config.Cache.Token.ExpiresIn)
	return dto.NewJWTResponse(token), nil
}
