package handler

import (
	"aspro/internal/domain/db/storage/mod"
	"aspro/pkg/client/postgresql/model"
	"aspro/pkg/client/postgresql/model/filter"
	"aspro/pkg/client/postgresql/model/sort"
	"context"
	"fmt"
)

func (h Handler) All(ctx context.Context, filterOptions filter.Options, sortOptions sort.Options) ([]model.Person, error) {
	options := mod.NewSortOptions(sortOptions.Field, sortOptions.Order)
	var userId int
	all, err := h.services.RecieveAll(ctx, userId, options)
	if err != nil {
		return nil, fmt.Errorf("failed to get all products due to error: %v", err)
	}
	return all, err
}
