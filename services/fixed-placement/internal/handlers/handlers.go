package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type Handler struct {

}

func NewHandler() *Handler {
	return &Handler{}
}

func (h *Handler) GetProducts(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"message": "Products endpoint",
	})
}

func (h *Handler) GetWarehouses(c *gin.Context) {

	c.JSON(http.StatusOK, gin.H{
		"message": "Warehouses endpoint",
	})
}

func (h *Handler) GetInventory(c *gin.Context) {

	c.JSON(http.StatusOK, gin.H{
		"message": "Inventory endpoint",
	})
}

func (h *Handler) RegisterRoutes(router *gin.Engine) {
	api := router.Group("/api/v1")
	{
		api.GET("/products", h.GetProducts)
		api.GET("/warehouses", h.GetWarehouses)
		api.GET("/inventory", h.GetInventory)
	}
} 