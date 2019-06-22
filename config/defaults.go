package config

import (
	"github.com/tryffel/market/modules/util"
	"time"
)

var TokenExpire = true
var tokenExpiration = util.Interval(time.Hour * 24 * 14)
var dbDefaultType = DbPostgres
var DbPostgres = "postgresql"
var DbMysql = "mysql"
var DbSqlite = "sqlite"
var tokenKeyLength = 80
var LogMainFile = "market.log"
var LogSqlFile = "sql.log"

func (c *Config) addDefaults() {
	if c.Tokens.Interval == 0 {
		c.Tokens.Interval = tokenExpiration
		c.Tokens.Expire = true
	}

	if c.Api.ListenTo == "" {
		c.Api.ListenTo = "127.0.0.1:8080"
	}

	if c.Database.Type == "" {
		c.Database.Type = "postgresql"
	}

	if c.Tokens.Key == "" {
		c.Tokens.Key = util.RandomKey(80)
	}

}
