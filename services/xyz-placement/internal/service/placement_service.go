package service

import (
	"context"
	"fmt"
	"sort"

	"warehouse/services/xyz-placement/internal/domain"
	"warehouse/services/xyz-placement/internal/repository"
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


	mr, err := s.repo.GetItemMr(ctx, req.ItemID)
	if err != nil {
		return nil, fmt.Errorf("error getting item Mr: %w", err)
	}

	var xyzCategory string
	if mr < 0.1 {
		xyzCategory = "X" 
	} else if mr < 0.25 {
		xyzCategory = "Y" 
	} else {
		xyzCategory = "Z"
	}


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


	mr, err := s.repo.GetItemMr(ctx, req.ItemID)
	if err != nil {
		return nil, fmt.Errorf("error getting item Mr: %w", err)
	}


	var xyzCategory string
	if mr < 0.1 {
		xyzCategory = "X"
	} else if mr < 0.25 {
		xyzCategory = "Y"
	} else {
		xyzCategory = "Z"
	}


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


		if err := s.repo.CreatePlacementLog(ctx, chosenSlotID, req.ItemID, req.BatchID, "xyz_placement"); err != nil {

			fmt.Printf("Error creating placement log: %v\n", err)
		}

		if err := s.repo.CreatePlacementResponse(ctx, requestID, true, chosenSlotID, "xyz_placement", 1.0, fmt.Sprintf("Item placed in slot %s in zone %s (XYZ category %s, Mr: %.2f)", chosenSlotID, targetZoneType, xyzCategory, mr)); err != nil {

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