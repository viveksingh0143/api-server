package domain

import (
	"errors"
	"time"
)

type InventoryType string

const (
	STOCK_IN       InventoryType = "STOCK IN"
	STOCK_OUT      InventoryType = "STOCK OUT"
	STOCK_RESERVED InventoryType = "STOCK RESERVED"
)

type Inventory struct {
	ID         int64         `json:"id"`
	Status     InventoryType `json:"status"`
	PalletID   *int64        `json:"pallet_id"`
	BinID      *int64        `json:"bin_id"`
	RackID     *int64        `json:"rack_id"`
	StoreID    *int64        `json:"store_id"`
	ProductID  int64         `json:"product_id"`
	Batch      string        `json:"batch"`
	Machine    string        `json:"machine"`
	Shift      string        `json:"shift"`
	Supervisor string        ``
	Quantity   float64       `json:"quantity"`
	Unit       string        `json:"unit"`
	StockInAt  time.Time     `json:"stockin_at"`
	StockOutAt *time.Time    `json:"stockout_at"`
}

func NewInventoryWithDefaults() *Inventory {
	return &Inventory{
		Status: STOCK_IN,
	}
}

func (i *Inventory) ValidateStatus() error {
	switch i.Status {
	case STOCK_IN, STOCK_OUT, STOCK_RESERVED:
		return nil
	default:
		return errors.New("invalid inventory status")
	}
}

type InventoryFormRawMaterial struct {
	Product_id string `json:"product_id"`
	Quantity   int64  `json:"quantity"`
	Pallet     string `json:"pallet"`
}
