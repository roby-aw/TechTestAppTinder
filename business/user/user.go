package user

import (
	"mime/multipart"
	"roby-backend-golang/utils"
)

type AuthLogin struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=8"`
}

type ResponseLogin struct {
	Email   string      `json:"email"`
	Package []Package   `json:"package"`
	Token   utils.Token `json:"token"`
}

type User struct {
	ID       string    `json:"id"`
	FullName string    `form:"fullname" validate:"required" json:"fullname"`
	Email    string    `form:"email" validate:"required,email" json:"email"`
	Password string    `json:"-" form:"password" validate:"required"`
	PhotoUrl string    `json:"photo_url"`
	Package  []string  `json:"-" bson:"package,omitempty"`
	Packages []Package `json:"packages"`
}

type ResponseRandomUser struct {
	ID       string    `json:"id"`
	FullName string    `json:"fullname" form:"fullname" validate:"required"`
	Email    string    `json:"email" form:"email" validate:"required,email"`
	PhotoUrl string    `json:"photo_url"`
	Packages []Package `json:"packages"`
}

type Register struct {
	FullName string                `form:"fullname" validate:"required"`
	Email    string                `form:"email" validate:"required,email"`
	Password string                `form:"password" validate:"required"`
	File     *multipart.FileHeader `form:"file"`
	PhotoUrl string                `json:"photo_url"`
}

type LastRandom struct {
	LastRandom string `json:"last_random"`
}

type SwipeUser struct {
	IDSwipe string `json:"id_swipe" validate:"required"`
	Swipe   string `json:"swipe" validate:"oneof=like pass"`
}

type Package struct {
	ID          string `json:"id" bson:"_id"`
	PackageName string `json:"package_name" bson:"package_name"`
	Description string `json:"description" bson:"description"`
}

type Purchase struct {
	ID string `json:"id" bson:"_id"`
}
