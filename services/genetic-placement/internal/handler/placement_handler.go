package handler

import (
	"fmt"
	"log"
	"net/http"

	"warehouse/services/genetic-placement/internal/domain"
	"warehouse/services/genetic-placement/internal/service"

	"github.com/gin-gonic/gin"
)

// PlacementHandler handles HTTP requests for genetic placement
type PlacementHandler struct {
	service *service.PlacementService
}

// NewPlacementHandler creates a new instance of PlacementHandler
func NewPlacementHandler(service *service.PlacementService) *PlacementHandler {
	return &PlacementHandler{service: service}
}

// RegisterRoutes registers the routes for the handler
func (h *PlacementHandler) RegisterRoutes(router *gin.Engine) {
	router.POST("/api/v1/genetic-placement", h.ProcessPlacementRequest)
}

// ProcessPlacementRequest handles incoming placement requests (analyze or place)
func (h *PlacementHandler) ProcessPlacementRequest(c *gin.Context) {
	var req domain.PlaceRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		log.Printf("Error binding JSON: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": fmt.Sprintf("Invalid request format: %v", err)})
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
		log.Printf("Unknown command: %s", req.Command)
		c.JSON(http.StatusBadRequest, gin.H{"error": fmt.Sprintf("Unknown command: %s", req.Command)})
		return
	}

	if err != nil {
		log.Printf("Service error: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("Service error: %v", err)})
		return
	}

	c.JSON(http.StatusOK, response)
} 