package service

import (
	"net/http"
	"strings"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"

	"github.com/syukur91/ischool-monitor/api/schema"
	"github.com/syukur91/ischool-monitor/pkg/apierror"
	"github.com/syukur91/ischool-monitor/pkg/query"
)

// SiswaService ...
type SiswaService struct {
	db *sqlx.DB
}

// NewSiswaService ...
func NewSiswaService(db *sqlx.DB) *SiswaService {
	return &SiswaService{db: db}
}

// CreateSiswa ...
func (s *SiswaService) CreateSiswa(request *schema.CreateSiswaRequest) (*schema.SiswaResponse, error) {
	if request.Nama == "" {
		return nil, apierror.NewError(http.StatusBadRequest, http.StatusBadRequest, "Siswa nama is not set", errors.New("createsiswa: siswa nama is not set"))
	}

	if request.IDKelas == 0 {
		return nil, apierror.NewError(http.StatusBadRequest, http.StatusBadRequest, "Siswa id kelas is not set", errors.New("createsiswa: siswa id kelas is not set"))
	}

	if request.IDWaliKelas == 0 {
		return nil, apierror.NewError(http.StatusBadRequest, http.StatusBadRequest, "Siswa id wali kelas is not set", errors.New("createsiswa: siswa id wali kelas is not set"))
	}

	if request.Tingkat == 0 {
		return nil, apierror.NewError(http.StatusBadRequest, http.StatusBadRequest, "Siswa tingkat is not set", errors.New("createsiswa: siswa tingkat is not set"))
	}

	if request.Alamat == "" {
		return nil, apierror.NewError(http.StatusBadRequest, http.StatusBadRequest, "Siswa alamat is not set", errors.New("createsiswa: siswa alamat is not set"))
	}

	tx, err := s.db.Begin()
	if err != nil {
		return nil, apierror.NewError(http.StatusInternalServerError, http.StatusInternalServerError, "Database transaction failed", errors.Wrap(err, "createsiswa: begin transaction failed"))
	}

	id := 0
	var createdAt time.Time

	{
		stmt, err := tx.Prepare(`
			INSERT INTO public.siswa (nama, id_kelas, id_wali_kelas, tingkat, alamat )
			VALUES($1, $2, $3, $4, $5)
			RETURNING id, created_at;
		`)

		if err != nil {
			tx.Rollback()
			return nil, apierror.NewError(http.StatusInternalServerError, http.StatusInternalServerError, "Database transaction failed", errors.Wrap(err, "createsiswa: prepare insert statement failed"))
		}
		defer stmt.Close()

		err = stmt.QueryRow(request.Nama, request.IDKelas, request.IDWaliKelas, request.Tingkat, request.Alamat).Scan(&id, &createdAt)
		if err != nil {
			tx.Rollback()

			if strings.Index(err.Error(), "duplicate key value violates unique constraint \"siswa_name_unique\"") > -1 {
				return nil, apierror.NewError(http.StatusBadRequest, http.StatusBadRequest, "Siswa with same nama already exists. Use different nama", errors.Wrap(err, "createsiswa: Siswa with same nama already exists"))
			}

			return nil, apierror.NewError(http.StatusInternalServerError, http.StatusInternalServerError, "Database transaction failed", errors.Wrap(err, "createsiswa: exec insert statement failed"))
		}
	}

	err = tx.Commit()
	if err != nil {
		return nil, apierror.NewError(http.StatusInternalServerError, http.StatusInternalServerError, "Database transaction failed", errors.Wrap(err, "createsiswa: commit transaction failed"))
	}

	return &schema.SiswaResponse{
		ID:          id,
		IDKelas:     request.IDKelas,
		IDWaliKelas: request.IDWaliKelas,
		Tingkat:     request.Tingkat,
		Alamat:      request.Alamat,
		CreatedAt:   &createdAt,
	}, nil
}

