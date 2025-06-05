package repository

import (
	"context"
	"warehouse/services/greedy-placement/internal/domain"
)

// Repository определяет интерфейс для работы с хранилищем данных для жадного алгоритма
type Repository interface {
	// Проверяет существование товара
	ItemExists(ctx context.Context, itemID string) (bool, error)
	// Проверяет существование партии
	BatchExists(ctx context.Context, batchID string) (bool, error)
	// Получает все доступные ячейки, отсортированные по расстоянию от выхода
	GetAllAvailableSlotsOrderedByDistance(ctx context.Context) ([]domain.Slot, error)
	// Проверяет занятость ячейки (может быть не нужна)
	// IsSlotOccupied(ctx context.Context, slotID string) (bool, error)
	// Создает запрос на размещение
	CreatePlacementRequest(ctx context.Context, req *domain.PlaceRequest) (int, error)
	// Обновляет статус занятости ячейки
	UpdateSlotOccupation(ctx context.Context, slotID string, isOccupied bool) error
	// Создает запись в логе размещений
	CreatePlacementLog(ctx context.Context, slotID, itemID, batchID, algorithm string) error
	// Создает ответ системы на размещение
	CreatePlacementResponse(ctx context.Context, requestID int, success bool, slotID, algorithm string, score float64, comment string) error
} 