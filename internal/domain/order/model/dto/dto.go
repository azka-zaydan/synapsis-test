package dto

import (
	"time"

	"github.com/azka-zaydan/synapsis-test/internal/domain/order/model"
	"github.com/gofrs/uuid"
	"github.com/guregu/null"
)

type OrderResponse struct {
	ID            uuid.UUID     `json:"id"`
	UserID        uuid.UUID     `json:"userId"`
	PaymentID     uuid.NullUUID `json:"paymentId,omitempty"`
	TotalPrice    float64       `json:"totalPrice"`
	Status        int           `json:"status"`
	OrderAt       time.Time     `json:"orderAt"`
	PaymentAt     null.Time     `json:"paymentAt"`
	CompletedAt   null.Time     `json:"completedAt"`
	CreatedBy     uuid.UUID     `json:"createdBy"`
	MetaCreatedAt time.Time     `json:"metaCreatedAt"`
	UpdatedBy     uuid.UUID     `json:"updatedBy"`
	MetaUpdatedAt time.Time     `json:"metaUpdatedAt"`
	DeletedBy     uuid.NullUUID `json:"deletedBy"`
	MetaDeletedAt null.Time     `json:"metaDeletedAt"`
}

type OrderDetailResponse struct {
	ID                   uuid.UUID     `json:"id"`
	OrderID              uuid.UUID     `json:"orderId"`
	ProductID            uuid.UUID     `json:"productId"`
	TotalItems           int           `json:"totalItems"`
	SubtotalProductPrice float64       `json:"subtotalProductPrice"`
	CreatedBy            uuid.UUID     `json:"createdBy"`
	MetaCreatedAt        time.Time     `json:"metaCreatedAt"`
	UpdatedBy            uuid.UUID     `json:"updatedBy"`
	MetaUpdatedAt        time.Time     `json:"metaUpdatedAt"`
	DeletedBy            uuid.NullUUID `json:"deletedBy"`
	MetaDeletedAt        null.Time     `json:"metaDeletedAt"`
}

type CreateOrderReqeust struct {
	UserID     string  `json:"userId"`
	TotalPrice float64 `json:"totalPrice"`
}

type CreateOrderDetailRequest struct {
	OrderID              string  `json:"orderId"`
	ProductID            string  `json:"productId"`
	TotalItems           int     `json:"totalItems"`
	SubtotalProductPrice float64 `json:"subtotalProductPrice"`
}

func (d *CreateOrderDetailRequest) ToModel() (res model.OrderDetail, err error) {
	id, err := uuid.NewV4()
	if err != nil {
		return
	}
	orderId, err := uuid.FromString(d.OrderID)
	if err != nil {
		return
	}
	productId, err := uuid.FromString(d.ProductID)
	if err != nil {
		return
	}

	return model.OrderDetail{
		ID:                   id,
		OrderID:              orderId,
		ProductID:            productId,
		TotalItems:           d.TotalItems,
		SubtotalProductPrice: d.SubtotalProductPrice,
		CreatedBy:            orderId,
		UpdatedBy:            orderId,
	}, nil
}

func (d *CreateOrderReqeust) ToModel() (res model.Order, err error) {
	id, err := uuid.NewV4()
	if err != nil {
		return
	}
	userId, err := uuid.FromString(d.UserID)
	if err != nil {
		return
	}
	return model.Order{
		ID:         id,
		UserID:     userId,
		TotalPrice: d.TotalPrice,
		Status:     int(model.OrderPlacedStatus),
		OrderAt:    time.Now(),
		CreatedBy:  userId,
		UpdatedBy:  userId,
	}, nil
}

func NewOrderResponse(order model.Order) OrderResponse {
	return OrderResponse{
		ID:            order.ID,
		UserID:        order.UserID,
		PaymentID:     order.PaymentID,
		TotalPrice:    order.TotalPrice,
		Status:        order.Status,
		OrderAt:       order.OrderAt,
		PaymentAt:     order.PaymentAt,
		CompletedAt:   order.CompletedAt,
		CreatedBy:     order.CreatedBy,
		MetaCreatedAt: order.MetaCreatedAt,
		UpdatedBy:     order.UpdatedBy,
		MetaUpdatedAt: order.MetaUpdatedAt,
		DeletedBy:     order.DeletedBy,
		MetaDeletedAt: order.MetaDeletedAt,
	}
}

func NewOrderDetailResponse(orderDetail model.OrderDetail) OrderDetailResponse {
	return OrderDetailResponse{
		ID:                   orderDetail.ID,
		OrderID:              orderDetail.OrderID,
		ProductID:            orderDetail.ProductID,
		TotalItems:           orderDetail.TotalItems,
		SubtotalProductPrice: orderDetail.SubtotalProductPrice,
		CreatedBy:            orderDetail.CreatedBy,
		MetaCreatedAt:        orderDetail.MetaCreatedAt,
		UpdatedBy:            orderDetail.UpdatedBy,
		MetaUpdatedAt:        orderDetail.MetaUpdatedAt,
		DeletedBy:            orderDetail.DeletedBy,
		MetaDeletedAt:        orderDetail.MetaDeletedAt,
	}
}
