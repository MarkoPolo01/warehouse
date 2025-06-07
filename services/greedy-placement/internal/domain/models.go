package domain

// PlaceRequest представляет запрос на размещение товара
type PlaceRequest struct {
	ItemID   string  `json:"item_id"`
	BatchID  string  `json:"batch_id"`
	Quantity int     `json:"quantity"`
	Command  string  `json:"command"` // "analyze" или "place"

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