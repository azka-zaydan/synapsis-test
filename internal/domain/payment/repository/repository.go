package repository

import (
	"context"
	"fmt"

	"github.com/azka-zaydan/synapsis-test/infras"
	"github.com/azka-zaydan/synapsis-test/internal/domain/payment/model"
	"github.com/azka-zaydan/synapsis-test/shared/logger"
)

//go:generate go run github.com/golang/mock/mockgen -source repository.go -destination mock/repository.go -package payment_repo_mock
type PaymentRepository interface {
	CreatePayment(ctx context.Context, payment *model.Payment) (err error)
	GetPaymentByOrderID(ctx context.Context, orderId string) (res model.Payment, err error)
	UpdatePayment(ctx context.Context, payment *model.Payment) (err error)
}

type PaymentRepositoryMySQL struct {
	DB *infras.MySQLConn
}

func ProvidePaymentRepositoryMySQL(db *infras.MySQLConn) *PaymentRepositoryMySQL {
	return &PaymentRepositoryMySQL{
		DB: db,
	}
}

func (repo *PaymentRepositoryMySQL) CreatePayment(ctx context.Context, payment *model.Payment) (err error) {
	_, err = repo.DB.Write.NamedExecContext(ctx, paymentInsertQuery, payment)
	if err != nil {
		logger.ErrorWithStack(err)
		return
	}
	return
}

func (repo *PaymentRepositoryMySQL) GetPaymentByOrderID(ctx context.Context, orderId string) (res model.Payment, err error) {
	err = repo.DB.Read.GetContext(ctx, &res, fmt.Sprintf("%s WHERE order_id = ?", paymentSelectQuery), orderId)
	if err != nil {
		logger.ErrorWithStack(err)
		return
	}
	return
}

func (repo *PaymentRepositoryMySQL) UpdatePayment(ctx context.Context, payment *model.Payment) (err error) {
	_, err = repo.DB.Write.NamedExecContext(ctx, paymentUpdateQuery, payment)
	if err != nil {
		logger.ErrorWithStack(err)
		return
	}
	return
}

var (
	paymentInsertQuery = `
	INSERT INTO payment (
		id,
		user_id,
		payment_method,
		order_id,
		total_price,
		status,
		payment_at,
		created_by,
		updated_by
	) VALUES (
		:id,
		:user_id,
		:payment_method,
		:order_id,
		:total_price,
		:status,
		:payment_at,
		:created_by,
		:updated_by
	)`
	paymentSelectQuery = `
	SELECT
		id,
		user_id,
		payment_method,
		order_id,
		total_price,
		status,
		payment_at,
		created_by,
		updated_by
	FROM payment`
	paymentUpdateQuery = `
	UPDATE payment SET
		user_id = :user_id,
		payment_method = :payment_method,
		order_id = :order_id,
		total_price = :total_price,
		status = :status,
		payment_at = :payment_at,
		updated_by = :updated_by
	WHERE id = :id
`
)
