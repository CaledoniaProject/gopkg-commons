package commons

import (
	"bytes"
	"fmt"
	"html/template"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"github.com/alexedwards/scs/redisstore"
	"github.com/alexedwards/scs/v2"
	"github.com/gomodule/redigo/redis"
	"github.com/minio/minio-go"
	"github.com/mitchellh/go-homedir"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"gopkg.in/yaml.v2"
	"gorm.io/gorm"
)

type Initializer interface {
	Init() error
}

type StandardInitializer interface {
	StandardInit() error
}

type MinioConfig struct {
	Endpoint  string `yaml:"endpoint"`
	AccessKey string `yaml:"accessKey"`
	SecretKey string `yaml:"secretKey"`
}

type RedisConfig struct {
	DSN string `yaml:"dsn"`
}

type SitemapConfig struct {
	BaseURL        string `yaml:"baseURL"`
	SitemapBaseURL string `yaml:"sitemapBaseURL"`
	Output         string `yaml:"output"`
	Interval       int    `yaml:"interval"`
}

type EmailTemplate struct {
	Template *template.Template
}

type EmailTemplates map[string]*EmailTemplate

func (e EmailTemplates) Execute(templateName string, data any) (string, error) {
	var (
		buffer bytes.Buffer
	)

	emailTemplate, ok := e[templateName]
	if !ok {
		return "", fmt.Errorf("missing template %s", templateName)
	}

	if err := emailTemplate.Template.Execute(&buffer, data); err != nil {
		return "", err
	}

	return buffer.String(), nil
}

func NewEmailTemplates(folder string) (EmailTemplates, error) {
	var (
		results = make(EmailTemplates)
	)

	if files, err := os.ReadDir(folder); err != nil {
		return nil, err
	} else {
		for _, filp := range files {
			if filp.IsDir() {
				continue
			}

			filename := filepath.Join(folder, filp.Name())
			if tmpl, err := template.ParseFiles(filename); err != nil {
				return nil, errors.Wrapf(err, "parse %s as template", filename)
			} else {
				results[filp.Name()] = &EmailTemplate{Template: tmpl}
			}
		}
	}

	return results, nil
}

type StandardConfig struct {
	Name                 string                    `yaml:"name"`
	Environ              string                    `yaml:"environ"`
	Mode                 string                    `yaml:"mode"`
	Proxy                string                    `yaml:"proxy"`
	Listen               string                    `yaml:"listen"`
	LogLevel             string                    `yaml:"logLevel"`
	LogDir               string                    `yaml:"logDir"`
	CookieDomain         string                    `yaml:"cookieDomain"`
	BaseURL              string                    `yaml:"baseURL"`
	ClientTokens         []string                  `yaml:"clientTokens"`
	MySQLConfig          *DatabaseConfig           `yaml:"mysql"`
	RedisConfig          *RedisConfig              `yaml:"redis"`
	MinioConfig          *MinioConfig              `yaml:"minio"`
	Sitemap              *SitemapConfig            `yaml:"sitemap"`
	OAuthConfigList      []*OAuthConfigBlock       `yaml:"oauth"`
	EmailTemplatesFolder string                    `yaml:"emailTemplatesFolder"`
	EmailTemplates       *EmailTemplates           ``
	SqlDB                *gorm.DB                  ``
	RedisPool            *redis.Pool               ``
	MinioClient          *minio.Client             ``
	RequestLogger        *logrus.Logger            ``
	TokenAuthenticator   *SimpleTokenAuthenticator ``
	SessionManager       *scs.SessionManager       ``
}

func (s *StandardConfig) RedisKeyName(key, value string) string {
	return fmt.Sprintf(`%s@%s-%s-%s`, s.Name, s.Environ, key, value)
}

func (s *StandardConfig) StandardInit() error {
	// logging directory
	if s.LogDir != "" {
		if tmp, err := homedir.Expand(s.LogDir); err != nil {
			return errors.Wrapf(err, "expand homedir")
		} else {
			s.LogDir = tmp
		}
	}

	// proxy
	if s.Proxy != "" {
		if err := SetGlobalHTTPProxy(s.Proxy); err != nil {
			return errors.Wrapf(err, "set global proxy")
		}
	}

	// database
	if s.MySQLConfig != nil {
		// MySQL
		if sqlDB, err := NewMySQLDatabase(s.MySQLConfig); err != nil {
			return errors.Wrapf(err, "create mysql database")
		} else {
			s.SqlDB = sqlDB
		}
	}

	// redis
	if s.RedisConfig != nil {
		if pool, err := NewRedisPool(s.RedisConfig.DSN); err != nil {
			return errors.Wrapf(err, "create redis pool")
		} else {
			s.RedisPool = pool
		}
	}

	// server side
	if s.Mode == "server" {
		// token authenticator
		if len(s.ClientTokens) != 0 {
			s.TokenAuthenticator = NewSimpleTokenAuthenticator(s.ClientTokens)
		}

		// session manager
		s.SessionManager = scs.New()
		s.SessionManager.Lifetime = 7 * 24 * time.Hour
		s.SessionManager.Cookie.HttpOnly = true
		s.SessionManager.Cookie.SameSite = http.SameSiteLaxMode
		if s.CookieDomain != "" {
			s.SessionManager.Cookie.Domain = s.CookieDomain
		}
		if s.Environ == "online" {
			s.SessionManager.Cookie.Secure = true
		}

		if s.RedisPool != nil {
			s.SessionManager.Store = redisstore.New(s.RedisPool)
		}

		// request logger
		if s.LogDir != "" {
			requestLogger, _ := GetRotatingFileLogger(s.LogDir + "/request.log")
			requestLogger.Formatter = GetCleanJSONFormatter()
			s.RequestLogger = requestLogger
		}

		// email
		if s.EmailTemplatesFolder != "" {
			if templates, err := NewEmailTemplates(s.EmailTemplatesFolder); err != nil {
				return err
			} else {
				s.EmailTemplates = &templates
			}
		}
	}

	return nil
}

func NewAppConfig(filename string, config any) error {
	if data, err := os.ReadFile(filename); err != nil {
		return err
	} else if err := yaml.Unmarshal(data, config); err != nil {
		return err
	} else {
		// standard initialize
		if standardInit, ok := config.(StandardInitializer); !ok {
			return fmt.Errorf("%T is not a StandardInitializer", config)
		} else if err := standardInit.StandardInit(); err != nil {
			return err
		}

		// custom initialize
		if methodInit, ok := config.(Initializer); ok {
			if err := methodInit.Init(); err != nil {
				return err
			}
		}
	}

	return nil
}
