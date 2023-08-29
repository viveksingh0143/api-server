package domain

import (
	"errors"
	"time"

	"github.com/vamika-digital/wms-api-server/internal/utility/customtypes"
)

type ProductType string

const (
	RAW_MATERIAL_TYPE        ProductType = "RAW Material"
	FINISHED_GOODS_TYPE      ProductType = "Finished Goods"
	SEMI_FINISHED_GOODS_TYPE ProductType = "Semi Finished Goods"
)

type Product struct {
	ID            int64                      `json:"id"`
	Type          ProductType                `json:"type"`
	Code          customtypes.NullableString `json:"code"`
	RawCode       customtypes.NullableString `json:"raw_code"`
	Name          customtypes.NullableString `json:"name"`
	Description   customtypes.NullableString `json:"description"`
	Unit          customtypes.NullableString `json:"unit"`
	Status        customtypes.NullableString `json:"status"`
	CreatedAt     time.Time                  `json:"created_at"`
	UpdatedAt     time.Time                  `json:"updated_at"`
	LastUpdatedBy customtypes.NullableString `json:"last_updated_by"`
}

func NewProductWithDefaults() Product {
	return Product{
		Type:   RAW_MATERIAL_TYPE,
		Status: "active",
	}
}

func (p *Product) ValidateType() error {
	switch p.Type {
	case RAW_MATERIAL_TYPE, SEMI_FINISHED_GOODS_TYPE, FINISHED_GOODS_TYPE:
		return nil
	default:
		return errors.New("invalid product type")
	}
}
