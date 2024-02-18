package controller

import (
	"github.com/go-playground/validator/v10"
	"github.com/patrickmn/go-cache"
	"github.com/praadit/dating-apps/service"
	"github.com/uptrace/bun"
)

type Controller struct {
	db       *bun.DB
	memCache *cache.Cache
	service  *service.Service
	validate *validator.Validate
}

func NewController(db *bun.DB, memCache *cache.Cache, validate *validator.Validate) *Controller {
	service := service.NewService(db, memCache)

	return &Controller{
		db:       db,
		memCache: memCache,
		service:  service,
		validate: validate,
	}
}
