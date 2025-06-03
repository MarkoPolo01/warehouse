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
	// Загрузка конфигурации
	cfg := config.LoadConfig()

	// Инициализация подключения к базе данных
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

	// Инициализация зависимостей
	repo := repository.NewPostgresRepository(db)
	placementService := service.NewPlacementService(repo)
	placementHandler := handler.NewPlacementHandler(placementService)

	// Настройка маршрутизатора
	router := gin.Default()
	placementHandler.RegisterRoutes(router)

	// Запуск сервера
	serverAddr := ":" + cfg.ServerPort
	log.Printf("Free Placement Service запущен на %s", serverAddr)
	if err := router.Run(serverAddr); err != nil {
		log.Fatalf("Ошибка запуска сервера: %v", err)
	}
} 