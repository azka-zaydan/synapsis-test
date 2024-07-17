package repository

import (
	"context"

	"github.com/azka-zaydan/synapsis-test/infras"
	"github.com/azka-zaydan/synapsis-test/internal/domain/order/model"
)

type OrderRepository interface {
	CreateOrder(ctx context.Context, Order *model.Order) (err error)
	GetOrderByOrderID(ctx context.Context, orderId string) (res model.Order, err error)
	UpdateOrder(ctx context.Context, Order *model.Order) (err error)
}

type OrderRepositoryMySQL struct {
	DB *infras.MySQLConn
}

func ProvideOrderRepositoryMySQL(db *infras.MySQLConn) *OrderRepositoryMySQL {
	return &OrderRepositoryMySQL{
		DB: db,
	}
}

func (repo *OrderRepositoryMySQL) CreateOrder(ctx context.Context, Order *model.Order) (err error) {
	return
}

func (repo *OrderRepositoryMySQL) GetOrderByOrderID(ctx context.Context, orderId string) (res model.Order, err error) {
	return
}

func (repo *OrderRepositoryMySQL) UpdateOrder(ctx context.Context, Order *model.Order) (err error) {
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
)
