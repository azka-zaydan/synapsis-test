package dto

import (
	"github.com/azka-zaydan/synapsis-test/internal/domain/product/model"
	"github.com/gofrs/uuid"
)

type ProductCreateRequest struct {
	CategoryID  string  `json:"categoryId"`
	Name        string  `json:"name"`
	Description string  `json:"description"`
	Price       float64 `json:"price"`
	Stock       int     `json:"stock"`
	CreatedBy   string  `json:"createdBy"`
}

func (d *ProductCreateRequest) ToModel() (res model.Product, err error) {
	id, err := uuid.NewV4()
	if err != nil {
		return
	}
	catId, err := uuid.FromString(d.CategoryID)
	if err != nil {
		return
	}
	return model.Product{
		ID:          id,
		CategoryID:  catId,
		Name:        d.Name,
		Description: d.Description,
		Price:       d.Price,
		Stock:       d.Stock,
		CreatedBy:   d.CreatedBy,
		UpdatedBy:   d.CreatedBy,
	}, nil
}
