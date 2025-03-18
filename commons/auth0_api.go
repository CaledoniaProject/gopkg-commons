package commons

import (
	"context"

	"github.com/auth0/go-auth0/management"
)

type Auth0APIConfig struct {
	ClientID     string `yaml:"clientId" json:"clientId"`
	ClientSecret string `yaml:"clientSecret" json:"clientSecret"`
	Domain       string `yaml:"domain" json:"domain"`
}

func (a *Auth0APIConfig) Auth0DeleteUser(externalId string) error {
	if authAPI, err := management.New(
		a.Domain,
		management.WithClientCredentials(context.Background(), a.ClientID, a.ClientSecret),
	); err != nil {
		return err
	} else if err := authAPI.User.Delete(context.Background(), externalId); err != nil {
		return err
	}

	return nil
}
