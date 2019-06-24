package config

import "github.com/tryffel/market/modules/util"

const (
	ApiV1BasePath       = "/api/v1"
	AuthorizationHeader = "Authorization"
)

// Config has complete application configuration
type Config struct {
	Database  Database    `yaml:"database"`
	Minio     Minio       `yaml:"minio"`
	Api       ApiEndpoint `yaml:"api_endpoint"`
	Logging   Logging     `yaml:"logging"`
	Tokens    Tokens      `yaml:"tokens"`
	S3Gateway S3Gateway   `yaml:"gateway_s3"`
}

type ApiEndpoint struct {
	ListenTo string `yaml:"listen_to"`
	BaseUrl  string `yaml:"base_url"`
}

type Logging struct {
	LogLevel  string `yaml:"log_level"`
	Directory string `yaml:"directory"`
	LogSql    bool   `yaml:"log_sql"`
	LogStd    bool   `yaml:"log_std"`
}

type Tokens struct {
	Key      string        `yaml:"token_key"`
	Expire   bool          `yaml:"expire"`
	Interval util.Interval `yaml:"expiration_interval"`
}

type Database struct {
	Type     string `yaml:"type"`
	Host     string `yaml:"host"`
	Port     string `yaml:"port"`
	Database string `yaml:"database"`
	Username string `yaml:"username"`
	Password string `yaml:"password"`
	Ssl      string `yaml:"ssl"`
}

type Minio struct {
	Url       string `yaml:"url"`
	Bucket    string `yaml:"bucket"`
	AccessKey string `yaml:"access_key"`
	SecretKey string `yaml:"secret_key"`
}

type S3Gateway struct {
	Enabled  bool   `yaml:"enabled"`
	ListenTo string `yaml:"listen_to"`
}
