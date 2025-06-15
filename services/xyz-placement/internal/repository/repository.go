package repository

import (
	"context"
	"warehouse/services/xyz-placement/internal/domain"
)


type Repository interface {
	
	ItemExists(ctx context.Context, itemID string) (bool, error)

	BatchExists(ctx context.Context, batchID string) (bool, error)

	GetItemMr(ctx context.Context, itemID string) (float64, error)

	GetAvailableSlots(ctx context.Context, zoneType string) ([]domain.Slot, error)

	CreatePlacementRequest(ctx context.Context, req *domain.PlaceRequest) (int, error)
	
	UpdateSlotOccupation(ctx context.Context, slotID string, isOccupied bool) error
	
	CreatePlacementLog(ctx context.Context, slotID, itemID, batchID, algorithm string) error

	CreatePlacementResponse(ctx context.Context, requestID int, success bool, slotID, algorithm string, score float64, comment string) error
} 