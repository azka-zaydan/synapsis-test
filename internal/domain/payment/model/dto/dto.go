package dto

import (
	"time"

	"github.com/azka-zaydan/synapsis-test/internal/domain/payment/model"
	"github.com/gofrs/uuid"
	"github.com/guregu/null"
)

type PaymentResponse struct {
	ID            string      `json:"id"`
	UserID        string      `json:"userId"`
	PaymentMethod string      `json:"paymentMethod"`
	OrderID       string      `json:"orderId"`
	TotalPrice    float64     `json:"totalPrice"`
	Status        int         `json:"status"`
	PaymentAt     null.Time   `json:"paymentAt"`
	CreatedBy     string      `json:"createdBy"`
	MetaCreatedAt time.Time   `json:"metaCreatedAt"`
	UpdatedBy     string      `json:"updatedBy"`
	MetaUpdatedAt time.Time   `json:"metaUpdatedAt"`
	DeletedBy     null.String `json:"deletedBy"`
	MetaDeletedAt null.Time   `json:"metaDeletedAt"`
}

type PayRequest struct {
	OrderID string `json:"orderId"`
}

type CreatePaymentRequest struct {
	UserID        string  `json:"user_id"`
	PaymentMethod string  `json:"payment_method"`
	OrderID       string  `json:"order_id"`
	TotalPrice    float64 `json:"total_price"`
}

func (d *CreatePaymentRequest) ToModel() (res model.Payment, err error) {
	id, err := uuid.NewV4()
	if err != nil {
		return
	}
	userId, err := uuid.FromString(d.UserID)
	if err != nil {
		return
	}
	orderId, err := uuid.FromString(d.OrderID)
	if err != nil {
		return
	}
	return model.Payment{
		ID:            id,
		UserID:        userId,
		OrderID:       orderId,
		PaymentMethod: d.PaymentMethod,
		TotalPrice:    d.TotalPrice,
		Status:        int(model.Unpaid),
		CreatedBy:     userId,
		UpdatedBy:     userId,
	}, nil
}

func NewPaymentResponse(payment model.Payment) PaymentResponse {
	return PaymentResponse{
		ID:            payment.ID.String(),
		UserID:        payment.UserID.String(),
		PaymentMethod: payment.PaymentMethod,
		OrderID:       payment.OrderID.String(),
		TotalPrice:    payment.TotalPrice,
		Status:        payment.Status,
		PaymentAt:     payment.PaymentAt,
		CreatedBy:     payment.CreatedBy.String(),
		MetaCreatedAt: payment.MetaCreatedAt,
		UpdatedBy:     payment.UpdatedBy.String(),
		MetaUpdatedAt: payment.MetaUpdatedAt,
	}
}
