package controllers

import (
	"github.com/go-playground/validator/v10"
	"github.com/patrickmn/go-cache"
	"github.com/praadit/dating-apps/services"
	"github.com/uptrace/bun"
)

type Controller struct {
	db       *bun.DB
	memCache *cache.Cache
	service  *services.Service
	validate *validator.Validate
}

func NewController(db *bun.DB, memCache *cache.Cache, validate *validator.Validate) *Controller {
	service := services.NewService(db, memCache)

	return &Controller{
		db:       db,
		memCache: memCache,
		service:  service,
		validate: validate,
	}
}
