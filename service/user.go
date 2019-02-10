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

// UserService ...
type UserService struct {
	db *sqlx.DB
}

// NewUserService ...
func NewUserService(db *sqlx.DB) *UserService {
	return &UserService{db: db}
}

// CreateUser ...
func (s *UserService) CreateUser(request *schema.CreateUserRequest) (*schema.UserResponse, error) {
	if request.Nama == "" {
		return nil, apierror.NewError(http.StatusBadRequest, http.StatusBadRequest, "User nama is not set", errors.New("createuser: user nama is not set"))
	}

	if request.Alamat == "" {
		return nil, apierror.NewError(http.StatusBadRequest, http.StatusBadRequest, "User alamat is not set", errors.New("createuser: user alamat is not set"))
	}

	if request.Password == "" {
		return nil, apierror.NewError(http.StatusBadRequest, http.StatusBadRequest, "User password is not set", errors.New("createuser: user password is not set"))
	}

	if request.Telepon == "" {
		return nil, apierror.NewError(http.StatusBadRequest, http.StatusBadRequest, "User telepon is not set", errors.New("createuser: user telepon is not set"))
	}

	tx, err := s.db.Begin()
	if err != nil {
		return nil, apierror.NewError(http.StatusInternalServerError, http.StatusInternalServerError, "Database transaction failed", errors.Wrap(err, "createuser: begin transaction failed"))
	}

	id := 0
	var createdAt time.Time

	{
		stmt, err := tx.Prepare(`
			INSERT INTO public.user (nama, alamat, password, telepon)
			VALUES($1, $2, $3, $4)
			RETURNING id, created_at;
		`)

		if err != nil {
			tx.Rollback()
			return nil, apierror.NewError(http.StatusInternalServerError, http.StatusInternalServerError, "Database transaction failed", errors.Wrap(err, "createuser: prepare insert statement failed"))
		}
		defer stmt.Close()

		err = stmt.QueryRow(request.Nama, request.Alamat, request.Password, request.Telepon).Scan(&id, &createdAt)
		if err != nil {
			tx.Rollback()

			if strings.Index(err.Error(), "duplicate key value violates unique constraint \"user_name_unique\"") > -1 {
				return nil, apierror.NewError(http.StatusBadRequest, http.StatusBadRequest, "User with same nama already exists. Use different nama", errors.Wrap(err, "createuser: User with same nama already exists"))
			}

			return nil, apierror.NewError(http.StatusInternalServerError, http.StatusInternalServerError, "Database transaction failed", errors.Wrap(err, "createuser: exec insert statement failed"))
		}
	}

	err = tx.Commit()
	if err != nil {
		return nil, apierror.NewError(http.StatusInternalServerError, http.StatusInternalServerError, "Database transaction failed", errors.Wrap(err, "createuser: commit transaction failed"))
	}

	return &schema.UserResponse{
		ID:        id,
		Nama:      request.Nama,
		Alamat:    request.Alamat,
		Password:  request.Password,
		Telepon:   request.Telepon,
		CreatedAt: &createdAt,
	}, nil
}

// GetUser ...
func (s *UserService) GetUser(id string) (*schema.UserResponse, error) {
	if id == "" {
		return nil, apierror.NewError(http.StatusBadRequest, http.StatusBadRequest, "User id is not set", errors.New("getuser: user id is not set"))
	}

	tx, err := s.db.Beginx()
	if err != nil {
		return nil, apierror.NewError(http.StatusInternalServerError, http.StatusInternalServerError, "Database transaction failed", errors.Wrap(err, "getuser: begin transaction failed"))
	}

	user := schema.UserResponse{}
	{
		err := tx.Get(&user, `
			SELECT id,nama,alamat,password,telepon,created_at,updated_at 
			FROM public.users 
			WHERE id=$1;`,
			id)

		if err != nil {
			tx.Rollback()

			if strings.Index(err.Error(), "sql: no rows in result set") > -1 {
				return nil, apierror.NewError(http.StatusNotFound, http.StatusNotFound, "User with id: "+id+" is not exists", errors.Wrap(err, "getuser: user with id: "+id+" is not exists"))
			}

			return nil, apierror.NewError(http.StatusInternalServerError, http.StatusInternalServerError, "Database transaction failed", errors.Wrap(err, "getuser: get data failed"))
		}
	}

	err = tx.Commit()
	if err != nil {
		return nil, apierror.NewError(http.StatusInternalServerError, http.StatusInternalServerError, "Database transaction failed", errors.Wrap(err, "getuser: commit transaction failed"))
	}

	return &user, nil
}

