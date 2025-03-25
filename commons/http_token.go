package commons

import (
	"net/http"

	"github.com/sirupsen/logrus"
)

type SimpleTokenAuthenticator struct {
	Tokens map[string]bool
}

func NewSimpleTokenAuthenticator(tokens []string) *SimpleTokenAuthenticator {
	auth := &SimpleTokenAuthenticator{
		Tokens: map[string]bool{},
	}

	for _, token := range tokens {
		auth.Tokens[token] = true
	}

	return auth
}

func (s *SimpleTokenAuthenticator) Auth(r *http.Request) error {
	if token := r.Header.Get("token"); token == "" {
		return ErrNoKey
	} else if _, ok := s.Tokens[token]; !ok {
		return ErrAuthFailed
	}

	return nil
}

func SimpleTokenMiddleWare(
	next func(w http.ResponseWriter, r *http.Request),
	auth *SimpleTokenAuthenticator,
	logger *logrus.Logger) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if err := auth.Auth(r); err != nil {
			logger.Errorf("token auth failed: %v", err)
			w.WriteHeader(401)
		} else {
			http.HandlerFunc(next).ServeHTTP(w, r)
		}
	})
}
