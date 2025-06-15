package domain

// PlaceRequest представляет запрос на размещение товара
type PlaceRequest struct {
	ItemID   string  `json:"item_id"`
	BatchID  string  `json:"batch_id"`
	Quantity int     `json:"quantity"`
	Command  string  `json:"command"`


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

type PlaceResponse struct {
	Success bool    `json:"success"`
	SlotID  string  `json:"slot_id,omitempty"` 
	Comment string  `json:"comment"`
	Score   float64 `json:"score"`
}


type Item struct {
	ItemID   string  `json:"item_id"`	
	Mr       float64 `json:"mr"`
}


type Slot struct {
	SlotID         string `json:"slot_id"`
	IsOccupied     bool   `json:"is_occupied"`
	ZoneType       string `json:"zone_type"`
	DistanceFromExit int  `json:"distance_from_exit"`
} 