// ListUsers ...
func (s *UserService) ListUsers(gridParams *query.GridParams) ([]schema.UserResponse, int, error) {

	tx, err := s.db.Beginx()
	if err != nil {
		return nil, 0, apierror.NewError(http.StatusInternalServerError, http.StatusInternalServerError, "Database transaction failed", errors.Wrap(err, "listuser: begin transaction failed"))
	}

	users := []schema.UserResponse{}
	total := 0
	{
		dataStatement := "SELECT id,nama,alamat,password,telepon,created_at,updated_at FROM public.users"
		dataQuery, dataParams := query.FullQuery(gridParams, "", nil)
		err := tx.Select(&users, dataStatement+dataQuery, dataParams...)
		if err != nil {
			tx.Rollback()
			return nil, 0, apierror.NewError(http.StatusInternalServerError, http.StatusInternalServerError, "Database transaction failed", errors.Wrap(err, "listuser: get data failed"))
		}

		countStatement := "SELECT count(*) FROM public.users"
		countQuery, countParams := query.FilterQuery(gridParams, "", nil)
		err = tx.QueryRow(countStatement+countQuery, countParams...).Scan(&total)
		if err != nil {
			tx.Rollback()
			return nil, 0, apierror.NewError(http.StatusInternalServerError, http.StatusInternalServerError, "Database transaction failed", errors.Wrap(err, "listuser: get count failed"))
		}
	}

	err = tx.Commit()
	if err != nil {
		return nil, 0, apierror.NewError(http.StatusInternalServerError, http.StatusInternalServerError, "Database transaction failed", errors.Wrap(err, "getuser: commit transaction failed"))
	}

	return users, total, nil
}

// UpdateUser ...
func (s *UserService) UpdateUser(id string, request *schema.UpdateUserRequest) (*schema.UserResponse, error) {
	if id == "" {
		return nil, apierror.NewError(http.StatusBadRequest, http.StatusBadRequest, "User id is not set", errors.New("updateuser: user id is not set"))
	}

	tx, err := s.db.Beginx()
	if err != nil {
		return nil, apierror.NewError(http.StatusInternalServerError, http.StatusInternalServerError, "Database transaction failed", errors.Wrap(err, "updateuser: begin transaction failed"))
	}

	// get existing user
	user := schema.UserResponse{}
	{
		err := tx.Get(&user, `
			SELECT id,nama,alamat,password,telepon,created_at,updated_at 
			FROM public.users 
			WHERE id=$1;`,
			id)

		if err != nil {
			tx.Rollback()

			if strings.Index(err.Error(), "sql: no rows in result set") > -1 {
				return nil, apierror.NewError(http.StatusNotFound, http.StatusNotFound, "User with id: "+id+" is not exists", errors.Wrap(err, "updateuser: user with id: "+id+" is not exists"))
			}

			return nil, apierror.NewError(http.StatusInternalServerError, http.StatusInternalServerError, "Database transaction failed", errors.Wrap(err, "updateuser: get data failed"))
		}
	}

	// update user
	var updatedAt time.Time
	{
		// only update if not empty
		if request.Nama != "" {
			user.Nama = request.Nama
		}

		if request.Alamat != "" {
			user.Alamat = request.Alamat
		}

		if request.Password != "" {
			user.Password = request.Password
		}

		if request.Telepon != "" {
			user.Telepon = request.Telepon
		}

		err := tx.QueryRow(`
			UPDATE public.users SET nama=$1,alamat=$2,password=$3,telepon=$4,updated_at=DEFAULT
			WHERE id=$5 returning updated_at `,
			user.Nama, user.Alamat, user.Password, user.Telepon, id).Scan(&updatedAt)

		if err != nil {
			tx.Rollback()

			return nil, apierror.NewError(http.StatusInternalServerError, http.StatusInternalServerError, "Database transaction failed", errors.Wrap(err, "updateuser: update data failed"))
		}
	}

	err = tx.Commit()
	if err != nil {
		return nil, apierror.NewError(http.StatusInternalServerError, http.StatusInternalServerError, "Database transaction failed", errors.Wrap(err, "updateuser: commit transaction failed"))
	}

	return &schema.UserResponse{
		ID:        user.ID,
		Nama:      user.Nama,
		Alamat:    user.Alamat,
		Password:  user.Password,
		Telepon:   user.Telepon,
		CreatedAt: user.CreatedAt,
		UpdatedAt: &updatedAt,
	}, nil
}

// DeleteUser ...
func (s *UserService) DeleteUser(id string) error {
	if id == "" {
		return apierror.NewError(http.StatusBadRequest, http.StatusBadRequest, "User id is not set", errors.New("deleteuser: user id is not set"))
	}

	tx, err := s.db.Beginx()
	if err != nil {
		return apierror.NewError(http.StatusInternalServerError, http.StatusInternalServerError, "Database transaction failed", errors.Wrap(err, "deleteuser: begin transaction failed"))
	}

	var rows int64
	{
		result, err := tx.Exec(`
			DELETE FROM public.users 
			WHERE id=$1`,
			id)
		rows, _ = result.RowsAffected()

		if err != nil {
			tx.Rollback()
			return apierror.NewError(http.StatusInternalServerError, http.StatusInternalServerError, "Database transaction failed", errors.Wrap(err, "deleteuser: delete data failed"))
		}
	}

	err = tx.Commit()
	if err != nil {
		return apierror.NewError(http.StatusInternalServerError, http.StatusInternalServerError, "Database transaction failed", errors.Wrap(err, "deleteuser: commit transaction failed"))
	}

	if rows == 0 {
		return apierror.NewError(http.StatusNotFound, http.StatusNotFound, "User with id: "+id+" is not exists", errors.Wrap(err, "deleteuser: user with id: "+id+" is not exists"))
	}

	return nil
}
