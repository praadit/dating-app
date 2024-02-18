package models

import (
	"context"
	"time"

	"github.com/uptrace/bun"
)

type Match struct {
	bun.BaseModel `bun:"table:matches"`
	Id            int       `bun:",pk,autoincrement"`
	UserId        int       `json:"user_id"`
	UserMatchID   int       `json:"user_match_id"`
	Liked         bool      `json:"liked"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
}

var _ bun.BeforeUpdateHook = (*Match)(nil)

func (*Match) BeforeUpdate(ctx context.Context, query *bun.UpdateQuery) error {
	query.Column("updated_at")
	return nil
}

var _ bun.BeforeAppendModelHook = (*Match)(nil)

func (m *Match) BeforeAppendModel(ctx context.Context, query bun.Query) error {
	switch query.(type) {
	case *bun.InsertQuery:
		m.CreatedAt = time.Now()
		m.UpdatedAt = time.Now()
	case *bun.UpdateQuery:
		m.UpdatedAt = time.Now()
	}
	return nil
}
