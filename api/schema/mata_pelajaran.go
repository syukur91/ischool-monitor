package service

import (
	"net/http"
	"strings"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"

	"gitlab.com/nextid/tenant-api/api/schema"
	"gitlab.com/nextid/tenant-api/pkg/apierror"
	"gitlab.com/nextid/tenant-api/pkg/query"
)

// Mata_PelajaranService ...
type Mata_PelajaranService struct {
	db *sqlx.DB
}

// NewMata_PelajaranService ...
func NewMata_PelajaranService(db *sqlx.DB) *Mata_PelajaranService {
	return &Mata_PelajaranService{db: db}
}

// CreateMata_Pelajaran ...
func (s *Mata_PelajaranService) CreateMata_Pelajaran(request *schema.CreateMata_PelajaranRequest) (*schema.Mata_PelajaranResponse, error) {
	if request.Name == "" {
		return nil, apierror.NewError(http.StatusBadRequest, http.StatusBadRequest, "Mata_Pelajaran name is not set", errors.New("createmata_pelajaran: mata_pelajaran name is not set"))
	}

	if request.Description == "" {
		return nil, apierror.NewError(http.StatusBadRequest, http.StatusBadRequest, "Mata_Pelajaran description is not set", errors.New("createmata_pelajaran: mata_pelajaran description is not set"))
	}

	tx, err := s.db.Begin()
	if err != nil {
		return nil, apierror.NewError(http.StatusInternalServerError, http.StatusInternalServerError, "Database transaction failed", errors.Wrap(err, "createmata_pelajaran: begin transaction failed"))
	}

	id := 0
	var createdAt time.Time

	{
		stmt, err := tx.Prepare(`
			INSERT INTO public.mata_pelajarans (name, description)
			VALUES($1, $2)
			RETURNING id, created_at;
		`)

		if err != nil {
			tx.Rollback()
			return nil, apierror.NewError(http.StatusInternalServerError, http.StatusInternalServerError, "Database transaction failed", errors.Wrap(err, "createmata_pelajaran: prepare insert statement failed"))
		}
		defer stmt.Close()

		err = stmt.QueryRow(request.Name, request.Description).Scan(&id, &createdAt)
		if err != nil {
			tx.Rollback()

			if strings.Index(err.Error(), "duplicate key value violates unique constraint \"mata_pelajaran_name_unique\"") > -1 {
				return nil, apierror.NewError(http.StatusBadRequest, http.StatusBadRequest, "Mata_Pelajaran with same name already exists. Use different name", errors.Wrap(err, "createmata_pelajaran: Mata_Pelajaran with same name already exists"))
			}

			return nil, apierror.NewError(http.StatusInternalServerError, http.StatusInternalServerError, "Database transaction failed", errors.Wrap(err, "createmata_pelajaran: exec insert statement failed"))
		}
	}

	err = tx.Commit()
	if err != nil {
		return nil, apierror.NewError(http.StatusInternalServerError, http.StatusInternalServerError, "Database transaction failed", errors.Wrap(err, "createmata_pelajaran: commit transaction failed"))
	}

	return &schema.Mata_PelajaranResponse{
		ID:          id,
		Name:        request.Name,
		Description: request.Description,
		CreatedAt:   &createdAt,
	}, nil
}

// GetMata_Pelajaran ...
func (s *Mata_PelajaranService) GetMata_Pelajaran(id string) (*schema.Mata_PelajaranResponse, error) {
	if id == "" {
		return nil, apierror.NewError(http.StatusBadRequest, http.StatusBadRequest, "Mata_Pelajaran id is not set", errors.New("getmata_pelajaran: mata_pelajaran id is not set"))
	}

	tx, err := s.db.Beginx()
	if err != nil {
		return nil, apierror.NewError(http.StatusInternalServerError, http.StatusInternalServerError, "Database transaction failed", errors.Wrap(err, "getmata_pelajaran: begin transaction failed"))
	}

	mata_pelajaran := schema.Mata_PelajaranResponse{}
	{
		err := tx.Get(&mata_pelajaran, `
			SELECT id,name,description,created_at,updated_at 
			FROM public.mata_pelajarans 
			WHERE id=$1;`,
			id)

		if err != nil {
			tx.Rollback()

			if strings.Index(err.Error(), "sql: no rows in result set") > -1 {
				return nil, apierror.NewError(http.StatusNotFound, http.StatusNotFound, "Mata_Pelajaran with id: "+id+" is not exists", errors.Wrap(err, "getmata_pelajaran: mata_pelajaran with id: "+id+" is not exists"))
			}

			return nil, apierror.NewError(http.StatusInternalServerError, http.StatusInternalServerError, "Database transaction failed", errors.Wrap(err, "getmata_pelajaran: get data failed"))
		}
	}

	err = tx.Commit()
	if err != nil {
		return nil, apierror.NewError(http.StatusInternalServerError, http.StatusInternalServerError, "Database transaction failed", errors.Wrap(err, "getmata_pelajaran: commit transaction failed"))
	}

	return &mata_pelajaran, nil
}

