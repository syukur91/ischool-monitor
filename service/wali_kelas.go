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

// Wali_KelasService ...
type Wali_KelasService struct {
	db *sqlx.DB
}

// NewWali_KelasService ...
func NewWali_KelasService(db *sqlx.DB) *Wali_KelasService {
	return &Wali_KelasService{db: db}
}

// CreateWali_Kelas ...
func (s *Wali_KelasService) CreateWali_Kelas(request *schema.CreateWali_KelasRequest) (*schema.Wali_KelasResponse, error) {
	if request.Nama == "" {
		return nil, apierror.NewError(http.StatusBadRequest, http.StatusBadRequest, "Wali_Kelas nama is not set", errors.New("createwali_kelas: wali_kelas nama is not set"))
	}

	if request.Alamat == "" {
		return nil, apierror.NewError(http.StatusBadRequest, http.StatusBadRequest, "Wali_Kelas alamat is not set", errors.New("createwali_kelas: wali_kelas alamat is not set"))
	}

	if request.Telpon == "" {
		return nil, apierror.NewError(http.StatusBadRequest, http.StatusBadRequest, "Wali_Kelas telpon is not set", errors.New("createwali_kelas: wali_kelas telpon is not set"))
	}

	tx, err := s.db.Begin()
	if err != nil {
		return nil, apierror.NewError(http.StatusInternalServerError, http.StatusInternalServerError, "Database transaction failed", errors.Wrap(err, "createwali_kelas: begin transaction failed"))
	}

	id := 0
	var createdAt time.Time

	{
		stmt, err := tx.Prepare(`
			INSERT INTO public.wali_kelas (nama, alamat, telpon)
			VALUES($1, $2, $3)
			RETURNING id, created_at;
		`)

		if err != nil {
			tx.Rollback()
			return nil, apierror.NewError(http.StatusInternalServerError, http.StatusInternalServerError, "Database transaction failed", errors.Wrap(err, "createwali_kelas: prepare insert statement failed"))
		}
		defer stmt.Close()

		err = stmt.QueryRow(request.Nama, request.Alamat, request.Telpon).Scan(&id, &createdAt)
		if err != nil {
			tx.Rollback()

			if strings.Index(err.Error(), "duplicate key value violates unique constraint \"wali_kelas_name_unique\"") > -1 {
				return nil, apierror.NewError(http.StatusBadRequest, http.StatusBadRequest, "Wali_Kelas with same nama already exists. Use different nama", errors.Wrap(err, "createwali_kelas: Wali_Kelas with same nama already exists"))
			}

			return nil, apierror.NewError(http.StatusInternalServerError, http.StatusInternalServerError, "Database transaction failed", errors.Wrap(err, "createwali_kelas: exec insert statement failed"))
		}
	}

	err = tx.Commit()
	if err != nil {
		return nil, apierror.NewError(http.StatusInternalServerError, http.StatusInternalServerError, "Database transaction failed", errors.Wrap(err, "createwali_kelas: commit transaction failed"))
	}

	return &schema.Wali_KelasResponse{
		ID:        id,
		Nama:      request.Nama,
		Alamat:    request.Alamat,
		Telpon:    request.Telpon,
		CreatedAt: &createdAt,
	}, nil
}

// GetWali_Kelas ...
func (s *Wali_KelasService) GetWali_Kelas(id string) (*schema.Wali_KelasResponse, error) {
	if id == "" {
		return nil, apierror.NewError(http.StatusBadRequest, http.StatusBadRequest, "Wali_Kelas id is not set", errors.New("getwali_kelas: wali_kelas id is not set"))
	}

	tx, err := s.db.Beginx()
	if err != nil {
		return nil, apierror.NewError(http.StatusInternalServerError, http.StatusInternalServerError, "Database transaction faileds", errors.Wrap(err, "getwali_kelas: begin transaction failed"))
	}

	wali_kelas := schema.Wali_KelasResponse{}
	{
		err := tx.Get(&wali_kelas, `
			SELECT id,nama,alamat,telpon,created_at,updated_at 
			FROM public.wali_kelas
			WHERE id=$1;`,
			id)

		if err != nil {
			tx.Rollback()

			if strings.Index(err.Error(), "sql: no rows in result set") > -1 {
				return nil, apierror.NewError(http.StatusNotFound, http.StatusNotFound, "Wali_Kelas with id: "+id+" is not exists", errors.Wrap(err, "getwali_kelas: wali_kelas with id: "+id+" is not exists"))
			}

			return nil, apierror.NewError(http.StatusInternalServerError, http.StatusInternalServerError, "Database transaction failedss", errors.Wrap(err, "getwali_kelas: get data failed"))
		}
	}

	err = tx.Commit()
	if err != nil {
		return nil, apierror.NewError(http.StatusInternalServerError, http.StatusInternalServerError, "Database transaction failedsss", errors.Wrap(err, "getwali_kelas: commit transaction failed"))
	}

	return &wali_kelas, nil
}

