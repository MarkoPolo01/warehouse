package repository

import (
	"context"
	"warehouse/services/abc-placement/internal/domain"
)

// Repository defines the interface for data operations
type Repository interface {
	// Check if an item exists
	ItemExists(ctx context.Context, itemID string) (bool, error)
	// Check if a batch exists
	BatchExists(ctx context.Context, batchID string) (bool, error)
	// Get item details for ABC analysis
	GetItemDetails(ctx context.Context, itemID string) (*domain.Item, error)
	// Get available slots based on zone type and distance from exit
	GetAvailableSlots(ctx context.Context, zoneType string) ([]domain.Slot, error)
	// Check if a slot is occupied
	IsSlotOccupied(ctx context.Context, slotID string) (bool, error)
	// Create a placement request entry
	CreatePlacementRequest(ctx context.Context, req *domain.PlaceRequest) (int, error)
	// Create a placement log entry
	CreatePlacementLog(ctx context.Context, slotID, itemID, batchID, algorithm string) error
	// Update slot occupation status
	UpdateSlotOccupation(ctx context.Context, slotID string, isOccupied bool) error
	// Create a placement response entry
	CreatePlacementResponse(ctx context.Context, requestID int, success bool, slotID, algorithm string, score float64, comment string) error
} 