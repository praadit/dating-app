package test

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/brianvoe/gofakeit"
	"github.com/praadit/dating-apps/models"
	"github.com/praadit/dating-apps/request"
	"github.com/praadit/dating-apps/utils"
	"github.com/stretchr/testify/assert"
	"github.com/uptrace/bun"
)

func generateFemaleUser() *models.User {
	pic := "picture"
	return &models.User{
		Active:   true,
		Name:     gofakeit.Name(),
		Email:    fmt.Sprintf("mail+%d@mail.com", time.Now().UnixNano()),
		Gender:   "f",
		Picture:  &pic,
		Password: "123",
		Benefits: models.UserBenefits{
			BaseSwipe: 10,
			IsPremium: false,
		},
	}
}
func generateMaleUser() *models.User {
	pic := "picture"
	return &models.User{
		Active:   true,
		Name:     gofakeit.Name(),
		Email:    fmt.Sprintf("mail+%d@mail.com", time.Now().UnixNano()),
		Gender:   "m",
		Picture:  &pic,
		Password: "123",
		Benefits: models.UserBenefits{
			BaseSwipe: 10,
			IsPremium: false,
		},
	}
}
func generateListOfUsers(count int, gender string) ([]int, map[int]*models.User) {
	ctx := context.Background()
	availUserId := []int{}
	mapAvailUser := map[int]*models.User{}

	for i := 0; i < count; i++ {
		var newUser *models.User
		if gender == "m" {
			newUser = generateMaleUser()
		} else if gender == "f" {
			newUser = generateFemaleUser()
		}

		BunDB.NewInsert().Model(newUser).Exec(ctx)
		if newUser.Id != 0 {
			availUserId = append(availUserId, newUser.Id)
			mapAvailUser[newUser.Id] = newUser
		}
	}

	return availUserId, mapAvailUser
}

func TestExplore(t *testing.T) {
	ctx := context.Background()

	newMale := generateMaleUser()
	BunDB.NewInsert().Model(newMale).Exec(ctx)
	assert.NotEqual(t, 0, newMale)
	defer func() {
		BunDB.NewDelete().Model(newMale).WherePK().Exec(ctx)
	}()

	availFemaleId, mapAvailFemale := generateListOfUsers(4, "f")
	defer func() {
		BunDB.NewDelete().Model((*models.User)(nil)).Where("id in (?)", bun.In(availFemaleId)).Exec(ctx)
	}()

	selUser, err := Service.Explore(ctx, newMale)
	assert.Nil(t, err)
	assert.NotEqual(t, newMale.Id, selUser.Id)

	// check female data
	assert.Equal(t, mapAvailFemale[selUser.Id].Id, selUser.Id)
	assert.Equal(t, mapAvailFemale[selUser.Id].Name, selUser.Name)
	assert.Equal(t, mapAvailFemale[selUser.Id].Email, selUser.Email)
	assert.Equal(t, mapAvailFemale[selUser.Id].Gender, selUser.Gender)
	assert.Equal(t, mapAvailFemale[selUser.Id].Picture, selUser.Picture)
	assert.Equal(t, mapAvailFemale[selUser.Id].Active, selUser.Active)

	assert.True(t, utils.Contains(availFemaleId, selUser.Id))
}
func TestExplore_NoProfile(t *testing.T) {
	ctx := context.Background()

	newMale := generateMaleUser()
	BunDB.NewInsert().Model(newMale).Exec(ctx)
	assert.NotEqual(t, 0, newMale)
	defer func() {
		BunDB.NewDelete().Model(newMale).WherePK().Exec(ctx)
	}()

	selUser, err := Service.Explore(ctx, newMale)
	assert.NotNil(t, err)
	assert.Equal(t, "No profile can be shown", err.Error())
	assert.Nil(t, selUser)
}

