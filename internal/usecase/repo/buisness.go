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

type BusinessRepo struct {
	pg     *postgres.Postgres
	config *config.Config
	logger *logger.Logger
}

// NewBusinessRepo creates a new instance of BusinessRepo
func NewBusinessRepo(pg *postgres.Postgres, config *config.Config, logger *logger.Logger) *BusinessRepo {
	return &BusinessRepo{
		pg:     pg,
		config: config,
		logger: logger,
	}
}

// Create adds a new business entry into the database
func (r *BusinessRepo) Create(ctx context.Context, req entity.Business) (entity.Business, error) {
	req.ID = uuid.NewString()

	query, args, err := r.pg.Builder.Insert("businesses").
		Columns(`id, owner_id, name, description, category, address, contact_info, photos`).
		Values(req.ID, req.OwnerID, req.Name, req.Description, req.Category, req.Address, req.ContactInfo, req.Photos).ToSql()
	if err != nil {
		return entity.Business{}, err
	}

	_, err = r.pg.Pool.Exec(ctx, query, args...)
	if err != nil {
		return entity.Business{}, err
	}

	return req, nil
}

func (r *BusinessRepo) GetSingle(ctx context.Context, req entity.Id) (entity.Business, error) {
	response := entity.Business{}
	var (
		createdAt, updatedAt time.Time
	)

	queryBuilder := r.pg.Builder.
		Select(`id, owner_id, name, description, category, address, contact_info, photos, created_at, updated_at`).
		From("businesses")

	switch {
	case req.ID != "":
		queryBuilder = queryBuilder.Where("id = ?", req.ID)

	default:
		return entity.Business{}, fmt.Errorf("GetSingle - invalid request")
	}

	query, args, err := queryBuilder.ToSql()
	if err != nil {
		return entity.Business{}, err
	}

	err = r.pg.Pool.QueryRow(ctx, query, args...).
		Scan(&response.ID, &response.OwnerID, &response.Name, &response.Description, &response.Category,
			&response.Address, &response.ContactInfo, &response.Photos, &createdAt, &updatedAt)
	if err != nil {
		return entity.Business{}, err
	}

	response.CreatedAt = createdAt.Format(time.RFC3339)
	response.UpdatedAt = updatedAt.Format(time.RFC3339)
	return response, nil
}

func (r *BusinessRepo) GetList(ctx context.Context, req entity.GetListFilter) (entity.BusinessList, error) {
	var response = entity.BusinessList{}
	var createdAt, updatedAt time.Time

	// Start building the SQL query
	queryBuilder := r.pg.Builder.
		Select(`id, owner_id, name, description, category, address, contact_info, photos, created_at, updated_at`).
		From("businesses")

	// Apply filters (if any)
	if req.Filters != nil {
		for _, filter := range req.Filters {
			if filter.Column == "owner_id" {
				if filter.Type == "eq" && filter.Value != "" {
					// Add filter for owner_id
					queryBuilder = queryBuilder.Where("owner_id = ?", filter.Value)
				}
				if filter.Type == "eq" && filter.Value == "" {
					// If owner_id is empty, return all businesses
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

	// Prepare the SQL query for fetching businesses
	query, args, err := queryBuilder.ToSql()
	if err != nil {
		return response, err
	}

	// Execute the query for businesses
	rows, err := r.pg.Pool.Query(ctx, query, args...)
	if err != nil {
		return response, err
	}
	defer rows.Close()

	// Scan the business records into response.Items
	for rows.Next() {
		var item entity.Business
		err = rows.Scan(&item.ID, &item.OwnerID, &item.Name, &item.Description, &item.Category,
			&item.Address, &item.ContactInfo, &item.Photos, &createdAt, &updatedAt)
		if err != nil {
			return response, err
		}

		item.CreatedAt = createdAt.Format(time.RFC3339)
		item.UpdatedAt = updatedAt.Format(time.RFC3339)

		response.Items = append(response.Items, item)
	}

	// Get the total count of businesses based on filters (if any)
	countQueryBuilder := r.pg.Builder.Select("COUNT(1)").From("businesses")
	if req.Filters != nil {
		for _, filter := range req.Filters {
			if filter.Column == "owner_id" && filter.Type == "eq" && filter.Value != "" {
				// Add the same filter to the COUNT query
				countQueryBuilder = countQueryBuilder.Where("owner_id = ?", filter.Value)
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



func (r *BusinessRepo) Update(ctx context.Context, req entity.Business) (entity.Business, error) {
	mp := make(map[string]interface{})

	if req.Name != "" {
		mp["name"] = req.Name
	}
	if req.Description != "" {
		mp["description"] = req.Description
	}
	if req.Category != "" {
		mp["category"] = req.Category
	}
	if req.Address != "" {
		mp["address"] = req.Address
	}
	if req.ContactInfo != "" {
		mp["contact_info"] = req.ContactInfo
	}
	if req.Photos != "" {
		mp["photos"] = req.Photos
	}

	mp["updated_at"] = "now()"

	if len(mp) == 0 {
		return entity.Business{}, errors.New("no fields to update")
	}

	query, args, err := r.pg.Builder.Update("businesses").SetMap(mp).Where("id = ?", req.ID).ToSql()
	if err != nil {
		return entity.Business{}, err
	}

	_, err = r.pg.Pool.Exec(ctx, query, args...)
	if err != nil {
		return entity.Business{}, err
	}

	res, err := r.GetSingle(ctx, entity.Id{ID: req.ID})
	if err != nil {
		return entity.Business{}, err
	}

	return res, nil
}

// Delete deletes a business by its I
func (r *BusinessRepo) Delete(ctx context.Context, req entity.Id) error {
	query, args, err := r.pg.Builder.Delete("businesses").Where("id = ?", req.ID).ToSql()
	if err != nil {
		return err
	}

	_, err = r.pg.Pool.Exec(ctx, query, args...)
	if err != nil {
		return err
	}

	return nil
}
