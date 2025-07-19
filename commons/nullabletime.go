package commons

import (
	"database/sql/driver"
	"encoding/json"
	"encoding/xml"
	"fmt"
	"strings"
	"time"
)

type NullableTime time.Time

func (nt *NullableTime) UnmarshalXML(d *xml.Decoder, start xml.StartElement) error {
	var s string
	if err := d.DecodeElement(&s, &start); err != nil || s == "" {
		return err
	}
	for _, layout := range []string{
		"2006-01-02",
		"2006-01-02T15:04:05Z07:00",
		"2006-01-02T15:04:05Z",
	} {
		if t, err := time.Parse(layout, s); err == nil {
			*nt = NullableTime(t)
			return nil
		}
	}

	return fmt.Errorf("invalid time format: %s", s)
}

func (nt *NullableTime) UnmarshalJSON(data []byte) error {
	stringVal := string(data)
	if stringVal == `""` || stringVal == "null" {
		return nil
	} else if !strings.Contains(stringVal, "Z") {
		stringVal = strings.TrimSuffix(stringVal, `"`) + `Z"`
	}

	return json.Unmarshal([]byte(stringVal), (*time.Time)(nt))
}

func (nt *NullableTime) IsZero() bool {
	return time.Time(*nt).IsZero()
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
