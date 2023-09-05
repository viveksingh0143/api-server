package domain

import (
	"time"

	"github.com/vamika-digital/wms-api-server/internal/utility/customtypes"
)

type Address struct {
	Address1 customtypes.NullableString `json:"address1"`
	Address2 customtypes.NullableString `json:"address2"`
	State    customtypes.NullableString `json:"state"`
	Country  customtypes.NullableString `json:"country"`
	Pincode  customtypes.NullableString `json:"pincode"`
}

type Customer struct {
	ID              int64                      `json:"id"`
	Code            customtypes.NullableString `json:"code"`
	Name            customtypes.NullableString `json:"name"`
	ContactPerson   customtypes.NullableString `json:"contact_person"`
	BillingAddress  Address                    `json:"billing_address"`
	ShippingAddress Address                    `json:"shipping_address"`
	Status          customtypes.NullableString `json:"status"`
	CreatedAt       time.Time                  `json:"created_at"`
	UpdatedAt       time.Time                  `json:"updated_at"`
	LastUpdatedBy   customtypes.NullableString `json:"last_updated_by"`
}

func NewCustomerWithDefaults() *Customer {
	return &Customer{
		Status:          "active",
		BillingAddress:  Address{},
		ShippingAddress: Address{},
	}
}
