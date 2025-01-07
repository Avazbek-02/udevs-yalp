package repo

import (
	"context"
	"database/sql"
	"time"

	"github.com/golanguzb70/system-design-openbudget/config"
	"github.com/golanguzb70/system-design-openbudget/internal/entity"
	"github.com/golanguzb70/system-design-openbudget/pkg/logger"
	"github.com/golanguzb70/system-design-openbudget/pkg/postgres"
	"github.com/google/uuid"
)

type SessionRepo struct {
	pg     *postgres.Postgres
	config *config.Config
	logger *logger.Logger
}

// New -.
func NewSessionRepo(pg *postgres.Postgres, config *config.Config, logger *logger.Logger) *SessionRepo {
	return &SessionRepo{
		pg:     pg,
		config: config,
		logger: logger,
	}
}

func (r *SessionRepo) Create(ctx context.Context, req entity.Session) (entity.Session, error) {
	req.ID = uuid.NewString()
	expireDate := sql.NullTime{}
	expiresat, err := time.Parse(time.RFC3339, req.ExpiresAt)
	if err == nil {
		expireDate.Time = expiresat
		expireDate.Valid = true
	}

	qeury, args, err := r.pg.Builder.Insert("auth_session").
		Columns(`id, user_id, ip_address, user_agent, is_active, expires_at`).
		Values(req.ID, req.UserID, req.IPAddress, req.UserAgent, req.IsActive, expireDate).ToSql()
	if err != nil {
		return entity.Session{}, err
	}

	_, err = r.pg.Pool.Exec(ctx, qeury, args...)
	if err != nil {
		return entity.Session{}, err
	}

	return req, nil
}

func (r *SessionRepo) GetSingle(ctx context.Context, req entity.Id) (entity.Session, error) {
	response := entity.Session{}
	var (
		createdAt, updatedAt, lastActiveAt time.Time
		expiresAt                          sql.NullTime
	)

	qeuryBuilder := r.pg.Builder.
		Select(`id, user_id, ip_address, user_agent, is_active, expires_at, last_active_at, created_at, updated_at`).
		From("auth_session").Where("id = ?", req.ID)

	qeury, args, err := qeuryBuilder.ToSql()
	if err != nil {
		return entity.Session{}, err
	}

	err = r.pg.Pool.QueryRow(ctx, qeury, args...).
		Scan(&response.ID, &response.UserID, &response.IPAddress, &response.UserAgent,
			&response.IsActive, &expiresAt, &lastActiveAt, &createdAt, &updatedAt)
	if err != nil {
		return entity.Session{}, err
	}

	response.CreatedAt = createdAt.Format(time.RFC3339)
	response.UpdatedAt = updatedAt.Format(time.RFC3339)
	response.LastActiveAt = lastActiveAt.Format(time.RFC3339)
	if expiresAt.Valid {
		response.ExpiresAt = expiresAt.Time.Format(time.RFC3339)
	}

	return response, nil
}

func (r *SessionRepo) GetList(ctx context.Context, req entity.GetListFilter) (entity.SessionList, error) {
	var (
		response = entity.SessionList{}
	)

	qeuryBuilder := r.pg.Builder.
		Select(`id, user_id, ip_address, user_agent, is_active, expires_at, last_active_at, created_at, updated_at`).
		From("auth_session")

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
		var (
			createdAt, updatedAt, lastActiveAt time.Time
			expiresAt                          sql.NullTime
			item                               entity.Session
		)
		err = rows.Scan(&item.ID, &item.UserID, &item.IPAddress, &item.UserAgent,
			&item.IsActive, &expiresAt, &lastActiveAt, &createdAt, &updatedAt)
		if err != nil {
			return response, err
		}

		item.CreatedAt = createdAt.Format(time.RFC3339)
		item.UpdatedAt = updatedAt.Format(time.RFC3339)
		item.LastActiveAt = lastActiveAt.Format(time.RFC3339)
		if expiresAt.Valid {
			item.ExpiresAt = expiresAt.Time.Format(time.RFC3339)
		}

		response.Items = append(response.Items, item)
	}

	countQuery, args, err := r.pg.Builder.Select("COUNT(1)").From("auth_session").Where(where).ToSql()
	if err != nil {
		return response, err
	}

	err = r.pg.Pool.QueryRow(ctx, countQuery, args...).Scan(&response.Count)
	if err != nil {
		return response, err
	}

	return response, nil
}

func (r *SessionRepo) Update(ctx context.Context, req entity.Session) (entity.Session, error) {
	mp := map[string]interface{}{
		"ip_address":     req.IPAddress,
		"is_active":      req.IsActive,
		"last_active_at": "now()",
		"updated_at":     "now()",
	}

	qeury, args, err := r.pg.Builder.Update("auth_session").SetMap(mp).Where("id = ?", req.ID).ToSql()
	if err != nil {
		return entity.Session{}, err
	}

	_, err = r.pg.Pool.Exec(ctx, qeury, args...)
	if err != nil {
		return entity.Session{}, err
	}

	return req, nil
}

func (r *SessionRepo) Delete(ctx context.Context, req entity.Id) error {
	qeury, args, err := r.pg.Builder.Delete("auth_session").Where("id = ?", req.ID).ToSql()
	if err != nil {
		return err
	}

	_, err = r.pg.Pool.Exec(ctx, qeury, args...)
	if err != nil {
		return err
	}

	return nil
}

func (r *SessionRepo) UpdateField(ctx context.Context, req entity.UpdateFieldRequest) (entity.RowsEffected, error) {
	mp := map[string]interface{}{}
	response := entity.RowsEffected{}

	for _, item := range req.Items {
		mp[item.Column] = item.Value
	}

	qeury, args, err := r.pg.Builder.Update("auth_session").SetMap(mp).Where(PrepareFilter(req.Filter)).ToSql()
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
