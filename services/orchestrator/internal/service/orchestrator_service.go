package service

import (
	"context"
	"fmt"
	"sync"

	"warehouse/services/orchestrator/internal/client"
	"warehouse/services/orchestrator/internal/config"
	"warehouse/services/orchestrator/internal/domain"
)

// OrchestratorService реализует логику оркестрации микросервисов
type OrchestratorService struct {
	config  *config.Config
	clients map[string]*client.PlacementClient
}

// NewOrchestratorService создает новый экземпляр OrchestratorService
func NewOrchestratorService(cfg *config.Config) *OrchestratorService {
	clients := make(map[string]*client.PlacementClient)
	for serviceID, serviceCfg := range cfg.Services {
		clients[serviceID] = client.NewPlacementClient(serviceCfg.URL)
	}

	return &OrchestratorService{
		config:  cfg,
		clients: clients,
	}
}

// AnalyzePlacement анализирует возможность размещения через все микросервисы
func (s *OrchestratorService) AnalyzePlacement(ctx context.Context, req *domain.PlacementRequest) (*domain.OrchestratorResponse, error) {
	var (
		wg      sync.WaitGroup
		mu      sync.Mutex
		results []domain.ServiceResult
	)

	// Запускаем анализ во всех микросервисах параллельно
	for serviceID, placementClient := range s.clients {
		wg.Add(1)
		go func(serviceID string, placementClient *client.PlacementClient) {
			defer wg.Done()

			resp, err := placementClient.AnalyzePlacement(ctx, req)
			if err != nil {
				// В случае ошибки добавляем отрицательный результат
				mu.Lock()
				results = append(results, domain.ServiceResult{
					ServiceName: s.config.Services[serviceID].Name,
					Response: domain.PlacementResponse{
						Success: false,
						Comment: "Ошибка сервиса: " + err.Error(),
						Score:   0,
					},
				})
				mu.Unlock()
				return
			}

			mu.Lock()
			results = append(results, domain.ServiceResult{
				ServiceName: s.config.Services[serviceID].Name,
				Response:    *resp,
			})
			mu.Unlock()
		}(serviceID, placementClient)
	}

	wg.Wait()

	// Выбираем лучший результат
	bestResult := s.selectBestResult(results)

	return &domain.OrchestratorResponse{
		Success:    bestResult.Response.Success,
		SlotID:     bestResult.Response.SlotID,
		Comment:    bestResult.Response.Comment,
		Score:      bestResult.Response.Score,
		Algorithm:  bestResult.ServiceName,
		AllResults: results,
	}, nil
}

// PlaceItem размещает товар используя выбранный алгоритм
func (s *OrchestratorService) PlaceItem(ctx context.Context, req *domain.PlacementRequest) (*domain.OrchestratorResponse, error) {
	// Сначала анализируем размещение
	analysis, err := s.AnalyzePlacement(ctx, req)
	if err != nil {
		return nil, err
	}

	if !analysis.Success {
		return analysis, nil
	}

	// Находим клиент для выбранного алгоритма
	var selectedClient *client.PlacementClient
	for serviceID, serviceCfg := range s.config.Services {
		if serviceCfg.Name == analysis.Algorithm {
			selectedClient = s.clients[serviceID]
			break
		}
	}

	if selectedClient == nil {
		return nil, fmt.Errorf("не найден клиент для алгоритма %s", analysis.Algorithm)
	}

	// Выполняем размещение
	resp, err := selectedClient.PlaceItem(ctx, req)
	if err != nil {
		return nil, err
	}

	return &domain.OrchestratorResponse{
		Success:    resp.Success,
		SlotID:     resp.SlotID,
		Comment:    resp.Comment,
		Score:      resp.Score,
		Algorithm:  analysis.Algorithm,
		AllResults: analysis.AllResults,
	}, nil
}

// selectBestResult выбирает лучший результат из всех полученных
func (s *OrchestratorService) selectBestResult(results []domain.ServiceResult) domain.ServiceResult {
	var bestResult domain.ServiceResult
	bestScore := -1.0

	for _, result := range results {
		if result.Response.Success && result.Response.Score > bestScore {
			bestScore = result.Response.Score
			bestResult = result
		}
	}

	return bestResult
} 