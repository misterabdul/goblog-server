package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type CategoryModel struct {
	UID       primitive.ObjectID `bson:"_id" json:"id,omitempty"`
	Slug      string             `json:"slug"`
	Name      string             `json:"name"`
	CreatedAt interface{}        `json:"createdAt"`
	UpdatedAt interface{}        `json:"updatedAt"`
	DeletedAt interface{}        `json:"deletedAt"`
}

type CategoryCommonModel struct {
	UID  primitive.ObjectID `bson:"_id" json:"id,omitempty"`
	Slug string             `json:"slug"`
	Name string             `json:"name"`
}

func (category *CategoryModel) ToCommonModel() (commonModel CategoryCommonModel) {
	return CategoryCommonModel{
		UID:  category.UID,
		Slug: category.Slug,
		Name: category.Name}
}
