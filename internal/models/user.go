package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type UserModel struct {
	UID       primitive.ObjectID `bson:"_id" json:"id,omitempty"`
	FirstName string             `json:"firstName"`
	LastName  string             `json:"lastName"`
	Username  string             `json:"username"`
	Email     string             `json:"email"`
	Password  string             `json:"password"`
	Roles     []UserRoles        `json:"roles"`
	CreatedAt interface{}        `json:"createdAt"`
	UpdatedAt interface{}        `json:"updatedAt"`
	DeletedAt interface{}        `json:"deletedAt"`
}

type UserCommonModel struct {
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
	Username  string `json:"username"`
	Email     string `json:"email"`
}

// 0 => SuperAdmin
//
// 1 => Admin
//
// 2 => Editor
type UserRoles struct {
	Level int                `json:"level"`
	Name  string             `json:"name"`
	Since primitive.DateTime `json:"since"`
}

func CreateUserCommonModel(user UserModel) UserCommonModel {
	return UserCommonModel{
		FirstName: user.FirstName,
		LastName:  user.LastName,
		Username:  user.Username,
		Email:     user.Email,
	}
}
