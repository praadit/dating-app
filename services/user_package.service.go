package services

import (
	"context"

	"github.com/praadit/dating-apps/models"
	"github.com/praadit/dating-apps/utils"
)

func (s *Service) UserPacakge(ctx context.Context, userId int, packageId int) (*models.UserPackage, error) {
	userPackage := &models.UserPackage{}
	query := s.db.NewSelect().Model(userPackage).Where("user_id = ?", userId).Where("package_id = ?", packageId)
	if _, err := utils.SqlPanicFilter(query.Scan(ctx), "Failed to get package", "Package not found"); err != nil {
		return nil, err
	}

	return userPackage, nil
}
