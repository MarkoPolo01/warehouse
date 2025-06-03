package handler

import (
	"net/http"

	"warehouse/services/free-placement/internal/domain"
	"warehouse/services/free-placement/internal/service"

	"github.com/gin-gonic/gin"
)

// PlacementHandler обрабатывает HTTP-запросы для размещения товаров
type PlacementHandler struct {
	service *service.PlacementService
}

// NewPlacementHandler создает новый экземпляр PlacementHandler
func NewPlacementHandler(service *service.PlacementService) *PlacementHandler {
	return &PlacementHandler{service: service}
}

// RegisterRoutes регистрирует маршруты для обработчика
func (h *PlacementHandler) RegisterRoutes(router *gin.Engine) {
	router.POST("/process-placement", h.ProcessPlacementRequest)
}

// ProcessPlacementRequest обрабатывает общий запрос на размещение/анализ
func (h *PlacementHandler) ProcessPlacementRequest(c *gin.Context) {
	var req domain.PlaceRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var response *domain.PlaceResponse
	var err error

	switch req.Command {
	case "analyze":
		response, err = h.service.AnalyzePlacement(c.Request.Context(), &req)
	case "place":
		response, err = h.service.PlaceItem(c.Request.Context(), &req)
	default:
		c.JSON(http.StatusBadRequest, gin.H{"error": "Неизвестная команда: " + req.Command})
		return
	}

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Внутренняя ошибка сервера"})
		return
	}

	c.JSON(http.StatusOK, response)
} 