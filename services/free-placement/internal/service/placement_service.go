package service

import (
	"context"
	"warehouse/services/free-placement/internal/domain"
	"warehouse/services/free-placement/internal/repository"
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
		return nil, err
	}

	batchExists, err := s.repo.BatchExists(ctx, req.BatchID)
	if err != nil {
		return nil, err
	}

	if !itemExists || !batchExists {
		return &domain.PlaceResponse{
			Success: false,
			Comment: "Товар или партия не найдены",
			Score:   0,
		}, nil
	}


	slotID, err := s.repo.GetFirstFreeSlot(ctx)
	if err != nil {
		return nil, err
	}

	if slotID == "" {
		return &domain.PlaceResponse{
			Success: false,
			Comment: "Нет свободных ячеек для размещения",
			Score:   0,
		}, nil
	}

	return &domain.PlaceResponse{
		Success: true,
		SlotID:  slotID,
		Comment: "Найдена свободная ячейка для размещения",
		Score:   0.9,
	}, nil
}


func (s *PlacementService) PlaceItem(ctx context.Context, req *domain.PlaceRequest) (*domain.PlaceResponse, error) {

	slotID, err := s.repo.GetFirstFreeSlot(ctx)
	if err != nil {
		return nil, err
	}

	if slotID == "" {
		return &domain.PlaceResponse{
			Success: false,
			Comment: "Нет свободных ячеек для размещения",
			Score:   0,
		}, nil
	}


	isOccupied, err := s.repo.IsSlotOccupied(ctx, slotID)
	if err != nil {
		return nil, err
	}

	if isOccupied {
		return &domain.PlaceResponse{
			Success: false,
			SlotID:  slotID,
			Comment: "Ячейка уже занята",
			Score:   0.1,
		}, nil
	}


	requestID, err := s.repo.CreatePlacementRequest(ctx, req)
	if err != nil {
		return nil, err
	}


	if err := s.repo.UpdateSlotOccupation(ctx, slotID, true); err != nil {
		return nil, err
	}


	if err := s.repo.CreatePlacementLog(ctx, slotID, req.ItemID, req.BatchID, "free_placement"); err != nil {
		return nil, err
	}


	if err := s.repo.CreatePlacementResponse(ctx, requestID, true, slotID, "free_placement", 1.0, "Товар успешно размещён в свободной ячейке"); err != nil {
		return nil, err
	}

	return &domain.PlaceResponse{
		Success: true,
		SlotID:  slotID,
		Comment: "Товар успешно размещён в свободной ячейке",
		Score:   1.0,
	}, nil
} 