package services

import (
	"github.com/patrickmn/go-cache"
	"github.com/uptrace/bun"
)

type Service struct {
	db       *bun.DB
	memCache *cache.Cache
}

func NewService(db *bun.DB, memCache *cache.Cache) *Service {
	return &Service{
		db:       db,
		memCache: memCache,
	}
}
