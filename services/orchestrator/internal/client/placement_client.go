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
	// Преобразуем расширенный запрос в формат для микросервисов
	serviceReq := struct {
		ItemID   string  `json:"item_id"`
		BatchID  string  `json:"batch_id"`
		Quantity int     `json:"quantity"`
		Command  string  `json:"command"`

		// Параметры товара
		Weight        float64 `json:"weight"`
		Volume        float64 `json:"volume"`
		TurnoverRate  float64 `json:"turnover_rate"`
		DemandRate    float64 `json:"demand_rate"`
		Seasonality   float64 `json:"seasonality"`
		ABCClass      string  `json:"abc_class"`
		XYZClass      string  `json:"xyz_class"`
		IsHeavy       bool    `json:"is_heavy"`
		IsFragile     bool    `json:"is_fragile"`
		IsHazardous   bool    `json:"is_hazardous"`
		StorageTemp   float64 `json:"storage_temp"`
		StorageHumidity float64 `json:"storage_humidity"`

		// Параметры склада
		WarehouseLoad float64 `json:"warehouse_load"`
		HasFixedSlot  bool    `json:"has_fixed_slot"`
		FastAccessZone bool   `json:"fast_access_zone"`
	}{
		ItemID:   req.ItemID,
		BatchID:  req.BatchID,
		Quantity: req.Quantity,
		Command:  "analyze",

		// Параметры товара
		Weight:         req.Weight,
		Volume:         req.Volume,
		TurnoverRate:   req.TurnoverRate,
		DemandRate:     req.DemandRate,
		Seasonality:    req.Seasonality,
		ABCClass:       req.ABCClass,
		XYZClass:       req.XYZClass,
		IsHeavy:        req.IsHeavy,
		IsFragile:      req.IsFragile,
		IsHazardous:    req.IsHazardous,
		StorageTemp:    req.StorageTemp,
		StorageHumidity: req.StorageHumidity,

		// Параметры склада
		WarehouseLoad:  req.WarehouseLoad,
		HasFixedSlot:   req.HasFixedSlot,
		FastAccessZone: req.FastAccessZone,
	}

	startTime := time.Now()
	resp, err := c.sendRequest(ctx, "/api/v1/placement", serviceReq)
	if err != nil {
		return nil, err
	}

	// Добавляем время ответа
	resp.ResponseTimeMs = time.Since(startTime).Milliseconds()

	// Добавляем дополнительные факторы из запроса
	resp.HasFixedSlot = req.HasFixedSlot
	resp.HighWarehouseLoad = req.WarehouseLoad > 0.8
	resp.HighTurnover = req.TurnoverRate > 0.8
	resp.HeavyItem = req.IsHeavy
	resp.NoPlacementHistory = false // TODO: получать из БД
	resp.FastAccessZone = req.FastAccessZone
	resp.XYZCompliant = req.XYZClass == "X" // X - стабильный спрос

	return resp, nil
}

// PlaceItem отправляет запрос на размещение товара
func (c *PlacementClient) PlaceItem(ctx context.Context, req *domain.PlacementRequest) (*domain.PlacementResponse, error) {
	// Преобразуем расширенный запрос в формат для микросервисов
	serviceReq := struct {
		ItemID   string  `json:"item_id"`
		BatchID  string  `json:"batch_id"`
		Quantity int     `json:"quantity"`
		Command  string  `json:"command"`

		// Параметры товара
		Weight        float64 `json:"weight"`
		Volume        float64 `json:"volume"`
		TurnoverRate  float64 `json:"turnover_rate"`
		DemandRate    float64 `json:"demand_rate"`
		Seasonality   float64 `json:"seasonality"`
		ABCClass      string  `json:"abc_class"`
		XYZClass      string  `json:"xyz_class"`
		IsHeavy       bool    `json:"is_heavy"`
		IsFragile     bool    `json:"is_fragile"`
		IsHazardous   bool    `json:"is_hazardous"`
		StorageTemp   float64 `json:"storage_temp"`
		StorageHumidity float64 `json:"storage_humidity"`

		// Параметры склада
		WarehouseLoad float64 `json:"warehouse_load"`
		HasFixedSlot  bool    `json:"has_fixed_slot"`
		FastAccessZone bool   `json:"fast_access_zone"`
	}{
		ItemID:   req.ItemID,
		BatchID:  req.BatchID,
		Quantity: req.Quantity,
		Command:  "place",

		// Параметры товара
		Weight:         req.Weight,
		Volume:         req.Volume,
		TurnoverRate:   req.TurnoverRate,
		DemandRate:     req.DemandRate,
		Seasonality:    req.Seasonality,
		ABCClass:       req.ABCClass,
		XYZClass:       req.XYZClass,
		IsHeavy:        req.IsHeavy,
		IsFragile:      req.IsFragile,
		IsHazardous:    req.IsHazardous,
		StorageTemp:    req.StorageTemp,
		StorageHumidity: req.StorageHumidity,

		// Параметры склада
		WarehouseLoad:  req.WarehouseLoad,
		HasFixedSlot:   req.HasFixedSlot,
		FastAccessZone: req.FastAccessZone,
	}

	return c.sendRequest(ctx, "/api/v1/placement", serviceReq)
}

// sendRequest отправляет запрос к микросервису
func (c *PlacementClient) sendRequest(ctx context.Context, path string, req interface{}) (*domain.PlacementResponse, error) {
	jsonData, err := json.Marshal(req)
	if err != nil {
		return nil, fmt.Errorf("ошибка сериализации запроса: %w", err)
	}

	// Используем базовый URL без дополнительного пути
	request, err := http.NewRequestWithContext(ctx, "POST", c.baseURL, bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, fmt.Errorf("ошибка создания запроса: %w", err)
	}

	request.Header.Set("Content-Type", "application/json")

	response, err := c.client.Do(request)
	if err != nil {
		return nil, fmt.Errorf("ошибка отправки запроса: %w", err)
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("ошибка сервера: %d", response.StatusCode)
	}

	var resp domain.PlacementResponse
	if err := json.NewDecoder(response.Body).Decode(&resp); err != nil {
		return nil, fmt.Errorf("ошибка десериализации ответа: %w", err)
	}

	return &resp, nil
} 