package model

import (
	"database/sql/driver"
	"fmt"
)

type NullableString string

func (n NullableString) Value() (driver.Value, error) {
	if n == "" {
		return nil, nil
	}
	return string(n), nil
}
func (n NullableString) String() string {
	return string(n)
}

func (n *NullableString) Scan(value interface{}) error {

	switch v := value.(type) {
	case nil:
		*n = ""
	case []byte:
		*n = NullableString(v)
	case string:
		*n = NullableString(v)
	case *string:
		*n = NullableString(*v)
	default:
		return fmt.Errorf("unsupported type for NullableString: %t", value)
	}

	return nil
}
