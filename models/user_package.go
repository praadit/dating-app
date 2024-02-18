package models

import (
	"context"
	"time"

	"github.com/uptrace/bun"
)

type UserPackage struct {
	bun.BaseModel `bun:"table:user_packages"`
	Id            int        `bun:",pk,autoincrement"`
	UserId        int        `json:"user_id"`
	PackageId     int        `json:"package_id"`
	ExpiredAt     *time.Time `json:"expired_at"`
	CreatedAt     time.Time  `json:"created_at"`
	UpdatedAt     time.Time  `json:"updated_at"`
}

var _ bun.BeforeUpdateHook = (*UserPackage)(nil)

func (*UserPackage) BeforeUpdate(ctx context.Context, query *bun.UpdateQuery) error {
	query.Column("updated_at")
	return nil
}

var _ bun.BeforeAppendModelHook = (*UserPackage)(nil)

func (m *UserPackage) BeforeAppendModel(ctx context.Context, query bun.Query) error {
	switch query.(type) {
	case *bun.InsertQuery:
		m.CreatedAt = time.Now()
		m.UpdatedAt = time.Now()
	case *bun.UpdateQuery:
		m.UpdatedAt = time.Now()
	}
	return nil
}
