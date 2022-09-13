package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type UserModel struct {
	UID       primitive.ObjectID `bson:"_id" json:"id,omitempty"`
	Username  string             `json:"username"`
	Email     string             `json:"email"`
	FirstName string             `json:"firstName"`
	LastName  string             `json:"lastName"`
	Password  string             `json:"password"`
	Roles     []UserRole         `json:"roles"`
	CreatedAt interface{}        `json:"createdAt"`
	UpdatedAt interface{}        `json:"updatedAt"`
	DeletedAt interface{}        `json:"deletedAt"`
}

type UserCommonModel struct {
	UID       primitive.ObjectID `bson:"_id" json:"id,omitempty"`
	Username  string             `json:"username"`
	Email     string             `json:"email"`
	FirstName string             `json:"firstName"`
	LastName  string             `json:"lastName"`
}

// 0 => SuperAdmin
//
// 1 => Admin
//
// 2 => Editor
//
// 3 => Writer
type UserRole struct {
	Level int                `json:"level"`
	Name  string             `json:"name"`
	Since primitive.DateTime `json:"since"`
}

func (user *UserModel) ToCommonModel() (commonModel UserCommonModel) {
	return UserCommonModel{
		UID:       user.UID,
		FirstName: user.FirstName,
		LastName:  user.LastName,
		Username:  user.Username,
		Email:     user.Email}
}
