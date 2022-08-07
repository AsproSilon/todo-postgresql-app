package service

import (
	"aspro/internal/domain/repository"
	"aspro/pkg/client/postgresql/model"
	"context"
)

type DeviceService struct {
	repo       repository.Device
	personRepo repository.Person
}

func NewDeviceService(repo repository.Device, personRepo repository.Person) *DeviceService {
	return &DeviceService{repo: repo, personRepo: personRepo}
}

func (s *DeviceService) Create(ctx context.Context, deviceId, userId, personId int, device model.Device) (int, error) {
	_, err := s.personRepo.RecieveById(ctx, userId)
	if err != nil {
		return 0, err
	}
	return s.repo.Create(ctx, userId, personId, deviceId, device)
}

func (s *DeviceService) FindAll(ctx context.Context, userId, personId int) (device []model.Device, err error) {
	return s.repo.FindAll(ctx, userId, personId)
}

func (s *DeviceService) FindById(ctx context.Context, userId, deviceId int) (model.Device, error) {
	return s.repo.FindById(ctx, userId, deviceId)
}

func (s *DeviceService) Update(ctx context.Context, userId, deviceId int) error {
	return s.repo.Update(ctx, userId, deviceId)
}

func (s *DeviceService) Delete(ctx context.Context, userId, deviceId int) error {
	return s.repo.Delete(ctx, userId, deviceId)
}
