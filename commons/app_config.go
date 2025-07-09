package commons

import (
	"fmt"
	"net/http"
	"os"
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

type StandardConfig struct {
	Environ            string                    `yaml:"environ"`
	Mode               string                    `yaml:"mode"`
	Listen             string                    `yaml:"listen"`
	LogLevel           string                    `yaml:"logLevel"`
	LogDir             string                    `yaml:"logDir"`
	CookieDomain       string                    `yaml:"cookieDomain"`
	ClientTokens       []string                  `yaml:"clientTokens"`
	MySQLConfig        *DatabaseConfig           `yaml:"mysql"`
	RedisConfig        *RedisConfig              `yaml:"redis"`
	MinioConfig        *MinioConfig              `yaml:"minio"`
	Sitemap            *SitemapConfig            `yaml:"sitemap"`
	OAuthConfigList    []*OAuthConfigBlock       `yaml:"oauth"`
	SqlDB              *gorm.DB                  ``
	RedisPool          *redis.Pool               ``
	MinioClient        *minio.Client             ``
	RequestLogger      *logrus.Logger            ``
	TokenAuthenticator *SimpleTokenAuthenticator ``
	SessionManager     *scs.SessionManager       ``
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
		s.SessionManager.Lifetime = 24 * time.Hour
		s.SessionManager.Cookie.SameSite = http.SameSiteLaxMode
		if s.CookieDomain != "" {
			s.SessionManager.Cookie.Domain = s.CookieDomain
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
	}

	return nil
}

func NewAppConfig(filename string, config any, migrateItems []any) error {
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
