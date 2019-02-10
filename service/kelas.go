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

// KelasService ...
type KelasService struct {
	db *sqlx.DB
}

// NewKelasService ...
func NewKelasService(db *sqlx.DB) *KelasService {
	return &KelasService{db: db}
}

// CreateKelas ...
func (s *KelasService) CreateKelas(request *schema.CreateKelasRequest) (*schema.KelasResponse, error) {
	if request.Nama == "" {
		return nil, apierror.NewError(http.StatusBadRequest, http.StatusBadRequest, "Kelas nama is not set", errors.New("createkelas: kelas nama is not set"))
	}

	if request.Tingkat == 0 {
		return nil, apierror.NewError(http.StatusBadRequest, http.StatusBadRequest, "Kelas tingkat is not set", errors.New("createkelas: kelas tingkat is not set"))
	}

	tx, err := s.db.Begin()
	if err != nil {
		return nil, apierror.NewError(http.StatusInternalServerError, http.StatusInternalServerError, "Database transaction failed", errors.Wrap(err, "createkelas: begin transaction failed"))
	}

	id := 0
	var createdAt time.Time

	{
		stmt, err := tx.Prepare(`
			INSERT INTO public.kelas (nama, tingkat)
			VALUES($1, $2)
			RETURNING id, created_at;
		`)

		if err != nil {
			tx.Rollback()
			return nil, apierror.NewError(http.StatusInternalServerError, http.StatusInternalServerError, "Database transaction failed", errors.Wrap(err, "createkelas: prepare insert statement failed"))
		}
		defer stmt.Close()

		err = stmt.QueryRow(request.Nama, request.Tingkat).Scan(&id, &createdAt)
		if err != nil {
			tx.Rollback()

			if strings.Index(err.Error(), "duplicate key value violates unique constraint \"kelas_name_unique\"") > -1 {
				return nil, apierror.NewError(http.StatusBadRequest, http.StatusBadRequest, "Kelas with same nama already exists. Use different nama", errors.Wrap(err, "createkelas: Kelas with same nama already exists"))
			}

			return nil, apierror.NewError(http.StatusInternalServerError, http.StatusInternalServerError, "Database transaction failed", errors.Wrap(err, "createkelas: exec insert statement failed"))
		}
	}

	err = tx.Commit()
	if err != nil {
		return nil, apierror.NewError(http.StatusInternalServerError, http.StatusInternalServerError, "Database transaction failed", errors.Wrap(err, "createkelas: commit transaction failed"))
	}

	return &schema.KelasResponse{
		ID:        id,
		Nama:      request.Nama,
		Tingkat:   request.Tingkat,
		CreatedAt: &createdAt,
	}, nil
}

// GetKelas ...
func (s *KelasService) GetKelas(id string) (*schema.KelasResponse, error) {
	if id == "" {
		return nil, apierror.NewError(http.StatusBadRequest, http.StatusBadRequest, "Kelas id is not set", errors.New("getkelas: kelas id is not set"))
	}

	tx, err := s.db.Beginx()
	if err != nil {
		return nil, apierror.NewError(http.StatusInternalServerError, http.StatusInternalServerError, "Database transaction failed", errors.Wrap(err, "getkelas: begin transaction failed"))
	}

	kelas := schema.KelasResponse{}
	{
		err := tx.Get(&kelas, `
			SELECT id,nama,tingkat,created_at,updated_at 
			FROM public.kelas 
			WHERE id=$1;`,
			id)

		if err != nil {
			tx.Rollback()

			if strings.Index(err.Error(), "sql: no rows in result set") > -1 {
				return nil, apierror.NewError(http.StatusNotFound, http.StatusNotFound, "Kelas with id: "+id+" is not exists", errors.Wrap(err, "getkelas: kelas with id: "+id+" is not exists"))
			}

			return nil, apierror.NewError(http.StatusInternalServerError, http.StatusInternalServerError, "Database transaction failed", errors.Wrap(err, "getkelas: get data failed"))
		}
	}

	err = tx.Commit()
	if err != nil {
		return nil, apierror.NewError(http.StatusInternalServerError, http.StatusInternalServerError, "Database transaction failed", errors.Wrap(err, "getkelas: commit transaction failed"))
	}

	return &kelas, nil
}

