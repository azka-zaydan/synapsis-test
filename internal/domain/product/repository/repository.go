package repository

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/azka-zaydan/synapsis-test/infras"
	"github.com/azka-zaydan/synapsis-test/internal/domain/product/model"
	"github.com/azka-zaydan/synapsis-test/shared/failure"
	"github.com/azka-zaydan/synapsis-test/shared/logger"
	"github.com/rs/zerolog/log"
)

type ProductRepository interface {
	GetProductByFilter(ctx context.Context, filter *model.Filter) (res []model.Product, totalData int, err error)
	CreateProduct(ctx context.Context, data *model.Product) (err error)
}

type ProductRepositoryMySQL struct {
	DB *infras.MySQLConn
}

func ProvideProductRepositoryMySQL(db *infras.MySQLConn) *ProductRepositoryMySQL {
	return &ProductRepositoryMySQL{
		DB: db,
	}
}

func (repo *ProductRepositoryMySQL) GetProductByFilter(ctx context.Context, filter *model.Filter) (res []model.Product, totalData int, err error) {
	query, args, err := repo.buildSQLQuery(countProductQuery, filter)
	if err != nil {
		err = failure.BadRequest(err)
		log.Error().Err(err).Msg("[GetProductByFilter] failed buildSQLQuery")
		return
	}

	err = repo.DB.Read.GetContext(ctx, &totalData, query, args...)
	if err != nil {
		log.Error().Err(err).Msg("[GetProductByFilter] failed counting total data")
		return
	}

	query, args, err = repo.buildSQLQuery(productSelectQuery, filter)
	if err != nil {
		err = failure.BadRequest(err)
		log.Error().Err(err).Msg("[GetProductByFilter] failed buildSQLQuery")
		return
	}

	err = repo.DB.Read.SelectContext(ctx, &res, query, args...)
	if err != nil {
		log.Error().Err(err).Msg("[GetProductByFilter] failed getting data")
		return
	}
	return
}

func (repo *ProductRepositoryMySQL) CreateProduct(ctx context.Context, data *model.Product) (err error) {
	_, err = repo.DB.Write.NamedExecContext(ctx, productInsertQuery, data)
	if err != nil {
		logger.ErrorWithStack(err)
		return
	}
	return
}

func (repo *ProductRepositoryMySQL) buildSQLQuery(baseQuery string, filter *model.Filter) (string, []interface{}, error) {
	var conditions []string
	var args []interface{}

	for _, f := range filter.FilterField {
		var condition string
		switch f.Operator {
		case "eq":
			condition = fmt.Sprintf("%s = ?", f.Field)
		case "not":
			condition = fmt.Sprintf("%s != ?", f.Field)
		case "like":
			condition = fmt.Sprintf("%s LIKE ?", f.Field)
			f.Value = fmt.Sprintf("%%%v%%", f.Value)
		default:
			return "", nil, errors.New("invalid operator: " + f.Operator)
		}
		conditions = append(conditions, condition)
		args = append(args, f.Value)
	}

	if len(conditions) > 0 {
		baseQuery = fmt.Sprintf("%s WHERE %s", baseQuery, strings.Join(conditions, " AND "))
	}

	if filter.PageSize > 0 {
		offset := (filter.Page - 1) * filter.PageSize
		baseQuery = fmt.Sprintf("%s LIMIT %d OFFSET %d", baseQuery, filter.PageSize, offset)
	}

	return baseQuery, args, nil
}

var (
	productSelectQuery = `SELECT 
        id, 
        category_id, 
        name, 
        description, 
        price, 
        stock, 
        created_by, 
        meta_created_at, 
        updated_by, 
        meta_updated_at, 
        deleted_by, 
        meta_deleted_at 
    FROM product `

	countProductQuery = `SELECT
        COUNT(id)
    FROM product`

	productInsertQuery = `
	INSERT INTO product (
		id, 
		category_id, 
		name, 
		description, 
		price, 
		stock, 
		created_by, 
		updated_by
	) VALUES (
		:id, 
		:category_id, 
		:name, 
		:description, 
		:price, 
		:stock, 
		:created_by, 
		:updated_by
	)`
)
