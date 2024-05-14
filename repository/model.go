package repository

import (
	"roby-backend-golang/business/user"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type User struct {
	ID       primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	Email    string             `json:"email" bson:"email,omitempty"`
	Password string             `json:"password" bson:"password,omitempty"`
	Fullname string             `json:"fullname" bson:"fullname,omitempty"`
	PhotoUrl string             `json:"photo_url" bson:"photo_url,omitempty"`
	Package  []string           `json:"package" bson:"package,omitempty"`
	Packages []user.Package     `json:"packages" bson:"packages,omitempty"`
}

type Package struct {
	ID          primitive.ObjectID `bson:"_id,omitempty" json:"_id"`
	PackageName string             `bson:"package_name,omitempty" binding:"required" json:"package_name"`
	Description string             `bson:"description,omitempty" binding:"required" json:"description"`
}

type Role struct {
	ID          primitive.ObjectID `bson:"_id,omitempty" json:"_id"`
	Rolename    string             `bson:"rolename,omitempty" binding:"required" json:"rolename"`
	Rolelabel   string             `bson:"rolelabel,omitempty" binding:"required" json:"rolelabel"`
	Description string             `bson:"description,omitempty" binding:"required" json:"description"`
}

type RegisterUser struct {
	ID        primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	Fullname  string             `json:"fullname" bson:"fullname,omitempty"`
	Package   string             `json:"package" bson:"package,omitempty"`
	Email     string             `json:"email" bson:"email,omitempty"`
	Type      string             `json:"type" bson:"type,omitempty"`
	Password  string             `json:"password" bson:"password,omitempty"`
	PhotoUrl  string             `json:"photo_url" bson:"photo_url,omitempty"`
	CreatedAt time.Time          `json:"created_at" bson:"created_at,omitempty"`
	UpdatedAt time.Time          `json:"updated_at" bson:"updated_at,omitempty"`
}

type FilterQuery bson.M

func NewFilterQuery() FilterQuery {
	return FilterQuery(bson.M{})
}

func (q FilterQuery) SetDeleteAtExists(exist bool) FilterQuery {
	q["deleted_at"] = bson.M{"$exists": exist}
	return q
}

func (q FilterQuery) SetRole() FilterQuery {
	q["$lookup"] = bson.M{
		"from":         "role",
		"localField":   "role_id",
		"foreignField": "id",
		"as":           "role",
	}
	return q
}

func (q FilterQuery) SetID(id primitive.ObjectID) FilterQuery {
	q["_id"] = id
	return q
}

func (q FilterQuery) SetEmail(email string) FilterQuery {
	q["email"] = email
	return q
}

func (q FilterQuery) SetDeleteAt() FilterQuery {
	q["deleted_at"] = time.Now()
	return q
}

func (q FilterQuery) SetUpdateAt() FilterQuery {
	q["updated_at"] = time.Now()
	return q
}

func (q FilterQuery) SetUpdate(data interface{}) FilterQuery {
	q["$set"] = data
	return q
}

func (q FilterQuery) SetSortDescending(field string) FilterQuery {
	q["$sort"] = bson.M{field: -1}
	return q
}

func (q FilterQuery) SetSortAscending(field string) FilterQuery {
	q["$sort"] = bson.M{field: 1}
	return q
}

func (q FilterQuery) SetRoleID(id int) FilterQuery {
	q["role_id"] = id
	return q
}

func (q FilterQuery) SetPassword(password string) FilterQuery {
	q["password"] = password
	return q
}

func (q FilterQuery) SetFullname(fullname string) FilterQuery {
	q["fullname"] = fullname
	return q
}
