package commons

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"
)

// gorm interface
type StringSlice []string

func (s *StringSlice) Scan(value interface{}) error {
	if data, ok := value.([]byte); !ok {
		return fmt.Errorf("can't unmarshal type %T", value)
	} else {
		return json.Unmarshal(data, s)
	}
}

func (s StringSlice) Value() (driver.Value, error) {
	data, err := json.Marshal(s)
	if err != nil {
		return nil, err
	}

	return string(data), nil
}
