package main

import (
	"log"

	"warehouse/services/orchestrator/internal/config"
	"warehouse/services/orchestrator/internal/handler"
	"warehouse/services/orchestrator/internal/service"

	"github.com/gin-gonic/gin"
)

func main() {
	
	cfg := config.NewConfig()


	orchestratorService := service.NewOrchestratorService(cfg)


	orchestratorHandler := handler.NewOrchestratorHandler(orchestratorService)


	router := gin.Default()
	orchestratorHandler.RegisterRoutes(router)


	log.Println("Orchestrator Service запущен на :8086")
	if err := router.Run(":8086"); err != nil {
		log.Fatalf("Ошибка запуска сервера: %v", err)
	}
} 