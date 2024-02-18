package service

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/patrickmn/go-cache"
	"github.com/praadit/dating-apps/constant"
	"github.com/praadit/dating-apps/models"
	modelhelper "github.com/praadit/dating-apps/models/model_helper"
	"github.com/praadit/dating-apps/requests"
	"github.com/praadit/dating-apps/response"
	"github.com/praadit/dating-apps/utils"
	"github.com/uptrace/bun"
)

func (s *Service) Packages(ctx context.Context, pagination *requests.Pagination) (*response.PackagesResponse, error) {
	pacakges := []*models.Package{}
	query := s.db.NewSelect().Model(&pacakges).Where("active = true")

	// set filter here if exist

	count, err := query.Count(ctx)
	if err != nil {
		return nil, utils.FilterError(err, "Failed to count items")
	}

	if pagination.Order != "" {
		if pagination.Order != constant.OrderAsc && pagination.Order != constant.OrderDesc {
			pagination.Order = constant.OrderAsc
		}
	}
	query = query.Order(fmt.Sprintf("name %s", pagination.Order)).Offset(pagination.PerPage * (pagination.Page - 1)).Limit(pagination.PerPage)
	if _, err := utils.SqlPanicFilter(query.Scan(ctx), "Failed to get package", "Package not found"); err != nil {
		return nil, err
	}

	packageResponses := []models.PackageResponse{}
	for _, pack := range pacakges {
		packageResponses = append(packageResponses, pack.ToResponse())
	}

	return &response.PackagesResponse{
		Pagination: response.PaginationResponse{
			Total:   count,
			Page:    pagination.Page,
			PerPage: pagination.PerPage,
		},
		Packages: packageResponses,
	}, nil
}

func (s *Service) Package(ctx context.Context, id int) (*models.PackageResponse, error) {
	pack, err := s.getPackageById(ctx, id)
	if err != nil {
		return nil, err
	}

	if !pack.Active {
		return nil, errors.New("Package not active")
	}

	res := pack.ToResponse()
	return &res, nil
}

func (s *Service) getPackageById(ctx context.Context, id int) (*models.Package, error) {
	cacheKey := fmt.Sprintf("package:%d", id)
	pacakge := &models.Package{
		Id: id,
	}
	if val, ok := s.memCache.Get(cacheKey); !ok {
		query := s.db.NewSelect().Model(pacakge).WherePK()
		if _, err := utils.SqlPanicFilter(query.Scan(ctx), "Failed to get package", "Package not found"); err != nil {
			return nil, err
		}
		s.memCache.Set(cacheKey, pacakge, cache.DefaultExpiration)
	} else {
		pacakge = val.(*models.Package)
	}

	return pacakge, nil
}

func (s *Service) Buy(ctx context.Context, user *models.User, req *requests.BuyPackage) error {
	pack, err := s.getPackageById(ctx, req.PackageID)
	if err != nil {
		return err
	}

	var expire *time.Time = nil
	now := time.Now()
	if pack.ActiveDays > 0 {
		now = now.AddDate(0, 0, pack.ActiveDays)
		if pack.ActiveDays == 1 {
			now = time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, now.Nanosecond(), now.Location())
		}
	}

	userBens, err := utils.StructToMap(user.Benefits)
	if err != nil {
		return utils.FilterError(err, "Failed to buy package")
	}

	packBenKey := utils.GetMapKeys(pack.Benefits)
	if err != nil {
		return utils.FilterError(err, "Package has invalid benefits")
	}

	if pack.Type == constant.PackageTypeOnetime {
		if _, err := s.UserPacakge(ctx, user.Id, pack.Id); err == nil {
			return errors.New("Cannot buy the badge package twice")
		}

		for _, key := range packBenKey {
			userBens[key] = pack.Benefits[key]
		}
	} else if pack.Type == constant.PackageTypeAddOn {
		for _, key := range packBenKey {
			if val, ok := userBens[key]; ok {
				userBens[key] = int(val.(float64) + pack.Benefits[key].(float64))
			} else {
				userBens[key] = pack.Benefits[key].(int)
			}
		}
	}

	userpack := &models.UserPackage{
		UserId:    user.Id,
		PackageId: pack.Id,
		ExpiredAt: expire,
	}
	user.Benefits = modelhelper.CreateUserBens(userBens)

	errTx := s.db.RunInTx(ctx, nil, func(ctx context.Context, tx bun.Tx) error {
		if _, err := tx.NewInsert().Model(userpack).Exec(ctx); err != nil {
			return err
		}
		if _, err := tx.NewUpdate().Model(user).WherePK().Column("benefits").Exec(ctx); err != nil {
			return err
		}

		return nil
	})
	if errTx != nil {
		return utils.FilterError(errTx, "Failed to buy package")
	}

	return nil
}
