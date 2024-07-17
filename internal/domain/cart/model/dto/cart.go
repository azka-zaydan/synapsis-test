package dto

import (
	"time"

	"github.com/azka-zaydan/synapsis-test/internal/domain/cart/model"
	"github.com/gofrs/uuid"
	"github.com/guregu/null"
)

type CartResponse struct {
	ID            string      `json:"id"`
	UserID        string      `json:"userId"`
	TotalItems    int         `json:"totalItems"`
	TotalPrice    float64     `json:"totalPrice"`
	CreatedBy     string      `json:"createdBy"`
	MetaCreatedAt time.Time   `json:"metaCreatedAt"`
	UpdatedBy     string      `json:"updatedBy"`
	MetaUpdatedAt time.Time   `json:"metaUpdatedAt"`
	DeletedBy     null.String `json:"deletedBy"`
	MetaDeletedAt null.Time   `json:"metaDeletedAt"`
}

type ListItemsResponse struct {
	Cart  CartResponse
	Items []CartItemResponse
}

func NewCartResponse(cart model.Cart) CartResponse {
	return CartResponse{
		ID:            cart.ID.String(),
		UserID:        cart.UserID.String(),
		TotalItems:    cart.TotalItems,
		TotalPrice:    cart.TotalPrice,
		CreatedBy:     cart.CreatedBy.String(),
		MetaCreatedAt: cart.MetaCreatedAt,
		UpdatedBy:     cart.UpdatedBy.String(),
		MetaUpdatedAt: cart.MetaUpdatedAt,
		DeletedBy:     cart.DeletedBy,
		MetaDeletedAt: cart.MetaDeletedAt,
	}
}

func NewListItemsResponse(cart model.Cart, items []model.CartItem) ListItemsResponse {
	var itemsRes []CartItemResponse
	cartRes := NewCartResponse(cart)

	for _, v := range items {
		itemsRes = append(itemsRes, NewCartItemResponse(v))
	}

	return ListItemsResponse{
		Cart:  cartRes,
		Items: itemsRes,
	}
}

type CreateCartRequest struct {
	UserID string `json:"userId"`
}

func (d *CreateCartRequest) ToModel() (res model.Cart, err error) {
	id, err := uuid.NewV4()
	if err != nil {
		return
	}
	userId, err := uuid.FromString(d.UserID)
	if err != nil {
		return
	}
	return model.Cart{
		ID:         id,
		UserID:     userId,
		TotalPrice: 0,
		TotalItems: 0,
		CreatedBy:  userId,
		UpdatedBy:  userId,
	}, nil
}
