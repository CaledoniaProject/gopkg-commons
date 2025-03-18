package commons

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
)

func RootHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("content-type", "text/html")
	http.Error(w, "No content", http.StatusOK)
}
func GetPostBody(w http.ResponseWriter, r *http.Request) ([]byte, error) {
	var (
		contentType       = r.Header.Get("content-type")
		maxRead     int64 = 10 * 1024 * 1024
	)

	if strings.Contains(contentType, "multipart/form-data") {
		if err := r.ParseMultipartForm(maxRead); err != nil {
			return nil, err
		}

		results := []string{}
		for key, values := range r.MultipartForm.Value {
			results = append(results, fmt.Sprintf("%s=%s", key, url.QueryEscape(values[0])))
		}

		return []byte(strings.Join(results, "&")), nil
	} else if strings.Contains(contentType, "application/json") {
		postBody, err := io.ReadAll(io.LimitReader(r.Body, maxRead))
		if err != nil {
			return nil, err
		} else {
			r.Body.Close()
			r.Body = io.NopCloser(bytes.NewReader(postBody))
			return postBody, err
		}
	}

	return nil, nil
}
