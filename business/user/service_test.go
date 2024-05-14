package user_test

import (
	"errors"
	"mime/multipart"
	businessUser "roby-backend-golang/business/user"
	"roby-backend-golang/config"
	repoUser "roby-backend-golang/repository/user"
	"roby-backend-golang/utils"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestLogin(t *testing.T) {
	asserting := assert.New(t)
	t.Run("Valid Test", func(t *testing.T) {
		resSample := businessUser.ResponseLogin{
			Email: "test@mail.com",
			Token: utils.Token{
				AccessToken:         "123",
				AccessTokenExpired:  123,
				RefreshToken:        "123",
				RefreshTokenExpired: 123,
			},
		}
		auth := businessUser.AuthLogin{
			Email:    "test@mail.com",
			Password: "12345678",
		}
		user := businessUser.User{
			ID:       "123",
			Email:    "test@mail.com",
			Password: "$2a$10$mfK4MlwOhHnvphtBNp0G0u/E6QjVHBk3ks0C.BnOMnRKI5Ue2J4SW",
			FullName: "test",
		}
		// mocking
		repoMock := &repoUser.UserMock{Mock: &mock.Mock{}}
		service := businessUser.NewService(repoMock, &config.AppConfig{})
		repoMock.On("FindUserByEmail", auth.Email).Return(user, nil)
		repoMock.On("GenerateTokenAuth", user.ID, user.Email).Return(&resSample.Token, nil)
		// repoMock.On("Login", mock.AnythingOfType("AuthLogin")).Return(resSample, nil)

		res, err := service.Login(auth)
		asserting.NoError(err)
		asserting.NotNil(res)
	})

	t.Run("Wrong Email Test", func(t *testing.T) {
		resSample := businessUser.ResponseLogin{
			Email: "test@mail.com",
			Token: utils.Token{
				AccessToken:         "123",
				AccessTokenExpired:  123,
				RefreshToken:        "123",
				RefreshTokenExpired: 123,
			},
		}
		auth := businessUser.AuthLogin{
			Email:    "test@mail.com",
			Password: "12345678",
		}
		user := businessUser.User{
			ID:       "123",
			Email:    "test@mail.com",
			Password: "$2a$10$mfK4MlwOhHnvphtBNp0G0u/E6QjVHBk3ks0C.BnOMnRKI5Ue2J4SW",
			FullName: "test",
		}
		// mocking
		repoMock := &repoUser.UserMock{Mock: &mock.Mock{}}
		service := businessUser.NewService(repoMock, &config.AppConfig{})
		repoMock.On("FindUserByEmail", auth.Email).Return(user, errors.New("wrong email"))
		repoMock.On("GenerateTokenAuth", user.ID, user.Email).Return(&resSample.Token, nil)
		// repoMock.On("Login", mock.AnythingOfType("AuthLogin")).Return(resSample, nil)

		_, err := service.Login(auth)
		asserting.Error(err)
	})

	t.Run("Validation Test", func(t *testing.T) {

		auth := businessUser.AuthLogin{
			Email:    "test@mail.com",
			Password: "",
		}
		user := businessUser.User{
			ID:       "123",
			Email:    "test@mail.com",
			Password: "$2a$10$mfK4MlwOhHnvphtBNp0G0u/E6QjVHBk3ks0C.BnOMnRKI5Ue2J4SW",
			FullName: "test",
		}
		// mocking
		repoMock := &repoUser.UserMock{Mock: &mock.Mock{}}
		service := businessUser.NewService(repoMock, &config.AppConfig{})
		repoMock.On("FindUserByEmail", auth.Email).Return(user, errors.New("wrong email"))
		repoMock.On("GenerateTokenAuth", user.ID, user.Email).Return(nil, nil)

		_, err := service.Login(auth)
		asserting.Error(err)
	})

	t.Run("Validation Test", func(t *testing.T) {

		auth := businessUser.AuthLogin{
			Email:    "test@mail.com",
			Password: "12345678",
		}
		user := businessUser.User{
			ID:       "123",
			Email:    "test@mail.com",
			Password: "$2a$10$/E6QjVHBk3ks0C.BnOMnRKI5Ue2J4SW",
			FullName: "test",
		}
		// mocking
		repoMock := &repoUser.UserMock{Mock: &mock.Mock{}}
		service := businessUser.NewService(repoMock, &config.AppConfig{})
		repoMock.On("FindUserByEmail", auth.Email).Return(user, nil)
		repoMock.On("GenerateTokenAuth", user.ID, user.Email).Return(nil, nil)

		_, err := service.Login(auth)
		asserting.Error(err)
	})

	t.Run("Generate Token Error Test", func(t *testing.T) {

		auth := businessUser.AuthLogin{
			Email:    "test@mail.com",
			Password: "12345678",
		}
		user := businessUser.User{
			ID:       "123",
			Email:    "test@mail.com",
			Password: "$2a$10$mfK4MlwOhHnvphtBNp0G0u/E6QjVHBk3ks0C.BnOMnRKI5Ue2J4SW",
			FullName: "test",
		}
		// mocking
		repoMock := &repoUser.UserMock{Mock: &mock.Mock{}}
		service := businessUser.NewService(repoMock, &config.AppConfig{})
		repoMock.On("FindUserByEmail", auth.Email).Return(user, nil)
		repoMock.On("GenerateTokenAuth", user.ID, user.Email).Return(&utils.Token{}, errors.New("error generate token"))

		_, err := service.Login(auth)
		asserting.Error(err)
	})

}

func TestRegisterUser(t *testing.T) {
	t.Run("Valid Test", func(t *testing.T) {
		asserting := assert.New(t)
		multipart := multipart.FileHeader{
			Filename: "test.jpeg",
			Size:     1,
		}
		inputUser := businessUser.Register{
			Email:    "test@mail.com",
			Password: "12345678",
			FullName: "test",
			File:     &multipart,
		}

		result := businessUser.User{
			ID:       "123",
			Email:    "test@mail.com",
			Password: "$2a$10$mfK4MlwOhHnvphtBNp0G0u/E6QjVHBk3ks0C.BnOMnRKI5Ue2J4SW",
		}

		repoMock := &repoUser.UserMock{Mock: &mock.Mock{}}
		service := businessUser.NewService(repoMock, &config.AppConfig{})
		repoMock.On("FindUserByEmail", inputUser.Email).Return(result, errors.New("email not found"))
		repoMock.On("UploadImageS3", mock.Anything).Return("url", nil)
		repoMock.On("CreateUser", mock.Anything).Return(nil)

		err := service.RegisterUser(inputUser)
		asserting.Nil(err)

	})
	t.Run("Wrong Validation Test", func(t *testing.T) {
		asserting := assert.New(t)
		multipart := multipart.FileHeader{
			Filename: "test.jpeg",
			Size:     1,
		}
		inputUser := businessUser.Register{
			Email:    "test@mail.com",
			Password: "",
			FullName: "test",
			File:     &multipart,
			PhotoUrl: "url",
		}

		repoMock := &repoUser.UserMock{Mock: &mock.Mock{}}
		service := businessUser.NewService(repoMock, &config.AppConfig{})
		repoMock.On("FindUserByEmail", inputUser.Email).Return(businessUser.User{}, errors.New("email already exist"))
		repoMock.On("UploadImageS3", &multipart).Return("url", nil)
		repoMock.On("CreateUser", mock.Anything).Return(nil)

		err := service.RegisterUser(inputUser)
		asserting.Error(err)

	})

	t.Run("Error Upload Image Test", func(t *testing.T) {
		asserting := assert.New(t)
		multipart := multipart.FileHeader{
			Filename: "test.jpeg",
			Size:     1,
		}
		inputUser := businessUser.Register{
			Email:    "test@mail.com",
			Password: "12345678",
			FullName: "test",
			File:     &multipart,
			PhotoUrl: "url",
		}

		repoMock := &repoUser.UserMock{Mock: &mock.Mock{}}
		service := businessUser.NewService(repoMock, &config.AppConfig{})
		repoMock.On("FindUserByEmail", inputUser.Email).Return(businessUser.User{}, errors.New("email already exist"))
		repoMock.On("UploadImageS3", &multipart).Return("", errors.New("error upload image"))
		repoMock.On("CreateUser", mock.Anything).Return(nil)

		err := service.RegisterUser(inputUser)
		asserting.Error(err)
	})

	t.Run("Error Create User Test", func(t *testing.T) {
		asserting := assert.New(t)
		multipart := multipart.FileHeader{
			Filename: "test.jpeg",
			Size:     1,
		}
		inputUser := businessUser.Register{
			Email:    "test@mail.com",
			Password: "12345678",
			FullName: "test",
			File:     &multipart,
		}

		result := businessUser.User{
			ID:       "123",
			Email:    "test@mail.com",
			Password: "$2a$10$mfK4MlwOhHnvphtBNp0G0u/E6QjVHBk3ks0C.BnOMnRKI5Ue2J4SW",
		}

		repoMock := &repoUser.UserMock{Mock: &mock.Mock{}}
		service := businessUser.NewService(repoMock, &config.AppConfig{})
		repoMock.On("FindUserByEmail", inputUser.Email).Return(result, errors.New("email not found"))
		repoMock.On("UploadImageS3", mock.Anything).Return("url", nil)
		repoMock.On("CreateUser", mock.Anything).Return(errors.New("error create user"))

		err := service.RegisterUser(inputUser)
		asserting.Error(err)
	})

	t.Run("Error Email Already Exist Test", func(t *testing.T) {
		asserting := assert.New(t)
		multipart := multipart.FileHeader{
			Filename: "test.jpeg",
			Size:     1,
		}
		inputUser := businessUser.Register{
			Email:    "test@mail.com",
			Password: "12345678",
			FullName: "test",
			File:     &multipart,
		}

		result := businessUser.User{
			ID:       "123",
			Email:    "test@mail.com",
			Password: "$2a$10$mfK4MlwOhHnvphtBNp0G0u/E6QjVHBk3ks0C.BnOMnRKI5Ue2J4SW",
		}

		repoMock := &repoUser.UserMock{Mock: &mock.Mock{}}
		service := businessUser.NewService(repoMock, &config.AppConfig{})
		repoMock.On("FindUserByEmail", inputUser.Email).Return(result, nil)
		repoMock.On("UploadImageS3", mock.Anything).Return("url", nil)
		repoMock.On("CreateUser", mock.Anything).Return(errors.New("error create user"))

		err := service.RegisterUser(inputUser)
		asserting.Error(err)
	})
}

func TestGetUserByID(t *testing.T) {
	t.Run("Valid Test", func(t *testing.T) {
		asserting := assert.New(t)
		user := businessUser.User{
			ID:       "123",
			Email:    "test@mail.com",
			Password: "$2a$10$mfK4MlwOhHnvphtBNp0G0u/E6QjVHBk3ks0C.BnOMnRKI5Ue2J4SW",
			FullName: "test",
		}
		repoMock := &repoUser.UserMock{Mock: &mock.Mock{}}
		service := businessUser.NewService(repoMock, &config.AppConfig{})
		repoMock.On("FindUserByID", user.ID).Return(user, nil)

		res, err := service.GetUserByID(user.ID)
		asserting.NoError(err)
		asserting.NotNil(res)
	})
}

func TestSwipeUser(t *testing.T) {
	t.Run("Valid Test", func(t *testing.T) {
		asserting := assert.New(t)
		user := businessUser.User{
			ID:    "123",
			Email: "test@mail.com",
		}
		swipe := businessUser.SwipeUser{
			IDSwipe: "1234",
			Swipe:   "like",
		}
		hour := 24 * time.Hour
		repoMock := &repoUser.UserMock{Mock: &mock.Mock{}}
		service := businessUser.NewService(repoMock, &config.AppConfig{})
		repoMock.On("GetMe", user.ID).Return(user, nil)
		repoMock.On("Get", "apptinder:allrandomid:123").Return("123", nil)
		repoMock.On("Set", "apptinder:allrandomid:123", "123,1234", hour).Return(nil)
		repoMock.On("GetRandomUser", mock.Anything).Return(businessUser.ResponseRandomUser{}, nil)
		repoMock.On("SwipeUser", user.ID, mock.Anything).Return(nil)

		err := service.SwipeUser(user.ID, swipe)
		asserting.NoError(err)
	})

	t.Run("Valid Test 2", func(t *testing.T) {
		asserting := assert.New(t)
		user := businessUser.User{
			ID:    "123",
			Email: "test@mail.com",
		}
		swipe := businessUser.SwipeUser{
			IDSwipe: "1234",
			Swipe:   "like",
		}
		hour := 24 * time.Hour
		repoMock := &repoUser.UserMock{Mock: &mock.Mock{}}
		service := businessUser.NewService(repoMock, &config.AppConfig{})
		repoMock.On("GetMe", user.ID).Return(user, nil)
		repoMock.On("Get", "apptinder:allrandomid:123").Return("123,2232", nil)
		repoMock.On("Set", "apptinder:allrandomid:123", "123,2232,1234", hour).Return(nil)
		repoMock.On("GetRandomUser", mock.Anything).Return(businessUser.ResponseRandomUser{}, nil)
		repoMock.On("SwipeUser", user.ID, mock.Anything).Return(nil)

		err := service.SwipeUser(user.ID, swipe)
		asserting.NoError(err)
	})

	t.Run("Test Error Set Redis", func(t *testing.T) {
		asserting := assert.New(t)
		user := businessUser.User{
			ID:    "123",
			Email: "test@mail.com",
		}
		swipe := businessUser.SwipeUser{
			IDSwipe: "1234",
			Swipe:   "like",
		}
		hour := 24 * time.Hour
		repoMock := &repoUser.UserMock{Mock: &mock.Mock{}}
		service := businessUser.NewService(repoMock, &config.AppConfig{})
		repoMock.On("GetMe", user.ID).Return(user, nil)
		repoMock.On("Get", "apptinder:allrandomid:123").Return("123,2232", nil)
		repoMock.On("Set", "apptinder:allrandomid:123", "123,2232,1234", hour).Return(errors.New("error set redis"))
		repoMock.On("GetRandomUser", mock.Anything).Return(businessUser.ResponseRandomUser{}, nil)
		repoMock.On("SwipeUser", user.ID, mock.Anything).Return(nil)

		err := service.SwipeUser(user.ID, swipe)
		asserting.Error(err)
	})

	t.Run("Limit 10 Swipe Test", func(t *testing.T) {
		asserting := assert.New(t)
		user := businessUser.User{
			ID:    "123",
			Email: "test@mail.com",
		}
		swipe := businessUser.SwipeUser{
			IDSwipe: "1234",
			Swipe:   "like",
		}
		hour := 24 * time.Hour
		repoMock := &repoUser.UserMock{Mock: &mock.Mock{}}
		service := businessUser.NewService(repoMock, &config.AppConfig{})
		repoMock.On("GetMe", user.ID).Return(user, nil)
		repoMock.On("Get", "apptinder:allrandomid:123").Return("123,221,21412,21412,214,214,124,124,214,214,214,21", nil)
		repoMock.On("Set", "apptinder:allrandomid:123", "123,1234", hour).Return(nil)
		repoMock.On("GetRandomUser", mock.Anything).Return(businessUser.ResponseRandomUser{}, nil)
		repoMock.On("SwipeUser", user.ID, mock.Anything).Return(nil)

		err := service.SwipeUser(user.ID, swipe)
		asserting.Error(err)
	})

	t.Run("PREMIUM PACKAGE Swipe Test", func(t *testing.T) {
		asserting := assert.New(t)
		packages := []businessUser.Package{
			{
				ID:          "123",
				PackageName: "premium",
				Description: "test",
			},
		}
		user := businessUser.User{
			ID:       "123",
			Email:    "test@mail.com",
			Packages: packages,
		}
		swipe := businessUser.SwipeUser{
			IDSwipe: "1234",
			Swipe:   "like",
		}
		hour := 24 * time.Hour
		repoMock := &repoUser.UserMock{Mock: &mock.Mock{}}
		service := businessUser.NewService(repoMock, &config.AppConfig{})
		repoMock.On("GetMe", user.ID).Return(user, nil)
		repoMock.On("Get", "apptinder:allrandomid:123").Return("123", nil)
		repoMock.On("Set", "apptinder:allrandomid:123", "123,1234", hour).Return(nil)
		repoMock.On("GetRandomUser", mock.Anything).Return(businessUser.ResponseRandomUser{}, nil)
		repoMock.On("SwipeUser", user.ID, mock.Anything).Return(nil)

		err := service.SwipeUser(user.ID, swipe)
		asserting.NoError(err)
	})

	t.Run("Id Swipe already Swipe Test", func(t *testing.T) {
		asserting := assert.New(t)

		swipe := businessUser.SwipeUser{
			IDSwipe: "",
			Swipe:   "like",
		}
		hour := 24 * time.Hour
		repoMock := &repoUser.UserMock{Mock: &mock.Mock{}}
		service := businessUser.NewService(repoMock, &config.AppConfig{})
		repoMock.On("Get", "apptinder:allrandomid:123").Return("1234", nil)
		repoMock.On("Set", "apptinder:allrandomid:123", "123,1234", hour).Return(nil)
		repoMock.On("GetRandomUser", mock.Anything).Return(businessUser.ResponseRandomUser{}, nil)

		err := service.SwipeUser(swipe.IDSwipe, swipe)
		asserting.Error(err)
	})

	t.Run("Error Get Me", func(t *testing.T) {
		asserting := assert.New(t)

		swipe := businessUser.SwipeUser{
			IDSwipe: "124",
			Swipe:   "like",
		}
		repoMock := &repoUser.UserMock{Mock: &mock.Mock{}}
		service := businessUser.NewService(repoMock, &config.AppConfig{})
		repoMock.On("GetMe", "1234").Return(businessUser.User{}, errors.New("error get me"))
		repoMock.On("Get", "apptinder:allrandomid:1234").Return("", nil)
		err := service.SwipeUser("1234", swipe)
		asserting.Error(err)
	})
	t.Run("Already Swipe Test", func(t *testing.T) {
		asserting := assert.New(t)
		user := businessUser.User{
			ID:    "123",
			Email: "test@mail.com",
		}
		swipe := businessUser.SwipeUser{
			IDSwipe: "1234",
			Swipe:   "like",
		}
		hour := 24 * time.Hour
		repoMock := &repoUser.UserMock{Mock: &mock.Mock{}}
		service := businessUser.NewService(repoMock, &config.AppConfig{})
		repoMock.On("GetMe", user.ID).Return(user, nil)
		repoMock.On("Get", "apptinder:allrandomid:123").Return("1234,24,124,214,214,214,21", nil)
		repoMock.On("Set", "apptinder:allrandomid:123", "123,1234", hour).Return(nil)
		repoMock.On("GetRandomUser", mock.Anything).Return(businessUser.ResponseRandomUser{}, nil)
		repoMock.On("SwipeUser", user.ID, mock.Anything).Return(nil)

		err := service.SwipeUser(user.ID, swipe)
		asserting.Error(err)
	})

	t.Run("Append Arr Test", func(t *testing.T) {
		asserting := assert.New(t)
		user := businessUser.User{
			ID:    "123",
			Email: "test@mail.com",
		}
		swipe := businessUser.SwipeUser{
			IDSwipe: "1234",
			Swipe:   "like",
		}
		hour := 24 * time.Hour
		repoMock := &repoUser.UserMock{Mock: &mock.Mock{}}
		service := businessUser.NewService(repoMock, &config.AppConfig{})
		repoMock.On("GetMe", user.ID).Return(user, nil)
		repoMock.On("Get", "apptinder:allrandomid:123").Return("", nil)
		repoMock.On("Set", "apptinder:allrandomid:123", "123,1234", hour).Return(nil)
		repoMock.On("GetRandomUser", mock.Anything).Return(businessUser.ResponseRandomUser{}, nil)
		repoMock.On("SwipeUser", user.ID, mock.Anything).Return(nil)

		err := service.SwipeUser(user.ID, swipe)
		asserting.NoError(err)
	})

}

func TestPurchasePackage(t *testing.T) {
	t.Run("Valid Test", func(t *testing.T) {
		asserting := assert.New(t)
		user := businessUser.User{
			ID:    "123",
			Email: "test@mail.com",
		}
		packages := "123"
		repoMock := &repoUser.UserMock{Mock: &mock.Mock{}}
		service := businessUser.NewService(repoMock, &config.AppConfig{})
		repoMock.On("PurchasePackage", user.ID, packages).Return(nil)
		repoMock.On("UpdatePackageUser", user.ID, mock.Anything).Return(nil)
		repoMock.On("GetPackageByID", packages).Return(businessUser.Package{}, nil)
		repoMock.On("GetMe", user.ID).Return(user, nil)

		err := service.PurchasePackage(user.ID, packages)
		asserting.NoError(err)
	})

	t.Run("Purchase Getme Error Test", func(t *testing.T) {
		asserting := assert.New(t)
		user := businessUser.User{
			ID:    "123",
			Email: "test@mail.com",
		}
		packages := "123"
		repoMock := &repoUser.UserMock{Mock: &mock.Mock{}}
		service := businessUser.NewService(repoMock, &config.AppConfig{})
		repoMock.On("PurchasePackage", user.ID, packages).Return(nil)
		repoMock.On("UpdatePackageUser", user.ID, mock.Anything).Return(nil)
		repoMock.On("GetPackageByID", packages).Return(businessUser.Package{}, nil)
		repoMock.On("GetMe", user.ID).Return(user, errors.New("error get me"))

		err := service.PurchasePackage(user.ID, packages)
		asserting.Error(err)
	})

	t.Run("already purchase Error Test", func(t *testing.T) {
		asserting := assert.New(t)
		user := businessUser.User{
			ID:      "123",
			Email:   "test@mail.com",
			Package: []string{"123"},
		}
		packages := "123"
		repoMock := &repoUser.UserMock{Mock: &mock.Mock{}}
		service := businessUser.NewService(repoMock, &config.AppConfig{})
		repoMock.On("PurchasePackage", user.ID, packages).Return(nil)
		repoMock.On("UpdatePackageUser", user.ID, mock.Anything).Return(nil)
		repoMock.On("GetPackageByID", packages).Return(businessUser.Package{}, nil)
		repoMock.On("GetMe", user.ID).Return(user, nil)

		err := service.PurchasePackage(user.ID, packages)
		asserting.Error(err)
	})

	t.Run("package not found Error Test", func(t *testing.T) {
		asserting := assert.New(t)
		user := businessUser.User{
			ID:    "123",
			Email: "test@mail.com",
		}
		packages := "123"
		repoMock := &repoUser.UserMock{Mock: &mock.Mock{}}
		service := businessUser.NewService(repoMock, &config.AppConfig{})
		repoMock.On("PurchasePackage", user.ID, packages).Return(nil)
		repoMock.On("UpdatePackageUser", user.ID, mock.Anything).Return(nil)
		repoMock.On("GetPackageByID", packages).Return(businessUser.Package{}, errors.New("package not found"))
		repoMock.On("GetMe", user.ID).Return(user, nil)

		err := service.PurchasePackage(user.ID, packages)
		asserting.Error(err)
	})

	t.Run("Update Package User Error Test", func(t *testing.T) {
		asserting := assert.New(t)
		user := businessUser.User{
			ID:    "123",
			Email: "test@mail.com",
		}
		packages := "123"
		repoMock := &repoUser.UserMock{Mock: &mock.Mock{}}
		service := businessUser.NewService(repoMock, &config.AppConfig{})
		repoMock.On("PurchasePackage", user.ID, packages).Return(nil)
		repoMock.On("UpdatePackageUser", user.ID, mock.Anything).Return(errors.New("error update package user"))
		repoMock.On("GetPackageByID", packages).Return(businessUser.Package{}, nil)
		repoMock.On("GetMe", user.ID).Return(user, nil)

		err := service.PurchasePackage(user.ID, packages)
		asserting.Error(err)
	})
}

func TestGetListPackage(t *testing.T) {
	t.Run("Valid Test", func(t *testing.T) {
		asserting := assert.New(t)
		packages := []businessUser.Package{
			{
				ID:          "123",
				PackageName: "test",
				Description: "test",
			},
		}
		repoMock := &repoUser.UserMock{Mock: &mock.Mock{}}
		service := businessUser.NewService(repoMock, &config.AppConfig{})
		repoMock.On("GetListPackage").Return(packages, nil)

		res, err := service.GetListPackage()
		asserting.NoError(err)
		asserting.NotNil(res)
	})
}

func TestGetMe(t *testing.T) {
	t.Run("Valid Test", func(t *testing.T) {
		asserting := assert.New(t)
		user := businessUser.User{
			ID:    "123",
			Email: "test@mail.com",
		}
		repoMock := &repoUser.UserMock{Mock: &mock.Mock{}}
		service := businessUser.NewService(repoMock, &config.AppConfig{})
		repoMock.On("GetMe", user.ID).Return(user, nil)

		res, err := service.GetMe(user.ID)
		asserting.NoError(err)

		asserting.NotNil(res)
	})
}

func TestGetPackageByID(t *testing.T) {
	t.Run("Valid Test", func(t *testing.T) {
		asserting := assert.New(t)
		packages := businessUser.Package{
			ID:          "123",
			PackageName: "test",
			Description: "test",
		}
		repoMock := &repoUser.UserMock{Mock: &mock.Mock{}}
		service := businessUser.NewService(repoMock, &config.AppConfig{})
		repoMock.On("GetPackageByID", packages.ID).Return(packages, nil)

		res, err := service.GetPackageByID(packages.ID)
		asserting.NoError(err)
		asserting.NotNil(res)
	})
}

func TestGetRandomUser(t *testing.T) {
	t.Run("Valid Test", func(t *testing.T) {
		asserting := assert.New(t)
		user := businessUser.User{
			ID:    "123",
			Email: "test@mail.com",
		}

		res := businessUser.ResponseRandomUser{
			ID:       "123",
			FullName: "test",
			Email:    "test@mail.com",
			PhotoUrl: "test",
		}

		repoMock := &repoUser.UserMock{Mock: &mock.Mock{}}
		service := businessUser.NewService(repoMock, &config.AppConfig{})
		repoMock.On("GetRandomUser", mock.Anything).Return(res, nil)
		repoMock.On("GetMe", user.ID).Return(user, nil)
		repoMock.On("Get", "apptinder:allrandomuser:123").Return("123", nil)
		repoMock.On("Set", "apptinder:allrandomuser:123", mock.Anything, mock.Anything).Return(nil)

		result, err := service.GetRandomUser(user.ID)
		asserting.NoError(err)
		asserting.NotNil(result)
	})

	t.Run("Error Get Random User Test", func(t *testing.T) {
		asserting := assert.New(t)
		user := businessUser.User{
			ID:    "123",
			Email: "test@mail.com",
		}

		res := businessUser.ResponseRandomUser{
			ID:       "123",
			FullName: "test",
			Email:    "test@mail.com",
			PhotoUrl: "test",
		}

		repoMock := &repoUser.UserMock{Mock: &mock.Mock{}}
		service := businessUser.NewService(repoMock, &config.AppConfig{})
		repoMock.On("GetRandomUser", mock.Anything).Return(res, errors.New("error get random user"))
		repoMock.On("GetMe", user.ID).Return(user, nil)
		repoMock.On("Get", "apptinder:allrandomuser:123").Return("123", nil)
		repoMock.On("Set", "apptinder:allrandomuser:123", mock.Anything, mock.Anything).Return(nil)

		_, err := service.GetRandomUser(user.ID)
		asserting.Error(err)
	})
	t.Run("Only 1 from val redis", func(t *testing.T) {
		asserting := assert.New(t)
		user := businessUser.User{
			ID:    "123",
			Email: "test@mail.com",
		}

		res := businessUser.ResponseRandomUser{
			ID:       "123",
			FullName: "test",
			Email:    "test@mail.com",
			PhotoUrl: "test",
		}
		repoMock := &repoUser.UserMock{Mock: &mock.Mock{}}
		service := businessUser.NewService(repoMock, &config.AppConfig{})
		repoMock.On("GetRandomUser", mock.Anything).Return(res, nil)
		repoMock.On("GetMe", user.ID).Return(user, nil)
		repoMock.On("Get", "apptinder:allrandomuser:123").Return("", nil)
		repoMock.On("Set", "apptinder:allrandomuser:123", mock.Anything, mock.Anything).Return(nil)

		result, err := service.GetRandomUser(user.ID)
		asserting.NoError(err)
		asserting.NotNil(result)
	})

	t.Run("Error set redis", func(t *testing.T) {
		asserting := assert.New(t)
		user := businessUser.User{
			ID:    "123",
			Email: "test@mail.com",
		}

		res := businessUser.ResponseRandomUser{
			ID:       "123",
			FullName: "test",
			Email:    "test@mail.com",
			PhotoUrl: "test",
		}
		repoMock := &repoUser.UserMock{Mock: &mock.Mock{}}
		service := businessUser.NewService(repoMock, &config.AppConfig{})
		repoMock.On("GetRandomUser", mock.Anything).Return(res, nil)
		repoMock.On("GetMe", user.ID).Return(user, nil)
		repoMock.On("Get", "apptinder:allrandomuser:123").Return("", nil)
		repoMock.On("Set", "apptinder:allrandomuser:123", mock.Anything, mock.Anything).Return(errors.New("error set redis"))

		_, err := service.GetRandomUser(user.ID)
		asserting.Error(err)
	})
}
