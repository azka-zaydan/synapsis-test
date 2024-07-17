package dto

import (
	"time"

	"github.com/azka-zaydan/synapsis-test/internal/domain/user/model"
	"github.com/gofrs/uuid"
	"github.com/guregu/null"
)

type CreateUserRequest struct {
	Username string `json:"username" validate:"required"`
	Password string `json:"password" validate:"required"`
}

type UserResponse struct {
	ID            string      `json:"id"`
	Username      string      `json:"username"`
	Password      string      `json:"passwordHash"`
	MetaCreatedAt time.Time   `json:"metaCreatedAt"`
	MetaUpdatedAt time.Time   `json:"metaUpdatedAt"`
	MetaDeletedAt null.Time   `json:"metaDeletedAt"`
	CreatedBy     string      `json:"createdBy"`
	UpdatedBy     string      `json:"updatedBy"`
	DeletedBy     null.String `json:"deletedBy"`
	CartID        null.String `json:"cartId,omitempty"`
}

func NewUserResponse(user model.User) UserResponse {
	return UserResponse{
		ID:            user.ID.String(),
		Username:      user.Username,
		Password:      user.Password,
		MetaCreatedAt: user.MetaCreatedAt,
		MetaUpdatedAt: user.MetaUpdatedAt,
		MetaDeletedAt: user.MetaDeletedAt,
		CreatedBy:     user.CreatedBy.String(),
		UpdatedBy:     user.UpdatedBy.String(),
		DeletedBy:     null.NewString(user.DeletedBy.UUID.String(), user.DeletedBy.Valid),
		CartID:        null.NewString(user.CartID.UUID.String(), user.CartID.Valid),
	}
}

func (d *CreateUserRequest) ToModel() model.User {
	id, _ := uuid.NewV4()

	return model.User{
		ID:        id,
		Username:  d.Username,
		Password:  d.Password,
		CreatedBy: id,
		UpdatedBy: id,
	}
}
