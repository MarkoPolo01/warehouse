package domain

// PlaceRequest represents a request to analyze or place an item
type PlaceRequest struct {
	ItemID   string `json:"item_id"`
	BatchID  string `json:"batch_id"`
	Quantity int    `json:"quantity"`
	Command  string `json:"command"` // "analyze" or "place"
}

// PlaceResponse represents the system's response to a placement request
type PlaceResponse struct {
	Success bool    `json:"success"`
	SlotID  string  `json:"slot_id,omitempty"`
	Comment string  `json:"comment"`
	Score   float64 `json:"score"`
}

// Item represents an item from the database with relevant characteristics
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

// Slot represents a warehouse slot with relevant information
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
	DistanceFromExit   int     `json:"distance_from_exit"` // for logistics
}

// PlacementCandidate represents a possible item-slot placement combination
type PlacementCandidate struct {
	Item    *Item
	Slot    *Slot
	Fitness float64 // The calculated fitness score for this candidate
}