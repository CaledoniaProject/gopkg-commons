package commons

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/sirupsen/logrus"
)

const (
	JsonStatusSuccess = iota
	JsonStatusInvalid
	JsonStatusFailed
	JsonStatusNoData
	JsonStatusUnauthorized
)

type JsonResponse struct {
	Code        int         `json:"code"`
	Total       int64       `json:"total"`
	ProcessTime string      `json:"process_time"`
	Message     string      `json:"message"`
	Data        interface{} `json:"data"`
}

type JsonResponse2 struct {
	Code    int              `json:"code"`
	Total   int              `json:"total"`
	Message string           `json:"message"`
	Data    *json.RawMessage `json:"data"`
}

func JsonGetProcessTime(r *http.Request) string {
	var (
		requestTime = r.Context().Value(ContextTypeRequestTimestamp).(time.Time)
		processTime = time.Since(requestTime)
	)

	return processTime.String()
}

func JsonSuccess(w http.ResponseWriter, r *http.Request, data interface{}) {
	w.Header().Add("content-type", "application/json")

	if err := json.NewEncoder(w).Encode(&JsonResponse{
		Code:        0,
		Message:     "succeed",
		Data:        data,
		ProcessTime: JsonGetProcessTime(r),
	}); err != nil {
		logrus.Errorf("json encode: %v", err)
	}
}

func JsonSuccessEx(w http.ResponseWriter, r *http.Request, data interface{}, total int64) {
	w.Header().Add("content-type", "application/json")

	if err := json.NewEncoder(w).Encode(&JsonResponse{
		Code:        0,
		Message:     "succeed",
		Data:        data,
		Total:       total,
		ProcessTime: JsonGetProcessTime(r),
	}); err != nil {
		logrus.Errorf("json encode: %v", err)
	}
}

func JsonFailEx(w http.ResponseWriter, r *http.Request, format string, params ...interface{}) {
	w.Header().Add("content-type", "application/json")

	if err := json.NewEncoder(w).Encode(&JsonResponse{
		Code:        JsonStatusInvalid,
		Message:     fmt.Sprintf(format, params...),
		ProcessTime: JsonGetProcessTime(r),
	}); err != nil {
		logrus.Errorf("json encode: %v", err)
	}
}
