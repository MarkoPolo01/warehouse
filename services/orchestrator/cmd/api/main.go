package main

import (
	"log"

	"warehouse/services/orchestrator/internal/config"
	"warehouse/services/orchestrator/internal/handler"
	"warehouse/services/orchestrator/internal/service"

	"github.com/gin-gonic/gin"
)

func main() {
	// Инициализация конфигурации
	cfg := config.NewConfig()

	// Инициализация сервиса
	orchestratorService := service.NewOrchestratorService(cfg)

	// Инициализация обработчика
	orchestratorHandler := handler.NewOrchestratorHandler(orchestratorService)

	// Настройка маршрутизатора
	router := gin.Default()
	orchestratorHandler.RegisterRoutes(router)

	// Запуск сервера
	log.Println("Orchestrator Service запущен на :8086")
	if err := router.Run(":8086"); err != nil {
		log.Fatalf("Ошибка запуска сервера: %v", err)
	}
} 