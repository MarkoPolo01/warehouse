package domain

// PlaceRequest представляет запрос на размещение товара
type PlaceRequest struct {
	ItemID   string `json:"item_id"`
	BatchID  string `json:"batch_id"`
	Quantity int    `json:"quantity"`
	Command  string `json:"command"` // "analyze" или "place"
}

// PlaceResponse представляет ответ на запрос размещения
type PlaceResponse struct {
	Success bool    `json:"success"`
	SlotID  string  `json:"slot_id"`
	Comment string  `json:"comment"`
	Score   float64 `json:"score"`
} 