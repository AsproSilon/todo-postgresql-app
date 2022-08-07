package service

import (
	"aspro/internal/domain/db/storage"
	"aspro/internal/domain/repository"
	"aspro/pkg/client/postgresql/model"
	"context"
)

type PersonService struct {
	repo repository.Person
}

func NewPersonService(repo repository.Person) *PersonService {
	return &PersonService{repo: repo}
}

func (s *PersonService) Create(ctx context.Context, person *model.Person) error {
	return s.repo.Create(ctx, person)
}

func (s *PersonService) RecieveAll(ctx context.Context, userId int, sortOptions storage.SortOptions) (p []model.Person, err error) {
	return s.repo.RecieveAll(ctx, userId, sortOptions)
}

func (s *PersonService) RecieveById(ctx context.Context, id int) (model.Person, error) {
	return s.repo.RecieveById(ctx, id)
}

func (s *PersonService) Update(ctx context.Context, person model.Person) error {
	return s.repo.Update(ctx, person)
}

func (s *PersonService) Delete(ctx context.Context, userId int, personId int) error {
	return s.repo.Delete(ctx, userId, personId)
}
