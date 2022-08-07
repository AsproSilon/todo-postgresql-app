package service

import (
	"aspro/internal/domain/db/storage"
	"aspro/internal/domain/repository"
	"aspro/pkg/client/postgresql/model"
	"context"
)

type Authorization interface {
	CreateUser(user model.User) (int, error)
	GenerateToken(username, password string) (string, error)
	ParseToken(token string) (int, error)
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

type Service struct {
	Authorization
	Person
	Device
}

func NewService(repos *repository.Repository) *Service {
	return &Service{
		Authorization: NewAuthorizationService(repos.Authorization),
		Person:        NewPersonService(repos.Person),
		Device:        NewDeviceService(repos.Device, repos.Person),
	}
}
