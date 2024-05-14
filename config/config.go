package config

import (
	"os"
	"sync"

	"github.com/joho/godotenv"
)

type AppConfig struct {
	App struct {
		Port string `toml:"port"`
	} `toml:"app"`
	Database struct {
		DBURL  string
		DBNAME string
		DBPASS string
		DBUSER string
		DBETC  string
		DBPORT string

		REDIS_HOST string
		REDIS_PASS string
	} `toml:"database"`
	AwsS3 struct {
		URL    string
		Access string
		Secret string
		Bucket string
		Zone   string
	}
	Secrettoken struct {
		Token string `toml:"token"`
	} `toml:"secrettoken"`
}

var lock = &sync.Mutex{}
var appConfig *AppConfig

func GetConfig() *AppConfig {
	_ = godotenv.Load()
	lock.Lock()
	defer lock.Unlock()

	if appConfig == nil {
		appConfig = initConfig()
	}
	return appConfig
}

func initConfig() *AppConfig {
	var finalConfig AppConfig
	finalConfig.App.Port = "8080"
	finalConfig.Database.DBNAME = os.Getenv("MONGO_DBNAME")
	finalConfig.Database.DBURL = os.Getenv("MONGO_HOST")
	finalConfig.Database.DBPASS = os.Getenv("MONGO_PASS")
	finalConfig.Database.DBUSER = os.Getenv("MONGO_USER")
	finalConfig.Database.DBETC = os.Getenv("MONGO_ETC")
	finalConfig.Database.DBPORT = os.Getenv("MONGO_PORT")

	finalConfig.Secrettoken.Token = os.Getenv("JWT_SECRET")

	finalConfig.Database.REDIS_HOST = os.Getenv("REDIS_HOST")
	finalConfig.Database.REDIS_PASS = os.Getenv("REDIS_PASS")

	finalConfig.AwsS3.URL = os.Getenv("AWS_S3_URL")
	finalConfig.AwsS3.Access = os.Getenv("AWS_S3_ACCESS")
	finalConfig.AwsS3.Secret = os.Getenv("AWS_S3_SECRET")
	finalConfig.AwsS3.Bucket = os.Getenv("AWS_S3_BUCKET")
	finalConfig.AwsS3.Zone = os.Getenv("AWS_S3_ZONE")

	return &finalConfig
}