// GetSiswa ...
func (s *SiswaService) GetSiswa(id string) (*schema.SiswaResponse, error) {
	if id == "" {
		return nil, apierror.NewError(http.StatusBadRequest, http.StatusBadRequest, "Siswa id is not set", errors.New("getsiswa: siswa id is not set"))
	}

	tx, err := s.db.Beginx()
	if err != nil {
		return nil, apierror.NewError(http.StatusInternalServerError, http.StatusInternalServerError, "Database transaction failed", errors.Wrap(err, "getsiswa: begin transaction failed"))
	}

	siswa := schema.SiswaResponse{}
	{
		err := tx.Get(&siswa, `
			SELECT id,nama,id_kelas,id_wali_kelas,tingkat,alamat,created_at,updated_at 
			FROM public.siswa 
			WHERE id=$1;`,
			id)

		if err != nil {
			tx.Rollback()

			if strings.Index(err.Error(), "sql: no rows in result set") > -1 {
				return nil, apierror.NewError(http.StatusNotFound, http.StatusNotFound, "Siswa with id: "+id+" is not exists", errors.Wrap(err, "getsiswa: siswa with id: "+id+" is not exists"))
			}

			return nil, apierror.NewError(http.StatusInternalServerError, http.StatusInternalServerError, "Database transaction failed", errors.Wrap(err, "getsiswa: get data failed"))
		}
	}

	err = tx.Commit()
	if err != nil {
		return nil, apierror.NewError(http.StatusInternalServerError, http.StatusInternalServerError, "Database transaction failed", errors.Wrap(err, "getsiswa: commit transaction failed"))
	}

	return &siswa, nil
}

// ListSiswas ...
func (s *SiswaService) ListSiswas(gridParams *query.GridParams) ([]schema.SiswaResponse, int, error) {

	tx, err := s.db.Beginx()
	if err != nil {
		return nil, 0, apierror.NewError(http.StatusInternalServerError, http.StatusInternalServerError, "Database transaction failed", errors.Wrap(err, "listsiswa: begin transaction failed"))
	}

	siswas := []schema.SiswaResponse{}
	total := 0
	{
		dataStatement := "SELECT id,nama,id_kelas,id_wali_kelas,tingkat,alamat,created_at,updated_at FROM public.siswa"
		dataQuery, dataParams := query.FullQuery(gridParams, "", nil)
		err := tx.Select(&siswas, dataStatement+dataQuery, dataParams...)
		if err != nil {
			tx.Rollback()
			return nil, 0, apierror.NewError(http.StatusInternalServerError, http.StatusInternalServerError, "Database transaction failed", errors.Wrap(err, "listsiswa: get data failed"))
		}

		countStatement := "SELECT count(*) FROM public.siswa"
		countQuery, countParams := query.FilterQuery(gridParams, "", nil)
		err = tx.QueryRow(countStatement+countQuery, countParams...).Scan(&total)
		if err != nil {
			tx.Rollback()
			return nil, 0, apierror.NewError(http.StatusInternalServerError, http.StatusInternalServerError, "Database transaction failed", errors.Wrap(err, "listsiswa: get count failed"))
		}
	}

	err = tx.Commit()
	if err != nil {
		return nil, 0, apierror.NewError(http.StatusInternalServerError, http.StatusInternalServerError, "Database transaction failed", errors.Wrap(err, "getsiswa: commit transaction failed"))
	}

	return siswas, total, nil
}

