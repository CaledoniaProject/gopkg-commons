package commons

import (
	"log"
	"os"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/pkg/errors"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"gorm.io/gorm/schema"
)

type DatabaseConfig struct {
	LogDir    string `yaml:"logDir"`
	GormDebug bool   `yaml:"gormDebug"`
	DSN       string `yaml:"dsn"`
}

func NewMySQLDatabase(config *DatabaseConfig) (*gorm.DB, error) {
	gormConfig := &gorm.Config{
		AllowGlobalUpdate: true,
		Logger:            logger.Default.LogMode(logger.Silent),
		NamingStrategy: &schema.NamingStrategy{
			NoLowerCase:   true,
			SingularTable: true,
		},
		NowFunc: func() time.Time {
			return time.Now().Local()
		},
	}

	if config.GormDebug {
		if config.LogDir == "" {
			return nil, errors.New("gorm debug enabled, but provided empty logDir")
		}

		if err := os.MkdirAll(config.LogDir, 0755); err != nil {
			return nil, err
		}

		fp, err := os.Create(config.LogDir + "/gorm.log")
		if err != nil {
			return nil, errors.Wrapf(err, "create gorm log file")
		}

		gormConfig.Logger = logger.New(
			log.New(fp, "\r\n", log.LstdFlags),
			logger.Config{
				SlowThreshold:             time.Duration(5) * time.Second, // Slow SQL threshold
				LogLevel:                  logger.Info,                    // Log level
				IgnoreRecordNotFoundError: true,                           // Ignore ErrRecordNotFound error for logger
				Colorful:                  false,                          // Disable color
			},
		)
	}

	sqlDB, err := gorm.Open(mysql.Open(config.DSN), gormConfig)
	if err != nil {
		return nil, errors.Wrap(err, "validate database")
	}

	db, err := sqlDB.DB()
	if err != nil {
		return nil, errors.Wrap(err, "get database")
	} else if err = db.Ping(); err != nil {
		return nil, errors.Wrap(err, "ping database")
	}

	db.SetConnMaxLifetime(time.Minute * 3)
	db.SetMaxOpenConns(100)
	db.SetMaxIdleConns(50)

	return sqlDB, nil
}
