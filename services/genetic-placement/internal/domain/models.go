package domain

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
	ItemID             string  `json:"item_id"`
	Name               string  `json:"name"`
	ItemType           string  `json:"item_type"`
	Weight             float64 `json:"weight"`
	Length             float64 `json:"length"`
	Width              float64 `json:"width"`
	Height             float64 `json:"height"`
	StorageConditions  string  `json:"storage_conditions"`
	LabelType          string  `json:"label_type"`
	Turnover           float64 `json:"turnover"` // for ABC
	Mr                 float64 `json:"mr"`       // for XYZ
}


type Slot struct {
	SlotID             string  `json:"slot_id"`
	LocationDescription string  `json:"location_description"`
	MaxWeight          float64 `json:"max_weight"`
	MaxLength          float64 `json:"max_length"`
	MaxWidth           float64 `json:"max_width"`
	MaxHeight          float64 `json:"max_height"`
	StorageConditions  string  `json:"storage_conditions"`
	IsOccupied         bool    `json:"is_occupied"`
	ZoneType           string  `json:"zone_type"`
	Level              int     `json:"level"`
	DistanceFromExit   int     `json:"distance_from_exit"`
}

type PlacementCandidate struct {
	Item    *Item
	Slot    *Slot
	Fitness float64
}