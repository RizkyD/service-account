package repository

import (
	"account-service/internal/model"
	"context"
	"errors"
	"fmt"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
)

type NasabahRepository interface {
	Daftar(ctx context.Context, nasabah *model.Nasabah) error
	GetNasabahForUpdate(ctx context.Context, tx pgx.Tx, id int) (*model.Nasabah, error)
	UpdateSaldoNasabah(ctx context.Context, tx pgx.Tx, updateSaldoNasabahDTO *model.UpdateSaldoNasabahDTO) error
	GetSaldo(ctx context.Context, id int) (*model.Nasabah, error)
}

type nasabahRepository struct {
	db *pgxpool.Pool
}

func NewNasabahRepository(db *pgxpool.Pool) NasabahRepository {
	return &nasabahRepository{db: db}
}

func (r *nasabahRepository) Daftar(ctx context.Context, nasabah *model.Nasabah) error {
	query := `INSERT INTO nasabah (nama, nik, no_hp) VALUES ($1, $2, $3) RETURNING id, created_at, updated_at`
	err := r.db.QueryRow(ctx, query, nasabah.Name, nasabah.NIK, nasabah.Phone).Scan(&nasabah.ID, &nasabah.CreatedAt, &nasabah.UpdatedAt)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == "23505" {
			return errors.New("nik atau nomor telepon telah digunakan")
		}
		return err
	}
	return nil
}

func (r *nasabahRepository) GetNasabahForUpdate(ctx context.Context, tx pgx.Tx, id int) (*model.Nasabah, error) {
	query := `SELECT id, saldo FROM nasabah WHERE id = $1 FOR UPDATE`
	row := tx.QueryRow(ctx, query, id)

	nasabah := &model.Nasabah{}
	err := row.Scan(&nasabah.ID, &nasabah.Saldo)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, errors.New("nasabah tidak ditemukan")
		}
		return nil, fmt.Errorf("error saat mengakses data nasabah untuk transaksi: %w", err)
	}
	return nasabah, nil
}

func (r *nasabahRepository) UpdateSaldoNasabah(ctx context.Context, tx pgx.Tx, updateSaldoNasabahDTO *model.UpdateSaldoNasabahDTO) error {
	query := fmt.Sprintf(`UPDATE nasabah SET saldo = $1, updated_at = CURRENT_TIMESTAMP WHERE id = $2`)
	_, err := tx.Exec(ctx, query, updateSaldoNasabahDTO.Saldo, updateSaldoNasabahDTO.ID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return errors.New("nasabah tidak ditemukan")
		}
		return fmt.Errorf("error saat melakukan update saldo: %w", err)
	}

	return nil
}

func (r *nasabahRepository) GetSaldo(ctx context.Context, id int) (*model.Nasabah, error) {
	query := `SELECT saldo FROM nasabah WHERE id = $1`
	row := r.db.QueryRow(ctx, query, id)

	var nasabah model.Nasabah
	err := row.Scan(&nasabah.Saldo)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, errors.New("nasabah tidak ditemukan")
		}
		return nil, err
	}
	return &nasabah, nil
}
