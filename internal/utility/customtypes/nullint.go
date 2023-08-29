package customtypes

import (
	"database/sql/driver"
	"fmt"
)

type NullableInt64 int64

func (s *NullableInt64) Scan(value interface{}) error {
	if value == nil {
		*s = 0
		return nil
	}
	asInt64, ok := value.(int64)
	if !ok {
		return fmt.Errorf("int64: expected int64, got %T", value)
	}
	*s = NullableInt64(asInt64)
	return nil
}

func (s NullableInt64) Value() (driver.Value, error) {
	if s <= 0 { // if nil or empty string
		return nil, nil
	}
	return s, nil
}
