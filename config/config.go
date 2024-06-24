package config

import (
	"fmt"
	"time"

	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

type Config struct {
	ServiceHost        string   `mapstructure:"service_host" json:"service_host"`
	ServiceEndpointV   string   `mapstructure:"service_endpoint_v" json:"service_endpoint_v"`
	ServiceEnvironment string   `mapstructure:"service_environment" json:"service_environment"`
	ServicePort        string   `mapstructure:"service_port" json:"service_port"`
	Database           Database `mapstructure:"database" json:"database"`
	Redis
	SecretKeyJWT string `mapstructure:"secret_key_jwt" json:"secret_key_jwt"`
}

func NewConfig() *Config {
	return &Config{
		ServiceHost:        viper.GetString(`APP_HOST`),
		ServiceEndpointV:   viper.GetString(`APP_ENDPOINT_V`),
		ServiceEnvironment: viper.GetString(`APP_ENVIRONMENT`),
		ServicePort:        viper.GetString(`APP_PORT`),
		Database:           LoadConfigDB(),
		SecretKeyJWT:       viper.GetString(`JWT_SECRET_KEY`),
		Redis: Redis{
			Host:     viper.GetString(`REDIS_HOST`),
			Password: viper.GetString(`REDIS_PASSWORD`),
			Db:       viper.GetInt(`REDIS_DB`),
			Port:     viper.GetString(`REDIS_PORT`),
		},
	}
}

type Database struct {
	ConnectionString string
	Flavor           string
	MaxIdleConns     int
	MaxOpenConns     int
	ConnMaxLifetime  time.Duration
	Location         string
	Timeout          time.Duration
}

type Redis struct {
	Host     string
	Password string
	Db       int
	Port     string
}

// Load all config for the system
func LoadConfigDB() Database {

	// Read the configuration file
	if err := viper.ReadInConfig(); err != nil {
		logrus.Errorf(err.Error())
	}

	if viper.GetBool(`debug`) {
		fmt.Println("Service RUN on DEBUG mode")
	}

	dbHost := viper.GetString(`DB_POSTGRES_HOST`)
	dbUser := viper.GetString(`DB_POSTGRES_USER`)
	dbName := viper.GetString(`DB_POSTGRES_NAME`)
	dbPassword := viper.GetString(`DB_POSTGRES_PASSWORD`)
	dbPort := viper.GetString(`DB_POSTGRES_PORT`)

	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable TimeZone=Asia/Jakarta", dbHost, dbUser, dbPassword, dbName, dbPort)
	cfg := Database{
		ConnectionString: dsn,
		MaxIdleConns:     viper.GetInt(`DB_POSTGRES_MAX_IDLE_CONNS`),
		MaxOpenConns:     viper.GetInt(`DB_POSTGRES_MAX_OPEN_CONNS`),
		ConnMaxLifetime:  viper.GetDuration(`DB_POSTGRES_MAX_LIFE_TIME`),
		Timeout:          viper.GetDuration(`DB_POSTGRES_TIMEOUT`),
	}

	return cfg
}
