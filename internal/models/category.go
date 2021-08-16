package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type CategoryModel struct {
	UID       primitive.ObjectID `bson:"_id" json:"id,omitempty"`
	Name      string             `json:"name"`
	Slug      string             `json:"slug"`
	CreatedAt interface{}        `json:"createdAt"`
	UpdatedAt interface{}        `json:"updatedAt"`
	DeletedAt interface{}        `json:"deletedAt"`
}

type CategoryCommonModel struct {
	Name string `json:"name"`
	Slug string `json:"slug"`
}
