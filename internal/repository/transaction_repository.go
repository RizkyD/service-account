package repository

import (
	"account-service/internal/model"
	"context"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type TransactionRepository interface {
	Create(ctx context.Context, tx *model.Transaction) error
	GetByID(ctx context.Context, id int64) (*model.Transaction, error)
	GetAll(ctx context.Context) ([]model.Transaction, error)
	Update(ctx context.Context, transaction *model.Transaction) error
	Delete(ctx context.Context, id int64) error
}

type transactionRepository struct {
	db *pgxpool.Pool
}

func NewTransactionRepository(db *pgxpool.Pool) TransactionRepository {
	return &transactionRepository{db: db}
}

func (r *transactionRepository) Create(ctx context.Context, t *model.Transaction) error {
	query := `INSERT INTO transactions (customer_id, transaction_at, transaction_amount) VALUES ($1, $2, $3) RETURNING id`
	err := r.db.QueryRow(ctx, query, t.CustomerID, t.TransactionAt, t.TransactionAmount).Scan(&t.ID)
	if err != nil {
		return fmt.Errorf("error inserting transaction: %w", err)
	}
	return nil
}

func (r *transactionRepository) GetByID(ctx context.Context, id int64) (*model.Transaction, error) {
	query := `SELECT id, customer_id, transaction_at, transaction_amount FROM transactions WHERE id = $1`
	row := r.db.QueryRow(ctx, query, id)
	var t model.Transaction
	err := row.Scan(&t.ID, &t.CustomerID, &t.TransactionAt, &t.TransactionAmount)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		return nil, fmt.Errorf("error fetching transaction by id: %w", err)
	}
	return &t, nil
}

func (r *transactionRepository) GetAll(ctx context.Context) ([]model.Transaction, error) {
	query := `SELECT id, customer_id, transaction_at, transaction_amount FROM transactions`
	rows, err := r.db.Query(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("error querying transactions: %w", err)
	}
	defer rows.Close()

	var result []model.Transaction
	for rows.Next() {
		var t model.Transaction
		if err = rows.Scan(&t.ID, &t.CustomerID, &t.TransactionAt, &t.TransactionAmount); err != nil {
			return nil, fmt.Errorf("error scanning transaction: %w", err)
		}
		result = append(result, t)
	}
	return result, nil
}

func (r *transactionRepository) Update(ctx context.Context, t *model.Transaction) error {
	query := `UPDATE transactions SET customer_id=$1, transaction_at=$2, transaction_amount=$3 WHERE id=$4`
	cmdTag, err := r.db.Exec(ctx, query, t.CustomerID, t.TransactionAt, t.TransactionAmount, t.ID)
	if err != nil {
		return fmt.Errorf("error updating transaction: %w", err)
	}
	if cmdTag.RowsAffected() == 0 {
		return pgx.ErrNoRows
	}
	return nil
}

func (r *transactionRepository) Delete(ctx context.Context, id int64) error {
	query := `DELETE FROM transactions WHERE id=$1`
	cmdTag, err := r.db.Exec(ctx, query, id)
	if err != nil {
		return fmt.Errorf("error deleting transaction: %w", err)
	}
	if cmdTag.RowsAffected() == 0 {
		return pgx.ErrNoRows
	}
	return nil
}
