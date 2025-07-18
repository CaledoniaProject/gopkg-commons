package commons

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"
	"strings"
	"time"
)

type NullableTime time.Time

func (m *NullableTime) UnmarshalJSON(data []byte) error {
	stringVal := string(data)
	if stringVal == `""` || stringVal == "null" {
		return nil
	} else if !strings.Contains(stringVal, "Z") {
		stringVal = strings.TrimSuffix(stringVal, `"`) + `Z"`
	}

	return json.Unmarshal([]byte(stringVal), (*time.Time)(m))
}

func (m *NullableTime) IsZero() bool {
	return time.Time(*m).IsZero()
}

// gorm

func (m NullableTime) Value() (driver.Value, error) {
	t := time.Time(m)
	if t.IsZero() {
		return nil, nil
	}
	return t, nil
}

func (m *NullableTime) Scan(value interface{}) error {
	if value == nil {
		*m = NullableTime{}
		return nil
	}

	switch v := value.(type) {
	case time.Time:
		*m = NullableTime(v)
		return nil
	default:
		return fmt.Errorf("cannot scan type %T into NullableTime", value)
	}
}
