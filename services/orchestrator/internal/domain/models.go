package domain

// PlacementRequest представляет запрос на размещение товара
type PlacementRequest struct {
	ItemID   string `json:"item_id"`
	BatchID  string `json:"batch_id"`
	Quantity int    `json:"quantity"`
	Command  string `json:"command"` // "analyze" или "place"
}

// PlacementResponse представляет ответ от микросервиса размещения
type PlacementResponse struct {
	Success bool    `json:"success"`
	SlotID  string  `json:"slot_id"`
	Comment string  `json:"comment"`
	Score   float64 `json:"score"`
}

// OrchestratorResponse представляет итоговый ответ оркестратора
type OrchestratorResponse struct {
	Success     bool             `json:"success"`
	SlotID      string          `json:"slot_id"`
	Comment     string          `json:"comment"`
	Score       float64         `json:"score"`
	Algorithm   string          `json:"algorithm"`
	AllResults  []ServiceResult `json:"all_results"`
}

// ServiceResult представляет результат от конкретного микросервиса
type ServiceResult struct {
	ServiceName string          `json:"service_name"`
	Response    PlacementResponse `json:"response"`
} 