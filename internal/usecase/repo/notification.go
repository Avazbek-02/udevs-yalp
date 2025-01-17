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

type NotificationRepo struct {
	pg     *postgres.Postgres
	config *config.Config
	logger *logger.Logger
}

func NewNotificationRepo(pg *postgres.Postgres, config *config.Config, logger *logger.Logger) *NotificationRepo {
	return &NotificationRepo{
		pg:     pg,
		config: config,
		logger: logger,
	}
}

func (r *NotificationRepo) Create(ctx context.Context, req entity.Notification) (entity.Notification, error) {
	req.ID = uuid.NewString()

	query, args, err := r.pg.Builder.Insert("notifications").
		Columns(`id,owner_id, user_id,email, message, status`).
		Values(req.ID, req.OwnerId, req.UserID, req.Email, req.Message, req.Status).ToSql()
	if err != nil {
		return entity.Notification{}, err
	}

	_, err = r.pg.Pool.Exec(ctx, query, args...)
	if err != nil {
		return entity.Notification{}, err
	}

	return req, nil
}

func (r *NotificationRepo) GetSingle(ctx context.Context, req entity.Id) (entity.Notification, error) {
	response := entity.Notification{}
	var createdAt time.Time

	queryBuilder := r.pg.Builder.
		Select(`id, user_id, message, status, created_at`).
		From("notifications")

	switch {
	case req.ID != "":
		queryBuilder = queryBuilder.Where("id = ?", req.ID)

	default:
		return entity.Notification{}, fmt.Errorf("GetSingle - invalid request")
	}

	query, args, err := queryBuilder.ToSql()
	if err != nil {
		return entity.Notification{}, err
	}

	err = r.pg.Pool.QueryRow(ctx, query, args...).
		Scan(&response.ID, &response.UserID, &response.Message, &response.Status, &createdAt)
	if err != nil {
		return entity.Notification{}, err
	}

	response.CreatedAt = createdAt.Format(time.RFC3339)
	return response, nil
}

func (r *NotificationRepo) GetList(ctx context.Context, req entity.GetListFilter) (entity.NotificationList, error) {
	var response = entity.NotificationList{}
	var createdAt time.Time

	// Start building the query
	queryBuilder := r.pg.Builder.
		Select(`id, user_id, message, status, created_at`).
		From("notifications")

	// Apply filters (if any)
	if req.Filters != nil {
		for _, filter := range req.Filters {
			if filter.Column == "user_id" {
				if filter.Type == "eq" && filter.Value != "" {
					// Add filter for user_id
					queryBuilder = queryBuilder.Where("user_id = ?", filter.Value)
				}
				if filter.Type == "eq" && filter.Value == "" {
					// If user_id is empty, return all notifications
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

	// Prepare the SQL query for fetching the notifications
	query, args, err := queryBuilder.ToSql()
	if err != nil {
		return response, err
	}

	// Execute the query for notifications
	rows, err := r.pg.Pool.Query(ctx, query, args...)
	if err != nil {
		return response, err
	}
	defer rows.Close()

	// Scan the notification records into response.Notifications
	for rows.Next() {
		var item entity.Notification
		err = rows.Scan(&item.ID, &item.UserID, &item.Message, &item.Status, &createdAt)
		if err != nil {
			return response, err
		}

		item.CreatedAt = createdAt.Format(time.RFC3339)
		response.Notifications = append(response.Notifications, item)
	}

	// Get the total count of notifications based on filters (if any)
	countQueryBuilder := r.pg.Builder.Select("COUNT(1)").From("notifications")
	if req.Filters != nil {
		for _, filter := range req.Filters {
			if filter.Column == "user_id" && filter.Type == "eq" && filter.Value != "" {
				// Add the same filter to the COUNT query
				countQueryBuilder = countQueryBuilder.Where("user_id = ?", filter.Value)
			}
		}
	}

	// Prepare and execute the COUNT query
	countQuery, args, err := countQueryBuilder.ToSql()
	if err != nil {
		return response, err
	}

	// Execute the count query
	err = r.pg.Pool.QueryRow(ctx, countQuery, args...).Scan(&response.TotalCount)
	if err != nil {
		return response, err
	}

	return response, nil
}

func (r *NotificationRepo) Update(ctx context.Context, req entity.Notification) (entity.Notification, error) {
	mp := make(map[string]interface{})
	userRes, err := r.GetSingle(ctx, entity.Id{ID: req.ID})
	if err != nil {
		return entity.Notification{}, err
	}

	if userRes.OwnerId != req.OwnerId && (req.OwnerRole != "admin" && req.OwnerRole != "superadmin") {
		return entity.Notification{}, fmt.Errorf("you are not allowed to update this notification")
	}
	if req.Message != "" {
		mp["message"] = req.Message
	}
	if req.Status != "" {
		mp["status"] = req.Status
	}

	mp["created_at"] = "CURRENT_TIMESTAMP"

	if len(mp) == 0 {
		return entity.Notification{}, errors.New("no fields to update")
	}

	query, args, err := r.pg.Builder.Update("notifications").SetMap(mp).Where("id = ?", req.ID).ToSql()
	if err != nil {
		return entity.Notification{}, err
	}

	_, err = r.pg.Pool.Exec(ctx, query, args...)
	if err != nil {
		return entity.Notification{}, err
	}

	res, err := r.GetSingle(ctx, entity.Id{ID: req.ID})
	if err != nil {
		return entity.Notification{}, err
	}

	return res, nil
}

func (r *NotificationRepo) Delete(ctx context.Context, req entity.Id) error {
	query, args, err := r.pg.Builder.Delete("notifications").Where("id = ?", req.ID).ToSql()
	if err != nil {
		return err
	}

	_, err = r.pg.Pool.Exec(ctx, query, args...)
	if err != nil {
		return err
	}

	return nil
}

func (r *NotificationRepo) UpdateStatus(ctx context.Context, req entity.Notification) (entity.Notification, error) {
	query, args, err := r.pg.Builder.Update("notifications").
		Set("status", "read").
		Where("id = ?", req.ID).
		ToSql()
	if err != nil {
		return entity.Notification{}, err
	}

	_, err = r.pg.Pool.Exec(ctx, query, args...)
	if err != nil {
		return entity.Notification{}, err
	}

	return r.GetSingle(ctx, entity.Id{ID: req.ID})
}
