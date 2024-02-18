package service

import (
	"context"
	"errors"
	"math/rand"

	"github.com/praadit/dating-apps/models"
	"github.com/praadit/dating-apps/request"
	"github.com/praadit/dating-apps/utils"
	"github.com/uptrace/bun"
)

func (s *Service) Explore(ctx context.Context, user *models.User) (*models.User, error) {
	history, err := s.getTodaySwipe(ctx, user.Id)
	if err != nil {
		return nil, err
	}

	swipedId := []int{}
	for _, his := range history {
		swipedId = append(swipedId, his.UserMatchID)
	}

	selectedUsers := []*models.User{}
	query := s.db.NewSelect().Model(&selectedUsers).Where("id != ?", user.Id).Where("gender != ? and active = true", user.Gender)
	if len(swipedId) > 0 {
		query = query.Where("id not in (?)", bun.In(swipedId))
	}
	if err := query.Scan(ctx); err != nil {
		return nil, utils.FilterError(err, "Failed to get potential profile")
	}

	if len(selectedUsers) < 1 {
		return nil, errors.New("No profile can be shown")
	} else if len(selectedUsers) == 1 {
		return selectedUsers[0], nil
	} else {
		idx := rand.Intn(len(selectedUsers)-0) + 0
		return selectedUsers[idx], nil
	}
}

func (s *Service) Swipe(ctx context.Context, user *models.User, req *request.SwipeRequest) error {
	if user.Id == req.UserId {
		return errors.New("You cannot swipe your self")
	}

	swipedUser, err := s.getUserById(ctx, req.UserId)
	if err != nil {
		return err
	}
	if !swipedUser.Active {
		return errors.New("User inactive")
	}

	if exist, err := s.db.NewSelect().Model((*models.Match)(nil)).Where("user_id = ? and user_match_id = ? and date(created_at) = current_date", user.Id, req.UserId).Exists(ctx); err != nil {
		return utils.FilterError(err, "Failed to validate request")
	} else {
		if exist {
			return errors.New("You can only swipe same profile once a day")
		}
	}

	history, err := s.getTodaySwipe(ctx, user.Id)
	if err != nil {
		return err
	}

	if !user.Benefits.IsPremium {
		if len(history) >= (user.Benefits.BaseSwipe) {
			return errors.New("You've reached swipe limit")
		}
	}

	newMatch := &models.Match{
		UserId:      user.Id,
		UserMatchID: req.UserId,
		Liked:       *req.Liked,
	}

	if _, err := s.db.NewInsert().Model(newMatch).Exec(ctx); err != nil {
		return utils.FilterError(err, "Fail to match with user")
	}

	return nil
}

func (s *Service) getTodaySwipe(ctx context.Context, userId int) ([]*models.Match, error) {
	history := []*models.Match{}
	if err := s.db.NewSelect().Model(&history).Where("user_id =  ?", userId).Where("date(created_at) = current_date").Scan(ctx); err != nil {
		return nil, utils.FilterError(err, "Something went wrong")
	}

	return history, nil
}
