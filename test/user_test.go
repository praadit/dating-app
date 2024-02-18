package test

import (
	"context"
	"testing"

	"github.com/praadit/dating-apps/models"
	"github.com/praadit/dating-apps/requests"
	"github.com/stretchr/testify/assert"
)

func TestLoginAndRegister(t *testing.T) {
	ctx := context.Background()
	testProfile := generateMaleUser()
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

	assert.Equal(t, testProfile.Name, newUser.Name)
	assert.Equal(t, testProfile.Email, newUser.Email)
	assert.Equal(t, testProfile.Gender, newUser.Gender)
	assert.Equal(t, testProfile.Picture, newUser.Picture)

	loginResponse, err := Service.Login(ctx, &requests.LoginRequest{
		Email:    testProfile.Email,
		Password: testProfile.Password,
	})

	assert.Nil(t, err)
	assert.NotNil(t, loginResponse.AccessToken)
	assert.NotNil(t, loginResponse.Expiry)
	assert.Equal(t, "Bearer", loginResponse.Type)
}

func TestRegister_DuplicateEmail(t *testing.T) {
	ctx := context.Background()
	testProfile := generateMaleUser()
	user := &models.User{
		Name:   testProfile.Name,
		Email:  testProfile.Email,
		Active: true,
	}

	BunDB.NewInsert().Model(user).Exec(ctx)
	assert.NotEqual(t, 0, user.Id)

	defer func() {
		BunDB.NewDelete().Model(user).WherePK().Exec(ctx)
	}()

	err := Service.SignupUser(ctx, &requests.SignupRequest{
		Name:            testProfile.Name,
		Email:           testProfile.Email,
		Password:        testProfile.Password,
		ConfirmPassword: testProfile.Password,
		Gender:          testProfile.Gender,
		Picture:         *testProfile.Picture,
	})

	assert.NotNil(t, err)
	assert.Equal(t, "Email already registered", err.Error())
}

func TestLogin_WrongPassword(t *testing.T) {
	ctx := context.Background()
	testProfile := generateMaleUser()
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

	assert.Equal(t, testProfile.Name, newUser.Name)
	assert.Equal(t, testProfile.Email, newUser.Email)
	assert.Equal(t, testProfile.Gender, newUser.Gender)
	assert.Equal(t, testProfile.Picture, newUser.Picture)

	loginResponse, err := Service.Login(ctx, &requests.LoginRequest{
		Email:    testProfile.Email,
		Password: "wrongPassword",
	})

	assert.NotNil(t, err)
	assert.Equal(t, "Email / Password doesnt match", err.Error())
	assert.Nil(t, loginResponse)
}

func TestLogin_WrongEmail(t *testing.T) {
	ctx := context.Background()
	testProfile := generateMaleUser()
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

	assert.Equal(t, testProfile.Name, newUser.Name)
	assert.Equal(t, testProfile.Email, newUser.Email)
	assert.Equal(t, testProfile.Gender, newUser.Gender)
	assert.Equal(t, testProfile.Picture, newUser.Picture)

	loginResponse, err := Service.Login(ctx, &requests.LoginRequest{
		Email:    "wrongEmail@mail.com",
		Password: testProfile.Password,
	})

	assert.NotNil(t, err)
	assert.Equal(t, "Email not found", err.Error())
	assert.Nil(t, loginResponse)
}
