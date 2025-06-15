package repository

import (
	"context"
	"warehouse/services/genetic-placement/internal/domain"
)


type Repository interface {

	ItemExists(ctx context.Context, itemID string) (bool, error)

	BatchExists(ctx context.Context, batchID string) (bool, error)

	GetItemDetails(ctx context.Context, itemID string) (*domain.Item, error)

	GetAllAvailableSlots(ctx context.Context) ([]domain.Slot, error)

	CreatePlacementRequest(ctx context.Context, req *domain.PlaceRequest) (int, error)

	UpdateSlotOccupation(ctx context.Context, slotID string, isOccupied bool) error

	CreatePlacementLog(ctx context.Context, slotID, itemID, batchID, algorithm string) error

	CreatePlacementResponse(ctx context.Context, requestID int, success bool, slotID, algorithm string, score float64, comment string) error
} 