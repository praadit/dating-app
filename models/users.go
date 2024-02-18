package models

import (
	"context"
	"time"

	"github.com/uptrace/bun"
)

type User struct {
	bun.BaseModel `bun:"table:users"`
	Id            int          `bun:",pk,autoincrement"`
	Active        bool         `json:"active"`
	Email         string       `json:"email"`
	Password      string       `json:"password"`
	Name          string       `json:"name"`
	Gender        string       `json:"gender"`
	Picture       *string      `json:"picture"`
	Benefits      UserBenefits `json:"benefits" bun:"type:jsonb"`
	CreatedAt     time.Time    `json:"created_at"`
	UpdatedAt     time.Time    `json:"updated_at"`
}

type UserBenefits struct {
	BaseSwipe int  `json:"base_swipe"`
	IsPremium bool `json:"is_premium"`
}

func (en *User) ToResponse() UserResponse {
	gender := "-"
	if en.Gender == "m" {
		gender = "Male"
	} else if en.Gender == "f" {
		gender = "Female"
	}
	return UserResponse{
		Id:         en.Id,
		Active:     en.Active,
		Name:       en.Name,
		Gender:     gender,
		Picture:    en.Picture,
		IsVerified: en.Benefits.IsPremium,
	}
}

type UserResponse struct {
	Id         int     `json:"id"`
	Active     bool    `json:"active"`
	Name       string  `json:"name"`
	Gender     string  `json:"gender"`
	Picture    *string `json:"picture"`
	IsVerified bool    `json:"is_verified"`
}

var _ bun.BeforeUpdateHook = (*User)(nil)

func (*User) BeforeUpdate(ctx context.Context, query *bun.UpdateQuery) error {
	query.Column("updated_at")
	return nil
}

var _ bun.BeforeAppendModelHook = (*User)(nil)

func (m *User) BeforeAppendModel(ctx context.Context, query bun.Query) error {
	switch query.(type) {
	case *bun.InsertQuery:
		m.CreatedAt = time.Now()
		m.UpdatedAt = time.Now()
	case *bun.UpdateQuery:
		m.UpdatedAt = time.Now()
	}
	return nil
}
