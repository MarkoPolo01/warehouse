package models

import "time"

type Product struct {
	ID          int64     `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	Quantity    int       `json:"quantity"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

type Warehouse struct {
	ID        int64     `json:"id"`
	Name      string    `json:"name"`
	Location  string    `json:"location"`
	Capacity  int       `json:"capacity"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type Inventory struct {
	ID          int64     `json:"id"`
	ProductID   int64     `json:"product_id"`
	WarehouseID int64     `json:"warehouse_id"`
	Quantity    int       `json:"quantity"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
} 