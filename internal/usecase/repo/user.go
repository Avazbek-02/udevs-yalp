package repo

import (
	"context"
	"fmt"
	"time"

	"github.com/Avazbek-02/udevslab-lesson6/config"
	"github.com/Avazbek-02/udevslab-lesson6/internal/entity"
	"github.com/Avazbek-02/udevslab-lesson6/pkg/logger"
	"github.com/Avazbek-02/udevslab-lesson6/pkg/postgres"
	"github.com/google/uuid"
)

type UserRepo struct {
	pg     *postgres.Postgres
	config *config.Config
	logger *logger.Logger
}

// New -.
func NewUserRepo(pg *postgres.Postgres, config *config.Config, logger *logger.Logger) *UserRepo {
	return &UserRepo{
		pg:     pg,
		config: config,
		logger: logger,
	}
}

func (r *UserRepo) Create(ctx context.Context, req entity.User) (entity.User, error) {
	req.ID = uuid.NewString()

	query, args, err := r.pg.Builder.Insert("users").
		Columns(`id, user_type, user_role, email, password, username, full_name, gender, profile_picture, bio, status`).
		Values(req.ID, req.UserType, req.UserRole, req.Email, req.Password, req.Username, req.FullName, req.Gender,
			req.Profile_picture, req.Bio, req.Status).ToSql()
	if err != nil {
		return entity.User{}, err
	}

	_, err = r.pg.Pool.Exec(ctx, query, args...)
	if err != nil {
		return entity.User{}, err
	}

	return req, nil
}

func (r *UserRepo) GetSingle(ctx context.Context, req entity.UserSingleRequest) (entity.User, error) {
	response := entity.User{}
	var createdAt, updatedAt time.Time

	queryBuilder := r.pg.Builder.
		Select(`id, user_type, user_role, email, password, username, full_name, gender, profile_picture, bio, status, created_at, updated_at`).
		From("users")

	switch {
	case req.ID != "":
		queryBuilder = queryBuilder.Where("id = ?", req.ID)
	case req.Email != "":
		queryBuilder = queryBuilder.Where("email = ?", req.Email)
	default:
		return entity.User{}, fmt.Errorf("GetSingle - invalid request")
	}

	query, args, err := queryBuilder.ToSql()
	if err != nil {
		return entity.User{}, err
	}

	err = r.pg.Pool.QueryRow(ctx, query, args...).
		Scan(&response.ID, &response.UserType, &response.UserRole, &response.Email, &response.Password, &response.Username,
			&response.FullName, &response.Gender, &response.Profile_picture, &response.Bio, &response.Status, &createdAt, &updatedAt)
	if err != nil {
		return entity.User{}, err
	}

	response.CreatedAt = createdAt.Format(time.RFC3339)
	response.UpdatedAt = updatedAt.Format(time.RFC3339)

	return response, nil
}

func (r *UserRepo) GetList(ctx context.Context, req entity.GetListFilter) (entity.UserList, error) {
	var (
		response             = entity.UserList{}
		createdAt, updatedAt time.Time
	)

	queryBuilder := r.pg.Builder.
		Select(`id, user_type, user_role, email, password, username, full_name, gender, profile_picture, bio, status, created_at, updated_at`).
		From("users")

	// Add filters if necessary, for example, by status or user_type=
	query, args, err := queryBuilder.ToSql()
	if err != nil {
		return response, err
	}

	rows, err := r.pg.Pool.Query(ctx, query, args...)
	if err != nil {
		return response, err
	}
	defer rows.Close()

	for rows.Next() {
		var item entity.User
		err = rows.Scan(&item.ID, &item.UserType, &item.UserRole, &item.Email, &item.Password, &item.Username,
			&item.FullName, &item.Gender, &item.Profile_picture, &item.Bio, &item.Status, &createdAt, &updatedAt)
		if err != nil {
			return response, err
		}

		item.CreatedAt = createdAt.Format(time.RFC3339)
		item.UpdatedAt = updatedAt.Format(time.RFC3339)

		response.Items = append(response.Items, item)
	}

	countQuery, args, err := r.pg.Builder.Select("COUNT(1)").From("users").ToSql()
	if err != nil {
		return response, err
	}

	err = r.pg.Pool.QueryRow(ctx, countQuery, args...).Scan(&response.Count)
	if err != nil {
		return response, err
	}

	return response, nil
}

func (r *UserRepo) Update(ctx context.Context, req entity.User) (entity.User, error) {
	mp := map[string]interface{}{
		"user_type":       req.UserType,
		"user_role":       req.UserRole,
		"email":           req.Email,
		"password":        req.Password,
		"username":        req.Username,
		"full_name":       req.FullName,
		"gender":          req.Gender,
		"profile_picture": req.Profile_picture,
		"bio":             req.Bio,
		"status":          req.Status,
	}

	query, args, err := r.pg.Builder.Update("users").SetMap(mp).Where("id = ?", req.ID).ToSql()
	if err != nil {
		return entity.User{}, err
	}

	_, err = r.pg.Pool.Exec(ctx, query, args...)
	if err != nil {
		return entity.User{}, err
	}

	return req, nil
}

func (r *UserRepo) Delete(ctx context.Context, req entity.Id) error {
	query, args, err := r.pg.Builder.Delete("users").Where("id = ?", req.ID).ToSql()
	if err != nil {
		return err
	}

	_, err = r.pg.Pool.Exec(ctx, query, args...)
	if err != nil {
		return err
	}

	return nil
}
