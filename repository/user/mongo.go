package user

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"math/rand"
	"mime/multipart"
	"path/filepath"
	businessUser "roby-backend-golang/business/user"
	"roby-backend-golang/config"
	"roby-backend-golang/repository"
	"roby-backend-golang/utils"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/go-redis/redis/v8"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/net/context"
)

type MongoDBRepository struct {
	colUser *mongo.Collection
	colPack *mongo.Collection
	conf    *config.AppConfig
	aws     *session.Session
	redis   *redis.Client
}

func NewMongoRepository(dbCon *utils.DatabaseConnection, conf *config.AppConfig) *MongoDBRepository {
	return &MongoDBRepository{
		colUser: dbCon.MongoDB.Collection("user"),
		colPack: dbCon.MongoDB.Collection("package"),
		conf:    conf,
		aws:     dbCon.AwsS3,
		redis:   dbCon.Redis,
	}
}

func (repo *MongoDBRepository) FindUserByEmail(email string) (businessUser.User, error) {
	var user repository.User
	var userBusiness businessUser.User
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	queryFilter := bson.A{
		bson.M{"$match": bson.M{"email": email}},
		bson.M{"$lookup": bson.M{
			"from":         "package",
			"localField":   "package",
			"foreignField": "_id",
			"as":           "packages",
		}},
		bson.M{
			"$limit": 1,
		},
	}

	cur, err := repo.colUser.Aggregate(ctx, queryFilter)
	if err != nil {
		return userBusiness, err
	}

	for cur.Next(ctx) {
		err = cur.Decode(&user)
		if err != nil {
			return userBusiness, err
		}
	}

	if user.ID.IsZero() {
		return userBusiness, errors.New("wrong email")
	}

	userBusiness.ID = user.ID.Hex()
	userBusiness.Email = user.Email
	userBusiness.Password = user.Password
	userBusiness.PhotoUrl = user.PhotoUrl
	userBusiness.FullName = user.Fullname
	userBusiness.Packages = user.Packages

	return userBusiness, nil
}

func (repo *MongoDBRepository) CreateUser(data businessUser.Register) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	passwd, err := utils.Hash(data.Password)
	if err != nil {
		return errors.New("failed to hash password")
	}

	insUser := repository.RegisterUser{
		ID:       primitive.NewObjectID(),
		Fullname: data.FullName,
		Email:    data.Email,
		Type:     "free",
		Password: string(passwd),
		PhotoUrl: data.PhotoUrl,
	}

	_, err = repo.colUser.InsertOne(ctx, insUser)
	if err != nil {
		return err
	}

	return nil
}

func (repo *MongoDBRepository) FindUserByID(id string) (businessUser.User, error) {
	var user repository.User
	var userBusiness businessUser.User
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return userBusiness, errors.New("invalid id")
	}

	queryFilter := repository.NewFilterQuery()
	queryFilter.SetID(objID)

	err = repo.colUser.FindOne(ctx, queryFilter).Decode(&user)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return userBusiness, errors.New("wrong id")
		}
		return userBusiness, err
	}
	userBusiness.ID = user.ID.Hex()
	userBusiness.Email = user.Email
	userBusiness.Password = user.Password
	return userBusiness, nil
}

func (repo *MongoDBRepository) UploadImageS3(file *multipart.FileHeader) (string, error) {
	ext := strings.ToLower(filepath.Ext(file.Filename))
	if ext != ".jpeg" && ext != ".webp" && ext != ".jpg" {
		return "", utils.HandleError(400, "file type not allowed")
	}

	// open file
	filearr, err := file.Open()
	if err != nil {
		return "", utils.HandleError(500, err.Error())
	}

	buf := bytes.NewBuffer(nil)

	// copy file to buffer
	if _, err := io.Copy(buf, filearr); err != nil {
		return "", utils.HandleError(500, err.Error())
	}

	// decode image
	_, err = utils.DecodeImage(buf.Bytes(), filepath.Ext(file.Filename))
	if err != nil {
		return "", utils.HandleError(500, err.Error())
	}
	filename := fmt.Sprintf("%d-%d.jpeg", time.Now().Unix(), rand.Intn(1000))

	_, err = utils.UploadImageS3(repo.conf, buf.Bytes(), file.Header.Get("Content-Type"), filename, repo.aws)
	if err != nil {
		return "", err
	}

	url := fmt.Sprintf("%s/%s/%s", repo.conf.AwsS3.URL, repo.conf.AwsS3.Bucket, filename)

	return url, nil
}

