package models

import (
	"context"
	"time"

	"github.com/uptrace/bun"
)

type Package struct {
	bun.BaseModel `bun:"table:packages"`
	Id            int                    `bun:",pk,autoincrement"`
	Active        bool                   `json:"active"`
	Name          string                 `json:"name"`
	Terms         string                 `json:"terms"`
	Description   string                 `json:"description"`
	ActiveDays    int                    `json:"active_days"`
	Type          string                 `json:"type"`
	Price         float64                `json:"price"`
	Benefits      map[string]interface{} `json:"benefits" bun:"type:jsonb"`
	CreatedAt     time.Time              `json:"created_at"`
	UpdatedAt     time.Time              `json:"updated_at"`
}

func (en *Package) ToResponse() PackageResponse {
	return PackageResponse{
		Id:          en.Id,
		Active:      en.Active,
		Name:        en.Name,
		Term:        en.Terms,
		Description: en.Description,
		CreatedAt:   en.CreatedAt,
		UpdatedAt:   en.UpdatedAt,
	}
}

type PackageResponse struct {
	Id          int       `json:"id"`
	Active      bool      `json:"active"`
	Name        string    `json:"name"`
	Term        string    `json:"terms"`
	Description string    `json:"description"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

var _ bun.BeforeUpdateHook = (*Package)(nil)

func (*Package) BeforeUpdate(ctx context.Context, query *bun.UpdateQuery) error {
	query.Column("updated_at")
	return nil
}

var _ bun.BeforeAppendModelHook = (*Package)(nil)

func (m *Package) BeforeAppendModel(ctx context.Context, query bun.Query) error {
	switch query.(type) {
	case *bun.InsertQuery:
		m.CreatedAt = time.Now()
		m.UpdatedAt = time.Now()
	case *bun.UpdateQuery:
		m.UpdatedAt = time.Now()
	}
	return nil
}
