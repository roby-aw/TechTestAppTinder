package main

import (
	"fmt"
	"os"
	"roby-backend-golang/app"
	"roby-backend-golang/config"
	"roby-backend-golang/utils"
)

// @title ID Backend Golang API Documentation
// @description This is a sample server for a ID Backend Golang API.
// @version 1.0.0
// @host localhost:8080
// @BasePath /v1
func main() {
	conf := config.GetConfig()
	dbCon := utils.NewConnectionDatabase(conf)
	server, port := app.Run(conf, dbCon)

	defer dbCon.CloseConnection()
	err := server.Listen(fmt.Sprintf(":%s", port))
	if err != nil {
		panic(err)
	}
	quit := make(chan os.Signal)
	<-quit
}
