package repo

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/Avazbek-02/udevslab-lesson6/config"
	"github.com/Avazbek-02/udevslab-lesson6/internal/entity"
	"github.com/Avazbek-02/udevslab-lesson6/pkg/logger"
	"github.com/Avazbek-02/udevslab-lesson6/pkg/postgres"
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

	qeury, args, err := r.pg.Builder.Insert("sessions").
		Columns(`id, user_id, ip_address, user_agent, is_active, expires_at, platform`).
		Values(req.ID, req.UserID, req.IPAddress, req.UserAgent, req.IsActive, expireDate, req.Platform).ToSql()
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
		createdAt, updatedAt    time.Time
		expiresAt, lastActiveAt sql.NullTime
	)
	qeuryBuilder := r.pg.Builder.
		Select(`id, user_id, ip_address, user_agent, is_active, expires_at, last_active_at, platform, created_at, updated_at`).
		From("sessions").Where("id = ?", req.ID)

	qeury, args, err := qeuryBuilder.ToSql()
	if err != nil {
		return entity.Session{}, err
	}

	err = r.pg.Pool.QueryRow(ctx, qeury, args...).
		Scan(&response.ID, &response.UserID, &response.IPAddress, &response.UserAgent,
			&response.IsActive, &expiresAt, &lastActiveAt, &response.Platform, &createdAt, &updatedAt)
	if err != nil {
		return entity.Session{}, err
	}

	response.CreatedAt = createdAt.Format(time.RFC3339)
	response.UpdatedAt = updatedAt.Format(time.RFC3339)
	if expiresAt.Valid {
		response.ExpiresAt = expiresAt.Time.Format(time.RFC3339)
	}

	if lastActiveAt.Valid {
		response.LastActiveAt = lastActiveAt.Time.Format(time.RFC3339)
	}

	return response, nil
}

func (r *SessionRepo) GetList(ctx context.Context, req entity.GetListFilter) (entity.SessionList, error) {
	var response entity.SessionList
	var createdAt, updatedAt time.Time

	// Start building the query
	queryBuilder := r.pg.Builder.
		Select(`id, user_id, ip_address, user_agent, is_active, expires_at, last_active_at, platform, created_at, updated_at`).
		From("sessions")

	// Apply filters to the query
	if req.Filters != nil {
		for _, filter := range req.Filters {
			switch filter.Column {
			case "user_id":
				if filter.Type == "eq" && filter.Value != "" {
					queryBuilder = queryBuilder.Where("user_id = ?", filter.Value)
				}
			case "ip_address":
				if filter.Type == "eq" && filter.Value != "" {
					queryBuilder = queryBuilder.Where("ip_address = ?", filter.Value)
				}
			case "is_active":
				if filter.Type == "eq" && filter.Value != "" {
					queryBuilder = queryBuilder.Where("is_active = ?", filter.Value)
				}
			}
		}
	}

	// Apply pagination (LIMIT and OFFSET)
	if req.Limit > 0 {
		queryBuilder = queryBuilder.Limit(uint64(req.Limit))
	}
	if req.Page > 0 {
		offset := (req.Page - 1) * req.Limit
		queryBuilder = queryBuilder.Offset(uint64(offset))
	}

	// Prepare and execute the SQL query
	query, args, err := queryBuilder.ToSql()
	if err != nil {
		return response, fmt.Errorf("error building query: %w", err)
	}

	rows, err := r.pg.Pool.Query(ctx, query, args...)
	if err != nil {
		return response, fmt.Errorf("error executing query: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		var item entity.Session
		var expiresAt, lastActiveAt sql.NullTime
		err = rows.Scan(&item.ID, &item.UserID, &item.IPAddress, &item.UserAgent,
			&item.IsActive, &expiresAt, &lastActiveAt, &item.Platform, &createdAt, &updatedAt)
		if err != nil {
			return response, fmt.Errorf("error scanning row: %w", err)
		}

		item.CreatedAt = createdAt.Format(time.RFC3339)
		item.UpdatedAt = updatedAt.Format(time.RFC3339)
		if expiresAt.Valid {
			item.ExpiresAt = expiresAt.Time.Format(time.RFC3339)
		}
		if lastActiveAt.Valid {
			item.LastActiveAt = lastActiveAt.Time.Format(time.RFC3339)
		}

		response.Items = append(response.Items, item)
	}

	// Count query to get the total number of records
	countQueryBuilder := r.pg.Builder.Select("COUNT(1)").From("sessions")
	if req.Filters != nil {
		for _, filter := range req.Filters {
			switch filter.Column {
			case "user_id":
				if filter.Type == "eq" && filter.Value != "" {
					countQueryBuilder = countQueryBuilder.Where("user_id = ?", filter.Value)
				}
			}
		}
	}

	// Prepare and execute the count query
	countQuery, args, err := countQueryBuilder.ToSql()
	if err != nil {
		return response, fmt.Errorf("error building count query: %w", err)
	}

	err = r.pg.Pool.QueryRow(ctx, countQuery, args...).Scan(&response.Count)
	if err != nil {
		return response, fmt.Errorf("error executing count query: %w", err)
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

	qeury, args, err := r.pg.Builder.Update("sessions").SetMap(mp).Where("id = ?", req.ID).ToSql()
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
	qeury, args, err := r.pg.Builder.Delete("sessions").Where("id = ?", req.ID).ToSql()
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

	qeury, args, err := r.pg.Builder.Update("sessions").SetMap(mp).Where(PrepareFilter(req.Filter)).ToSql()
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
