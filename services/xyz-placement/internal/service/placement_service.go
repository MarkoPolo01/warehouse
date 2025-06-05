package service

import (
	"context"
	"fmt"
	"sort"

	"warehouse/services/xyz-placement/internal/domain"
	"warehouse/services/xyz-placement/internal/repository"
)

// PlacementService implements the business logic for XYZ placement
type PlacementService struct {
	repo repository.Repository
}

// NewPlacementService creates a new instance of PlacementService
func NewPlacementService(repo repository.Repository) *PlacementService {
	return &PlacementService{repo: repo}
}

// AnalyzePlacement analyzes the best placement for an item based on XYZ analysis
func (s *PlacementService) AnalyzePlacement(ctx context.Context, req *domain.PlaceRequest) (*domain.PlaceResponse, error) {
	// Check if item and batch exist
	itemExists, err := s.repo.ItemExists(ctx, req.ItemID)
	if err != nil {
		return nil, fmt.Errorf("error checking item existence: %w", err)
	}
	batchExists, err := s.repo.BatchExists(ctx, req.BatchID)
	if err != nil {
		return nil, fmt.Errorf("error checking batch existence: %w", err)
	}

	if !itemExists || !batchExists {
		return &domain.PlaceResponse{
			Success: false,
			Comment: "Item or batch not found",
			Score:   0,
		}, nil
	}

	// Get item Mr for XYZ analysis
	mr, err := s.repo.GetItemMr(ctx, req.ItemID)
	if err != nil {
		return nil, fmt.Errorf("error getting item Mr: %w", err)
	}

	// Determine XYZ category based on Mr (coefficient of variation)
	var xyzCategory string
	// Assuming thresholds: X (мr < 0.1), Y (0.1 <= мr < 0.25), Z (мr >= 0.25)
	if mr < 0.1 {
		xyzCategory = "X" // Стабильный спрос
	} else if mr < 0.25 {
		xyzCategory = "Y" // Умеренно изменчивый спрос
	} else {
		xyzCategory = "Z" // Нестабильный спрос
	}

	// Determine target zone based on XYZ category
	var targetZoneType string
	switch xyzCategory {
	case "X":
		targetZoneType = "fast-access"
	case "Y":
		targetZoneType = "regular"
	case "Z":
		targetZoneType = "deep"
	default:
		targetZoneType = "regular" // Default to regular for unknown categories
	}

	// Find available slots in the target zone
	slots, err := s.repo.GetAvailableSlots(ctx, targetZoneType)
	if err != nil {
		return nil, fmt.Errorf("error getting available slots: %w", err)
	}

	// Sort slots by distance from exit (closest first)
	sort.Slice(slots, func(i, j int) bool {
		return slots[i].DistanceFromExit < slots[j].DistanceFromExit
	})

	// Select the first available slot (closest to exit)
	var suggestedSlotID string
	if len(slots) > 0 {
		suggestedSlotID = slots[0].SlotID
		return &domain.PlaceResponse{
			Success: true,
			SlotID:  suggestedSlotID,
			Comment: fmt.Sprintf("Suggested placement in zone %s (XYZ category %s, Mr: %.2f)", targetZoneType, xyzCategory, mr),
			Score:   0.9,
		}, nil
	}

	return &domain.PlaceResponse{
		Success: false,
		Comment: fmt.Sprintf("No available slots found in zone %s (XYZ category %s, Mr: %.2f)", targetZoneType, xyzCategory, mr),
		Score:   0,
	}, nil
}

// PlaceItem places an item in the best available slot based on XYZ analysis
func (s *PlacementService) PlaceItem(ctx context.Context, req *domain.PlaceRequest) (*domain.PlaceResponse, error) {
	// Create a placement request entry
	requestID, err := s.repo.CreatePlacementRequest(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("error creating placement request: %w", err)
	}

	// Check if item and batch exist
	itemExists, err := s.repo.ItemExists(ctx, req.ItemID)
	if err != nil {
		return nil, fmt.Errorf("error checking item existence: %w", err)
	}
	batchExists, err := s.repo.BatchExists(ctx, req.BatchID)
	if err != nil {
		return nil, fmt.Errorf("error checking batch existence: %w", err)
	}

	if !itemExists || !batchExists {
		return &domain.PlaceResponse{
			Success: false,
			Comment: "Item or batch not found",
			Score:   0,
		}, nil
	}

	// Get item Mr for XYZ analysis
	mr, err := s.repo.GetItemMr(ctx, req.ItemID)
	if err != nil {
		return nil, fmt.Errorf("error getting item Mr: %w", err)
	}

	// Determine XYZ category based on Mr
	var xyzCategory string
	if mr < 0.1 {
		xyzCategory = "X"
	} else if mr < 0.25 {
		xyzCategory = "Y"
	} else {
		xyzCategory = "Z"
	}

	// Determine target zone based on XYZ category
	var targetZoneType string
	switch xyzCategory {
	case "X":
		targetZoneType = "fast-access"
	case "Y":
		targetZoneType = "regular"
	case "Z":
		targetZoneType = "deep"
	default:
		targetZoneType = "regular"
	}

	// Find available slots in the target zone
	slots, err := s.repo.GetAvailableSlots(ctx, targetZoneType)
	if err != nil {
		return nil, fmt.Errorf("error getting available slots: %w", err)
	}

	// Sort slots by distance from exit (closest first)
	sort.Slice(slots, func(i, j int) bool {
		return slots[i].DistanceFromExit < slots[j].DistanceFromExit
	})

	// Select the first available slot (closest to exit)
	var chosenSlotID string
	if len(slots) > 0 {
		chosenSlotID = slots[0].SlotID

		// Update slot occupation
		if err := s.repo.UpdateSlotOccupation(ctx, chosenSlotID, true); err != nil {
			return nil, fmt.Errorf("error updating slot occupation: %w", err)
		}

		// Create placement log
		if err := s.repo.CreatePlacementLog(ctx, chosenSlotID, req.ItemID, req.BatchID, "xyz_placement"); err != nil {
			// Log the error but don't return it as placement was successful
			fmt.Printf("Error creating placement log: %v\n", err)
		}

		// Create placement response
		if err := s.repo.CreatePlacementResponse(ctx, requestID, true, chosenSlotID, "xyz_placement", 1.0, fmt.Sprintf("Item placed in slot %s in zone %s (XYZ category %s, Mr: %.2f)", chosenSlotID, targetZoneType, xyzCategory, mr)); err != nil {
			// Log the error but don't return it as placement was successful
			fmt.Printf("Error creating placement response: %v\n", err)
		}

		return &domain.PlaceResponse{
			Success: true,
			SlotID:  chosenSlotID,
			Comment: fmt.Sprintf("Item placed successfully in slot %s", chosenSlotID),
			Score:   1.0,
		}, nil
	}

	return &domain.PlaceResponse{
		Success: false,
		Comment: fmt.Sprintf("No available slots found in zone %s (XYZ category %s, Mr: %.2f) for placement", targetZoneType, xyzCategory, mr),
		Score:   0,
	}, nil
} 