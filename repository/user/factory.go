package user

import (
	"roby-backend-golang/business/user"
	"roby-backend-golang/config"
	"roby-backend-golang/utils"
)

func RepositoryFactory(dbCon *utils.DatabaseConnection, conf *config.AppConfig) user.Repository {
	adminRepo := NewMongoRepository(dbCon, conf)
	return adminRepo
}
