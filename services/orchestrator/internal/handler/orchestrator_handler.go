package handler

import (
	"net/http"

	"warehouse/services/orchestrator/internal/domain"
	"warehouse/services/orchestrator/internal/service"

	"github.com/gin-gonic/gin"
)

// OrchestratorHandler обрабатывает HTTP-запросы для оркестратора
type OrchestratorHandler struct {
	service *service.OrchestratorService
}

// NewOrchestratorHandler создает новый экземпляр OrchestratorHandler
func NewOrchestratorHandler(service *service.OrchestratorService) *OrchestratorHandler {
	return &OrchestratorHandler{service: service}
}

// RegisterRoutes регистрирует маршруты для обработчика
func (h *OrchestratorHandler) RegisterRoutes(router *gin.Engine) {
	router.POST("/place", h.PlaceItem)
}

// PlaceItem обрабатывает запрос на размещение товара
func (h *OrchestratorHandler) PlaceItem(c *gin.Context) {
	var req domain.PlacementRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Сначала анализируем размещение через все микросервисы
	analysis, err := h.service.AnalyzePlacement(c.Request.Context(), &req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка при анализе размещения: " + err.Error()})
		return
	}

	// Если анализ не успешен, возвращаем результаты анализа
	if !analysis.Success {
		c.JSON(http.StatusOK, analysis)
		return
	}

	// Выполняем размещение через выбранный алгоритм
	response, err := h.service.PlaceItem(c.Request.Context(), &req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка при размещении: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, response)
} 