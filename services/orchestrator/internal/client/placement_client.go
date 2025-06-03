package client

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"warehouse/services/orchestrator/internal/domain"
)

// PlacementClient представляет клиент для работы с микросервисами размещения
type PlacementClient struct {
	client  *http.Client
	baseURL string
}

// NewPlacementClient создает новый клиент
func NewPlacementClient(baseURL string) *PlacementClient {
	return &PlacementClient{
		client: &http.Client{
			Timeout: 5 * time.Second,
		},
		baseURL: baseURL,
	}
}

// AnalyzePlacement отправляет запрос на анализ размещения
func (c *PlacementClient) AnalyzePlacement(ctx context.Context, req *domain.PlacementRequest) (*domain.PlacementResponse, error) {
	req.Command = "analyze"
	return c.sendRequest(ctx, "/analyze", req)
}

// PlaceItem отправляет запрос на размещение товара
func (c *PlacementClient) PlaceItem(ctx context.Context, req *domain.PlacementRequest) (*domain.PlacementResponse, error) {
	req.Command = "place"
	return c.sendRequest(ctx, "/place", req)
}

// sendRequest отправляет запрос к микросервису
func (c *PlacementClient) sendRequest(ctx context.Context, path string, req *domain.PlacementRequest) (*domain.PlacementResponse, error) {
	body, err := json.Marshal(req)
	if err != nil {
		return nil, fmt.Errorf("ошибка маршалинга запроса: %v", err)
	}

	httpReq, err := http.NewRequestWithContext(ctx, "POST", c.baseURL+path, bytes.NewBuffer(body))
	if err != nil {
		return nil, fmt.Errorf("ошибка создания запроса: %v", err)
	}

	httpReq.Header.Set("Content-Type", "application/json")

	resp, err := c.client.Do(httpReq)
	if err != nil {
		return nil, fmt.Errorf("ошибка отправки запроса: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("неожиданный статус ответа: %d", resp.StatusCode)
	}

	var response domain.PlacementResponse
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return nil, fmt.Errorf("ошибка разбора ответа: %v", err)
	}

	return &response, nil
} 