// Package usecase implements application business logic. Each logic group in own file.
package usecase

import (
	"context"

	"github.com/golanguzb70/system-design-openbudget/internal/entity"
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
		UpdateField(ctx context.Context, req entity.UpdateFieldRequest) (entity.RowsEffected, error)
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
)
