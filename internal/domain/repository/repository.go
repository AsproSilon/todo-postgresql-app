package repository

import (
	"aspro/internal/domain/db/storage"
	"aspro/pkg/client/postgresql/model"
	"context"

	"github.com/jmoiron/sqlx"
)

type Authorization interface {
	CreateUser(user model.User) (int, error)
	GetUser(username, password string) (model.User, error)
}

type Person interface {
	Create(ctx context.Context, person *model.Person) error
	RecieveAll(ctx context.Context, userId int, sortOptions storage.SortOptions) (p []model.Person, err error)
	RecieveById(ctx context.Context, id int) (model.Person, error)
	Update(ctx context.Context, person model.Person) error
	Delete(ctx context.Context, userId int, personId int) error
}

type Device interface {
	Create(ctx context.Context, deviceId, userId, personId int, device model.Device) (int, error)
	FindAll(ctx context.Context, userId, personId int) (device []model.Device, err error)
	FindById(ctx context.Context, userId, deviceId int) (model.Device, error)
	Update(ctx context.Context, userId, deviceId int) error
	Delete(ctx context.Context, userId, deviceId int) error
}

type Repository struct {
	Authorization
	Person
	Device
}

func NewRepository(db *sqlx.DB) *Repository {
	return &Repository{
		Authorization: NewAuthorizationPostgres(db),
		Person:        NewPersonPostgres(db),
		Device:        NewDevicePostgres(db),
	}
}
