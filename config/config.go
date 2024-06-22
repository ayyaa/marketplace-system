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
}

func NewConfig() *Config {
	return &Config{
		ServiceHost:        viper.GetString(`server.app_host`),
		ServiceEndpointV:   viper.GetString(`server.endpoint_v`),
		ServiceEnvironment: viper.GetString(`server.environtment`),
		ServicePort:        viper.GetString(`server.port`),
		Database:           LoadConfigDB(),
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

// Load all config for the system
func LoadConfigDB() Database {

	// Read the configuration file
	if err := viper.ReadInConfig(); err != nil {
		logrus.Errorf(err.Error())
	}

	if viper.GetBool(`debug`) {
		fmt.Println("Service RUN on DEBUG mode")
	}

	dbHost := viper.GetString(`database.host`)
	dbUser := viper.GetString(`database.user`)
	dbName := viper.GetString(`database.name`)
	dbPassword := viper.GetString(`database.password`)
	dbPort := viper.GetString(`database.port`)

	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable TimeZone=Asia/Jakarta", dbHost, dbUser, dbPassword, dbName, dbPort)
	cfg := Database{
		ConnectionString: dsn,
		MaxIdleConns:     viper.GetInt(`database.max_idle_conn`),
		MaxOpenConns:     viper.GetInt(`database.max_open_conn`),
		ConnMaxLifetime:  viper.GetDuration(`database.max_life_time`),
		Timeout:          viper.GetDuration(`database.timeout`),
	}

	return cfg
}
