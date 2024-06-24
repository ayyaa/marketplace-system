package main

import (
	"log"
	"marketplace-system/config"
	"marketplace-system/database"
	"marketplace-system/handlers"
	"marketplace-system/lib/validator"
	repository "marketplace-system/repositories"
	"marketplace-system/routes"
	"marketplace-system/services"
	"os"

	"github.com/go-redis/redis/v8"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"gorm.io/gorm"
)

type Main struct {
	cfg      *config.Config
	database Database
	repo     *repository.Main
	service  *services.Main
	handler  *handlers.Main
	router   *echo.Echo
}

type Database struct {
	Postgres *gorm.DB
	Redis    *redis.Client
}

func New() *Main {
	return new(Main)
}

// @title Echo Swagger Marketplace System API
// @version 1.0
// @description This is a Marketplace System server.
// @termsOfService http://swagger.io/terms/

// @contact.name API Support
// @contact.url http://www.swagger.io/support
// @contact.email support@swagger.io

// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html

// @host localhost:8080
// @BasePath /v1
// @schemes http
func (m *Main) Init() (err error) {
	// Set the configuration file name and path
	viper.SetConfigName("config")
	viper.AddConfigPath(".")
	viper.SetConfigType("json")
	viper.AutomaticEnv()
	err = viper.ReadInConfig()
	if err != nil {
		m.cfg = &config.Config{
			ServiceHost:        "localhost",
			ServiceEndpointV:   "V1",
			ServiceEnvironment: "DEVELOPMENT",
			ServicePort:        "3002",
			Database:           config.LoadConfigDB(),
		}

	} else {

		m.cfg = config.NewConfig()
		m.database.Postgres, err = database.GetConnection(config.LoadConfigDB())
		if err != nil {
			return
		}

		m.database.Redis, err = database.GetConnectionRedis(m.cfg.Redis)
		if err != nil {
			return
		}
	}

	e := echo.New()

	// Create a new instance of the logger
	logger := logrus.New()

	// Set the desired log level (optional, default is logrus.InfoLevel)
	logger.SetLevel(logrus.DebugLevel)
	logger.SetOutput(os.Stdout)
	logger.SetFormatter(&logrus.TextFormatter{
		FullTimestamp: true,
	})

	// Middleware
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.Use(middleware.CORS())

	m.repo = repository.Init(repository.Options{
		Config:   m.cfg,
		Postgres: m.database.Postgres,
		Redis:    m.database.Redis,
	})

	m.service = services.Init(services.Options{
		Config:     m.cfg,
		Repository: m.repo,
		Postgres:   m.database.Postgres,
	})
	m.handler = handlers.Init(handlers.Options{
		Config:   m.cfg,
		Services: m.service,
	})

	m.router = e
	routes.ConfigureRouter(e, m.handler, m.cfg)
	return nil
}

func (m *Main) Run() (err error) {
	defer m.close()
	logrus.Infof("Server run on localhost%v", m.cfg.ServicePort)
	m.router.Start(":" + m.cfg.ServicePort)
	return
}

func (m *Main) close() {
	if m.database.Postgres != nil {
		if db, err := m.database.Postgres.DB(); err == nil {
			db.Close()
		}
	}
}

func main() {
	app := New()
	validator.InitValidator()
	err := app.Init()
	if err != nil {
		log.Fatalf("Error in initializing the application: %+v", err)
		return
	}

	err = app.Run()
	if err != nil {
		log.Fatalf("Error in running the application: %+v", err)
		return
	}

}
