package schema

import (
	"time"
)

// CreateSiswaRequest ...
type CreateSiswaRequest struct {
	Nama        string `json:"nama" validate:"required"`
	IDKelas     int    `json:"id_kelas" validate:"required"`
	IDWaliKelas int    `json:"id_wali_kelas" validate:"required"`
	Tingkat     int    `json:"tingkat" validate:"required"`
	Alamat      string `json:"alamat" validate:"required"`
}

// SiswaResponse ...
type SiswaResponse struct {
	ID          int        `json:"id" db:"id"`
	Nama        string     `json:"nama" db:"nama"`
	IDKelas     int        `json:"id_kelas" db:"id_kelas"`
	IDWaliKelas int        `json:"id_wali_kelas" db:"id_wali_kelas"`
	Tingkat     int        `json:"tingkat" db:"tingkat"`
	Alamat      string     `json:"alamat" db:"alamat"`
	CreatedAt   *time.Time `json:"created_at,omitempty" db:"created_at"`
	UpdatedAt   *time.Time `json:"updated_at,omitempty" db:"updated_at"`
}

// UpdateSiswaRequest ...
type UpdateSiswaRequest struct {
	Nama        string `json:"nama"`
	IDKelas     int    `json:"id_kelas"`
	IDWaliKelas int    `json:"id_wali_kelas"`
	Tingkat     int    `json:"tingkat"`
	Alamat      string `json:"alamat"`
}