// ListKelass ...
func (s *KelasService) ListKelass(gridParams *query.GridParams) ([]schema.KelasResponse, int, error) {

	tx, err := s.db.Beginx()
	if err != nil {
		return nil, 0, apierror.NewError(http.StatusInternalServerError, http.StatusInternalServerError, "Database transaction failed", errors.Wrap(err, "listkelas: begin transaction failed"))
	}

	kelass := []schema.KelasResponse{}
	total := 0
	{
		dataStatement := "SELECT id,nama,tingkat,created_at,updated_at FROM public.kelas"
		dataQuery, dataParams := query.FullQuery(gridParams, "", nil)
		err := tx.Select(&kelass, dataStatement+dataQuery, dataParams...)
		if err != nil {
			tx.Rollback()
			return nil, 0, apierror.NewError(http.StatusInternalServerError, http.StatusInternalServerError, "Database transaction faileds", errors.Wrap(err, "listkelas: get data failed"))
		}

		countStatement := "SELECT count(*) FROM public.kelas"
		countQuery, countParams := query.FilterQuery(gridParams, "", nil)
		err = tx.QueryRow(countStatement+countQuery, countParams...).Scan(&total)
		if err != nil {
			tx.Rollback()
			return nil, 0, apierror.NewError(http.StatusInternalServerError, http.StatusInternalServerError, "Database transaction failedss", errors.Wrap(err, "listkelas: get count failed"))
		}
	}

	err = tx.Commit()
	if err != nil {
		return nil, 0, apierror.NewError(http.StatusInternalServerError, http.StatusInternalServerError, "Database transaction failedsss", errors.Wrap(err, "getkelas: commit transaction failed"))
	}

	return kelass, total, nil
}

// UpdateKelas ...
func (s *KelasService) UpdateKelas(id string, request *schema.UpdateKelasRequest) (*schema.KelasResponse, error) {
	if id == "" {
		return nil, apierror.NewError(http.StatusBadRequest, http.StatusBadRequest, "Kelas id is not set", errors.New("updatekelas: kelas id is not set"))
	}

	tx, err := s.db.Beginx()
	if err != nil {
		return nil, apierror.NewError(http.StatusInternalServerError, http.StatusInternalServerError, "Database transaction failed", errors.Wrap(err, "updatekelas: begin transaction failed"))
	}

	// get existing kelas
	kelas := schema.KelasResponse{}
	{
		err := tx.Get(&kelas, `
			SELECT id,nama,tingkat,created_at,updated_at 
			FROM public.kelas 
			WHERE id=$1;`,
			id)

		if err != nil {
			tx.Rollback()

			if strings.Index(err.Error(), "sql: no rows in result set") > -1 {
				return nil, apierror.NewError(http.StatusNotFound, http.StatusNotFound, "Kelas with id: "+id+" is not exists", errors.Wrap(err, "updatekelas: kelas with id: "+id+" is not exists"))
			}

			return nil, apierror.NewError(http.StatusInternalServerError, http.StatusInternalServerError, "Database transaction failed", errors.Wrap(err, "updatekelas: get data failed"))
		}
	}

	// update kelas
	var updatedAt time.Time
	{
		// only update if not empty
		if request.Nama != "" {
			kelas.Nama = request.Nama
		}

		if request.Tingkat != 0 {
			kelas.Tingkat = request.Tingkat
		}

		err := tx.QueryRow(`
			UPDATE public.kelas SET nama=$1,tingkat=$2,updated_at=DEFAULT
			WHERE id=$3 returning updated_at `,
			kelas.Nama, kelas.Tingkat, id).Scan(&updatedAt)

		if err != nil {
			tx.Rollback()

			return nil, apierror.NewError(http.StatusInternalServerError, http.StatusInternalServerError, "Database transaction failed", errors.Wrap(err, "updatekelas: update data failed"))
		}
	}

	err = tx.Commit()
	if err != nil {
		return nil, apierror.NewError(http.StatusInternalServerError, http.StatusInternalServerError, "Database transaction failed", errors.Wrap(err, "updatekelas: commit transaction failed"))
	}

	return &schema.KelasResponse{
		ID:        kelas.ID,
		Nama:      kelas.Nama,
		Tingkat:   kelas.Tingkat,
		CreatedAt: kelas.CreatedAt,
		UpdatedAt: &updatedAt,
	}, nil
}

// DeleteKelas ...
func (s *KelasService) DeleteKelas(id string) error {
	if id == "" {
		return apierror.NewError(http.StatusBadRequest, http.StatusBadRequest, "Kelas id is not set", errors.New("deletekelas: kelas id is not set"))
	}

	tx, err := s.db.Beginx()
	if err != nil {
		return apierror.NewError(http.StatusInternalServerError, http.StatusInternalServerError, "Database transaction failed", errors.Wrap(err, "deletekelas: begin transaction failed"))
	}

	var rows int64
	{
		result, err := tx.Exec(`
			DELETE FROM public.kelas 
			WHERE id=$1`,
			id)
		rows, _ = result.RowsAffected()

		if err != nil {
			tx.Rollback()
			return apierror.NewError(http.StatusInternalServerError, http.StatusInternalServerError, "Database transaction failed", errors.Wrap(err, "deletekelas: delete data failed"))
		}
	}

	err = tx.Commit()
	if err != nil {
		return apierror.NewError(http.StatusInternalServerError, http.StatusInternalServerError, "Database transaction failed", errors.Wrap(err, "deletekelas: commit transaction failed"))
	}

	if rows == 0 {
		return apierror.NewError(http.StatusNotFound, http.StatusNotFound, "Kelas with id: "+id+" is not exists", errors.Wrap(err, "deletekelas: kelas with id: "+id+" is not exists"))
	}

	return nil
}
