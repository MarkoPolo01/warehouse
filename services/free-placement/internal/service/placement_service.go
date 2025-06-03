package service

import (
	"context"
	"warehouse/services/free-placement/internal/domain"
	"warehouse/services/free-placement/internal/repository"
)

// PlacementService реализует бизнес-логику размещения товаров
type PlacementService struct {
	repo repository.Repository
}

// NewPlacementService создает новый экземпляр PlacementService
func NewPlacementService(repo repository.Repository) *PlacementService {
	return &PlacementService{repo: repo}
}

// AnalyzePlacement анализирует возможность размещения товара
func (s *PlacementService) AnalyzePlacement(ctx context.Context, req *domain.PlaceRequest) (*domain.PlaceResponse, error) {
	// Проверяем существование товара и партии
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

	// Ищем первую свободную ячейку
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

// PlaceItem размещает товар в первую свободную ячейку
func (s *PlacementService) PlaceItem(ctx context.Context, req *domain.PlaceRequest) (*domain.PlaceResponse, error) {
	// Ищем первую свободную ячейку
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

	// Проверяем занятость ячейки
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

	// Создаем запрос на размещение
	requestID, err := s.repo.CreatePlacementRequest(ctx, req)
	if err != nil {
		return nil, err
	}

	// Обновляем статус ячейки
	if err := s.repo.UpdateSlotOccupation(ctx, slotID, true); err != nil {
		return nil, err
	}

	// Создаем запись в логе размещений
	if err := s.repo.CreatePlacementLog(ctx, slotID, req.ItemID, req.BatchID, "free_placement"); err != nil {
		return nil, err
	}

	// Создаем ответ системы
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