package repository

import (
	"context"
	"fmt"

	"github.com/azka-zaydan/synapsis-test/infras"
	"github.com/azka-zaydan/synapsis-test/internal/domain/order/model"
	"github.com/azka-zaydan/synapsis-test/shared/logger"
)

type OrderRepository interface {
	CreateOrder(ctx context.Context, order *model.Order) (err error)
	CreateOrderDetail(ctx context.Context, detail *model.OrderDetail) (err error)
	GetOrderByID(ctx context.Context, orderId string) (res model.Order, err error)
	UpdateOrder(ctx context.Context, order *model.Order) (err error)
	UpdateOrderDetail(ctx context.Context, order *model.OrderDetail) (err error)
	GetOrderDetailByID(ctx context.Context, orderDetailId string) (res model.Order, err error)
}

type OrderRepositoryMySQL struct {
	DB *infras.MySQLConn
}

func ProvideOrderRepositoryMySQL(db *infras.MySQLConn) *OrderRepositoryMySQL {
	return &OrderRepositoryMySQL{
		DB: db,
	}
}

func (repo *OrderRepositoryMySQL) CreateOrder(ctx context.Context, order *model.Order) (err error) {
	_, err = repo.DB.Write.NamedExecContext(ctx, orderInsertQuery, order)
	if err != nil {
		logger.ErrorWithStack(err)
		return
	}
	return
}

func (repo *OrderRepositoryMySQL) CreateOrderDetail(ctx context.Context, detail *model.OrderDetail) (err error) {
	_, err = repo.DB.Write.NamedExecContext(ctx, orderDetailInsertQuery, detail)
	if err != nil {
		logger.ErrorWithStack(err)
		return
	}
	return
}

func (repo *OrderRepositoryMySQL) GetOrderByID(ctx context.Context, orderId string) (res model.Order, err error) {
	err = repo.DB.Read.GetContext(ctx, &res, fmt.Sprintf("%s WHERE id = ?", orderSelectQuery), orderId)
	if err != nil {
		logger.ErrorWithStack(err)
		return
	}
	return
}

func (repo *OrderRepositoryMySQL) UpdateOrder(ctx context.Context, order *model.Order) (err error) {
	_, err = repo.DB.Write.NamedExecContext(ctx, orderUpdateQuery, order)
	if err != nil {
		logger.ErrorWithStack(err)
		return
	}
	return
}

func (repo *OrderRepositoryMySQL) UpdateOrderDetail(ctx context.Context, order *model.OrderDetail) (err error) {
	_, err = repo.DB.Write.NamedExecContext(ctx, orderDetailUpdateQuery, order)
	if err != nil {
		logger.ErrorWithStack(err)
		return
	}
	return
}

func (repo *OrderRepositoryMySQL) GetOrderDetailByID(ctx context.Context, orderDetailId string) (res model.Order, err error) {
	err = repo.DB.Read.GetContext(ctx, &res, fmt.Sprintf("%s WHERE id = ?", orderDetailSelectQuery), orderDetailId)
	if err != nil {
		logger.ErrorWithStack(err)
		return
	}
	return
}

var (
	orderInsertQuery = "INSERT INTO `order` (id,user_id,payment_id,total_price,status,order_at,payment_at,completed_at,created_by,updated_by) VALUES (:id,:user_id,:payment_id,:total_price,:status,:order_at,:payment_at,:completed_at,:created_by,:updated_by)"

	orderDetailInsertQuery = `
	INSERT INTO order_detail (
		id,
		order_id,
		product_id,
		total_items,
		subtotal_product_price,
		created_by,
		updated_by
	) VALUES (
		:id,
		:order_id,
		:product_id,
		:total_items,
		:subtotal_product_price,
		:created_by,
		:updated_by
	)`

	orderSelectQuery       = "SELECT id, user_id, payment_id, total_price, status, order_at, payment_at, completed_at, created_by, updated_by FROM `order`"
	orderDetailSelectQuery = `
	SELECT
		id,
		order_id,
		product_id,
		total_items,
		subtotal_product_price,
		created_by,
		updated_by
	FROM order_detail
`
	orderUpdateQuery = "UPDATE `order` SET user_id = :user_id, payment_id = :payment_id, total_price = :total_price, status = :status, order_at = :order_at, payment_at = :payment_at, completed_at = :completed_at, updated_by = :updated_by WHERE id = :id"

	orderDetailUpdateQuery = `
	UPDATE order_detail SET
		order_id = :order_id,
		product_id = :product_id,
		total_items = :total_items,
		subtotal_product_price = :subtotal_product_price,
		updated_by = :updated_by
	WHERE id = :id
`
)
