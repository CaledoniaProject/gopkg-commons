package commons

import (
	"context"
	"crypto/tls"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/sirupsen/logrus"
)

var (
	GlobalHTTPProxy *url.URL
)

type RequestOptions struct {
	URL                      string
	Method                   string
	Headers                  map[string]string
	Body                     io.Reader
	Proxy                    *url.URL
	CookieJar                http.CookieJar
	Username                 string
	Password                 string
	MaxRedirect              int
	CheckRedirect            func(req *http.Request, via []*http.Request, options *RequestOptions) error
	AllowRedirectToDirectory bool
	Context                  context.Context
	Timeout                  int
	ReadBody                 bool
	MaxBodyRead              int64
	MaxRetry                 int
	RetryInterval            int
	RetryBackOffFactor       int
	CheckError               func(no int, resp *http.Response, body []byte, err error) error
}

func SetGlobalHTTPProxy(proxyURL string) error {
	if urlObj, err := url.Parse(proxyURL); err != nil {
		return err
	} else {
		GlobalHTTPProxy = urlObj
		return nil
	}
}

func DefaultCheckHTTPError(no int, resp *http.Response, body []byte, err error) error {
	if err != nil {
		return err
	} else if resp.StatusCode != 200 {
		return fmt.Errorf("bad status code %d", resp.StatusCode)
	}

	return nil
}

func HttpRequest(options *RequestOptions) (resp *http.Response, body []byte, err error) {
	var (
		timeStart = time.Now()
		transport = &http.Transport{
			TLSClientConfig: &tls.Config{
				InsecureSkipVerify: true,
			},
		}
		client = &http.Client{
			Timeout: 10 * time.Second,
			CheckRedirect: func(req *http.Request, via []*http.Request) error {
				if options.CheckRedirect != nil {
					return options.CheckRedirect(req, via, options)
				}

				return nil
			},
		}
	)

	if options.Timeout != 0 {
		client.Timeout = time.Duration(options.Timeout) * time.Second
	}

	if options.CookieJar != nil {
		client.Jar = options.CookieJar
	}

	if options.MaxBodyRead == 0 {
		options.MaxBodyRead = 4096
	}

	if options.MaxRetry == 0 {
		options.MaxRetry = 1
	}

	if options.RetryInterval == 0 {
		options.RetryInterval = 1
	}

	if options.RetryBackOffFactor == 0 {
		options.RetryBackOffFactor = 1
	}

	if options.CheckError == nil {
		options.CheckError = DefaultCheckHTTPError
	}

	if options.Method == "" && options.Body != nil {
		options.Method = "POST"
	}

	req, err := http.NewRequest(options.Method, options.URL, options.Body)
	if err != nil {
		return
	}

	if options.Context != nil {
		req = req.WithContext(options.Context)
	}

	req.Header.Set("connection", "close")
	if options.Username != "" && options.Password != "" {
		req.SetBasicAuth(options.Username, options.Password)
	}

	if options.Headers != nil {
		for k, v := range options.Headers {
			req.Header.Set(k, v)

			// resolve host header issue
			if strings.ToLower(k) == "host" {
				req.Host = v
			}
		}
	}

	if options.Proxy != nil {
		transport.Proxy = http.ProxyURL(options.Proxy)
	} else if GlobalHTTPProxy != nil {
		transport.Proxy = http.ProxyURL(GlobalHTTPProxy)
	}

	client.Transport = transport

	for i := 1; i <= options.MaxRetry; i++ {
		logrus.Debugf("[%d/%d] requesting %s", i, options.MaxRetry, options.URL)

		resp, err = client.Do(req)
		if err == nil {
			if options.ReadBody || options.MaxBodyRead > 0 {
				defer resp.Body.Close()
				if options.MaxBodyRead != 0 {
					body, err = io.ReadAll(io.LimitReader(resp.Body, options.MaxBodyRead))
				} else {
					body, err = io.ReadAll(resp.Body)
				}
			}
		}

		if options.CheckError(i, resp, body, err) == nil {
			logrus.Debugf("[%d/%d] request done %s, took %.2f seconds", i, options.MaxRetry, options.URL, time.Since(timeStart).Seconds())
			break
		} else if i != options.MaxRetry {
			time.Sleep(time.Duration(options.RetryInterval + i*options.RetryBackOffFactor))
		}
	}

	return
}

func URLDownloadToByets(options *RequestOptions) ([]byte, error) {
	if resp, body, err := HttpRequest(options); err != nil {
		return nil, err
	} else if resp.StatusCode != 200 {
		return nil, err
	} else {
		return body, nil
	}
}
