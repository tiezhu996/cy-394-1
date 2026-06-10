package router

import (
	"fitnessapi/internal/config"
	"fitnessapi/internal/handler"
	"fitnessapi/internal/middleware"
	"fitnessapi/internal/repository"
	"fitnessapi/internal/service"
	"net/http"

	_ "fitnessapi/docs"
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

func New(db *gorm.DB, logger *zap.Logger, cfg config.Config) *gin.Engine {
	app := gin.New()
	app.Use(middleware.RequestLogger(logger), middleware.ErrorHandler(logger), middleware.JWTAuth(cfg))
	recordRepo := repository.NewRecordRepository(db)
	goalRepo := repository.NewGoalRepository(db)
	recordSvc := service.NewRecordService(recordRepo)
	goalSvc := service.NewGoalService(goalRepo)
	records := handler.NewRecordHandler(recordSvc)
	stats := handler.NewStatsHandler(recordSvc)
	rank := handler.NewRankHandler(recordSvc)
	goals := handler.NewGoalHandler(goalSvc, recordSvc, goalRepo)
	webhooks := handler.NewWebhookHandler(recordSvc)
	app.GET("/health", func(c *gin.Context) { c.JSON(http.StatusOK, gin.H{"status": "ok"}) })
	app.POST("/api/v1/login", handler.Login(cfg))
	app.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	api := app.Group("/api/v1")
	api.POST("/records", records.Create)
	api.GET("/records", records.List)
	api.PUT("/records/:id", records.Update)
	api.DELETE("/records/:id", records.Delete)
	api.GET("/sports/types", stats.Types)
	api.GET("/stats", stats.Stats)
	api.GET("/pr", stats.PR)
	api.GET("/rankings", rank.Rankings)
	api.POST("/goals", goals.Save)
	api.GET("/goals/progress", goals.Progress)
	api.POST("/webhooks/device", webhooks.Device)
	return app
}
