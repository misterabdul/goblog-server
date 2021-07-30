package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type UserModel struct {
	UID       primitive.ObjectID `bson:"_id" json:"id,omitempty"`
	FirstName string             `json:"firstName"`
	LastName  string             `json:"lastName"`
	Email     string             `json:"email"`
	Password  string             `json:"password"`
	CreatedAt interface{}        `json:"createdAt"`
	UpdatedAt interface{}        `json:"updatedAt"`
	DeletedAt interface{}        `json:"deletedAt"`
}
