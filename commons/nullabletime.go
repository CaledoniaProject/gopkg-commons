package commons

import (
	"encoding/json"
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