// ListMata_Pelajarans ...
func (s *Mata_PelajaranService) ListMata_Pelajarans(gridParams *query.GridParams) ([]schema.Mata_PelajaranResponse, int, error) {

	tx, err := s.db.Beginx()
	if err != nil {
		return nil, 0, apierror.NewError(http.StatusInternalServerError, http.StatusInternalServerError, "Database transaction failed", errors.Wrap(err, "listmata_pelajaran: begin transaction failed"))
	}

	mata_pelajarans := []schema.Mata_PelajaranResponse{}
	total := 0
	{
		dataStatement := "SELECT id,name,description,created_at,updated_at FROM public.mata_pelajarans"
		dataQuery, dataParams := query.FullQuery(gridParams, "", nil)
		err := tx.Select(&mata_pelajarans, dataStatement+dataQuery, dataParams...)
		if err != nil {
			tx.Rollback()
			return nil, 0, apierror.NewError(http.StatusInternalServerError, http.StatusInternalServerError, "Database transaction failed", errors.Wrap(err, "listmata_pelajaran: get data failed"))
		}

		countStatement := "SELECT count(*) FROM public.mata_pelajarans"
		countQuery, countParams := query.FilterQuery(gridParams, "", nil)
		err = tx.QueryRow(countStatement+countQuery, countParams...).Scan(&total)
		if err != nil {
			tx.Rollback()
			return nil, 0, apierror.NewError(http.StatusInternalServerError, http.StatusInternalServerError, "Database transaction failed", errors.Wrap(err, "listmata_pelajaran: get count failed"))
		}
	}

	err = tx.Commit()
	if err != nil {
		return nil, 0, apierror.NewError(http.StatusInternalServerError, http.StatusInternalServerError, "Database transaction failed", errors.Wrap(err, "getmata_pelajaran: commit transaction failed"))
	}

	return mata_pelajarans, total, nil
}

// UpdateMata_Pelajaran ...
func (s *Mata_PelajaranService) UpdateMata_Pelajaran(id string, request *schema.UpdateMata_PelajaranRequest) (*schema.Mata_PelajaranResponse, error) {
	if id == "" {
		return nil, apierror.NewError(http.StatusBadRequest, http.StatusBadRequest, "Mata_Pelajaran id is not set", errors.New("updatemata_pelajaran: mata_pelajaran id is not set"))
	}

	tx, err := s.db.Beginx()
	if err != nil {
		return nil, apierror.NewError(http.StatusInternalServerError, http.StatusInternalServerError, "Database transaction failed", errors.Wrap(err, "updatemata_pelajaran: begin transaction failed"))
	}

	// get existing mata_pelajaran
	mata_pelajaran := schema.Mata_PelajaranResponse{}
	{
		err := tx.Get(&mata_pelajaran, `
			SELECT id,name,description,created_at,updated_at 
			FROM public.mata_pelajarans 
			WHERE id=$1;`,
			id)

		if err != nil {
			tx.Rollback()

			if strings.Index(err.Error(), "sql: no rows in result set") > -1 {
				return nil, apierror.NewError(http.StatusNotFound, http.StatusNotFound, "Mata_Pelajaran with id: "+id+" is not exists", errors.Wrap(err, "updatemata_pelajaran: mata_pelajaran with id: "+id+" is not exists"))
			}

			return nil, apierror.NewError(http.StatusInternalServerError, http.StatusInternalServerError, "Database transaction failed", errors.Wrap(err, "updatemata_pelajaran: get data failed"))
		}
	}

	// update mata_pelajaran
	var updatedAt time.Time
	{
		// only update if not empty
		if request.Name != "" {
			mata_pelajaran.Name = request.Name
		}

		if request.Description != "" {
			mata_pelajaran.Description = request.Description
		}

		err := tx.QueryRow(`
			UPDATE public.mata_pelajarans SET name=$1,description=$2,updated_at=DEFAULT
			WHERE id=$3 returning updated_at `,
			mata_pelajaran.Name, mata_pelajaran.Description, id).Scan(&updatedAt)

		if err != nil {
			tx.Rollback()

			return nil, apierror.NewError(http.StatusInternalServerError, http.StatusInternalServerError, "Database transaction failed", errors.Wrap(err, "updatemata_pelajaran: update data failed"))
		}
	}

	err = tx.Commit()
	if err != nil {
		return nil, apierror.NewError(http.StatusInternalServerError, http.StatusInternalServerError, "Database transaction failed", errors.Wrap(err, "updatemata_pelajaran: commit transaction failed"))
	}

	return &schema.Mata_PelajaranResponse{
		ID:          mata_pelajaran.ID,
		Name:        mata_pelajaran.Name,
		Description: mata_pelajaran.Description,
		CreatedAt:   mata_pelajaran.CreatedAt,
		UpdatedAt:   &updatedAt,
	}, nil
}

// DeleteMata_Pelajaran ...
func (s *Mata_PelajaranService) DeleteMata_Pelajaran(id string) error {
	if id == "" {
		return apierror.NewError(http.StatusBadRequest, http.StatusBadRequest, "Mata_Pelajaran id is not set", errors.New("deletemata_pelajaran: mata_pelajaran id is not set"))
	}

	tx, err := s.db.Beginx()
	if err != nil {
		return apierror.NewError(http.StatusInternalServerError, http.StatusInternalServerError, "Database transaction failed", errors.Wrap(err, "deletemata_pelajaran: begin transaction failed"))
	}

	var rows int64
	{
		result, err := tx.Exec(`
			DELETE FROM public.mata_pelajarans 
			WHERE id=$1`,
			id)
		rows, _ = result.RowsAffected()

		if err != nil {
			tx.Rollback()
			return apierror.NewError(http.StatusInternalServerError, http.StatusInternalServerError, "Database transaction failed", errors.Wrap(err, "deletemata_pelajaran: delete data failed"))
		}
	}

	err = tx.Commit()
	if err != nil {
		return apierror.NewError(http.StatusInternalServerError, http.StatusInternalServerError, "Database transaction failed", errors.Wrap(err, "deletemata_pelajaran: commit transaction failed"))
	}

	if rows == 0 {
		return apierror.NewError(http.StatusNotFound, http.StatusNotFound, "Mata_Pelajaran with id: "+id+" is not exists", errors.Wrap(err, "deletemata_pelajaran: mata_pelajaran with id: "+id+" is not exists"))
	}

	return nil
}
