package repo

import (
	"context"
	"errors"
	"time"

	"github.com/Avazbek-02/udevslab-lesson6/config"
	"github.com/Avazbek-02/udevslab-lesson6/internal/entity"
	"github.com/Avazbek-02/udevslab-lesson6/pkg/logger"
	"github.com/Avazbek-02/udevslab-lesson6/pkg/postgres"
	"github.com/google/uuid"
)

type EventRepo struct {
	pg     *postgres.Postgres
	config *config.Config
	logger *logger.Logger
}

func NewEventRepo(pg *postgres.Postgres, config *config.Config, logger *logger.Logger) *EventRepo {
	return &EventRepo{
		pg:     pg,
		config: config,
		logger: logger,
	}
}

// Create adds a new event to the database
func (r *EventRepo) Create(ctx context.Context, req entity.Event) (entity.Event, error) {
	req.ID = uuid.NewString()
	query, args, err := r.pg.Builder.Insert("events").
		Columns("id, business_id, name, description, date, location").
		Values(req.ID, req.BusinessID, req.Name, req.Description, req.Date, req.Location).ToSql()
	if err != nil {
		return entity.Event{}, err
	}

	_, err = r.pg.Pool.Exec(ctx, query, args...)
	if err != nil {
		return entity.Event{}, err
	}

	return req, nil
}

// GetSingle retrieves a single event by ID
func (r *EventRepo) GetSingle(ctx context.Context, req entity.Id) (entity.Event, error) {
	var response entity.Event
	var createdAt time.Time

	query, args, err := r.pg.Builder.
		Select("id, business_id, name, description, date, location, created_at").
		From("events").
		Where("id = ?", req.ID).ToSql()
	if err != nil {
		return entity.Event{}, err
	}

	err = r.pg.Pool.QueryRow(ctx, query, args...).
		Scan(&response.ID, &response.BusinessID, &response.Name, &response.Description, &response.Date, &response.Location, &createdAt)
	if err != nil {
		return entity.Event{}, err
	}

	response.CreatedAt = createdAt.Format(time.RFC3339)
	return response, nil
}

// GetList retrieves a list of events based on filters
func (r *EventRepo) GetList(ctx context.Context, req entity.GetListFilter) (entity.EventList, error) {
	var response entity.EventList
	var createdAt time.Time

	queryBuilder := r.pg.Builder.
		Select("id, business_id, name, description, date, location, created_at").
		From("events")

	// Apply filters
	for _, filter := range req.Filters {
		if filter.Column == "business_id" && filter.Type == "eq" && filter.Value != "" {
			queryBuilder = queryBuilder.Where("business_id = ?", filter.Value)
		}
	}

	// Apply pagination
	if req.Limit > 0 {
		queryBuilder = queryBuilder.Limit(uint64(req.Limit))
	}
	if req.Page > 0 {
		offset := (req.Page - 1) * req.Limit
		queryBuilder = queryBuilder.Offset(uint64(offset))
	}

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
		var item entity.Event
		err = rows.Scan(&item.ID, &item.BusinessID, &item.Name, &item.Description, &item.Date, &item.Location, &createdAt)
		if err != nil {
			return response, err
		}

		item.CreatedAt = createdAt.Format(time.RFC3339)
		response.Events = append(response.Events, item)
	}

	// Count total events
	countQuery, args, err := r.pg.Builder.
		Select("COUNT(1)").
		From("events").
		ToSql()
	if err != nil {
		return response, err
	}

	err = r.pg.Pool.QueryRow(ctx, countQuery, args...).Scan(&response.Count)
	if err != nil {
		return response, err
	}

	return response, nil
}

// Update updates an event
func (r *EventRepo) Update(ctx context.Context, req entity.Event) (entity.Event, error) {
	mp := map[string]interface{}{}

	// Ma'lumotlarni qo'shish
	if req.Name != "" && req.Name != "string" {
		mp["name"] = req.Name
	}
	if req.Description != "" && req.Description != "string" {
		mp["description"] = req.Description
	}
	if req.Date != "" && req.Date != "string" {
		mp["date"] = req.Date
	}
	if req.Location != "" && req.Location != "string" {
		mp["location"] = req.Location
	}
	mp["created_at"] = "now()"

	if len(mp) == 0 {
		return entity.Event{}, errors.New("no fields to update")
	}

	query, args, err := r.pg.Builder.Update("events").SetMap(mp).Where("id = ?", req.ID).ToSql()
	if err != nil {
		return entity.Event{}, err
	}

	_, err = r.pg.Pool.Exec(ctx, query, args...)
	if err != nil {
		return entity.Event{}, err
	}
	return r.GetSingle(ctx, entity.Id{ID: req.ID})
}

// Delete removes an event by ID
func (r *EventRepo) Delete(ctx context.Context, id entity.Id) error {
	query, args, err := r.pg.Builder.Delete("events").Where("id = ?", id).ToSql()
	if err != nil {
		return err
	}

	_, err = r.pg.Pool.Exec(ctx, query, args...)
	return err
}

// AddParticipant adds a participant to an event
func (r *EventRepo) AddParticipant(ctx context.Context, req entity.EventParticipant) (entity.EventParticipant, error) {
	req.ID = uuid.NewString()
	query, args, err := r.pg.Builder.Insert("event_participants").
		Columns("id, event_id, user_id, joined_at").
		Values(req.ID, req.EventID, req.UserID, time.Now()).ToSql()
	if err != nil {
		return entity.EventParticipant{}, err
	}

	_, err = r.pg.Pool.Exec(ctx, query, args...)
	if err != nil {
		return entity.EventParticipant{}, err
	}

	return req, nil
}

// RemoveParticipant removes a participant from an event
func (r *EventRepo) RemoveParticipant(ctx context.Context, req entity.EventParticipant) error {
	query, args, err := r.pg.Builder.Delete("event_participants").Where("event_id = ? AND user_id = ?", req.EventID, req.UserID).ToSql()
	if err != nil {
		return err
	}

	_, err = r.pg.Pool.Exec(ctx, query, args...)
	return err
}

// GetParticipants retrieves participants of an event
func (r *EventRepo) GetParticipants(ctx context.Context, req entity.GetListFilter) (entity.EventParticipantList, error) {
	var response entity.EventParticipantList

	queryBuilder := r.pg.Builder.
		Select("id, event_id, user_id, joined_at").
		From("event_participants")

	// Apply filters
	for _, filter := range req.Filters {
		if filter.Column == "event_id" && filter.Type == "eq" && filter.Value != "" {
			queryBuilder = queryBuilder.Where("event_id = ?", filter.Value)
		}
	}

	// Apply pagination
	if req.Limit > 0 {
		queryBuilder = queryBuilder.Limit(uint64(req.Limit))
	}
	if req.Page > 0 {
		offset := (req.Page - 1) * req.Limit
		queryBuilder = queryBuilder.Offset(uint64(offset))
	}

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
		var participant entity.EventParticipant
		var joinedAt time.Time
		err = rows.Scan(&participant.ID, &participant.EventID, &participant.UserID, &joinedAt)
		if err != nil {
			return response, err
		}
		participant.JoinedAt = joinedAt.Format(time.RFC3339)
		response.Participants = append(response.Participants, participant)
	}

	// Count total participants
	countQuery, args, err := r.pg.Builder.Select("COUNT(1)").From("event_participants").ToSql()
	if err != nil {
		return response, err
	}

	err = r.pg.Pool.QueryRow(ctx, countQuery, args...).Scan(&response.Count)
	if err != nil {
		return response, err
	}

	return response, nil
}