func (repo *MongoDBRepository) GetRandomUser(id []string) (businessUser.ResponseRandomUser, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	var user repository.User
	var userBusiness businessUser.ResponseRandomUser
	var objArr []primitive.ObjectID
	for _, v := range id {
		if v != "" {
			objID, err := primitive.ObjectIDFromHex(v)
			if err != nil {
				return userBusiness, errors.New("invalid id")
			}
			objArr = append(objArr, objID)
		}
	}

	filter := bson.A{
		bson.M{"$match": bson.M{"_id": bson.M{"$nin": objArr}}},
		bson.M{"$sample": bson.M{"size": 1}},
		bson.M{"$lookup": bson.M{
			"from":         "package",
			"localField":   "package",
			"foreignField": "_id",
			"as":           "packages",
		}},
	}

	cur, err := repo.colUser.Aggregate(ctx, filter)
	if err != nil {
		return userBusiness, err
	}

	for cur.Next(ctx) {
		err = cur.Decode(&user)
		if err != nil {
			return userBusiness, err
		}
	}

	if user.ID.IsZero() {
		return userBusiness, errors.New("no user found, please wait for tomorrow")
	}

	userBusiness.ID = user.ID.Hex()
	userBusiness.Email = user.Email
	userBusiness.PhotoUrl = user.PhotoUrl
	userBusiness.FullName = user.Fullname
	userBusiness.Packages = user.Packages

	return userBusiness, nil
}

func (repo *MongoDBRepository) Set(key string, value interface{}, expiration time.Duration) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	return repo.redis.Set(ctx, key, value, expiration).Err()
}

func (repo *MongoDBRepository) Get(key string) (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	res, err := repo.redis.Get(ctx, key).Result()
	if err != nil {
		return "", err
	}

	return res, nil
}

func (repo *MongoDBRepository) Del(key string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	return repo.redis.Del(ctx, key).Err()
}

func (repo *MongoDBRepository) PurchasePackage(id string, packages []string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return errors.New("invalid id")
	}

	queryFilter := repository.NewFilterQuery()
	queryFilter.SetID(objID)

	update := bson.M{"$set": bson.M{"package": packages}}

	_, err = repo.colUser.UpdateOne(ctx, queryFilter, update)
	if err != nil {
		return err
	}

	return nil
}

func (repo *MongoDBRepository) GetListPackage() ([]businessUser.Package, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	var packages []businessUser.Package

	cur, err := repo.colPack.Find(ctx, bson.M{})
	if err != nil {
		return packages, err
	}

	for cur.Next(ctx) {
		var pack businessUser.Package
		err = cur.Decode(&pack)
		if err != nil {
			return packages, err
		}
		packages = append(packages, pack)
	}

	return packages, nil
}

func (repo *MongoDBRepository) GetMe(id string) (businessUser.User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	var user repository.User
	var userBusiness businessUser.User

	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return userBusiness, errors.New("invalid id")
	}

	filter := bson.A{
		bson.M{"$match": bson.M{"_id": objID}},
		bson.M{"$lookup": bson.M{
			"from":         "package",
			"localField":   "package",
			"foreignField": "_id",
			"as":           "packages",
		}},
	}

	cursor, err := repo.colUser.Aggregate(ctx, filter)
	if err != nil {
		return userBusiness, err
	}

	for cursor.Next(ctx) {
		err = cursor.Decode(&user)
		if err != nil {
			return userBusiness, err
		}
		userBusiness.ID = user.ID.Hex()
		userBusiness.Email = user.Email
		userBusiness.PhotoUrl = user.PhotoUrl
		userBusiness.FullName = user.Fullname
		userBusiness.Packages = user.Packages
		userBusiness.Package = user.Package
	}

	return userBusiness, nil
}

func (repo *MongoDBRepository) UpdatePackageUser(id string, idPackage []string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	var objArr []primitive.ObjectID
	for _, v := range idPackage {
		objID, err := primitive.ObjectIDFromHex(v)
		if err != nil {
			return errors.New("invalid id")
		}
		objArr = append(objArr, objID)
	}

	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return errors.New("invalid id")
	}

	queryFilter := repository.NewFilterQuery()
	queryFilter.SetID(objID)

	update := bson.M{"$set": bson.M{"package": objArr}}

	_, err = repo.colUser.UpdateOne(ctx, queryFilter, update)
	if err != nil {
		return err
	}

	return nil
}

func (repo *MongoDBRepository) GetPackageByID(id string) (businessUser.Package, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	var pack businessUser.Package

	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return pack, errors.New("invalid id")
	}

	queryFilter := repository.NewFilterQuery()
	queryFilter.SetID(objID)

	err = repo.colPack.FindOne(ctx, queryFilter).Decode(&pack)
	if err != nil {
		return pack, err
	}

	return pack, nil
}

func (repo *MongoDBRepository) GenerateTokenAuth(id, email string) (*utils.Token, error) {
	exp, token, err := utils.GenerateAccessTokenUser(id, email, repo.conf.Secrettoken.Token)
	if err != nil {
		return nil, err
	}
	exprefresh, refreshtoken, err := utils.GenerateRefreshTokenUser(id, email, repo.conf.Secrettoken.Token)
	if err != nil {
		return nil, err
	}
	var restoken = utils.Token{
		AccessToken:         token,
		AccessTokenExpired:  exp,
		RefreshToken:        refreshtoken,
		RefreshTokenExpired: exprefresh,
	}

	return &restoken, nil
}
