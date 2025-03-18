package commons

import (
	"context"
	"errors"
	"time"

	"github.com/coreos/go-oidc/v3/oidc"
	"golang.org/x/oauth2"
)

type Auth0Config struct {
	ClientID     string `yaml:"clientId" json:"clientId"`
	ClientSecret string `yaml:"clientSecret" json:"clientSecret"`
	Domain       string `yaml:"domain" json:"domain"`
	RedirectURL  string `yaml:"redirectUrl" json:"redirectUrl"`
}

type JWTUserClaims struct {
	Aud           string    `json:"aud"`
	Exp           float64   `json:"exp"` // expire time
	FamilyName    string    `json:"family_name"`
	GivenName     string    `json:"given_name"`
	Email         string    `json:"email"`
	EmailVerified bool      `json:"email_verified"`
	Iat           float64   `json:"iat"`
	Iss           string    `json:"iss"` // auth0 domain
	Name          string    `json:"name"`
	Nickname      string    `json:"nickname"`
	Picture       string    `json:"picture"`
	Sid           string    `json:"sid"`
	Sub           string    `json:"sub"` // account id，e.g google-oauth|UID
	UpdatedAt     time.Time `json:"updated_at"`
}

type Auth0Authenticator struct {
	*oidc.Provider
	oauth2.Config
}

func NewAuth0Authenticator(auth0Config *Auth0Config) (*Auth0Authenticator, error) {
	provider, err := oidc.NewProvider(context.Background(), "https://"+auth0Config.Domain+"/")
	if err != nil {
		return nil, err
	}

	// https://auth0.com/docs/secure/tokens/json-web-tokens/create-custom-claims
	conf := oauth2.Config{
		ClientID:     auth0Config.ClientID,
		ClientSecret: auth0Config.ClientSecret,
		RedirectURL:  auth0Config.RedirectURL,
		Endpoint:     provider.Endpoint(),
		Scopes:       []string{"openid profile email nickname family_name given_name picture"},
	}

	return &Auth0Authenticator{
		Provider: provider,
		Config:   conf,
	}, nil
}

func (a *Auth0Authenticator) VerifyIDToken(ctx context.Context, token *oauth2.Token) (*oidc.IDToken, error) {
	rawIDToken, ok := token.Extra("id_token").(string)
	if !ok {
		return nil, errors.New("no id_token field in oauth2 token")
	}

	oidcConfig := &oidc.Config{
		ClientID: a.ClientID,
	}

	return a.Verifier(oidcConfig).Verify(ctx, rawIDToken)
}
