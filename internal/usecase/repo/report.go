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

func NewReportRepo(pg *postgres.Postgres, config *config.Config, logger *logger.Logger) *ReportRepo {
	return &ReportRepo{
		pg:     pg,
		config: config,
		logger: logger,
	}
}

func (r *ReportRepo) Create(ctx context.Context, req entity.Report) (entity.Report, error) {
	req.ID = uuid.NewString()
	fmt.Println("::::::",req.UserID, req.BusinessID)
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

func (r *ReportRepo) GetList(ctx context.Context, req entity.GetListFilter) (entity.ReportList, error) {
	var response = entity.ReportList{}
	var createdAt time.Time

	// Start building the query
	queryBuilder := r.pg.Builder.
		Select(`id, user_id, business_id, reason, created_at`).
		From("reports")

	// Apply filters (if any)
	if req.Filters != nil {
		for _, filter := range req.Filters {
			if filter.Column == "business_id" {
				if filter.Type == "eq" && filter.Value != "" {
					// Add filter for business_id
					queryBuilder = queryBuilder.Where("business_id = ?", filter.Value)
				}
				if filter.Type == "eq" && filter.Value == "" {
					// If business_id is empty, return all reports
					break
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

	// Prepare the SQL query for fetching the reports
	query, args, err := queryBuilder.ToSql()
	if err != nil {
		return response, err
	}

	// Execute the query for reports
	rows, err := r.pg.Pool.Query(ctx, query, args...)
	if err != nil {
		return response, err
	}
	defer rows.Close()

	// Scan the report records into response.Reports
	for rows.Next() {
		var item entity.Report
		err = rows.Scan(&item.ID, &item.UserID, &item.BusinessID, &item.Reason, &createdAt)
		if err != nil {
			return response, err
		}

		item.CreatedAt = createdAt.Format(time.RFC3339)
		response.Reports = append(response.Reports, item)
	}

	// Now, count the reports based on the same filter (business_id)
	countQueryBuilder := r.pg.Builder.Select("COUNT(1)").From("reports")
	if req.Filters != nil {
		for _, filter := range req.Filters {
			if filter.Column == "business_id" && filter.Type == "eq" && filter.Value != "" {
				// Add the same filter to the COUNT query
				countQueryBuilder = countQueryBuilder.Where("business_id = ?", filter.Value)
			}
		}
	}

	// Prepare and execute the COUNT query
	countQuery, args, err := countQueryBuilder.ToSql()
	if err != nil {
		return response, err
	}

	// Execute the count query
	err = r.pg.Pool.QueryRow(ctx, countQuery, args...).Scan(&response.Count)
	if err != nil {
		return response, err
	}

	return response, nil
}

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