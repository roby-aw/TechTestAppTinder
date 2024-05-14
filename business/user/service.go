package user

import (
	"errors"
	"fmt"
	"mime/multipart"
	"roby-backend-golang/config"
	"roby-backend-golang/utils"
	"strings"
	"time"

	"github.com/go-playground/validator/v10"
	"golang.org/x/exp/slices"
)

type Repository interface {
	FindUserByID(id string) (User, error)
	FindUserByEmail(email string) (User, error)
	CreateUser(data Register) error
	GetRandomUser(id []string) (ResponseRandomUser, error)
	UploadImageS3(file *multipart.FileHeader) (string, error)
	PurchasePackage(id string, packages []string) error
	GetListPackage() ([]Package, error)
	GetMe(id string) (User, error)
	UpdatePackageUser(id string, idPackage []string) error
	GetPackageByID(id string) (Package, error)
	GenerateTokenAuth(id, email string) (*utils.Token, error)
	// Redis
	Set(key string, value interface{}, expiration time.Duration) error
	Get(key string) (string, error)
	Del(keys string) error
}
type Service interface {
	Login(auth AuthLogin) (*ResponseLogin, error)
	RegisterUser(data Register) error
	GetUserByID(id string) (User, error)
	GetRandomUser(id string) (ResponseRandomUser, error)
	SwipeUser(id string, input SwipeUser) error
	PurchasePackage(id, packages string) error
	GetListPackage() ([]Package, error)
	GetMe(id string) (User, error)
	GetPackageByID(id string) (Package, error)
}

type service struct {
	repository Repository
	validate   *validator.Validate
	conf       *config.AppConfig
}

func NewService(repository Repository, conf *config.AppConfig) Service {
	return &service{
		repository: repository,
		validate:   validator.New(),
		conf:       conf,
	}
}

func (s *service) Login(auth AuthLogin) (*ResponseLogin, error) {
	err := s.validate.Struct(&auth)
	if err != nil {
		return nil, utils.HandleErrorValidator(err)
	}

	user, err := s.repository.FindUserByEmail(auth.Email)
	if err != nil {
		return nil, err
	}

	err = utils.VerifyPassword(user.Password, auth.Password)
	if err != nil {
		fmt.Println(user.Password)
		return nil, errors.New("wrong password")
	}

	restoken, err := s.repository.GenerateTokenAuth(user.ID, user.Email)
	if err != nil {
		return nil, err
	}

	return &ResponseLogin{
		Email:   user.Email,
		Package: user.Packages,
		Token:   *restoken,
	}, nil
}

func (s *service) RegisterUser(data Register) error {
	err := s.validate.Struct(&data)
	if err != nil {
		return err
	}

	_, err = s.repository.FindUserByEmail(data.Email)
	if err == nil {
		return errors.New("email already exist")
	}

	data.PhotoUrl, err = s.repository.UploadImageS3(data.File)
	if err != nil {
		return err
	}

	err = s.repository.CreateUser(data)
	if err != nil {
		return err
	}
	return nil
}

func (s *service) GetUserByID(id string) (User, error) {
	return s.repository.FindUserByID(id)
}

func (s *service) GetRandomUser(id string) (ResponseRandomUser, error) {

	keyRedis := fmt.Sprintf("apptinder:allrandomuser:%s", id)

	// // first, check if data is exists in redis server get data from redis
	val, _ := s.repository.Get(keyRedis)

	strArr := strings.Split(val, ",")

	resUser, err := s.repository.GetRandomUser(append(strArr, id))
	if err != nil {
		return ResponseRandomUser{}, utils.HandleError(500, err.Error())
	}

	// if data is not exists in redis server, get data from mongodb
	now := time.Now()
	time12am := time.Date(now.Year(), now.Month(), now.Day(), 23, 59, 59, 59, time.Local)

	duration := time12am.Sub(now)

	if val != "" {
		val = fmt.Sprintf("%s,%s", val, resUser.ID)
	} else {
		val = resUser.ID
	}

	err = s.repository.Set(keyRedis, val, duration)
	if err != nil {
		return ResponseRandomUser{}, utils.HandleError(500, err.Error())
	}

	return resUser, nil
}

func (s *service) SwipeUser(id string, input SwipeUser) error {
	err := s.validate.Struct(&input)
	if err != nil {
		return utils.HandleErrorValidator(err)
	}

	keyRedis := fmt.Sprintf("apptinder:allrandomid:%s", id)

	// // first, check if data is exists in redis server get data from redis
	val, _ := s.repository.Get(keyRedis)

	res, err := s.repository.GetMe(id)
	if err != nil {
		return err
	}

	var premium bool
	for _, v := range res.Packages {
		if v.PackageName == "premium" {
			premium = true
		}
	}

	var strArr []string

	if val != "" {
		strArr = strings.Split(val, ",")
		if !premium {
			if len(strArr) >= 10 {
				return utils.HandleError(400, "please purchase premium packages that unlocks one premium feature")
			}
		}
		if utils.CheckArray(strArr, input.IDSwipe) {
			return utils.HandleError(400, "already swipe")
		}
	} else {
		strArr = append(strArr, id)
	}

	if len(strArr) > 1 {
		val = fmt.Sprintf("%s,%s", val, input.IDSwipe)
	} else {
		val = fmt.Sprintf("%s,%s", id, input.IDSwipe)
	}

	// if data is not exists in redis server, get data from mongodb
	err = s.repository.Set(keyRedis, val, 24*time.Hour)
	if err != nil {
		return utils.HandleError(500, err.Error())
	}

	return nil
}

func (s *service) PurchasePackage(id, packages string) error {
	res, err := s.repository.GetMe(id)
	if err != nil {
		return err
	}
	if slices.Contains(res.Package, packages) {
		return utils.HandleError(400, "already purchase package")
	}

	_, err = s.repository.GetPackageByID(packages)
	if err != nil {
		return utils.HandleError(400, "package not found")
	}

	res.Package = append(res.Package, packages)

	err = s.repository.UpdatePackageUser(id, res.Package)
	if err != nil {
		return utils.HandleError(500, err.Error())
	}
	return nil
}

func (s *service) GetListPackage() ([]Package, error) {
	return s.repository.GetListPackage()
}

func (s *service) GetMe(id string) (User, error) {
	return s.repository.GetMe(id)
}

func (s *service) GetPackageByID(id string) (Package, error) {
	return s.repository.GetPackageByID(id)
}
