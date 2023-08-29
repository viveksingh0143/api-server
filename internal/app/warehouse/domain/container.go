package domain

import (
	"errors"
	"time"

	"github.com/vamika-digital/wms-api-server/internal/utility/customtypes"
)

type ContainerType string

const (
	PALLET_TYPE ContainerType = "PALLET"
	BIN_TYPE    ContainerType = "BIN"
	RACK_TYPE   ContainerType = "RACK"
)

type Container struct {
	ID            int64                      `json:"id"`
	Type          ContainerType              `json:"type"`
	Code          customtypes.NullableString `json:"code"`
	Name          customtypes.NullableString `json:"name"`
	Address       customtypes.NullableString `json:"address"`
	Status        customtypes.NullableString `json:"status"`
	CreatedAt     time.Time                  `json:"created_at"`
	UpdatedAt     time.Time                  `json:"updated_at"`
	LastUpdatedBy customtypes.NullableString `json:"last_updated_by"`
}

func NewContainerWithDefaults() Container {
	return Container{
		Type:   PALLET_TYPE,
		Status: "active",
	}
}

func (p *Container) ValidateType() error {
	switch p.Type {
	case PALLET_TYPE, BIN_TYPE, RACK_TYPE:
		return nil
	default:
		return errors.New("invalid product type")
	}
}
