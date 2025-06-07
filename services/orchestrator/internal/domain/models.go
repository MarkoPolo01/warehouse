package domain

// PlacementRequest представляет запрос на размещение товара
type PlacementRequest struct {
	// Основные параметры
	ItemID   string `json:"item_id"`
	BatchID  string `json:"batch_id"`
	Quantity int    `json:"quantity"`

	// Параметры товара
	Weight        float64 `json:"weight"`         // вес товара
	Volume        float64 `json:"volume"`         // объем товара
	TurnoverRate  float64 `json:"turnover_rate"`  // коэффициент оборачиваемости
	DemandRate    float64 `json:"demand_rate"`    // коэффициент спроса
	Seasonality   float64 `json:"seasonality"`    // коэффициент сезонности
	ABCClass      string  `json:"abc_class"`      // класс ABC-анализа (A, B, C)
	XYZClass      string  `json:"xyz_class"`      // класс XYZ-анализа (X, Y, Z)
	IsHeavy       bool    `json:"is_heavy"`       // является ли товар тяжелым
	IsFragile     bool    `json:"is_fragile"`     // является ли товар хрупким
	IsHazardous   bool    `json:"is_hazardous"`   // является ли товар опасным
	StorageTemp   float64 `json:"storage_temp"`   // требуемая температура хранения
	StorageHumidity float64 `json:"storage_humidity"` // требуемая влажность хранения

	// Параметры склада
	WarehouseLoad float64 `json:"warehouse_load"` // текущая загрузка склада (0-1)
	HasFixedSlot  bool    `json:"has_fixed_slot"` // есть ли фиксированная ячейка
	FastAccessZone bool   `json:"fast_access_zone"` // требуется ли размещение в fast-access зоне
}

// PlacementResponse представляет ответ от микросервиса размещения
type PlacementResponse struct {
	Success bool    `json:"success"`
	SlotID  string  `json:"slot_id"`
	Comment string  `json:"comment"`
	Score   float64 `json:"score"`
	
	// Факторы для расчета итоговой оценки
	BaseScore          float64 `json:"base_score"`           // S0 - базовая оценка (0.90, 0.75, 0.60, 0.50, 0.00)
	ResponseTimeMs     int64   `json:"response_time_ms"`     // для расчета Δt
	DistanceToExit     float64 `json:"distance_to_exit"`     // для расчета Δd
	HasFixedSlot      bool    `json:"has_fixed_slot"`       // +0.07
	HighWarehouseLoad bool    `json:"high_warehouse_load"`  // +0.05
	HighTurnover      bool    `json:"high_turnover"`        // +0.05
	HeavyItem         bool    `json:"heavy_item"`           // +0.03
	NoPlacementHistory bool   `json:"no_placement_history"` // +0.02
	FastAccessZone    bool    `json:"fast_access_zone"`     // +0.03
	XYZCompliant      bool    `json:"xyz_compliant"`        // +0.02
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