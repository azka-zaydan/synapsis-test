package repository

import (
	"context"
	"fmt"

	"github.com/azka-zaydan/synapsis-test/infras"
	"github.com/azka-zaydan/synapsis-test/internal/domain/cart/model"
	"github.com/azka-zaydan/synapsis-test/shared/logger"
)

type CartRepository interface {
	CreateCart(ctx context.Context, cart *model.Cart) (err error)
	GetCartByUserID(ctx context.Context, userId string) (res model.Cart, err error)
	GetCartItemsByCartID(ctx context.Context, cartId string) (res []model.CartItem, err error)
	CreateCartItem(ctx context.Context, item *model.CartItem) (err error)
	UpdateCart(ctx context.Context, cart *model.Cart) (err error)
	UpdateCartItem(ctx context.Context, item *model.CartItem) (err error)
	DeleteCartItem(ctx context.Context, cartID string) (err error)
}

type CartRepositoryMySQL struct {
	DB *infras.MySQLConn
}

func ProvideCartRepositoryMySQL(db *infras.MySQLConn) *CartRepositoryMySQL {
	return &CartRepositoryMySQL{
		DB: db,
	}
}

func (repo *CartRepositoryMySQL) CreateCart(ctx context.Context, cart *model.Cart) (err error) {
	_, err = repo.DB.Write.NamedExecContext(ctx, cartInsertQuery, cart)
	if err != nil {
		logger.ErrorWithStack(err)
		return
	}
	return
}

func (repo *CartRepositoryMySQL) GetCartByUserID(ctx context.Context, userId string) (res model.Cart, err error) {
	err = repo.DB.Read.GetContext(ctx, &res, fmt.Sprintf("%s WHERE user_id = ?", cartSelectQuery), userId)
	if err != nil {
		logger.ErrorWithStack(err)
		return
	}
	return
}

func (repo *CartRepositoryMySQL) GetCartItemsByCartID(ctx context.Context, cartId string) (res []model.CartItem, err error) {
	err = repo.DB.Read.SelectContext(ctx, &res, cartItemSelectQuery, cartId)
	if err != nil {
		logger.ErrorWithStack(err)
		return
	}
	return
}

func (repo *CartRepositoryMySQL) CreateCartItem(ctx context.Context, item *model.CartItem) (err error) {
	_, err = repo.DB.Write.NamedExecContext(ctx, cartItemInsertQuery, item)
	if err != nil {
		logger.ErrorWithStack(err)
		return
	}
	return
}

func (repo *CartRepositoryMySQL) UpdateCartItem(ctx context.Context, item *model.CartItem) (err error) {
	_, err = repo.DB.Write.NamedExecContext(ctx, cartItemUpdateQuery, item)
	if err != nil {
		logger.ErrorWithStack(err)
		return
	}
	return
}

func (repo *CartRepositoryMySQL) UpdateCart(ctx context.Context, cart *model.Cart) (err error) {
	_, err = repo.DB.Write.NamedExecContext(ctx, cartUpdateQuery, cart)
	if err != nil {
		logger.ErrorWithStack(err)
		return
	}
	return
}

func (repo *CartRepositoryMySQL) DeleteCartItem(ctx context.Context, cartID string) (err error) {
	_, err = repo.DB.Write.ExecContext(ctx, cartItemDeleteQuery, cartID)
	if err != nil {
		logger.ErrorWithStack(err)
		return
	}
	return
}

var (
	cartInsertQuery = `
	INSERT INTO cart (id, user_id, total_items, total_price, created_by, updated_by)
	VALUES (:id, :user_id, :total_items, :total_price, :created_by, :updated_by)`
	cartSelectQuery = `
	SELECT id, user_id, total_items, total_price, created_by, meta_created_at, updated_by, meta_updated_at, deleted_by, meta_deleted_at
	FROM cart`
	cartItemSelectQuery = `
	SELECT id, cart_id, product_id, quantity, total_price, created_by, meta_created_at, updated_by, meta_updated_at, deleted_by, meta_deleted_at
	FROM cart_item
	WHERE cart_id = ?`
	cartItemInsertQuery = `
	INSERT INTO cart_item (id, cart_id, product_id, quantity, total_price, created_by, updated_by)
	VALUES (:id, :cart_id, :product_id, :quantity, :total_price, :created_by, :updated_by)`
	cartUpdateQuery = `
	UPDATE cart 
	SET user_id = :user_id, total_items = :total_items, total_price = :total_price, updated_by = :updated_by
	WHERE id = :id`
	cartItemUpdateQuery = `
	UPDATE cart_item 
	SET cart_id = :cart_id, product_id = :product_id, quantity = :quantity, total_price = :total_price, updated_by = :updated_by
	WHERE id = :id`
	cartItemDeleteQuery = `
	DELETE FROM cart_item 
	WHERE id = ?`
)
