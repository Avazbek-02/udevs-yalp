package repo

import (
	"context"
	"fmt"
	"time"

	"github.com/golanguzb70/system-design-openbudget/config"
	"github.com/golanguzb70/system-design-openbudget/internal/entity"
	"github.com/golanguzb70/system-design-openbudget/pkg/logger"
	"github.com/golanguzb70/system-design-openbudget/pkg/postgres"
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

	qeury, args, err := r.pg.Builder.Insert("auth_user").
		Columns(`id, full_name, phone_number, username, password, user_type, status`).
		Values(req.ID, req.FullName, req.PhoneNumber, req.Username, req.Password, req.UserType, req.Status).ToSql()
	if err != nil {
		return entity.User{}, err
	}

	_, err = r.pg.Pool.Exec(ctx, qeury, args...)
	if err != nil {
		return entity.User{}, err
	}

	return req, nil
}

func (r *UserRepo) GetSingle(ctx context.Context, req entity.UserSingleRequest) (entity.User, error) {
	response := entity.User{}
	var (
		createdAt, updatedAt time.Time
	)

	qeuryBuilder := r.pg.Builder.
		Select(`id, full_name, phone_number, username, password, user_type, status, created_at, updated_at`).
		From("auth_user")

	switch {
	case req.ID != "":
		qeuryBuilder = qeuryBuilder.Where("id = ?", req.ID)
	case req.PhoneNumber != "" && req.UserType != "":
		qeuryBuilder = qeuryBuilder.Where("phone_number = ? AND user_type = ?", req.PhoneNumber, req.UserType)
	case req.UserName != "":
		qeuryBuilder = qeuryBuilder.Where("username = ?", req.UserName)
	default:
		return entity.User{}, fmt.Errorf("GetSingle - invalid request")
	}

	qeury, args, err := qeuryBuilder.ToSql()
	if err != nil {
		return entity.User{}, err
	}

	err = r.pg.Pool.QueryRow(ctx, qeury, args...).
		Scan(&response.ID, &response.FullName, &response.PhoneNumber, &response.Username, &response.Password,
			&response.UserType, &response.Status, &createdAt, &updatedAt)
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

	qeuryBuilder := r.pg.Builder.
		Select(`id, full_name, phone_number, username, password, user_type, status, created_at, updated_at`).
		From("auth_user")

	qeuryBuilder, where := PrepareGetListQuery(qeuryBuilder, req)
	qeury, args, err := qeuryBuilder.ToSql()
	if err != nil {
		return response, err
	}

	rows, err := r.pg.Pool.Query(ctx, qeury, args...)
	if err != nil {
		return response, err
	}
	defer rows.Close()

	for rows.Next() {
		var item entity.User
		err = rows.Scan(&item.ID, &item.FullName, &item.PhoneNumber, &item.Username, &item.Password,
			&item.UserType, &item.Status, &createdAt, &updatedAt)
		if err != nil {
			return response, err
		}

		item.CreatedAt = createdAt.Format(time.RFC3339)
		item.UpdatedAt = updatedAt.Format(time.RFC3339)

		response.Items = append(response.Items, item)
	}

	countQuery, args, err := r.pg.Builder.Select("COUNT(1)").From("auth_user").Where(where).ToSql()
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
		"full_name":    req.FullName,
		"phone_number": req.PhoneNumber,
		"username":     req.Username,
		"status":       req.Status,
	}

	if req.Password != "" {
		mp["password"] = req.Password
	}

	qeury, args, err := r.pg.Builder.Update("auth_user").SetMap(mp).Where("id = ?", req.ID).ToSql()
	if err != nil {
		return entity.User{}, err
	}

	_, err = r.pg.Pool.Exec(ctx, qeury, args...)
	if err != nil {
		return entity.User{}, err
	}

	return req, nil
}

func (r *UserRepo) Delete(ctx context.Context, req entity.Id) error {
	qeury, args, err := r.pg.Builder.Delete("auth_user").Where("id = ?", req.ID).ToSql()
	if err != nil {
		return err
	}

	_, err = r.pg.Pool.Exec(ctx, qeury, args...)
	if err != nil {
		return err
	}

	return nil
}

func (r *UserRepo) UpdateField(ctx context.Context, req entity.UpdateFieldRequest) (entity.RowsEffected, error) {
	mp := map[string]interface{}{}
	response := entity.RowsEffected{}

	for _, item := range req.Items {
		mp[item.Column] = item.Value
	}

	qeury, args, err := r.pg.Builder.Update("auth_user").SetMap(mp).Where(PrepareFilter(req.Filter)).ToSql()
	if err != nil {
		return response, err
	}

	n, err := r.pg.Pool.Exec(ctx, qeury, args...)
	if err != nil {
		return response, err
	}

	response.RowsEffected = int(n.RowsAffected())

	return response, nil
}