// UpdateSiswa ...
func (s *SiswaService) UpdateSiswa(id string, request *schema.UpdateSiswaRequest) (*schema.SiswaResponse, error) {
	if id == "" {
		return nil, apierror.NewError(http.StatusBadRequest, http.StatusBadRequest, "Siswa id is not set", errors.New("updatesiswa: siswa id is not set"))
	}

	tx, err := s.db.Beginx()
	if err != nil {
		return nil, apierror.NewError(http.StatusInternalServerError, http.StatusInternalServerError, "Database transaction failed", errors.Wrap(err, "updatesiswa: begin transaction failed"))
	}

	// get existing siswa
	siswa := schema.SiswaResponse{}
	{
		err := tx.Get(&siswa, `
			SELECT id,nama,id_kelas,id_wali_kelas,tingkat,alamat,created_at,updated_at 
			FROM public.siswa 
			WHERE id=$1;`,
			id)

		if err != nil {
			tx.Rollback()

			if strings.Index(err.Error(), "sql: no rows in result set") > -1 {
				return nil, apierror.NewError(http.StatusNotFound, http.StatusNotFound, "Siswa with id: "+id+" is not exists", errors.Wrap(err, "updatesiswa: siswa with id: "+id+" is not exists"))
			}

			return nil, apierror.NewError(http.StatusInternalServerError, http.StatusInternalServerError, "Database transaction failed", errors.Wrap(err, "updatesiswa: get data failed"))
		}
	}

	// update siswa
	var updatedAt time.Time
	{
		// only update if not empty
		if request.Nama != "" {
			siswa.Nama = request.Nama
		}

		if request.IDKelas != 0 {
			siswa.IDKelas = request.IDKelas
		}

		if request.IDWaliKelas != 0 {
			siswa.IDWaliKelas = request.IDWaliKelas
		}

		if request.Tingkat != 0 {
			siswa.Tingkat = request.Tingkat
		}

		if request.Alamat != "" {
			siswa.Alamat = request.Alamat
		}

		err := tx.QueryRow(`
			UPDATE public.siswa SET nama=$1,id_kelas=$2,id_wali_kelas=$3,tingkat=$4,alamat=$5,updated_at=DEFAULT
			WHERE id=$6 returning updated_at `,
			siswa.Nama, siswa.IDKelas, siswa.IDWaliKelas, siswa.Tingkat, siswa.Alamat, id).Scan(&updatedAt)

		if err != nil {
			tx.Rollback()

			return nil, apierror.NewError(http.StatusInternalServerError, http.StatusInternalServerError, "Database transaction failed", errors.Wrap(err, "updatesiswa: update data failed"))
		}
	}

	err = tx.Commit()
	if err != nil {
		return nil, apierror.NewError(http.StatusInternalServerError, http.StatusInternalServerError, "Database transaction failed", errors.Wrap(err, "updatesiswa: commit transaction failed"))
	}

	return &schema.SiswaResponse{
		ID:          siswa.ID,
		Nama:        siswa.Nama,
		IDKelas:     siswa.IDKelas,
		IDWaliKelas: siswa.IDWaliKelas,
		Tingkat:     siswa.Tingkat,
		Alamat:      siswa.Alamat,
		CreatedAt:   siswa.CreatedAt,
		UpdatedAt:   &updatedAt,
	}, nil
}

// DeleteSiswa ...
func (s *SiswaService) DeleteSiswa(id string) error {
	if id == "" {
		return apierror.NewError(http.StatusBadRequest, http.StatusBadRequest, "Siswa id is not set", errors.New("deletesiswa: siswa id is not set"))
	}

	tx, err := s.db.Beginx()
	if err != nil {
		return apierror.NewError(http.StatusInternalServerError, http.StatusInternalServerError, "Database transaction failed", errors.Wrap(err, "deletesiswa: begin transaction failed"))
	}

	var rows int64
	{
		result, err := tx.Exec(`
			DELETE FROM public.siswa 
			WHERE id=$1`,
			id)
		rows, _ = result.RowsAffected()

		if err != nil {
			tx.Rollback()
			return apierror.NewError(http.StatusInternalServerError, http.StatusInternalServerError, "Database transaction failed", errors.Wrap(err, "deletesiswa: delete data failed"))
		}
	}

	err = tx.Commit()
	if err != nil {
		return apierror.NewError(http.StatusInternalServerError, http.StatusInternalServerError, "Database transaction failed", errors.Wrap(err, "deletesiswa: commit transaction failed"))
	}

	if rows == 0 {
		return apierror.NewError(http.StatusNotFound, http.StatusNotFound, "Siswa with id: "+id+" is not exists", errors.Wrap(err, "deletesiswa: siswa with id: "+id+" is not exists"))
	}

	return nil
}
