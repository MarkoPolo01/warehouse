package repository

import (
	"context"
	"database/sql"
	"fmt"

	"warehouse/services/abc-placement/internal/domain"
)

// PostgresRepository implements the Repository interface for PostgreSQL
type PostgresRepository struct {
	db *sql.DB
}

// NewPostgresRepository creates a new instance of PostgresRepository
func NewPostgresRepository(db *sql.DB) *PostgresRepository {
	return &PostgresRepository{db: db}
}

// ItemExists checks if an item exists in the database
func (r *PostgresRepository) ItemExists(ctx context.Context, itemID string) (bool, error) {
	var exists bool
	err := r.db.QueryRowContext(ctx, "SELECT EXISTS(SELECT 1 FROM items WHERE item_id = $1)", itemID).Scan(&exists)
	return exists, err
}

// BatchExists checks if a batch exists in the database
func (r *PostgresRepository) BatchExists(ctx context.Context, batchID string) (bool, error) {
	var exists bool
	err := r.db.QueryRowContext(ctx, "SELECT EXISTS(SELECT 1 FROM batches WHERE batch_id = $1)", batchID).Scan(&exists)
	return exists, err
}

// GetItemDetails retrieves item details for ABC analysis
func (r *PostgresRepository) GetItemDetails(ctx context.Context, itemID string) (*domain.Item, error) {
	var item domain.Item
	err := r.db.QueryRowContext(ctx, "SELECT item_id, turnover, item_type FROM items WHERE item_id = $1", itemID).Scan(&item.ItemID, &item.Turnover, &item.ItemType)
	if err == sql.ErrNoRows {
		return nil, nil // Item not found
	}
	return &item, err
}

// GetAvailableSlots retrieves available slots based on zone type and distance from exit
func (r *PostgresRepository) GetAvailableSlots(ctx context.Context, zoneType string) ([]domain.Slot, error) {
	rows, err := r.db.QueryContext(ctx, "SELECT slot_id, is_occupied, zone_type, distance_from_exit FROM slots WHERE is_occupied = false AND zone_type = $1 ORDER BY distance_from_exit ASC", zoneType)
	if err != nil {
		return nil, fmt.Errorf("failed to get available slots: %w", err)
	}
	defer rows.Close()

	var slots []domain.Slot
	for rows.Next() {
		var slot domain.Slot
		if err := rows.Scan(&slot.SlotID, &slot.IsOccupied, &slot.ZoneType, &slot.DistanceFromExit); err != nil {
			return nil, fmt.Errorf("failed to scan slot row: %w", err)
		}
		slots = append(slots, slot)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("row iteration error: %w", err)
	}

	return slots, nil
}

// IsSlotOccupied checks if a slot is occupied
func (r *PostgresRepository) IsSlotOccupied(ctx context.Context, slotID string) (bool, error) {
	var isOccupied bool
	err := r.db.QueryRowContext(ctx, "SELECT is_occupied FROM slots WHERE slot_id = $1", slotID).Scan(&isOccupied)
	if err == sql.ErrNoRows {
		return false, nil // Slot not found, consider it not occupied for placement purposes
	}
	return isOccupied, err
}

// CreatePlacementRequest creates a new placement request entry and returns its ID
func (r *PostgresRepository) CreatePlacementRequest(ctx context.Context, req *domain.PlaceRequest) (int, error) {
	var requestID int
	err := r.db.QueryRowContext(ctx,
		"INSERT INTO placement_requests (item_id, batch_id, quantity) VALUES ($1, $2, $3) RETURNING request_id",
		req.ItemID, req.BatchID, req.Quantity,
	).Scan(&requestID)
	return requestID, err
}

// CreatePlacementLog creates a log entry for a placement
func (r *PostgresRepository) CreatePlacementLog(ctx context.Context, slotID, itemID, batchID, algorithm string) error {
	_, err := r.db.ExecContext(ctx,
		"INSERT INTO placement_logs (slot_id, item_id, batch_id, algorithm) VALUES ($1, $2, $3, $4)",
		slotID, itemID, batchID, algorithm,
	)
	return err
}

// UpdateSlotOccupation updates the occupation status of a slot
func (r *PostgresRepository) UpdateSlotOccupation(ctx context.Context, slotID string, isOccupied bool) error {
	_, err := r.db.ExecContext(ctx, "UPDATE slots SET is_occupied = $1 WHERE slot_id = $2", isOccupied, slotID)
	return err
}

// CreatePlacementResponse creates a response entry for a placement request
func (r *PostgresRepository) CreatePlacementResponse(ctx context.Context, requestID int, success bool, slotID, algorithm string, score float64, comment string) error {
	_, err := r.db.ExecContext(ctx,
		"INSERT INTO placement_responses (request_id, success, slot_id, algorithm_used, score, comment) VALUES ($1, $2, $3, $4, $5, $6)",
		requestID, success, slotID, algorithm, score, comment,
	)
	return err
} 