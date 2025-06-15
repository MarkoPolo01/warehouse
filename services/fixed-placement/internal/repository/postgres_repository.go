package repository

import (
	"context"
	"database/sql"
	"warehouse/services/fixed-placement/internal/domain"
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

func (r *PostgresRepository) GetFixedSlot(ctx context.Context, itemID string) (string, error) {
	var slotID string
	err := r.db.QueryRowContext(ctx, "SELECT slot_id FROM item_slot_map WHERE item_id = $1", itemID).Scan(&slotID)
	if err == sql.ErrNoRows {
		return "", nil
	}
	return slotID, err
}

func (r *PostgresRepository) IsSlotOccupied(ctx context.Context, slotID string) (bool, error) {
	var isOccupied bool
	err := r.db.QueryRowContext(ctx, "SELECT is_occupied FROM slots WHERE slot_id = $1", slotID).Scan(&isOccupied)
	return isOccupied, err
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