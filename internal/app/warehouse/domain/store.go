package domain

import (
	"time"

	"github.com/vamika-digital/wms-api-server/internal/utility/customtypes"
)

type StoreOwner struct {
	ID      customtypes.NullableInt64  `json:"id"`
	Name    customtypes.NullableString `json:"name"`
	StaffID customtypes.NullableString `json:"staff_id"`
	Email   customtypes.NullableString `json:"email"`
}

type Store struct {
	ID            int64                      `json:"id"`
	Name          customtypes.NullableString `json:"name"`
	Location      customtypes.NullableString `json:"location"`
	Status        customtypes.NullableString `json:"status"`
	CreatedAt     time.Time                  `json:"created_at"`
	UpdatedAt     time.Time                  `json:"updated_at"`
	LastUpdatedBy customtypes.NullableString `json:"last_updated_by"`
	Owner         *StoreOwner                `json:"owner"`
}

func NewStoreWithDefaults() *Store {
	return &Store{
		Status: "active",
		Owner:  &StoreOwner{},
	}
}
