package domain

// PlaceRequest represents a request to analyze or place an item
type PlaceRequest struct {
	ItemID   string  `json:"item_id"`
	BatchID  string  `json:"batch_id"`
	Quantity int     `json:"quantity"`
	Command  string  `json:"command"` // "analyze" or "place"

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
}

// PlaceResponse represents the system's response to a placement request
type PlaceResponse struct {
	Success bool    `json:"success"`
	SlotID  string  `json:"slot_id,omitempty"` // Use omitempty for optional fields
	Comment string  `json:"comment"`
	Score   float64 `json:"score"`
}

// Item represents an item from the database, including turnover for ABC analysis
type Item struct {
	ItemID   string  `json:"item_id"`
	Turnover float64 `json:"turnover"`
	ItemType string  `json:"item_type"`
}

// Slot represents a warehouse slot with relevant placement information
type Slot struct {
	SlotID         string `json:"slot_id"`
	IsOccupied     bool   `json:"is_occupied"`
	ZoneType       string `json:"zone_type"`
	DistanceFromExit int  `json:"distance_from_exit"`
} 