// ListWali_Kelass ...
func (s *Wali_KelasService) ListWali_Kelass(gridParams *query.GridParams) ([]schema.Wali_KelasResponse, int, error) {

	tx, err := s.db.Beginx()
	if err != nil {
		return nil, 0, apierror.NewError(http.StatusInternalServerError, http.StatusInternalServerError, "Database transaction failed", errors.Wrap(err, "listwali_kelas: begin transaction failed"))
	}

	wali_kelass := []schema.Wali_KelasResponse{}
	total := 0
	{
		dataStatement := "SELECT id,nama,alamat,telpon,created_at,updated_at FROM public.wali_kelas"
		dataQuery, dataParams := query.FullQuery(gridParams, "", nil)
		err := tx.Select(&wali_kelass, dataStatement+dataQuery, dataParams...)
		if err != nil {
			tx.Rollback()
			return nil, 0, apierror.NewError(http.StatusInternalServerError, http.StatusInternalServerError, "Database transaction failed", errors.Wrap(err, "listwali_kelas: get data failed"))
		}

		countStatement := "SELECT count(*) FROM public.wali_kelas"
		countQuery, countParams := query.FilterQuery(gridParams, "", nil)
		err = tx.QueryRow(countStatement+countQuery, countParams...).Scan(&total)
		if err != nil {
			tx.Rollback()
			return nil, 0, apierror.NewError(http.StatusInternalServerError, http.StatusInternalServerError, "Database transaction failed", errors.Wrap(err, "listwali_kelas: get count failed"))
		}
	}

	err = tx.Commit()
	if err != nil {
		return nil, 0, apierror.NewError(http.StatusInternalServerError, http.StatusInternalServerError, "Database transaction failed", errors.Wrap(err, "getwali_kelas: commit transaction failed"))
	}

	return wali_kelass, total, nil
}

// UpdateWali_Kelas ...
func (s *Wali_KelasService) UpdateWali_Kelas(id string, request *schema.UpdateWali_KelasRequest) (*schema.Wali_KelasResponse, error) {
	if id == "" {
		return nil, apierror.NewError(http.StatusBadRequest, http.StatusBadRequest, "Wali_Kelas id is not set", errors.New("updatewali_kelas: wali_kelas id is not set"))
	}

	tx, err := s.db.Beginx()
	if err != nil {
		return nil, apierror.NewError(http.StatusInternalServerError, http.StatusInternalServerError, "Database transaction failed", errors.Wrap(err, "updatewali_kelas: begin transaction failed"))
	}

	// get existing wali_kelas
	wali_kelas := schema.Wali_KelasResponse{}
	{
		err := tx.Get(&wali_kelas, `
			SELECT id,nama,alamat,telpon,created_at,updated_at 
			FROM public.wali_kelas
			WHERE id=$1;`,
			id)

		if err != nil {
			tx.Rollback()

			if strings.Index(err.Error(), "sql: no rows in result set") > -1 {
				return nil, apierror.NewError(http.StatusNotFound, http.StatusNotFound, "Wali_Kelas with id: "+id+" is not exists", errors.Wrap(err, "updatewali_kelas: wali_kelas with id: "+id+" is not exists"))
			}

			return nil, apierror.NewError(http.StatusInternalServerError, http.StatusInternalServerError, "Database transaction failed", errors.Wrap(err, "updatewali_kelas: get data failed"))
		}
	}

	// update wali_kelas
	var updatedAt time.Time
	{
		// only update if not empty
		if request.Nama != "" {
			wali_kelas.Nama = request.Nama
		}

		if request.Alamat != "" {
			wali_kelas.Alamat = request.Alamat
		}

		if request.Telpon != "" {
			wali_kelas.Telpon = request.Telpon
		}

		err := tx.QueryRow(`
			UPDATE public.wali_kelas SET nama=$1,alamat=$2,telpon=$3,updated_at=DEFAULT
			WHERE id=$4 returning updated_at `,
			wali_kelas.Nama, wali_kelas.Alamat, wali_kelas.Telpon, id).Scan(&updatedAt)

		if err != nil {
			tx.Rollback()

			return nil, apierror.NewError(http.StatusInternalServerError, http.StatusInternalServerError, "Database transaction failedss", errors.Wrap(err, "updatewali_kelas: update data failed"))
		}
	}

	err = tx.Commit()
	if err != nil {
		return nil, apierror.NewError(http.StatusInternalServerError, http.StatusInternalServerError, "Database transaction failedsss", errors.Wrap(err, "updatewali_kelas: commit transaction failed"))
	}

	return &schema.Wali_KelasResponse{
		ID:        wali_kelas.ID,
		Nama:      wali_kelas.Nama,
		Alamat:    wali_kelas.Alamat,
		Telpon:    wali_kelas.Telpon,
		CreatedAt: wali_kelas.CreatedAt,
		UpdatedAt: &updatedAt,
	}, nil
}

// DeleteWali_Kelas ...
func (s *Wali_KelasService) DeleteWali_Kelas(id string) error {
	if id == "" {
		return apierror.NewError(http.StatusBadRequest, http.StatusBadRequest, "Wali_Kelas id is not set", errors.New("deletewali_kelas: wali_kelas id is not set"))
	}

	tx, err := s.db.Beginx()
	if err != nil {
		return apierror.NewError(http.StatusInternalServerError, http.StatusInternalServerError, "Database transaction failed", errors.Wrap(err, "deletewali_kelas: begin transaction failed"))
	}

	var rows int64
	{
		result, err := tx.Exec(`
			DELETE FROM public.wali_kelas
			WHERE id=$1`,
			id)
		rows, _ = result.RowsAffected()

		if err != nil {
			tx.Rollback()
			return apierror.NewError(http.StatusInternalServerError, http.StatusInternalServerError, "Database transaction failed", errors.Wrap(err, "deletewali_kelas: delete data failed"))
		}
	}

	err = tx.Commit()
	if err != nil {
		return apierror.NewError(http.StatusInternalServerError, http.StatusInternalServerError, "Database transaction failed", errors.Wrap(err, "deletewali_kelas: commit transaction failed"))
	}

	if rows == 0 {
		return apierror.NewError(http.StatusNotFound, http.StatusNotFound, "Wali_Kelas with id: "+id+" is not exists", errors.Wrap(err, "deletewali_kelas: wali_kelas with id: "+id+" is not exists"))
	}

	return nil
}
