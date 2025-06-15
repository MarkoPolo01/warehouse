package handler

import (
	"net/http"

	"warehouse/services/orchestrator/internal/domain"
	"warehouse/services/orchestrator/internal/service"

	"github.com/gin-gonic/gin"
)


type OrchestratorHandler struct {
	service *service.OrchestratorService
}


func NewOrchestratorHandler(service *service.OrchestratorService) *OrchestratorHandler {
	return &OrchestratorHandler{service: service}
}


func (h *OrchestratorHandler) RegisterRoutes(router *gin.Engine) {
	router.POST("/place", h.PlaceItem)
}


func (h *OrchestratorHandler) PlaceItem(c *gin.Context) {
	var req domain.PlacementRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}


	analysis, err := h.service.AnalyzePlacement(c.Request.Context(), &req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка при анализе размещения: " + err.Error()})
		return
	}

	
	if !analysis.Success {
		c.JSON(http.StatusOK, analysis)
		return
	}


	response, err := h.service.PlaceItem(c.Request.Context(), &req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка при размещении: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, response)
} 