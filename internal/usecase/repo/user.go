package repo

import (
	"context"
	"database/sql"
	"errors"
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

// New 
func NewUserRepo(pg *postgres.Postgres, config *config.Config, logger *logger.Logger) *UserRepo {
	return &UserRepo{
		pg:     pg,
		config: config,
		logger: logger,
	}
}

func (r *UserRepo) Create(ctx context.Context, req entity.User) (entity.User, error) {
	req.ID = uuid.NewString()

	qeury, args, err := r.pg.Builder.Insert("users").
		Columns(`id, full_name, email, bio, username, password_hash, user_type, user_role, status, avatar_id, gender`).
		Values(req.ID, req.FullName, req.Email, req.Bio, req.Username, req.Password, req.UserType, req.UserRole, req.Status, req.AvatarId, req.Gender).ToSql()
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
		avatarID             sql.NullString
	)

	qeuryBuilder := r.pg.Builder.
		Select(`id, full_name, email, bio, username, password_hash, user_type, user_role, status, avatar_id, gender, created_at, updated_at`).
		From("users")

	switch {
	case req.ID != "":
		qeuryBuilder = qeuryBuilder.Where("id = ?", req.ID)
	case req.Email != "":
		qeuryBuilder = qeuryBuilder.Where("email = ?", req.Email)
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
		Scan(&response.ID, &response.FullName, &response.Email, &response.Bio, &response.Username, &response.Password,
			&response.UserType, &response.UserRole, &response.Status, &avatarID, &response.Gender, &createdAt, &updatedAt)
	if err != nil {
		return entity.User{}, err
	}

	// Agar avatar_id NULL bo'lsa, bo'sh qator sifatida saqlanadi
	if avatarID.Valid {
		response.AvatarId = avatarID.String
	} else {
		response.AvatarId = ""
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
		Select(`id, full_name, email, bio, username, user_type, user_role, status, avatar_id, gender, created_at, updated_at`).
		From("users")

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
		err = rows.Scan(&item.ID, &item.FullName, &item.Email, &item.Bio, &item.Username,
			&item.UserType, &item.UserRole, &item.Status, &item.AvatarId, &item.Gender, &createdAt, &updatedAt)
		if err != nil {
			return response, err
		}

		item.CreatedAt = createdAt.Format(time.RFC3339)
		item.UpdatedAt = updatedAt.Format(time.RFC3339)

		response.Items = append(response.Items, item)
	}

	countQuery, args, err := r.pg.Builder.Select("COUNT(1)").From("users").Where(where).ToSql()
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
	mp := make(map[string]interface{})

	if req.FullName != "" && req.FullName != "string" {
		mp["full_name"] = req.FullName
	}
	if req.Username != "" && req.Username != "string" {
		mp["username"] = req.Username
	}
	if req.Status != "" && req.Status != "string" {
		mp["status"] = req.Status
	}
	if req.Email != "" && req.Email != "string" {
		mp["email"] = req.Email
	}
	if req.Bio != "" && req.Bio != "string" {
		mp["bio"] = req.Bio
	}
	if req.AvatarId != "" && req.AvatarId != "string" {
		mp["avatar_id"] = req.AvatarId
	}
	if req.Gender != "" && req.Gender != "string" {
		mp["gender"] = req.Gender
	}
	if req.UserRole != "" && req.UserRole != "string" {
		mp["user_role"] = req.UserRole
	}
	if req.Password != "" && req.Password != "string" {
		mp["password_hash"] = req.Password
	}

	// Doim yangilanishi kerak bo'lgan maydonlar
	mp["updated_at"] = "now()"

	// Xarita bo'sh bo'lsa, yangilashni davom ettirmaymiz
	if len(mp) == 0 {
		return entity.User{}, errors.New("no fields to update")
	}

	query, args, err := r.pg.Builder.Update("users").SetMap(mp).Where("id = ?", req.ID).ToSql()
	if err != nil {
		return entity.User{}, err
	}

	_, err = r.pg.Pool.Exec(ctx, query, args...)
	if err != nil {
		return entity.User{}, err
	}
	res, err := r.GetSingle(ctx, entity.UserSingleRequest{ID: req.ID})
	if err != nil {
		return entity.User{}, err
	}
	res.Password = " "
	return res, nil
}

func (r *UserRepo) Delete(ctx context.Context, req entity.Id) error {
	qeury, args, err := r.pg.Builder.Delete("users").Where("id = ?", req.ID).ToSql()
	if err != nil {
		return err
	}

	_, err = r.pg.Pool.Exec(ctx, qeury, args...)
	if err != nil {
		return err
	}

	return nil
}
