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

func JsonRequest(options *RequestOptions) (*json.RawMessage, error) {
	var (
		jsonResp JsonResponse2
	)

	if resp, body, err := HttpRequest(options); err != nil {
		return nil, err
	} else if resp.StatusCode != 200 {
		return nil, fmt.Errorf("bad http status %d", resp.StatusCode)
	} else if err := json.Unmarshal(body, &jsonResp); err != nil {
		return nil, err
	} else if jsonResp.Code != 0 {
		return nil, fmt.Errorf("code=%d, msg=%s", jsonResp.Code, jsonResp.Message)
	} else {
		return jsonResp.Data, nil
	}
}

func JsonRequestEx(options *RequestOptions, dataStruct any) error {
	var (
		jsonResp JsonResponse2
	)

	if resp, body, err := HttpRequest(options); err != nil {
		return err
	} else if resp.StatusCode != 200 {
		return fmt.Errorf("bad http status %d", resp.StatusCode)
	} else if err := json.Unmarshal(body, &jsonResp); err != nil {
		return err
	} else if jsonResp.Code != 0 {
		return fmt.Errorf("code=%d, msg=%s", jsonResp.Code, jsonResp.Message)
	} else if err := json.Unmarshal(*jsonResp.Data, dataStruct); err != nil {
		return err
	} else {
		return nil
	}
}

func JsonGetProcessTime(r *http.Request) string {
	if timeCtx := r.Context().Value(ContextTypeRequestTimestamp); timeCtx == nil {
		return ""
	} else if timeValue, ok := timeCtx.(time.Time); !ok {
		return ""
	} else {
		return time.Since(timeValue).String()
	}
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

func JsonFail(w http.ResponseWriter, r *http.Request, err error) {
	w.Header().Add("content-type", "application/json")

	if err := json.NewEncoder(w).Encode(&JsonResponse{
		Code:        JsonStatusInvalid,
		Message:     err.Error(),
		ProcessTime: JsonGetProcessTime(r),
	}); err != nil {
		logrus.Errorf("json encode: %v", err)
	}
}

func JsonSuccessEx(w http.ResponseWriter, r *http.Request, data interface{}, total int64) {
	w.Header().Add("content-type", "application/json")

	if err := json.NewEncoder(w).Encode(&JsonResponse{
		Code:        JsonStatusSuccess,
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
