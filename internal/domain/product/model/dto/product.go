package dto

import (
	"errors"
	"time"

	"github.com/azka-zaydan/synapsis-test/internal/domain/product/model"
	"github.com/guregu/null"
)

type ProductJSONField string

const (
	Id            ProductJSONField = "id"
	CategoryID    ProductJSONField = "categoryId"
	Name          ProductJSONField = "name"
	Description   ProductJSONField = "description"
	Price         ProductJSONField = "price"
	Stock         ProductJSONField = "stock"
	CreatedBy     ProductJSONField = "createdBy"
	MetaCreatedAt ProductJSONField = "metaCreatedAt"
	UpdatedBy     ProductJSONField = "updatedBy"
	MetaUpdatedAt ProductJSONField = "metaUpdatedAt"
	DeletedBy     ProductJSONField = "deletedBy"
	MetaDeletedAt ProductJSONField = "metaDeletedAt"
)

type ProductResponse struct {
	ID            string      `json:"id"`
	CategoryID    string      `json:"categoryId"`
	Name          string      `json:"name"`
	Description   string      `json:"description"`
	Price         float64     `json:"price"`
	Stock         int         `json:"stock"`
	CreatedBy     string      `json:"createdBy"`
	MetaCreatedAt time.Time   `json:"metaCreatedAt"`
	UpdatedBy     string      `json:"updatedBy"`
	MetaUpdatedAt time.Time   `json:"metaUpdatedAt"`
	DeletedBy     null.String `json:"deletedBy"`
	MetaDeletedAt null.Time   `json:"metaDeletedAt"`
}

type ProductFilterResponse struct {
	Data     []ProductResponse `json:"data"`
	Metadata Metadata          `json:"metadata"`
}

func NewProductResponse(prod model.Product) ProductResponse {
	return ProductResponse{
		ID:            prod.ID.String(),
		CategoryID:    prod.CategoryID.String(),
		Name:          prod.Name,
		Description:   prod.Description,
		Price:         prod.Price,
		Stock:         prod.Stock,
		CreatedBy:     prod.CreatedBy,
		MetaCreatedAt: prod.MetaCreatedAt,
		UpdatedBy:     prod.UpdatedBy,
		MetaUpdatedAt: prod.MetaUpdatedAt,
		DeletedBy:     prod.DeletedBy,
		MetaDeletedAt: prod.MetaDeletedAt,
	}
}

func TransformProductList(data []model.Product) (res []ProductResponse) {
	for _, v := range data {
		res = append(res, NewProductResponse(v))
	}
	return
}

func NewProductFilterResponse(data []model.Product, page, pageSize, totalData int) ProductFilterResponse {
	var totalPage int
	if totalData > 0 {
		totalPage = totalData / pageSize
	}
	dataResp := TransformProductList(data)
	return ProductFilterResponse{
		Data: dataResp,
		Metadata: Metadata{
			Page:      page,
			PageSize:  pageSize,
			TotalData: totalData,
			TotalPage: totalPage,
		},
	}
}

type Metadata struct {
	Page      int `json:"page"`
	PageSize  int `json:"pageSize"`
	TotalData int `json:"totalData"`
	TotalPage int `json:"totalPage"`
}

func ValidateAndSetDefaultFilter(filter *model.Filter) (err error) {
	if len(filter.FilterField) > 0 {
		for _, v := range filter.FilterField {
			err := validateFilterField(v.Field)
			if err != nil {
				return err
			}
			err = validateOperator(v.Operator)
			if err != nil {
				return err
			}
		}
	}

	if filter.Page == 0 {
		filter.Page = 1
	}

	if filter.PageSize == 0 {
		filter.PageSize = 10
	}

	return nil
}

func validateFilterField(field string) (err error) {
	switch field {
	case string(Id):
	case string(CategoryID):
	case string(Name):
	case string(Description):
	case string(Price):
	case string(Stock):
	case string(CreatedBy):
	case string(MetaCreatedAt):
	case string(UpdatedBy):
	case string(MetaUpdatedAt):
	case string(DeletedBy):
	case string(MetaDeletedAt):
	default:
		return errors.New("invalid filter field: " + field)
	}
	return
}

func TransformToDBField(filter *model.Filter) (err error) {
	if len(filter.FilterField) > 0 {
		for i, v := range filter.FilterField {
			switch v.Field {
			case string(Id):
				filter.FilterField[i].Field = string(model.Id)
			case string(CategoryID):
				filter.FilterField[i].Field = string(model.CategoryID)
			case string(Name):
				filter.FilterField[i].Field = string(model.Name)
			case string(Description):
				filter.FilterField[i].Field = string(model.Description)
			case string(Price):
				filter.FilterField[i].Field = string(model.Price)
			case string(Stock):
				filter.FilterField[i].Field = string(model.Stock)
			case string(CreatedBy):
				filter.FilterField[i].Field = string(model.CreatedBy)
			case string(MetaCreatedAt):
				filter.FilterField[i].Field = string(model.MetaCreatedAt)
			case string(UpdatedBy):
				filter.FilterField[i].Field = string(model.UpdatedBy)
			case string(MetaUpdatedAt):
				filter.FilterField[i].Field = string(model.MetaUpdatedAt)
			case string(DeletedBy):
				filter.FilterField[i].Field = string(model.DeletedBy)
			case string(MetaDeletedAt):
				filter.FilterField[i].Field = string(model.MetaDeletedAt)
			default:
				return errors.New("invalid filter field: " + v.Field)
			}
		}
	}
	return nil
}

func validateOperator(operator string) error {
	switch operator {
	case model.OperatorNot:
		return nil
	case model.OperatorEq:
		return nil
	case model.OperatorLike:
		return nil
	default:
		return errors.New("invalid operator")
	}
}
