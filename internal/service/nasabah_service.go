package service

import (
	"account-service/internal/model"
	"account-service/internal/repository"
	"context"
	"errors"
	"fmt"
	"github.com/jackc/pgx/v5/pgxpool"
	zlog "github.com/rs/zerolog/log"
)

type NasabahService interface {
	Daftar(ctx context.Context, nasabahDTO *model.DaftarNasabahDTO) (*model.Nasabah, error)
	UpdateSaldo(ctx context.Context, updateSaldoNasabahDTO *model.UpdateSaldoNasabahDTO, isDeposit bool) error
	GetSaldo(ctx context.Context, id int) (*model.Nasabah, error)
}

type nasabahService struct {
	nasabahRepo repository.NasabahRepository
	dbPool      *pgxpool.Pool
}

func NewNasabahService(repo repository.NasabahRepository, db *pgxpool.Pool) NasabahService {
	return &nasabahService{nasabahRepo: repo, dbPool: db}
}

func (s *nasabahService) Daftar(ctx context.Context, nasabahDTO *model.DaftarNasabahDTO) (*model.Nasabah, error) {
	nasabah := &model.Nasabah{
		Name:  nasabahDTO.Name,
		NIK:   nasabahDTO.NIK,
		Phone: nasabahDTO.PhoneNumber,
	}

	err := s.nasabahRepo.Daftar(ctx, nasabah)
	if err != nil {
		return nil, err
	}

	return nasabah, nil
}

func (s *nasabahService) UpdateSaldo(ctx context.Context, updateSaldoNasabahDTO *model.UpdateSaldoNasabahDTO, isDeposit bool) error {
	tx, err := s.dbPool.Begin(ctx)
	if err != nil {
		return fmt.Errorf("error saat memulai koneksi transaksi : %w", err)
	}

	defer func() {
		if err != nil {
			rbErr := tx.Rollback(ctx)
			if rbErr != nil {
				zlog.Error().Err(rbErr).Msg("error saat melakukan rollback")
			}
		}
	}()

	var nasabahFromDB *model.Nasabah
	nasabahFromDB, err = s.nasabahRepo.GetNasabahForUpdate(ctx, tx, updateSaldoNasabahDTO.ID)
	if err != nil {
		return err
	}
	if nasabahFromDB == nil {
		err = errors.New("nasabah tidak ditemukan")
		return err
	}

	if isDeposit {
		updateSaldoNasabahDTO.Saldo += nasabahFromDB.Saldo
	} else {
		updateSaldoNasabahDTO.Saldo = nasabahFromDB.Saldo - updateSaldoNasabahDTO.Saldo
	}

	if updateSaldoNasabahDTO.Saldo < 0 {
		err = errors.New("saldo anda tidak mencukupi")
		return err
	}

	err = s.nasabahRepo.UpdateSaldoNasabah(ctx, tx, updateSaldoNasabahDTO)
	if err != nil {
		return err
	}

	err = tx.Commit(ctx)
	if err != nil {
		return err
	}

	return nil
}

func (s *nasabahService) GetSaldo(ctx context.Context, id int) (*model.Nasabah, error) {
	dataNasabah, err := s.nasabahRepo.GetSaldo(ctx, id)
	if err != nil {
		return nil, err
	}
	if dataNasabah == nil {
		return nil, errors.New("nasabah tidak ditemukan")
	}
	return dataNasabah, nil
}
