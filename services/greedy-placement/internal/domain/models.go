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

// Item представляет товар из базы данных (для общих данных, если потребуются)
type Item struct {
	ItemID   string  `json:"item_id"`	
	// Добавьте другие поля Item, если они будут нужны для логики жадного алгоритма
	// Например, размеры, вес, условия хранения - хотя жадный алгоритм может их не использовать напрямую
}

// Slot представляет ячейку склада с соответствующей информацией
type Slot struct {
	SlotID         string `json:"slot_id"`	
	IsOccupied     bool   `json:"is_occupied"`
	ZoneType       string `json:"zone_type"`
	DistanceFromExit int  `json:"distance_from_exit"` // Ключевое поле для жадного алгоритма
	// Добавьте другие поля Slot, если они будут нужны (например, размеры, условия хранения)
} 