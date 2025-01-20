package repo

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/Avazbek-02/udevslab-lesson6/config"
	"github.com/Avazbek-02/udevslab-lesson6/internal/entity"
	"github.com/Avazbek-02/udevslab-lesson6/pkg/logger"
	"github.com/Avazbek-02/udevslab-lesson6/pkg/postgres"
	"github.com/google/uuid"
)

type ReviewRepo struct {
	pg     *postgres.Postgres
	config *config.Config
	logger *logger.Logger
}

func NewReviewRepo(pg *postgres.Postgres, config *config.Config, logger *logger.Logger) *ReviewRepo {
	return &ReviewRepo{
		pg:     pg,
		config: config,
		logger: logger,
	}
}

func (r *ReviewRepo) Create(ctx context.Context, req entity.Review) (entity.Review, error) {
	req.ID = uuid.NewString()
	photosJSON, err := json.Marshal(req.Photos)
	if err != nil {
		return entity.Review{}, fmt.Errorf("failed to marshal photos: %v", err)
	}
	fmt.Println(req)
	query, args, err := r.pg.Builder.Insert("reviews").
		Columns(`id, user_id, business_id, rating, feedback, photos`).
		Values(req.ID, req.UserID, req.BusinessID, req.Rating, req.Feedback, photosJSON).ToSql()
	if err != nil {
		return entity.Review{}, err
	}

	_, err = r.pg.Pool.Exec(ctx, query, args...)
	if err != nil {
		return entity.Review{}, err
	}
	
	return req, nil
}

func (r *ReviewRepo) GetSingle(ctx context.Context, req entity.Id) (entity.Review, error) {
	response := entity.Review{}
	var createdAt time.Time

	queryBuilder := r.pg.Builder.
		Select(`id, user_id, business_id, rating, feedback, photos, created_at`).
		From("reviews")

	switch {
	case req.ID != "":
		queryBuilder = queryBuilder.Where("id = ?", req.ID)

	default:
		return entity.Review{}, fmt.Errorf("GetSingle - invalid request")
	}

	query, args, err := queryBuilder.ToSql()
	if err != nil {
		return entity.Review{}, err
	}

	err = r.pg.Pool.QueryRow(ctx, query, args...).
		Scan(&response.ID, &response.UserID, &response.BusinessID, &response.Rating, &response.Feedback, &response.Photos, &createdAt)
	if err != nil {
		return entity.Review{}, err
	}

	response.CreatedAt = createdAt.Format(time.RFC3339)
	return response, nil
}

func (r *ReviewRepo) GetList(ctx context.Context, req entity.GetListFilter) (entity.ReviewList, error) {
	var response = entity.ReviewList{}
	var createdAt time.Time

	queryBuilder := r.pg.Builder.
		Select(`id, user_id, business_id, rating, feedback, photos, created_at`).
		From("reviews")

	// Apply filters
	if req.Filters != nil {
		for _, filter := range req.Filters {
			if filter.Column == "business_id" {
				if filter.Type == "eq" && filter.Value != "" {
					// Add filter for business_id
					queryBuilder = queryBuilder.Where("business_id = ?", filter.Value)
				}
				if filter.Type == "eq" && filter.Value == "" {
					// If business_id is empty, return all reviews
					break
				}
			}
		}
	}

	// Prepare the SQL query for fetching the reviews
	query, args, err := queryBuilder.ToSql()
	if err != nil {
		return response, err
	}

	// Execute the query for reviews
	rows, err := r.pg.Pool.Query(ctx, query, args...)
	if err != nil {
		return response, err
	}
	defer rows.Close()

	// Scan the review records into response.Items
	for rows.Next() {
		var item entity.Review
		err = rows.Scan(&item.ID, &item.UserID, &item.BusinessID, &item.Rating, &item.Feedback, &item.Photos, &createdAt)
		if err != nil {
			return response, err
		}

		item.CreatedAt = createdAt.Format(time.RFC3339)
		response.Items = append(response.Items, item)
	}

	// Now, count the reviews based on the same filter (business_id)
	countQueryBuilder := r.pg.Builder.Select("COUNT(1)").From("reviews")
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

func (r *ReviewRepo) Update(ctx context.Context, req entity.Review) (entity.Review, error) {
	mp := make(map[string]interface{})

	if req.Rating != 0 {
		mp["rating"] = req.Rating
	}
	if req.Feedback != "" {
		mp["feedback"] = req.Feedback
	}
	if req.Photos != "" {
		mp["photos"] = req.Photos
	}

	mp["created_at"] = "now()"

	if len(mp) == 0 {
		return entity.Review{}, errors.New("no fields to update")
	}

	query, args, err := r.pg.Builder.Update("reviews").SetMap(mp).Where("id = ?", req.ID).ToSql()
	if err != nil {
		return entity.Review{}, err
	}

	_, err = r.pg.Pool.Exec(ctx, query, args...)
	if err != nil {
		return entity.Review{}, err
	}

	res, err := r.GetSingle(ctx, entity.Id{ID: req.ID})
	if err != nil {
		return entity.Review{}, err
	}

	return res, nil
}

func (r *ReviewRepo) Delete(ctx context.Context, req entity.Id) error {
	query, args, err := r.pg.Builder.Delete("reviews").Where("id = ?", req.ID).ToSql()
	if err != nil {
		return err
	}

	_, err = r.pg.Pool.Exec(ctx, query, args...)
	if err != nil {
		return err
	}

	return nil
}