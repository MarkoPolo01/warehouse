package repository

import (
	"context"
	"database/sql"
	"fmt"

	"warehouse/services/greedy-placement/internal/domain"

	_ "github.com/lib/pq"
)


type PostgresRepository struct {
	db *sql.DB
}


func NewPostgresRepository(db *sql.DB) *PostgresRepository {
	return &PostgresRepository{db: db}
}

func (r *PostgresRepository) ItemExists(ctx context.Context, itemID string) (bool, error) {
	var exists bool
	err := r.db.QueryRowContext(ctx, "SELECT EXISTS(SELECT 1 FROM items WHERE item_id = $1)", itemID).Scan(&exists)
	return exists, err
}


func (r *PostgresRepository) BatchExists(ctx context.Context, batchID string) (bool, error) {
	var exists bool
	err := r.db.QueryRowContext(ctx, "SELECT EXISTS(SELECT 1 FROM batches WHERE batch_id = $1)", batchID).Scan(&exists)
	return exists, err
}


func (r *PostgresRepository) GetAllAvailableSlotsOrderedByDistance(ctx context.Context) ([]domain.Slot, error) {
	rows, err := r.db.QueryContext(ctx, "SELECT slot_id, is_occupied, zone_type, distance_from_exit FROM slots WHERE is_occupied = false ORDER BY distance_from_exit ASC")
	if err != nil {
		return nil, fmt.Errorf("failed to get all available slots: %w", err)
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


func (r *PostgresRepository) CreatePlacementRequest(ctx context.Context, req *domain.PlaceRequest) (int, error) {
	var requestID int
	err := r.db.QueryRowContext(ctx,
		"INSERT INTO placement_requests (item_id, batch_id, quantity) VALUES ($1, $2, $3) RETURNING request_id",
		req.ItemID, req.BatchID, req.Quantity,
	).Scan(&requestID)
	return requestID, err
}


func (r *PostgresRepository) UpdateSlotOccupation(ctx context.Context, slotID string, isOccupied bool) error {
	_, err := r.db.ExecContext(ctx, "UPDATE slots SET is_occupied = $1 WHERE slot_id = $2", isOccupied, slotID)
	return err
}


func (r *PostgresRepository) CreatePlacementLog(ctx context.Context, slotID, itemID, batchID, algorithm string) error {
	_, err := r.db.ExecContext(ctx,
		"INSERT INTO placement_logs (slot_id, item_id, batch_id, algorithm) VALUES ($1, $2, $3, $4)",
		slotID, itemID, batchID, algorithm,
	)
	return err
}


func (r *PostgresRepository) CreatePlacementResponse(ctx context.Context, requestID int, success bool, slotID, algorithm string, score float64, comment string) error {
	_, err := r.db.ExecContext(ctx,
		"INSERT INTO placement_responses (request_id, success, slot_id, algorithm_used, score, comment) VALUES ($1, $2, $3, $4, $5, $6)",
		requestID, success, slotID, algorithm, score, comment,
	)
	return err
} 