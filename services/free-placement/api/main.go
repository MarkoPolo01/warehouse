package main

import (
	"log"

	"warehouse/services/free-placement/internal/config"
	"warehouse/services/free-placement/internal/handler"
	"warehouse/services/free-placement/internal/repository"
	"warehouse/services/free-placement/internal/service"
	"warehouse/services/free-placement/pkg/database"

	"github.com/gin-gonic/gin"
)

func main() {
	cfg := config.LoadConfig()

	dbConfig := &database.Config{
		Host:     cfg.DBHost,
		Port:     cfg.DBPortInt,
		User:     cfg.DBUser,
		Password: cfg.DBPassword,
		DBName:   cfg.DBName,
		SSLMode:  "disable",
	}

	db, err := database.NewPostgresConnection(dbConfig)
	if err != nil {
		log.Fatalf("Ошибка подключения к базе данных: %v", err)
	}
	defer db.Close()

	repo := repository.NewPostgresRepository(db)
	placementService := service.NewPlacementService(repo)
	placementHandler := handler.NewPlacementHandler(placementService)


	router := gin.Default()
	placementHandler.RegisterRoutes(router)


	serverAddr := ":" + cfg.ServerPort
	log.Printf("Free Placement Service запущен на %s", serverAddr)
	if err := router.Run(serverAddr); err != nil {
		log.Fatalf("Ошибка запуска сервера: %v", err)
	}
} 