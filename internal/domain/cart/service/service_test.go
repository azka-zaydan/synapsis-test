package service_test

import (
	"context"
	"testing"

	"github.com/azka-zaydan/synapsis-test/configs"
	"github.com/azka-zaydan/synapsis-test/internal/domain/cart/model/dto"
	cart_repo_mock "github.com/azka-zaydan/synapsis-test/internal/domain/cart/repository/mock"
	"github.com/azka-zaydan/synapsis-test/internal/domain/cart/service"

	"github.com/go-redis/redismock/v9"
	"github.com/gofrs/uuid"
	"github.com/golang/mock/gomock"
	"github.com/rs/zerolog/log"
	"github.com/stretchr/testify/assert"
)

func TestCreateCart(t *testing.T) {
	id, _ := uuid.NewV4()
	userId, _ := uuid.NewV4()
	testCases := []struct {
		name           string
		expectedResult dto.CartResponse
		expectedError  error
		body           dto.CreateCartRequest
		mockRepo       func(repo *cart_repo_mock.MockCartRepository)
		redisMock      func(redis *redismock.ClientMock)
	}{
		{
			name: "Success",
			expectedResult: dto.CartResponse{
				ID:         id.String(),
				UserID:     userId.String(),
				TotalItems: 0,
				TotalPrice: 0,
				CreatedBy:  userId.String(),
				UpdatedBy:  userId.String(),
			},
			expectedError: nil,
			body: dto.CreateCartRequest{
				UserID: userId.String(),
			},
			mockRepo: func(repo *cart_repo_mock.MockCartRepository) {
				repo.EXPECT().CreateCart(gomock.Any(), gomock.Any()).Return(nil).AnyTimes()
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			mockRepo := cart_repo_mock.NewMockCartRepository(ctrl)
			// redis, redisMock := redismock.NewClientMock()

			// tc.redisMock(&redisMock)
			tc.mockRepo(mockRepo)
			cfg := &configs.Config{}
			serviceImpl := service.CartServiceImpl{
				Repo:   mockRepo,
				Config: cfg,
			}
			res, err := serviceImpl.CreateCart(context.Background(), tc.body)
			if err != nil {
				log.Err(err).Msg(err.Error())
				t.Fail()
			}
			res.ID = id.String()
			assert.Equal(t, tc.expectedResult, res)
		})
	}
}
