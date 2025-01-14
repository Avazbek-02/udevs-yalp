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

// NewNotificationRepo creates a new instance of NotificationRepo
func NewNotificationRepo(pg *postgres.Postgres, config *config.Config, logger *logger.Logger) *NotificationRepo {
	return &NotificationRepo{
		pg:     pg,
		config: config,
		logger: logger,
	}
}

// Create adds a new notification entry into the database
func (r *NotificationRepo) Create(ctx context.Context, req entity.Notification) (entity.Notification, error) {
	req.ID = uuid.NewString()

	query, args, err := r.pg.Builder.Insert("notifications").
		Columns(`id, user_id, message, status`).
		Values(req.ID, req.UserID, req.Message, req.Status).ToSql()
	if err != nil {
		return entity.Notification{}, err
	}

	_, err = r.pg.Pool.Exec(ctx, query, args...)
	if err != nil {
		return entity.Notification{}, err
	}

	return req, nil
}

// GetSingle retrieves a single notification by its ID
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

	queryBuilder := r.pg.Builder.
		Select(`id, user_id, message, status, created_at`).
		From("notifications")

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
		var item entity.Notification
		err = rows.Scan(&item.ID, &item.UserID, &item.Message, &item.Status, &createdAt)
		if err != nil {
			return response, err
		}

		item.CreatedAt = createdAt.Format(time.RFC3339)
		response.Notifications = append(response.Notifications, item)
	}

	countQuery, args, err := r.pg.Builder.Select("COUNT(1)").From("notifications").Where(where).ToSql()
	if err != nil {
		return response, err
	}

	err = r.pg.Pool.QueryRow(ctx, countQuery, args...).Scan(&response.TotalCount)
	if err != nil {
		return response, err
	}

	return response, nil
}

// Update updates the details of a notification
func (r *NotificationRepo) Update(ctx context.Context, req entity.Notification) (entity.Notification, error) {
	mp := make(map[string]interface{})

	if req.Message != "" {
		mp["message"] = req.Message
	}
	if req.Status != "" {
		mp["status"] = req.Status
	}

	mp["created_at"] = "now()"

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

// Delete deletes a notification by its ID
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

// UpdateStatus updates the status of a notification (e.g., 'read' or 'unread')
func (r *NotificationRepo) UpdateStatus(ctx context.Context, req entity.Notification) (entity.Notification, error) {
	query, args, err := r.pg.Builder.Update("notifications").
		Set("status", req.Status).
		Where("id = ?", req.ID).
		ToSql()
	if err != nil {
		return entity.Notification{}, err
	}

	_, err = r.pg.Pool.Exec(ctx, query, args...)
	if err != nil {
		return entity.Notification{}, err
	}

	// Return the updated notification
	return r.GetSingle(ctx, entity.Id{ID: req.ID})
}
