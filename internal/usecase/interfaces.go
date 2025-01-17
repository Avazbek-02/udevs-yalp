// Package usecase implements application business logic. Each logic group in own file.
package usecase

import (
	"context"

	"github.com/Avazbek-02/udevslab-lesson6/internal/entity"
)

//go:generate mockgen -source=interfaces.go -destination=./mocks_test.go -package=usecase_test

type (
	// UserRepo -.
	UserRepoI interface {
		Create(ctx context.Context, req entity.User) (entity.User, error)
		GetSingle(ctx context.Context, req entity.UserSingleRequest) (entity.User, error)
		GetList(ctx context.Context, req entity.GetListFilter) (entity.UserList, error)
		Update(ctx context.Context, req entity.User) (entity.User, error)
		Delete(ctx context.Context, req entity.Id) error
	}

	// SessionRepo -.
	SessionRepoI interface {
		Create(ctx context.Context, req entity.Session) (entity.Session, error)
		GetSingle(ctx context.Context, req entity.Id) (entity.Session, error)
		GetList(ctx context.Context, req entity.GetListFilter) (entity.SessionList, error)
		Update(ctx context.Context, req entity.Session) (entity.Session, error)
		Delete(ctx context.Context, req entity.Id) error
		UpdateField(ctx context.Context, req entity.UpdateFieldRequest) (entity.RowsEffected, error)
	}
	BusinessRepoI interface {
		Create(ctx context.Context, req entity.Business) (entity.Business, error)
		GetSingle(ctx context.Context, req entity.Id) (entity.Business, error)
		GetList(ctx context.Context, req entity.GetListFilter) (entity.BusinessList, error)
		Update(ctx context.Context, req entity.Business) (entity.Business, error)
		Delete(ctx context.Context, req entity.Id) error
	}
	NotificationRepoI interface {
		Create(ctx context.Context, req entity.Notification) (entity.Notification, error)
		GetSingle(ctx context.Context, req entity.Id) (entity.Notification, error)
		GetList(ctx context.Context, req entity.GetListFilter) (entity.NotificationList, error)
		Update(ctx context.Context, req entity.Notification) (entity.Notification, error)
		Delete(ctx context.Context, req entity.Id) error
		UpdateStatus(ctx context.Context, req entity.Notification) (entity.Notification, error)
	}
	ReviewRepoI interface {
		Create(ctx context.Context, req entity.Review) (entity.Review, error)
		GetSingle(ctx context.Context, req entity.Id) (entity.Review, error)
		GetList(ctx context.Context, req entity.GetListFilter) (entity.ReviewList, error)
		Update(ctx context.Context, req entity.Review) (entity.Review, error)
		Delete(ctx context.Context, req entity.Id) error
	}
	ReportRepoI interface {
		Create(ctx context.Context, req entity.Report) (entity.Report, error)
		GetSingle(ctx context.Context, req entity.Id) (entity.Report, error)
		GetList(ctx context.Context, req entity.GetListFilter) (entity.ReportList, error)
		Update(ctx context.Context, req entity.Report) (entity.Report, error)
		Delete(ctx context.Context, req entity.Id) error
	}
	EventRepoI interface {
		Create(ctx context.Context, req entity.Event) (entity.Event, error)
		GetSingle(ctx context.Context, req entity.Id) (entity.Event, error)
		GetList(ctx context.Context, req entity.GetListFilter) (entity.EventList, error)
		Update(ctx context.Context, req entity.Event) (entity.Event, error)
		Delete(ctx context.Context, req entity.Id) error
		AddParticipant(ctx context.Context, req entity.EventParticipant) (entity.EventParticipant, error)
		RemoveParticipant(ctx context.Context, req entity.EventParticipant) error
		GetParticipants(ctx context.Context, req entity.GetListFilter) (entity.EventParticipantList, error)
	}
)
