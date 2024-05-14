package user

import (
	"mime/multipart"
	businessUser "roby-backend-golang/business/user"
	"roby-backend-golang/utils"
	"time"

	"github.com/stretchr/testify/mock"
)

type UserMock struct {
	*mock.Mock
}

func (m *UserMock) FindUserByID(id string) (businessUser.User, error) {
	args := m.Called(id)
	return args.Get(0).(businessUser.User), args.Error(1)
}

func (m *UserMock) FindUserByEmail(email string) (businessUser.User, error) {
	args := m.Called(email)
	return args.Get(0).(businessUser.User), args.Error(1)
}

func (m *UserMock) CreateUser(data businessUser.Register) error {
	args := m.Called(data)
	return args.Error(0)
}

func (m *UserMock) GetRandomUser(id []string) (businessUser.ResponseRandomUser, error) {
	args := m.Called(id)
	return args.Get(0).(businessUser.ResponseRandomUser), args.Error(1)
}

func (m *UserMock) UploadImageS3(file *multipart.FileHeader) (string, error) {
	args := m.Called(file)
	return args.String(0), args.Error(1)
}

func (m *UserMock) PurchasePackage(id string, packages []string) error {
	args := m.Called(id, packages)
	return args.Error(0)
}

func (m *UserMock) GetListPackage() ([]businessUser.Package, error) {
	args := m.Called()
	return args.Get(0).([]businessUser.Package), args.Error(1)
}

func (m *UserMock) GetMe(id string) (businessUser.User, error) {
	args := m.Called(id)
	return args.Get(0).(businessUser.User), args.Error(1)
}

func (m *UserMock) UpdatePackageUser(id string, idPackage []string) error {
	args := m.Called(id, idPackage)
	return args.Error(0)
}

func (m *UserMock) GetPackageByID(id string) (businessUser.Package, error) {
	args := m.Called(id)
	return args.Get(0).(businessUser.Package), args.Error(1)
}

func (m *UserMock) Set(key string, value interface{}, expiration time.Duration) error {
	args := m.Called(key, value, expiration)
	return args.Error(0)
}

func (m *UserMock) Get(key string) (string, error) {
	args := m.Called(key)
	return args.String(0), args.Error(1)
}

func (m *UserMock) Del(keys string) error {
	args := m.Called(keys)
	return args.Error(0)
}

func (m *UserMock) GenerateTokenAuth(id, email string) (*utils.Token, error) {
	args := m.Called(id, email)
	return args.Get(0).(*utils.Token), args.Error(1)
}
