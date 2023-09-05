package domain

import (
	"time"

	"github.com/vamika-digital/wms-api-server/internal/utility/customtypes"
)

type Machine struct {
	ID            int64                      `json:"id"`
	Code          customtypes.NullableString `json:"code"`
	Name          customtypes.NullableString `json:"name"`
	Status        customtypes.NullableString `json:"status"`
	CreatedAt     time.Time                  `json:"created_at"`
	UpdatedAt     time.Time                  `json:"updated_at"`
	LastUpdatedBy customtypes.NullableString `json:"last_updated_by"`
}

func NewMachineWithDefaults() *Machine {
	return &Machine{
		Status: "active",
	}
}
