package model

import "time"

type Nasabah struct {
	ID        int       `json:"id"`
	Name      string    `json:"nama"`
	NIK       string    `json:"nik"`
	Phone     string    `json:"no_hp"`
	Saldo     float64   `json:"saldo"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type DaftarNasabahDTO struct {
	Name        string `json:"nama" validate:"required,min=1,max=256"`
	NIK         string `json:"nik" validate:"required,min=1,max=256"`
	PhoneNumber string `json:"no_hp" validate:"required,min=1,max=256"`
}

type UpdateSaldoNasabahDTO struct {
	ID    int     `json:"no_rekening" validate:"required"`
	Saldo float64 `json:"saldo" validate:"required,min=1"`
}
