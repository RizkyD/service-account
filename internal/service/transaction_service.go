package service

import (
	"account-service/internal/model"
	"account-service/internal/repository"
	"context"
	"errors"
	"github.com/jackc/pgx/v5"
)

type TransactionService interface {
	Create(ctx context.Context, dto *model.CreateTransactionDTO) (*model.Transaction, error)
	GetByID(ctx context.Context, id int64) (*model.Transaction, error)
	// GetAll retrieves transactions with the given limit and offset.
	GetAll(ctx context.Context, limit, offset int) ([]model.Transaction, error)
	Update(ctx context.Context, id int64, dto *model.UpdateTransactionDTO) error
	Delete(ctx context.Context, id int64) error
}

type transactionService struct {
	repo repository.TransactionRepository
}

func NewTransactionService(r repository.TransactionRepository) TransactionService {
	return &transactionService{repo: r}
}

func (s *transactionService) Create(ctx context.Context, dto *model.CreateTransactionDTO) (*model.Transaction, error) {
	t := &model.Transaction{
		CustomerID:        dto.CustomerID,
		TransactionAt:     dto.TransactionAt,
		TransactionAmount: dto.TransactionAmount,
	}
	if err := s.repo.Create(ctx, t); err != nil {
		return nil, err
	}
	return t, nil
}

func (s *transactionService) GetByID(ctx context.Context, id int64) (*model.Transaction, error) {
	tx, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if tx == nil {
		return nil, errors.New("transaction not found")
	}
	return tx, nil
}

func (s *transactionService) GetAll(ctx context.Context, limit, offset int) ([]model.Transaction, error) {
	return s.repo.GetAll(ctx, limit, offset)
}

func (s *transactionService) Update(ctx context.Context, id int64, dto *model.UpdateTransactionDTO) error {
	t := &model.Transaction{
		ID:                id,
		CustomerID:        dto.CustomerID,
		TransactionAt:     dto.TransactionAt,
		TransactionAmount: dto.TransactionAmount,
	}
	err := s.repo.Update(ctx, t)
	if err == pgx.ErrNoRows {
		return errors.New("transaction not found")
	}
	return err
}

func (s *transactionService) Delete(ctx context.Context, id int64) error {
	err := s.repo.Delete(ctx, id)
	if err == pgx.ErrNoRows {
		return errors.New("transaction not found")
	}
	return err
}
