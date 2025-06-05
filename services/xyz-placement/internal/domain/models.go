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
	SlotID  string  `json:"slot_id,omitempty"` // Use omitempty for optional fields
	Comment string  `json:"comment"`
	Score   float64 `json:"score"`
}

// Item представляет товар из базы данных, включая мr для XYZ-анализа
type Item struct {
	ItemID   string  `json:"item_id"`	
	Mr       float64 `json:"mr"` // Коэффициент вариации для XYZ
}

// Slot представляет ячейку склада с соответствующей информацией
type Slot struct {
	SlotID         string `json:"slot_id"`
	IsOccupied     bool   `json:"is_occupied"`
	ZoneType       string `json:"zone_type"`
	DistanceFromExit int  `json:"distance_from_exit"`
} 