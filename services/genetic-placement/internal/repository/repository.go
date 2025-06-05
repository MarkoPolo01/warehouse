package repository

import (
	"context"
	"warehouse/services/genetic-placement/internal/domain"
)

// Repository defines the interface for data operations for the genetic algorithm
type Repository interface {
	// Check if an item exists
	ItemExists(ctx context.Context, itemID string) (bool, error)
	// Check if a batch exists
	BatchExists(ctx context.Context, batchID string) (bool, error)
	// Get detailed information about an item
	GetItemDetails(ctx context.Context, itemID string) (*domain.Item, error)
	// Get all available slots with their details
	GetAllAvailableSlots(ctx context.Context) ([]domain.Slot, error)
	// Create a placement request entry
	CreatePlacementRequest(ctx context.Context, req *domain.PlaceRequest) (int, error)
	// Update slot occupation status
	UpdateSlotOccupation(ctx context.Context, slotID string, isOccupied bool) error
	// Create a placement log entry
	CreatePlacementLog(ctx context.Context, slotID, itemID, batchID, algorithm string) error
	// Create a placement response entry
	CreatePlacementResponse(ctx context.Context, requestID int, success bool, slotID, algorithm string, score float64, comment string) error
} 