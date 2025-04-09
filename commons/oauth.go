package commons

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/facebook"
	"golang.org/x/oauth2/github"
	"golang.org/x/oauth2/google"
	"golang.org/x/oauth2/microsoft"
)

var (
	OAuthProviderMicrosoft = "microsoft"
	OAuthProviderGoogle    = "google"
	OAuthProviderGitHub    = "github"
	OAuthProviderFacebook  = "facebook"
	OAuthProviderLinkedIn  = "linkedin"
	OAuthProviderApple     = "apple"
)

type OAuthUserInfo struct {
	Id          string
	Email       string
	DisplayName string
	Avatar      string
}

type GoogleUser struct {
	Id            string `json:"id"`
	Email         string `json:"email"`
	FamilyName    string `json:"family_name"`
	GivenName     string `json:"given_name"`
	Name          string `json:"name"`
	Picture       string `json:"picture"`
	VerifiedEmail bool   `json:"verified_email"`
}

func (g *GoogleUser) ToUserInfo() *OAuthUserInfo {
	return &OAuthUserInfo{
		Id:          g.Id,
		Email:       g.Email,
		DisplayName: g.Name,
		Avatar:      g.Picture,
	}
}

type MicrosoftUser struct {
	Id                string `json:"id"`
	Mail              string `json:"mail"`
	DisplayName       string `json:"displayName"`
	GivenName         string `json:"givenName"`
	Surname           string `json:"surname"`
	UserPrincipalName string `json:"userPrincipalName"`
}

func (m *MicrosoftUser) ToUserInfo() *OAuthUserInfo {
	return &OAuthUserInfo{
		Id:          m.Id,
		Email:       m.Mail,
		DisplayName: m.DisplayName,
	}
}

type OAuthProviderConfig struct {
	Scopes      []string
	Endpoint    oauth2.Endpoint
	UserInfoURL string
}

var providerConfigs = map[string]OAuthProviderConfig{
	OAuthProviderGoogle: {
		Scopes:      []string{"openid", "profile", "email"},
		Endpoint:    google.Endpoint,
		UserInfoURL: "https://www.googleapis.com/oauth2/v2/userinfo",
	},
	OAuthProviderGitHub: {
		Scopes:      []string{"read:user", "user:email"},
		Endpoint:    github.Endpoint,
		UserInfoURL: "https://api.github.com/user",
	},
	OAuthProviderFacebook: {
		Scopes:      []string{"public_profile", "email"},
		Endpoint:    facebook.Endpoint,
		UserInfoURL: "https://graph.facebook.com/me?fields=id,name,email",
	},
	OAuthProviderMicrosoft: {
		Scopes:      []string{"User.Read"},
		Endpoint:    microsoft.AzureADEndpoint("common"),
		UserInfoURL: "https://graph.microsoft.com/v1.0/me",
	},
	OAuthProviderLinkedIn: {
		Scopes: []string{"r_liteprofile", "r_emailaddress"},
		Endpoint: oauth2.Endpoint{
			AuthURL:  "https://www.linkedin.com/oauth/v2/authorization",
			TokenURL: "https://www.linkedin.com/oauth/v2/accessToken",
		},
		UserInfoURL: "https://api.linkedin.com/v2/me",
	},
	OAuthProviderApple: {
		Scopes: []string{"name", "email"},
		Endpoint: oauth2.Endpoint{
			AuthURL:  "https://appleid.apple.com/auth/authorize",
			TokenURL: "https://appleid.apple.com/auth/token",
		},
		UserInfoURL: "", // Apple OAuth requires JWT handling
	},
}

func GetOAuthConfigBlock(oauthConfigBlocks []*OAuthConfigBlock, provider string) *OAuthConfigBlock {
	for _, oauthConfigBlock := range oauthConfigBlocks {
		if oauthConfigBlock.Provider == provider {
			return oauthConfigBlock
		}
	}

	return nil
}

type OAuthConfigBlock struct {
	Provider     string `yaml:"provider" json:"provider"`
	ClientID     string `yaml:"clientId" json:"clientId"`
	ClientSecret string `yaml:"clientSecret" json:"clientSecret"`
	RedirectURL  string `yaml:"redirectUrl" json:"redirectUrl"`
}

func (o *OAuthConfigBlock) GetConfig() *oauth2.Config {
	providerConfig, exists := providerConfigs[o.Provider]
	if !exists {
		return nil
	}

	return &oauth2.Config{
		ClientID:     o.ClientID,
		ClientSecret: o.ClientSecret,
		RedirectURL:  o.RedirectURL,
		Scopes:       providerConfig.Scopes,
		Endpoint:     providerConfig.Endpoint,
	}
}

func (o *OAuthConfigBlock) GetUserInfo(ctx context.Context, provider string, code string) (*OAuthUserInfo, error) {
	var (
		oauth2Config = o.GetConfig()
	)

	token, err := oauth2Config.Exchange(ctx, code)
	if err != nil {
		return nil, err
	}

	providerConfig, exists := providerConfigs[provider]
	if !exists {
		return nil, fmt.Errorf("provider %s not found", provider)
	}

	if provider == OAuthProviderApple {
		return nil, fmt.Errorf("apple OAuth user info fetching requires JWT handling")
	}

	resp, err := oauth2Config.Client(ctx, token).Get(providerConfig.UserInfoURL)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch user info from %s: %v", providerConfig.UserInfoURL, err)
	}
	defer resp.Body.Close()

	// decode data
	if provider == "google" {
		googleUser := &GoogleUser{}

		if err := json.NewDecoder(resp.Body).Decode(&googleUser); err != nil {
			return nil, fmt.Errorf("decode %s user info: %v", provider, err)
		} else {
			return googleUser.ToUserInfo(), nil
		}
	} else if provider == "microsoft" {
		microsoftUser := &MicrosoftUser{}

		if err := json.NewDecoder(resp.Body).Decode(&microsoftUser); err != nil {
			return nil, fmt.Errorf("decode %s user info: %v", provider, err)
		} else {
			return microsoftUser.ToUserInfo(), nil
		}
	}

	return nil, errors.New("unprocessed")
}
