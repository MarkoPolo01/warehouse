package service

import (
	"context"
	"fmt"

	"warehouse/services/genetic-placement/internal/config" // Импортируем пакет config для весов
	"warehouse/services/genetic-placement/internal/domain"
	"warehouse/services/genetic-placement/internal/repository"
)

// PlacementService implements the business logic for genetic algorithm placement
type PlacementService struct {
	repo   repository.Repository
	config *config.Config // Добавляем конфигурацию для доступа к весам
}

// NewPlacementService creates a new instance of PlacementService
func NewPlacementService(repo repository.Repository, config *config.Config) *PlacementService {
	return &PlacementService{repo: repo, config: config}
}

// AnalyzePlacement analyzes the best placement for an item using a simplified genetic algorithm approach
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

	// Get item details
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

	// Get all available slots
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

	// Evaluate fitness for each available slot as a potential placement
	var bestCandidate *domain.PlacementCandidate
	maxFitness := -1.0 // Initialize with a value lower than any possible fitness

	for _, slot := range availableSlots {
		candidate := &domain.PlacementCandidate{
			Item: item,
			Slot: &slot,
		}
		candidate.Fitness = s.calculateFitness(candidate)

		// Select the best candidate based on fitness
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

// PlaceItem places an item using the best placement found by the genetic algorithm approach
func (s *PlacementService) PlaceItem(ctx context.Context, req *domain.PlaceRequest) (*domain.PlaceResponse, error) {
	// Create a placement request entry
	requestID, err := s.repo.CreatePlacementRequest(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("error creating placement request: %w", err)
	}

	// Use AnalyzePlacement to find the best slot
	analyzeResponse, err := s.AnalyzePlacement(ctx, req)
	if err != nil {
		// Log the placement response
		if err := s.repo.CreatePlacementResponse(ctx, requestID, false, "", "genetic_placement", 0, fmt.Sprintf("Error during analysis: %v", err)); err != nil {
			fmt.Printf("Error creating placement response for analysis error: %v\n", err)
		}
		return nil, fmt.Errorf("error during placement analysis: %w", err)
	}

	if !analyzeResponse.Success {
		// Log the placement response
		if err := s.repo.CreatePlacementResponse(ctx, requestID, false, "", "genetic_placement", 0, analyzeResponse.Comment); err != nil {
			fmt.Printf("Error creating placement response for unsuccessful analysis: %v\n", err)
		}
		return analyzeResponse, nil
	}

	// Proceed with placement in the suggested slot
	chosenSlotID := analyzeResponse.SlotID

	// Update slot occupation
	if err := s.repo.UpdateSlotOccupation(ctx, chosenSlotID, true); err != nil {
		// Log the placement response
		if err := s.repo.CreatePlacementResponse(ctx, requestID, false, chosenSlotID, "genetic_placement", 0, fmt.Sprintf("Error updating slot occupation: %v", err)); err != nil {
			fmt.Printf("Error creating placement response for occupation error: %v\n", err)
		}
		return nil, fmt.Errorf("error updating slot occupation: %w", err)
	}

	// Create placement log
	if err := s.repo.CreatePlacementLog(ctx, chosenSlotID, req.ItemID, req.BatchID, "genetic_placement"); err != nil {
		fmt.Printf("Error creating placement log: %v\n", err)
	}

	// Create placement response
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

// calculateFitness calculates the fitness score for a placement candidate
func (s *PlacementService) calculateFitness(candidate *domain.PlacementCandidate) float64 {
	// Fitness = (w1 * normalized_distance) + (w2 * size_compatibility) + (w3 * storage_conditions)
	// Higher fitness is better

	// 1. Normalized Distance: Closer is better
	maxPossibleDistance := 1000.0 // Placeholder, adjust as needed
	normalizedDistance := 1.0 - (float64(candidate.Slot.DistanceFromExit) / maxPossibleDistance)
	if normalizedDistance < 0 {
		normalizedDistance = 0
	}

	// 2. Size Compatibility: How well the item fits into the slot
	sizeCompatibility := 1.0
	if candidate.Item.Weight > candidate.Slot.MaxWeight {
		return 0 // Item too heavy
	}
	if candidate.Item.Length > candidate.Slot.MaxLength {
		return 0 // Item too long
	}
	if candidate.Item.Width > candidate.Slot.MaxWidth {
		return 0 // Item too wide
	}
	if candidate.Item.Height > candidate.Slot.MaxHeight {
		return 0 // Item too high
	}

	// Calculate volume ratio for size compatibility
	itemVolume := candidate.Item.Length * candidate.Item.Width * candidate.Item.Height
	slotVolume := candidate.Slot.MaxLength * candidate.Slot.MaxWidth * candidate.Slot.MaxHeight
	if slotVolume > 0 {
		sizeCompatibility = itemVolume / slotVolume
	} else {
		sizeCompatibility = 0
	}

	// 3. Storage Conditions: Binary score (1 if compatible, 0 if not)
	storageConditionsCompatibility := 0.0
	if candidate.Item.StorageConditions == candidate.Slot.StorageConditions {
		storageConditionsCompatibility = 1.0
	}

	// Calculate total fitness using weights from config
	fitness := (s.config.WeightDistance * normalizedDistance) +
		(s.config.WeightSize * sizeCompatibility) +
		(s.config.WeightStorageConditions * storageConditionsCompatibility)

	return fitness
} 