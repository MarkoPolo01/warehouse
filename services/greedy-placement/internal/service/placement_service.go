package service

import (
	"context"
	"fmt"

	"warehouse/services/greedy-placement/internal/domain"
	"warehouse/services/greedy-placement/internal/repository"
)


type PlacementService struct {
	repo repository.Repository
}


func NewPlacementService(repo repository.Repository) *PlacementService {
	return &PlacementService{repo: repo}
}


func (s *PlacementService) AnalyzePlacement(ctx context.Context, req *domain.PlaceRequest) (*domain.PlaceResponse, error) {
	// Check if item and batch exist (optional for greedy, but good practice)
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


	slots, err := s.repo.GetAllAvailableSlotsOrderedByDistance(ctx)
	if err != nil {
		return nil, fmt.Errorf("error getting available slots: %w", err)
	}


	var suggestedSlotID string
	if len(slots) > 0 {
		suggestedSlotID = slots[0].SlotID
		return &domain.PlaceResponse{
			Success: true,
			SlotID:  suggestedSlotID,
			Comment: fmt.Sprintf("Suggested placement in the closest available slot %s (Zone: %s, Distance: %d)", suggestedSlotID, slots[0].ZoneType, slots[0].DistanceFromExit),
			Score:   1.0,
		}, nil
	}

	return &domain.PlaceResponse{
		Success: false,
		Comment: "No available slots found",
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

		if err := s.repo.CreatePlacementResponse(ctx, requestID, false, "", "greedy_placement", 0, "Item or batch not found"); err != nil {
			fmt.Printf("Error creating placement response for not found item/batch: %v\n", err)
		}
		return &domain.PlaceResponse{
			Success: false,
			Comment: "Item or batch not found",
			Score:   0,
		}, nil
	}


	slots, err := s.repo.GetAllAvailableSlotsOrderedByDistance(ctx)
	if err != nil {
		// Log the placement response
		if err := s.repo.CreatePlacementResponse(ctx, requestID, false, "", "greedy_placement", 0, fmt.Sprintf("Error getting available slots: %v", err)); err != nil {
			fmt.Printf("Error creating placement response for slot error: %v\n", err)
		}
		return nil, fmt.Errorf("error getting available slots: %w", err)
	}


	var chosenSlotID string
	if len(slots) > 0 {
		chosenSlotID = slots[0].SlotID


		if err := s.repo.UpdateSlotOccupation(ctx, chosenSlotID, true); err != nil {

			if err := s.repo.CreatePlacementResponse(ctx, requestID, false, chosenSlotID, "greedy_placement", 0, fmt.Sprintf("Error updating slot occupation: %v", err)); err != nil {
				fmt.Printf("Error creating placement response for occupation error: %v\n", err)
			}
			return nil, fmt.Errorf("error updating slot occupation: %w", err)
		}


		if err := s.repo.CreatePlacementLog(ctx, chosenSlotID, req.ItemID, req.BatchID, "greedy_placement"); err != nil {

			fmt.Printf("Error creating placement log: %v\n", err)
		}

		// Create placement response
		if err := s.repo.CreatePlacementResponse(ctx, requestID, true, chosenSlotID, "greedy_placement", 1.0, fmt.Sprintf("Item placed in slot %s (Zone: %s, Distance: %d)", chosenSlotID, slots[0].ZoneType, slots[0].DistanceFromExit)); err != nil {

			fmt.Printf("Error creating placement response: %v\n", err)
		}

		return &domain.PlaceResponse{
			Success: true,
			SlotID:  chosenSlotID,
			Comment: fmt.Sprintf("Item placed successfully in slot %s", chosenSlotID),
			Score:   1.0,
		}, nil
	}


	if err := s.repo.CreatePlacementResponse(ctx, requestID, false, "", "greedy_placement", 0, "No available slots found for placement"); err != nil {
		fmt.Printf("Error creating placement response for no slots: %v\n", err)
	}

	return &domain.PlaceResponse{
		Success: false,
		Comment: "No available slots found for placement",
		Score:   0,
	}, nil
} 