package main

import (
	"fitnessapi/internal/config"
	"fitnessapi/internal/model"
	"fitnessapi/internal/router"
	"fitnessapi/internal/utils"
	"fmt"
)

// @title Fitness API
// @version 1.0
// @description 运动数据统计 API 服务
// @host localhost:19409
// @BasePath /
func main() {
	cfg := config.Load()
	logger := utils.NewLogger()
	db := config.MustOpenDB(cfg)
	if err := db.AutoMigrate(&model.User{}, &model.WorkoutRecord{}, &model.Goal{}); err != nil {
		logger.Fatal("migrate failed", utils.ZapError(err))
	}
	app := router.New(db, logger, cfg)
	logger.Info("server started", utils.ZapString("port", cfg.ServerPort))
	if err := app.Run(fmt.Sprintf(":%s", cfg.ServerPort)); err != nil {
		logger.Fatal("server stopped", utils.ZapError(err))
	}
}
