package storage

import (
	"fmt"
	"github.com/jinzhu/gorm"
	"github.com/sirupsen/logrus"
	"github.com/tryffel/market/config"
	"github.com/tryffel/market/modules/logger"
)

type Database struct {
	engine      *gorm.DB
	initialized bool
	db          string
	log         *logger.SqlLogger
}

// NewDatabase initializes new database
func NewDatabase(c *config.Config, logger *logger.SqlLogger) (*Database, error) {
	db := &Database{}

	conf := c.Database

	var engine *gorm.DB
	var err error

	db.log = logger

	switch c.Database.Type {
	case config.DbSqlite:
		engine, err = gorm.Open("sqlite3", c.Database.Database)
		db.db = c.Database.Database
	case config.DbPostgres:
		ssl := "disable"
		if c.Database.Ssl != "" {
			ssl = "require"
		}

		url := fmt.Sprintf("host=%s port=%s user=%s dbname=%s password=%s sslmode=%s",
			conf.Host, conf.Port, conf.Username, conf.Database, conf.Password, ssl)
		engine, err = gorm.Open("postgres", url)
		db.db = fmt.Sprintf("postgres://%s:%s/%s, ssl: %s", conf.Host, conf.Port, conf.Database, ssl)
	case config.DbMysql:
		engine, err = gorm.Open("mysql",
			fmt.Sprintf("%s:%s@%s:%s/%s?charset=utf8mb4&parseTime=True&loc=Local",
				conf.Username, conf.Password, conf.Host, conf.Port, conf.Database))
		db.db = fmt.Sprintf("mysql://%s:%s/%s", conf.Host, conf.Port, conf.Database)

	default:
		logrus.Fatal("Invalid database type configured!")
		panic("Invalid database configured")
		return nil, err
	}

	if err != nil {
		return db, err
	}

	db.engine = engine

	if err != nil {
		return db, err
	}
	db.engine.LogMode(c.Logging.LogSql)
	if c.Logging.LogSql {
		db.engine.SetLogger(db.log)
	}

	db.initialized = true
	return db, nil
}

func (db *Database) GetEngine() *gorm.DB {
	return db.engine
}

// Close closes database connection
func (db *Database) Close() error {
	//db.sqlLogger.Close()
	return db.engine.Close()
}
