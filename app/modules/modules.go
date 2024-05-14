package modules

import (
	"roby-backend-golang/api"
	userController "roby-backend-golang/api/user"
	userBusiness "roby-backend-golang/business/user"
	"roby-backend-golang/config"
	userRepository "roby-backend-golang/repository/user"
	"roby-backend-golang/utils"
)

func RegistrationModules(dbCon *utils.DatabaseConnection, conf *config.AppConfig) api.Controller {
	userPermitRepository := userRepository.RepositoryFactory(dbCon, conf)
	userPermitService := userBusiness.NewService(userPermitRepository, conf)
	userPermitController := userController.NewController(userPermitService)
	// Register controller
	controller := api.Controller{
		UserController: userPermitController,
	}

	return controller
}
