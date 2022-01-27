package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type PostModel struct {
	UID                primitive.ObjectID    `bson:"_id" json:"id,omitempty"`
	Slug               string                `json:"slug"`
	Title              string                `json:"title"`
	FeaturingImagePath string                `json:"featuringImagePath"`
	Description        string                `json:"description"`
	Categories         []CategoryCommonModel `json:"categories"`
	Tags               []string              `json:"tags"`
	Author             UserCommonModel       `json:"author"`
	CommentCount       int16                 `json:"commentCount"`
	PublishedAt        interface{}           `json:"publishedAt"`
	CreatedAt          interface{}           `json:"createdAt"`
	UpdatedAt          interface{}           `json:"updatedAt"`
	DeletedAt          interface{}           `json:"deletedAt"`
}

type PostContentModel struct {
	UID     primitive.ObjectID `bson:"_id" json:"id,omitempty"`
	Content string             `json:"content"`
}
