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
	router.POST("/analyze", h.AnalyzePlacement)
	router.POST("/place", h.PlaceItem)
}

// AnalyzePlacement обрабатывает запрос на анализ размещения
func (h *OrchestratorHandler) AnalyzePlacement(c *gin.Context) {
	var req domain.PlacementRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	response, err := h.service.AnalyzePlacement(c.Request.Context(), &req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Внутренняя ошибка сервера"})
		return
	}

	c.JSON(http.StatusOK, response)
}

// PlaceItem обрабатывает запрос на размещение товара
func (h *OrchestratorHandler) PlaceItem(c *gin.Context) {
	var req domain.PlacementRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	response, err := h.service.PlaceItem(c.Request.Context(), &req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Внутренняя ошибка сервера"})
		return
	}

	c.JSON(http.StatusOK, response)
} 