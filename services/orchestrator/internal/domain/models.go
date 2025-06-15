package domain

// PlacementRequest представляет запрос на размещение товара
type PlacementRequest struct {
	// Основные параметры
	ItemID   string `json:"item_id"`
	BatchID  string `json:"batch_id"`
	Quantity int    `json:"quantity"`

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


	WarehouseLoad float64 `json:"warehouse_load"` 
	HasFixedSlot  bool    `json:"has_fixed_slot"` 
	FastAccessZone bool   `json:"fast_access_zone"` 
}


type PlacementResponse struct {
	Success bool    `json:"success"`
	SlotID  string  `json:"slot_id"`
	Comment string  `json:"comment"`
	Score   float64 `json:"score"`
	

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


type OrchestratorResponse struct {
	Success     bool             `json:"success"`
	SlotID      string          `json:"slot_id"`
	Comment     string          `json:"comment"`
	Score       float64         `json:"score"`
	Algorithm   string          `json:"algorithm"`
	AllResults  []ServiceResult `json:"all_results"`
}

type ServiceResult struct {
	ServiceName string          `json:"service_name"`
	Response    PlacementResponse `json:"response"`
} 