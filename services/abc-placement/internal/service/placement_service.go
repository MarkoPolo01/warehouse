package service

import (
	"context"
	"fmt"
	"sort"

	"warehouse/services/abc-placement/internal/domain"
	"warehouse/services/abc-placement/internal/repository"
)

// PlacementService implements the business logic for ABC placement
type PlacementService struct {
	repo repository.Repository
}

// NewPlacementService creates a new instance of PlacementService
func NewPlacementService(repo repository.Repository) *PlacementService {
	return &PlacementService{repo: repo}
}

// AnalyzePlacement analyzes the best placement for an item based on ABC analysis
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

	// Get item details for ABC analysis
	item, err := s.repo.GetItemDetails(ctx, req.ItemID)
	if err != nil {
		return nil, fmt.Errorf("error getting item details: %w", err)
	}
	
	if item == nil {
		return &domain.PlaceResponse{
			Success: false,
			Comment: "Item details not found",
			Score:   0,
		}, nil
	}

	// Determine ABC category based on turnover (using 80/15/5 rule as a reference)
	var abcCategory string
	if item.Turnover >= 0.8 {
		abcCategory = "A"
	} else if item.Turnover >= 0.15 {
		abcCategory = "B"
	} else {
		abcCategory = "C"
	}

	// Determine target zone based on ABC category
	var targetZoneType string
	switch abcCategory {
	case "A":
		targetZoneType = "fast-access"
	case "B":
		targetZoneType = "regular"
	case "C":
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
			Comment: fmt.Sprintf("Suggested placement in zone %s (ABC category %s)", targetZoneType, abcCategory),
			Score:   0.9,
		}, nil
	}

	return &domain.PlaceResponse{
		Success: false,
		Comment: fmt.Sprintf("No available slots found in zone %s (ABC category %s)", targetZoneType, abcCategory),
		Score:   0,
	}, nil
}

// PlaceItem places an item in the best available slot based on ABC analysis
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

	// Get item details for ABC analysis
	item, err := s.repo.GetItemDetails(ctx, req.ItemID)
	if err != nil {
		return nil, fmt.Errorf("error getting item details: %w", err)
	}

	if item == nil {
		return &domain.PlaceResponse{
			Success: false,
			Comment: "Item details not found",
			Score:   0,
		}, nil
	}

	// Determine ABC category based on turnover
	var abcCategory string
	if item.Turnover >= 0.8 {
		abcCategory = "A"
	} else if item.Turnover >= 0.15 {
		abcCategory = "B"
	} else {
		abcCategory = "C"
	}

	// Determine target zone based on ABC category
	var targetZoneType string
	switch abcCategory {
	case "A":
		targetZoneType = "fast-access"
	case "B":
		targetZoneType = "regular"
	case "C":
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
		if err := s.repo.CreatePlacementLog(ctx, chosenSlotID, req.ItemID, req.BatchID, "abc_placement"); err != nil {
			// Log the error but don't return it as placement was successful
			fmt.Printf("Error creating placement log: %v\n", err)
		}

		// Create placement response
		if err := s.repo.CreatePlacementResponse(ctx, requestID, true, chosenSlotID, "abc_placement", 1.0, fmt.Sprintf("Item placed in slot %s in zone %s (ABC category %s)", chosenSlotID, targetZoneType, abcCategory)); err != nil {
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
		Comment: fmt.Sprintf("No available slots found in zone %s (ABC category %s) for placement", targetZoneType, abcCategory),
		Score:   0,
	}, nil
} 