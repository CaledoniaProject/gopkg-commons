package commons

import (
	"fmt"
	"log"
	"net/url"
	"reflect"
	"strings"
	"time"
)

// 编码为form参数
func EncodeStructToURLParams(model interface{}) string {
	var (
		values = reflect.ValueOf(model).Elem()
		types  = reflect.TypeOf(model).Elem()
		result = url.Values{}
	)

	for i := 0; i < types.NumField(); i++ {
		paramName := types.Field(i).Tag.Get("param")
		if paramName == "" || paramName == "-" {
			continue
		}

		joinWith := types.Field(i).Tag.Get("join")
		if joinWith == "" {
			joinWith = ","
		}

		timeLayout := types.Field(i).Tag.Get("layout")
		if timeLayout == "" {
			timeLayout = "2006-01-02"
		}

		switch newVal := values.Field(i).Interface().(type) {
		case string:
			result.Add(paramName, newVal)
		case []string:
			result.Add(paramName, strings.Join(newVal, joinWith))
		case int:
			result.Add(paramName, fmt.Sprintf("%d", newVal))
		case int64:
			result.Add(paramName, fmt.Sprintf("%d", newVal))
		case time.Time:
			result.Add(paramName, newVal.Format(timeLayout))
		default:
			log.Fatalf("unsupported field type: %v", newVal)
		}
	}

	return result.Encode()
}
