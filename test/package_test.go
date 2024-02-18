package test

import (
	"context"
	"testing"

	"github.com/brianvoe/gofakeit"
	"github.com/praadit/dating-apps/constant"
	"github.com/praadit/dating-apps/models"
	"github.com/praadit/dating-apps/requests"
	"github.com/stretchr/testify/assert"
)

func generatePackage() *models.Package {
	return &models.Package{
		Active:      true,
		Name:        gofakeit.Name(),
		Terms:       gofakeit.Word(),
		Description: gofakeit.Word(),
		ActiveDays:  0,
		Type:        constant.PackageTypeOnetime,
		Price:       1000,
		Benefits: map[string]interface{}{
			"is_premium": true,
		},
	}
}

func TestGetPackageByID(t *testing.T) {
	ctx := context.Background()

	newPack := generatePackage()
	BunDB.NewInsert().Model(newPack).Exec(ctx)
	defer func() {
		BunDB.NewDelete().Model(newPack).WherePK().Exec(ctx)
	}()

	assert.NotEqual(t, 0, newPack.Id)
	pack, err := Service.Package(ctx, newPack.Id)

	assert.Nil(t, err)
	assert.Equal(t, newPack.Id, pack.Id)
	assert.Equal(t, newPack.Name, pack.Name)
	assert.Equal(t, newPack.Active, pack.Active)
	assert.Equal(t, newPack.Terms, pack.Terms)
	assert.Equal(t, newPack.Description, pack.Description)
}

func TestGetPackagesList(t *testing.T) {
	ctx := context.Background()

	newPack := generatePackage()
	BunDB.NewInsert().Model(newPack).Exec(ctx)
	defer func() {
		BunDB.NewDelete().Model(newPack).WherePK().Exec(ctx)
	}()

	assert.NotEqual(t, 0, newPack.Id)

	pack, err := Service.Packages(ctx, &requests.Pagination{
		Page:    1,
		PerPage: 10,
		Order:   "asc",
	})

	assert.Nil(t, err)
	assert.True(t, len(pack.Packages) > 1)
}

func TestGetPackageByID_NotFound(t *testing.T) {
	ctx := context.Background()

	pack, err := Service.Package(ctx, -1)

	assert.Nil(t, pack)
	assert.NotNil(t, err)
	assert.Equal(t, "Package not found", err.Error())
}

func TestBuyPackage(t *testing.T) {
	ctx := context.Background()
	testProfile := generateMaleUser()

	// register user
	err := Service.SignupUser(ctx, &requests.SignupRequest{
		Name:            testProfile.Name,
		Email:           testProfile.Email,
		Password:        testProfile.Password,
		ConfirmPassword: testProfile.Password,
		Gender:          testProfile.Gender,
		Picture:         *testProfile.Picture,
	})
	defer func() {
		BunDB.NewDelete().Model((*models.User)(nil)).Where("email = ?", testProfile.Email).Exec(ctx)
	}()
	assert.Nil(t, err)
	newUser := &models.User{}
	BunDB.NewSelect().Model(newUser).Where("email = ?", testProfile.Email).Scan(ctx)
	assert.NotEqual(t, 0, newUser.Id)

	// register new
	newPack := generatePackage()
	BunDB.NewInsert().Model(newPack).Exec(ctx)
	defer func() {
		BunDB.NewDelete().Model(newPack).WherePK().Exec(ctx)
	}()
	assert.NotEqual(t, 0, newPack.Id)

	// buy package testcase
	err = Service.Buy(ctx, newUser, &requests.BuyPackage{
		PackageID: newPack.Id,
	})
	assert.Nil(t, err)

	userpack := &models.UserPackage{}
	BunDB.NewSelect().Model(userpack).Where("user_id = ? and package_id = ?", newUser.Id, newPack.Id).Scan(ctx)
	defer func() {
		BunDB.NewDelete().Model(userpack).WherePK().Exec(ctx)
	}()
	assert.NotEqual(t, 0, userpack.Id)

	newBenefitUser := &models.User{Id: newUser.Id}
	BunDB.NewSelect().Model(newBenefitUser).WherePK().Scan(ctx)
	assert.Equal(t, true, newBenefitUser.Benefits.IsPremium)
}
