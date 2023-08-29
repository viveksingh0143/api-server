package customtypes

import (
	"database/sql/driver"
	"fmt"
)

type NullableString string

func (s *NullableString) Scan(value interface{}) error {
	if value == nil {
		*s = ""
		return nil
	}
	asBytes, ok := value.([]byte)
	if !ok {
		return fmt.Errorf("string: expected []byte, got %T", value)
	}
	*s = NullableString(string(asBytes))
	return nil
}

func (s NullableString) Value() (driver.Value, error) {
	if len(s) == 0 { // if nil or empty string
		return nil, nil
	}
	return string(s), nil
}
