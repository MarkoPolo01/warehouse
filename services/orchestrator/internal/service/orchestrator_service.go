package service

import (
	"context"
	"fmt"
	"sync"

	"warehouse/services/orchestrator/internal/client"
	"warehouse/services/orchestrator/internal/config"
	"warehouse/services/orchestrator/internal/domain"
)


type OrchestratorService struct {
	config  *config.Config
	clients map[string]*client.PlacementClient
}


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

func (s *OrchestratorService) AnalyzePlacement(ctx context.Context, req *domain.PlacementRequest) (*domain.OrchestratorResponse, error) {
	var (
		wg      sync.WaitGroup
		mu      sync.Mutex
		results []domain.ServiceResult
	)


	for serviceID, placementClient := range s.clients {
		wg.Add(1)
		go func(serviceID string, placementClient *client.PlacementClient) {
			defer wg.Done()

			resp, err := placementClient.AnalyzePlacement(ctx, req)
			if err != nil {
	
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


func (s *OrchestratorService) PlaceItem(ctx context.Context, req *domain.PlacementRequest) (*domain.OrchestratorResponse, error) {
	
	analysis, err := s.AnalyzePlacement(ctx, req)
	if err != nil {
		return nil, err
	}

	if !analysis.Success {
		return analysis, nil
	}


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


func (s *OrchestratorService) selectBestResult(results []domain.ServiceResult) domain.ServiceResult {
	var bestResult domain.ServiceResult
	bestScore := -1.0

	for _, result := range results {
		if !result.Response.Success {
			continue
		}


		score := result.Response.BaseScore


		switch {
		case result.Response.ResponseTimeMs < 300:
			score += 0.03
		case result.Response.ResponseTimeMs > 1000:
			score -= 0.05
		}


		switch {
		case result.Response.DistanceToExit < 5:
			score += 0.05
		case result.Response.DistanceToExit > 15:
			score -= 0.03
		}


		if result.Response.HasFixedSlot {
			score += 0.07
		}
		if result.Response.HighWarehouseLoad {
			score += 0.05
		}
		if result.Response.HighTurnover {
			score += 0.05
		}
		if result.Response.HeavyItem {
			score += 0.03
		}
		if result.Response.NoPlacementHistory {
			score += 0.02
		}
		if result.Response.FastAccessZone {
			score += 0.03
		}
		if result.Response.XYZCompliant {
			score += 0.02
		}

	
		if score > bestScore {
			bestScore = score
			bestResult = result
	
			bestResult.Response.Score = score
		}
	}

	return bestResult
} 