package handler

import (
	"aspro/internal/apperror"
	"aspro/internal/domain/service"
	"aspro/pkg/client/postgresql/model/filter"
	"aspro/pkg/client/postgresql/model/sort"
	"net/http"
	"strings"

	"github.com/spf13/viper"
)

const (
	URL = "/api/heartbeat"
)

type Handler struct {
	services *service.Service
}

func NewHandler(services *service.Service) *Handler {
	return &Handler{services: services}
}

// type of handle http request
type HandlerFunc interface {
	HandlerFunc(method, path string, handler http.HandlerFunc)
}

// add the routes for metricHandler
func (h *Handler) Register(router HandlerFunc) {
	router.HandlerFunc(http.MethodGet, URL, filter.Middleware(sort.Sorting(apperror.Middleware(h.Heartbeat), "created_at", sort.DESC), viper.GetInt("default_limit")))
}

func (h *Handler) Heartbeat(w http.ResponseWriter, r *http.Request) error {

	filterOptions := r.Context().Value(sort.OptionsContextKey).(filter.Options)

	name := r.URL.Query().Get("name")
	if name != "" {
		filterOptions.AddField("name", filter.OperatorLike, name, filter.DataTypeStr)
	}

	age := r.URL.Query().Get("age")
	if age != "" {
		operator := filter.OperatorEq
		value := age
		if strings.Index(age, ":") != -1 {
			split := strings.Split(age, ":")
			operator = split[0]
			value = split[1]
		}
		filterOptions.AddField("age", operator, value, filter.DataTypeInt)
	}

	createdAt := r.URL.Query().Get("created_at")
	if createdAt != "" {
		var operator string
		if strings.Index(createdAt, ":") != -1 {
			operator = filter.OperatorBetween
		} else {
			operator = filter.OperatorEq
		}
		filterOptions.AddField("created_at", operator, createdAt, filter.DataTypeInt)
	}
	return nil
}
