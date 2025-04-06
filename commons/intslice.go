package commons

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"
)

// gorm interface
type IntSlice []int

func (s *IntSlice) Scan(value interface{}) error {
	if data, ok := value.([]byte); !ok {
		return fmt.Errorf("can't unmarshal type %T", value)
	} else {
		return json.Unmarshal(data, s)
	}
}

func (s IntSlice) Value() (driver.Value, error) {
	data, err := json.Marshal(s)
	if err != nil {
		return nil, err
	}

	return string(data), nil
}
