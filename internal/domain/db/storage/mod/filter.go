package mod

import (
	"aspro/internal/domain/db/storage"
	"aspro/pkg/client/postgresql/model/filter"
)

type filterOptions struct {
	limit  int
	fields []filter.Field
}

func NewOptions(options filter.Options) storage.FilterOptions {
	return &filterOptions{limit: options.Limit(), fields: options.Fields()}
}
