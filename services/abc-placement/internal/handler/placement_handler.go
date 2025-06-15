package handler

import (
	"net/http"

	"warehouse/services/abc-placement/internal/domain"
	"warehouse/services/abc-placement/internal/service"

	"github.com/gin-gonic/gin"
)
type PlacementHandler struct {
	service *service.PlacementService
}

func NewPlacementHandler(service *service.PlacementService) *PlacementHandler {
	return &PlacementHandler{service: service}
}

func (h *PlacementHandler) RegisterRoutes(router *gin.Engine) {
	router.POST("/api/v1/abc-placement", h.ProcessPlacementRequest)
}

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
		c.JSON(http.StatusBadRequest, gin.H{"error": "Unknown command: " + req.Command})
		return
	}

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		return
	}

	c.JSON(http.StatusOK, response)
}
