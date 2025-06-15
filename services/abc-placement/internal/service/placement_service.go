package service

import (
	"context"
	"fmt"
	"sort"

	"warehouse/services/abc-placement/internal/domain"
	"warehouse/services/abc-placement/internal/repository"
)


type PlacementService struct {
	repo repository.Repository
}


func NewPlacementService(repo repository.Repository) *PlacementService {
	return &PlacementService{repo: repo}
}


func (s *PlacementService) AnalyzePlacement(ctx context.Context, req *domain.PlaceRequest) (*domain.PlaceResponse, error) {

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


	var abcCategory string
	if item.Turnover >= 0.8 {
		abcCategory = "A"
	} else if item.Turnover >= 0.15 {
		abcCategory = "B"
	} else {
		abcCategory = "C"
	}


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


	slots, err := s.repo.GetAvailableSlots(ctx, targetZoneType)
	if err != nil {
		return nil, fmt.Errorf("error getting available slots: %w", err)
	}


	sort.Slice(slots, func(i, j int) bool {
		return slots[i].DistanceFromExit < slots[j].DistanceFromExit
	})


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


func (s *PlacementService) PlaceItem(ctx context.Context, req *domain.PlaceRequest) (*domain.PlaceResponse, error) {
	
	requestID, err := s.repo.CreatePlacementRequest(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("error creating placement request: %w", err)
	}


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


	var abcCategory string
	if item.Turnover >= 0.8 {
		abcCategory = "A"
	} else if item.Turnover >= 0.15 {
		abcCategory = "B"
	} else {
		abcCategory = "C"
	}

	
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

	slots, err := s.repo.GetAvailableSlots(ctx, targetZoneType)
	if err != nil {
		return nil, fmt.Errorf("error getting available slots: %w", err)
	}


	sort.Slice(slots, func(i, j int) bool {
		return slots[i].DistanceFromExit < slots[j].DistanceFromExit
	})


	var chosenSlotID string
	if len(slots) > 0 {
		chosenSlotID = slots[0].SlotID
		
	
		if err := s.repo.UpdateSlotOccupation(ctx, chosenSlotID, true); err != nil {
			return nil, fmt.Errorf("error updating slot occupation: %w", err)
		}


		if err := s.repo.CreatePlacementLog(ctx, chosenSlotID, req.ItemID, req.BatchID, "abc_placement"); err != nil {

			fmt.Printf("Error creating placement log: %v\n", err)
		}


		if err := s.repo.CreatePlacementResponse(ctx, requestID, true, chosenSlotID, "abc_placement", 1.0, fmt.Sprintf("Item placed in slot %s in zone %s (ABC category %s)", chosenSlotID, targetZoneType, abcCategory)); err != nil {

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