func TestSwipe(t *testing.T) {
	ctx := context.Background()

	newMale := generateMaleUser()
	BunDB.NewInsert().Model(newMale).Exec(ctx)
	assert.NotEqual(t, 0, newMale)
	defer func() {
		BunDB.NewDelete().Model(newMale).WherePK().Exec(ctx)
	}()

	availFemaleId, _ := generateListOfUsers(1, "f")
	defer func() {
		BunDB.NewDelete().Model((*models.User)(nil)).Where("id in (?)", bun.In(availFemaleId)).Exec(ctx)
	}()

	liked := false
	err := Service.Swipe(ctx, newMale, &request.SwipeRequest{
		UserId: availFemaleId[0],
		Liked:  &liked,
	})
	assert.Nil(t, err)

	match := &models.Match{}
	BunDB.NewSelect().Model(match).Where("user_id = ?", newMale.Id).Where("user_match_id = ?", availFemaleId[0]).Scan(ctx)
	defer func() {
		BunDB.NewDelete().Model(match).WherePK().Exec(ctx)
	}()
	assert.NotEqual(t, 0, match.Id)
	assert.Equal(t, liked, match.Liked)
}
func TestSwipe_SwipeMySelf(t *testing.T) {
	ctx := context.Background()

	newMale := generateMaleUser()
	BunDB.NewInsert().Model(newMale).Exec(ctx)
	assert.NotEqual(t, 0, newMale)
	defer func() {
		BunDB.NewDelete().Model(newMale).WherePK().Exec(ctx)
	}()

	liked := false
	err := Service.Swipe(ctx, newMale, &request.SwipeRequest{
		UserId: newMale.Id,
		Liked:  &liked,
	})
	assert.NotNil(t, err)
	assert.Equal(t, "You cannot swipe your self", err.Error())
}
func TestSwipe_InactiveUser(t *testing.T) {
	ctx := context.Background()

	newMale := generateMaleUser()
	BunDB.NewInsert().Model(newMale).Exec(ctx)
	assert.NotEqual(t, 0, newMale)
	defer func() {
		BunDB.NewDelete().Model(newMale).WherePK().Exec(ctx)
	}()

	newFemale := generateFemaleUser()
	newFemale.Active = false
	BunDB.NewInsert().Model(newFemale).Exec(ctx)
	assert.NotEqual(t, 0, newFemale)
	defer func() {
		BunDB.NewDelete().Model(newFemale).WherePK().Exec(ctx)
	}()

	liked := false
	err := Service.Swipe(ctx, newMale, &request.SwipeRequest{
		UserId: newFemale.Id,
		Liked:  &liked,
	})
	assert.NotNil(t, err)
	assert.Equal(t, "User inactive", err.Error())
}
func TestSwipe_OnceADay(t *testing.T) {
	ctx := context.Background()

	newMale := generateMaleUser()
	BunDB.NewInsert().Model(newMale).Exec(ctx)
	assert.NotEqual(t, 0, newMale)
	defer func() {
		BunDB.NewDelete().Model(newMale).WherePK().Exec(ctx)
	}()

	availFemaleId, _ := generateListOfUsers(1, "f")
	defer func() {
		BunDB.NewDelete().Model((*models.User)(nil)).Where("id in (?)", bun.In(availFemaleId)).Exec(ctx)
	}()

	match := &models.Match{
		UserId:      newMale.Id,
		UserMatchID: availFemaleId[0],
		Liked:       false,
	}
	BunDB.NewInsert().Model(match).Exec(ctx)
	defer func() {
		BunDB.NewDelete().Model(match).WherePK().Exec(ctx)
	}()

	liked := false
	err := Service.Swipe(ctx, newMale, &request.SwipeRequest{
		UserId: availFemaleId[0],
		Liked:  &liked,
	})
	assert.NotNil(t, err)
	assert.Equal(t, "You can only swipe same profile once a day", err.Error())
}
func TestSwipe_SwipeLimit(t *testing.T) {
	ctx := context.Background()

	newMale := generateMaleUser()
	newMale.Benefits.BaseSwipe = 1
	BunDB.NewInsert().Model(newMale).Exec(ctx)
	assert.NotEqual(t, 0, newMale)
	defer func() {
		BunDB.NewDelete().Model(newMale).WherePK().Exec(ctx)
	}()

	availFemaleId, _ := generateListOfUsers(2, "f")
	defer func() {
		BunDB.NewDelete().Model((*models.User)(nil)).Where("id in (?)", bun.In(availFemaleId)).Exec(ctx)
	}()

	match := &models.Match{
		UserId:      newMale.Id,
		UserMatchID: availFemaleId[0],
		Liked:       false,
	}
	BunDB.NewInsert().Model(match).Exec(ctx)
	defer func() {
		BunDB.NewDelete().Model(match).WherePK().Exec(ctx)
	}()

	liked := false
	err := Service.Swipe(ctx, newMale, &request.SwipeRequest{
		UserId: availFemaleId[1],
		Liked:  &liked,
	})
	assert.NotNil(t, err)
	assert.Equal(t, "You've reached swipe limit", err.Error())
}
func TestSwipe_SwipeLimit_Premium(t *testing.T) {
	ctx := context.Background()

	newMale := generateMaleUser()
	newMale.Benefits.BaseSwipe = 1
	newMale.Benefits.IsPremium = true
	BunDB.NewInsert().Model(newMale).Exec(ctx)
	assert.NotEqual(t, 0, newMale)
	defer func() {
		BunDB.NewDelete().Model(newMale).WherePK().Exec(ctx)
	}()

	availFemaleId, _ := generateListOfUsers(2, "f")
	defer func() {
		BunDB.NewDelete().Model((*models.User)(nil)).Where("id in (?)", bun.In(availFemaleId)).Exec(ctx)
	}()

	match := &models.Match{
		UserId:      newMale.Id,
		UserMatchID: availFemaleId[0],
		Liked:       false,
	}
	BunDB.NewInsert().Model(match).Exec(ctx)
	defer func() {
		BunDB.NewDelete().Model(match).WherePK().Exec(ctx)
	}()

	liked := false
	err := Service.Swipe(ctx, newMale, &request.SwipeRequest{
		UserId: availFemaleId[1],
		Liked:  &liked,
	})
	assert.Nil(t, err)

	newMatch := &models.Match{}
	BunDB.NewSelect().Model(newMatch).Where("user_id = ?", newMale.Id).Where("user_match_id = ?", availFemaleId[1]).Scan(ctx)
	defer func() {
		BunDB.NewDelete().Model(newMatch).WherePK().Exec(ctx)
	}()
	assert.NotEqual(t, 0, newMatch.Id)
	assert.Equal(t, liked, newMatch.Liked)
}
