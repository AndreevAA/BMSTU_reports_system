package main

import (
	"reports_system/cmd/server"
	_ "reports_system/docs"
	"reports_system/internal/handlers/account"
	"reports_system/internal/handlers/label"
	"reports_system/internal/handlers/report"
	"reports_system/internal/mapper"
	"reports_system/internal/repository"
	"reports_system/internal/service"
	"reports_system/internal/session"
	"reports_system/pkg/client/psqlclient"
	"reports_system/pkg/logging"

	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

// @title RS
// @version 1.0
// @description API Server for reports-taking applications
// @termsOfService  http://swagger.io/terms/

// @host localhost:8080
// @BasePath /

// @securityDefinitions.apikey ApiKeyAuth
// @in header
// @name Authorization
func main() {
	logging.Init()
	logger := logging.GetLogger()

	cfg := session.GetConfig()

	client, err := psqlclient.NewClient(cfg.DB)
	if err != nil {
		logger.Fatal(err)
	}

	logger.Info("Create new gin router")
	router := gin.New()

	router.GET("api/v1/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	logger.Info("initializing repository")
	repos := repository.New(client, logger)

	logger.Info("initializing services")
	services := service.New(repos, logger)
	mappers := mapper.New(logger)

	accountHandler := account.NewHandler(logger, services.Account, mappers.Account)
	accountHandler.Register(router)

	reportsHandler := report.NewHandler(logger, services.Report, mappers.Report)
	reportsHandler.Register(router)

	labelsHandler := label.NewHandler(logger, services.Label, mappers.Label)
	labelsHandler.Register(router)

	server.Run(cfg, router, logger)
}
