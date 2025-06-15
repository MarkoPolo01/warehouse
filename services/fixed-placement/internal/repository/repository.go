package repository

import (
	"context"
	"warehouse/services/fixed-placement/internal/domain"
)

type Repository interface {

	ItemExists(ctx context.Context, itemID string) (bool, error)

	BatchExists(ctx context.Context, batchID string) (bool, error)

	GetFixedSlot(ctx context.Context, itemID string) (string, error)
	IsSlotOccupied(ctx context.Context, slotID string) (bool, error)
	CreatePlacementRequest(ctx context.Context, req *domain.PlaceRequest) (int, error)

	UpdateSlotOccupation(ctx context.Context, slotID string, isOccupied bool) error

	CreatePlacementLog(ctx context.Context, slotID, itemID, batchID, algorithm string) error
	
	CreatePlacementResponse(ctx context.Context, requestID int, success bool, slotID, algorithm string, score float64, comment string) error
} 