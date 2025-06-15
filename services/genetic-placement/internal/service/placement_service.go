package service

import (
	"context"
	"fmt"

	"warehouse/services/genetic-placement/internal/config"
	"warehouse/services/genetic-placement/internal/domain"
	"warehouse/services/genetic-placement/internal/repository"
)


type PlacementService struct {
	repo   repository.Repository
	config *config.Config 
}


func NewPlacementService(repo repository.Repository, config *config.Config) *PlacementService {
	return &PlacementService{repo: repo, config: config}
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

	availableSlots, err := s.repo.GetAllAvailableSlots(ctx)
	if err != nil {
		return nil, fmt.Errorf("error getting available slots: %w", err)
	}

	if len(availableSlots) == 0 {
		return &domain.PlaceResponse{
			Success: false,
			Comment: "No available slots found",
			Score:   0,
		}, nil
	}


	var bestCandidate *domain.PlacementCandidate
	maxFitness := -1.0

	for _, slot := range availableSlots {
		candidate := &domain.PlacementCandidate{
			Item: item,
			Slot: &slot,
		}
		candidate.Fitness = s.calculateFitness(candidate)


		if candidate.Fitness > maxFitness {
			maxFitness = candidate.Fitness
			bestCandidate = candidate
		}
	}

	if bestCandidate != nil {
		return &domain.PlaceResponse{
			Success: true,
			SlotID:  bestCandidate.Slot.SlotID,
			Comment: fmt.Sprintf("Suggested placement in slot %s with fitness %.2f", bestCandidate.Slot.SlotID, bestCandidate.Fitness),
			Score:   bestCandidate.Fitness,
		}, nil
	}

	return &domain.PlaceResponse{
		Success: false,
		Comment: "Could not determine best placement",
		Score:   0,
	}, nil
}


func (s *PlacementService) PlaceItem(ctx context.Context, req *domain.PlaceRequest) (*domain.PlaceResponse, error) {

	requestID, err := s.repo.CreatePlacementRequest(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("error creating placement request: %w", err)
	}


	analyzeResponse, err := s.AnalyzePlacement(ctx, req)
	if err != nil {

		if err := s.repo.CreatePlacementResponse(ctx, requestID, false, "", "genetic_placement", 0, fmt.Sprintf("Error during analysis: %v", err)); err != nil {
			fmt.Printf("Error creating placement response for analysis error: %v\n", err)
		}
		return nil, fmt.Errorf("error during placement analysis: %w", err)
	}

	if !analyzeResponse.Success {

		if err := s.repo.CreatePlacementResponse(ctx, requestID, false, "", "genetic_placement", 0, analyzeResponse.Comment); err != nil {
			fmt.Printf("Error creating placement response for unsuccessful analysis: %v\n", err)
		}
		return analyzeResponse, nil
	}


	chosenSlotID := analyzeResponse.SlotID


	if err := s.repo.UpdateSlotOccupation(ctx, chosenSlotID, true); err != nil {
		// Log the placement response
		if err := s.repo.CreatePlacementResponse(ctx, requestID, false, chosenSlotID, "genetic_placement", 0, fmt.Sprintf("Error updating slot occupation: %v", err)); err != nil {
			fmt.Printf("Error creating placement response for occupation error: %v\n", err)
		}
		return nil, fmt.Errorf("error updating slot occupation: %w", err)
	}

	
	if err := s.repo.CreatePlacementLog(ctx, chosenSlotID, req.ItemID, req.BatchID, "genetic_placement"); err != nil {
		fmt.Printf("Error creating placement log: %v\n", err)
	}


	if err := s.repo.CreatePlacementResponse(ctx, requestID, true, chosenSlotID, "genetic_placement", analyzeResponse.Score, fmt.Sprintf("Item placed successfully in slot %s", chosenSlotID)); err != nil {
		fmt.Printf("Error creating placement response: %v\n", err)
	}

	return &domain.PlaceResponse{
		Success: true,
		SlotID:  chosenSlotID,
		Comment: fmt.Sprintf("Item placed successfully in slot %s", chosenSlotID),
		Score:   analyzeResponse.Score,
	}, nil
}

func (s *PlacementService) calculateFitness(candidate *domain.PlacementCandidate) float64 {

	maxPossibleDistance := 1000.0
	normalizedDistance := 1.0 - (float64(candidate.Slot.DistanceFromExit) / maxPossibleDistance)
	if normalizedDistance < 0 {
		normalizedDistance = 0
	}


	sizeCompatibility := 1.0
	if candidate.Item.Weight > candidate.Slot.MaxWeight {
		return 0
	}
	if candidate.Item.Length > candidate.Slot.MaxLength {
		return 0 
	}
	if candidate.Item.Width > candidate.Slot.MaxWidth {
		return 0 
	}
	if candidate.Item.Height > candidate.Slot.MaxHeight {
		return 0
	}


	itemVolume := candidate.Item.Length * candidate.Item.Width * candidate.Item.Height
	slotVolume := candidate.Slot.MaxLength * candidate.Slot.MaxWidth * candidate.Slot.MaxHeight
	if slotVolume > 0 {
		sizeCompatibility = itemVolume / slotVolume
	} else {
		sizeCompatibility = 0
	}

	
	storageConditionsCompatibility := 0.0
	if candidate.Item.StorageConditions == candidate.Slot.StorageConditions {
		storageConditionsCompatibility = 1.0
	}


	fitness := (s.config.WeightDistance * normalizedDistance) +
		(s.config.WeightSize * sizeCompatibility) +
		(s.config.WeightStorageConditions * storageConditionsCompatibility)

	return fitness
} 