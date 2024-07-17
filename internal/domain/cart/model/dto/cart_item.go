package dto

import (
	"time"

	"github.com/azka-zaydan/synapsis-test/internal/domain/cart/model"
	"github.com/gofrs/uuid"
	"github.com/guregu/null"
)

type CartItemResponse struct {
	ID            string      `json:"id"`
	CartID        string      `json:"cartId"`
	ProductID     string      `json:"productId"`
	TotalPrice    float64     `json:"totalPrice"`
	Quantity      int         `json:"quantity"`
	CreatedBy     string      `json:"createdBy"`
	MetaCreatedAt time.Time   `json:"metaCreatedAt"`
	UpdatedBy     string      `json:"updatedBy"`
	MetaUpdatedAt time.Time   `json:"metaUpdatedAt"`
	DeletedBy     null.String `json:"deletedBy"`
	MetaDeletedAt null.Time   `json:"metaDeletedAt"`
}

func NewCartItemResponse(item model.CartItem) CartItemResponse {
	return CartItemResponse{
		ID:            item.ID.String(),
		CartID:        item.CartID.String(),
		ProductID:     item.ProductID.String(),
		TotalPrice:    item.TotalPrice,
		Quantity:      item.Quantity,
		CreatedBy:     item.CreatedBy.String(),
		MetaCreatedAt: item.MetaCreatedAt,
		UpdatedBy:     item.UpdatedBy.String(),
		MetaUpdatedAt: item.MetaUpdatedAt,
		DeletedBy:     item.DeletedBy,
		MetaDeletedAt: item.MetaDeletedAt,
	}
}

type AddItemsRequest []ItemRequest

type DeleteItemsRequest []DeleteItemRequest

type ItemRequest struct {
	ProductID string  `json:"productId"`
	Quantity  int     `json:"quantity"`
	Price     float64 `json:"price"`
}

type CartItemCreateRequest struct {
	CartID    string  `json:"cartId"`
	ProductID string  `json:"productId"`
	Quantity  int     `json:"quantity"`
	Price     float64 `json:"price"`
}

type DeleteItemRequest struct {
	ItemId   string  `json:"itemId"`
	Quantity int     `json:"quantity"`
	Price    float64 `json:"price"`
}

func (d *CartItemCreateRequest) ToModel() (res model.CartItem, err error) {
	id, err := uuid.NewV4()
	if err != nil {
		return
	}
	cartId, err := uuid.FromString(d.CartID)
	if err != nil {
		return
	}
	productId, err := uuid.FromString(d.ProductID)
	if err != nil {
		return
	}
	return model.CartItem{
		ID:         id,
		CartID:     cartId,
		ProductID:  productId,
		Quantity:   d.Quantity,
		TotalPrice: d.Price,
		CreatedBy:  cartId,
		UpdatedBy:  cartId,
	}, nil
}
