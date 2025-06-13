package model

import "time"

// Transaction entity
// used for table transactions
// id bigserial, customer_id bigint, transaction_at timestamp, transaction_amount float8

// Transaction struct maps to table transactions
// To ensure JSON serialization, we add struct tags.
type Transaction struct {
	ID                int64     `json:"id"`
	CustomerID        int64     `json:"customer_id"`
	TransactionAt     time.Time `json:"transaction_at"`
	TransactionAmount float64   `json:"transaction_amount"`
}

type CreateTransactionDTO struct {
	CustomerID        int64     `json:"customer_id" validate:"required"`
	TransactionAt     time.Time `json:"transaction_at" validate:"required"`
	TransactionAmount float64   `json:"transaction_amount" validate:"required"`
}

type UpdateTransactionDTO struct {
	CustomerID        int64     `json:"customer_id" validate:"required"`
	TransactionAt     time.Time `json:"transaction_at" validate:"required"`
	TransactionAmount float64   `json:"transaction_amount" validate:"required"`
}
