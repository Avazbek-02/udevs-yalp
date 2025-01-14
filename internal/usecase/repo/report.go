package repo

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/Avazbek-02/udevslab-lesson6/config"
	"github.com/Avazbek-02/udevslab-lesson6/internal/entity"
	"github.com/Avazbek-02/udevslab-lesson6/pkg/logger"
	"github.com/Avazbek-02/udevslab-lesson6/pkg/postgres"
	"github.com/google/uuid"
)

type ReportRepo struct {
	pg     *postgres.Postgres
	config *config.Config
	logger *logger.Logger
}

// NewReportRepo creates a new instance of ReportRepo
func NewReportRepo(pg *postgres.Postgres, config *config.Config, logger *logger.Logger) *ReportRepo {
	return &ReportRepo{
		pg:     pg,
		config: config,
		logger: logger,
	}
}

// Create adds a new report entry into the database
func (r *ReportRepo) Create(ctx context.Context, req entity.Report) (entity.Report, error) {
	req.ID = uuid.NewString()

	query, args, err := r.pg.Builder.Insert("reports").
		Columns(`id, user_id, business_id, reason`).
		Values(req.ID, req.UserID, req.BusinessID, req.Reason).ToSql()
	if err != nil {
		return entity.Report{}, err
	}

	_, err = r.pg.Pool.Exec(ctx, query, args...)
	if err != nil {
		return entity.Report{}, err
	}

	return req, nil
}

// GetSingle retrieves a single report by its ID
func (r *ReportRepo) GetSingle(ctx context.Context, req entity.Id) (entity.Report, error) {
	response := entity.Report{}
	var createdAt time.Time

	queryBuilder := r.pg.Builder.
		Select(`id, user_id, business_id, reason, created_at`).
		From("reports")

	switch {
	case req.ID != "":
		queryBuilder = queryBuilder.Where("id = ?", req.ID)

	default:
		return entity.Report{}, fmt.Errorf("GetSingle - invalid request")
	}

	query, args, err := queryBuilder.ToSql()
	if err != nil {
		return entity.Report{}, err
	}

	err = r.pg.Pool.QueryRow(ctx, query, args...).
		Scan(&response.ID, &response.UserID, &response.BusinessID, &response.Reason, &createdAt)
	if err != nil {
		return entity.Report{}, err
	}

	response.CreatedAt = createdAt.Format(time.RFC3339)
	return response, nil
}

// GetList retrieves a list of reports based on filters
func (r *ReportRepo) GetList(ctx context.Context, req entity.GetListFilter) (entity.ReportList, error) {
	var response = entity.ReportList{}
	var createdAt time.Time

	queryBuilder := r.pg.Builder.
		Select(`id, user_id, business_id, reason, created_at`).
		From("reports")

	queryBuilder, where := PrepareGetListQuery(queryBuilder, req)

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
		var item entity.Report
		err = rows.Scan(&item.ID, &item.UserID, &item.BusinessID, &item.Reason, &createdAt)
		if err != nil {
			return response, err
		}

		item.CreatedAt = createdAt.Format(time.RFC3339)
		response.Reports = append(response.Reports, item)
	}

	countQuery, args, err := r.pg.Builder.Select("COUNT(1)").From("reports").Where(where).ToSql()
	if err != nil {
		return response, err
	}

	err = r.pg.Pool.QueryRow(ctx, countQuery, args...).Scan(&response.Count)
	if err != nil {
		return response, err
	}

	return response, nil
}

// Update updates the details of a report
func (r *ReportRepo) Update(ctx context.Context, req entity.Report) (entity.Report, error) {
	mp := make(map[string]interface{})

	if req.Reason != "" {
		mp["reason"] = req.Reason
	}

	mp["created_at"] = "now()"

	if len(mp) == 0 {
		return entity.Report{}, errors.New("no fields to update")
	}

	query, args, err := r.pg.Builder.Update("reports").SetMap(mp).Where("id = ?", req.ID).ToSql()
	if err != nil {
		return entity.Report{}, err
	}

	_, err = r.pg.Pool.Exec(ctx, query, args...)
	if err != nil {
		return entity.Report{}, err
	}

	res, err := r.GetSingle(ctx, entity.Id{ID: req.ID})
	if err != nil {
		return entity.Report{}, err
	}

	return res, nil
}

// Delete deletes a report by its ID
func (r *ReportRepo) Delete(ctx context.Context, req entity.Id) error {
	query, args, err := r.pg.Builder.Delete("reports").Where("id = ?", req.ID).ToSql()
	if err != nil {
		return err
	}

	_, err = r.pg.Pool.Exec(ctx, query, args...)
	if err != nil {
		return err
	}

	return nil
}
