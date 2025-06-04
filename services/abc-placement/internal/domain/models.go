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