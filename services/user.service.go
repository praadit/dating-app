package services

import (
	"context"
	"errors"
	"fmt"
	"strconv"

	"github.com/golang-jwt/jwt"
	"github.com/patrickmn/go-cache"
	"github.com/praadit/dating-apps/constants"
	"github.com/praadit/dating-apps/models"
	"github.com/praadit/dating-apps/requests"
	"github.com/praadit/dating-apps/response"
	"github.com/praadit/dating-apps/utils"
	"github.com/uptrace/bun"
)

func (s *Service) AuthUser(ctx context.Context) (*models.User, error) {
	claim := ctx.Value(constants.CTX_CLAIM)
	if claim == nil {
		return nil, utils.ERR_INVALID_TOKEN()
	}

	claims := &jwt.StandardClaims{}
	if val, ok := claim.(*jwt.StandardClaims); ok {
		claims = val
	} else {
		return nil, utils.ERR_INVALID_TOKEN()
	}

	id, err := strconv.Atoi(claims.Subject)
	if err != nil {
		return nil, utils.FilterError(err, "Id invalid")
	}
	user, err := s.getUserById(ctx, id)
	if err != nil {
		return nil, err
	}

	return user, nil
}

func (s *Service) getUserById(ctx context.Context, id int) (*models.User, error) {
	cacheKey := fmt.Sprintf("user:%d", id)
	user := &models.User{
		Id: id,
	}
	if val, ok := s.memCache.Get(cacheKey); !ok {
		query := s.db.NewSelect().Model(user).WherePK()
		if _, err := utils.SqlPanicFilter(query.Scan(ctx), "Failed to get user", "User not found"); err != nil {
			return nil, err
		}
		s.memCache.Set(cacheKey, user, cache.DefaultExpiration)
	} else {
		user = val.(*models.User)
	}

	return user, nil
}

func (s *Service) Login(ctx context.Context, req *requests.LoginRequest) (*response.LoginResponse, error) {
	user := &models.User{}
	query := s.db.NewSelect().Model(user).Where("email = ?", req.Email)
	if _, err := utils.SqlPanicFilter(query.Scan(ctx), "Failed to get user", "Email not found"); err != nil {
		return nil, err
	}

	if ok := utils.CheckPasswordHash(req.Password, user.Password); !ok {
		return nil, errors.New("Email / Password doesnt match")
	}
	token, expire, err := utils.GenerateLoginJwt(user.Id)
	if err != nil {
		return nil, utils.FilterError(err, "Failed to generate jwt")
	}

	return &response.LoginResponse{
		AccessToken: token,
		Expiry:      expire,
		Type:        "Bearer",
	}, nil
}

func (s *Service) SignupUser(ctx context.Context, req *requests.SignupRequest) error {
	if exist, err := s.db.NewSelect().Model((*models.User)(nil)).Where("email = ?", req.Email).Exists(ctx); err != nil {
		return utils.FilterError(err, "Fail to validate email user")
	} else {
		if exist {
			return errors.New("Email already registered")
		}
	}

	hashed, err := utils.HashPassword(req.Password)
	if err != nil {
		return utils.ERR_UNKNOWN(err)
	}

	if req.Gender != constants.GenderMale && req.Gender != constants.GenderFemale {
		return errors.New("Gender invalid. Please use 'm' for male and 'f' for female")
	}

	newUser := &models.User{
		Active:   true,
		Email:    req.Email,
		Password: hashed,
		Name:     req.Name,
		Gender:   req.Gender,
		Picture:  &req.Picture,
		Benefits: models.UserBenefits{
			BaseSwipe: 10,
		},
	}

	errTx := s.db.RunInTx(ctx, nil, func(ctx context.Context, tx bun.Tx) error {
		if _, err := tx.NewInsert().Model(newUser).Exec(ctx); err != nil {
			return err
		}
		return nil
	})
	if errTx != nil {
		return utils.FilterError(errTx, "Failed to register user")
	}
	return nil
}
