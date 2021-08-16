package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type PostModel struct {
	UID         primitive.ObjectID    `bson:"_id" json:"id,omitempty"`
	Slug        string                `json:"slug"`
	Title       string                `json:"title"`
	Categories  []CategoryCommonModel `json:"categories"`
	Tags        []string              `json:"tags"`
	Content     string                `json:"content"`
	Author      UserCommonModel       `json:"author"`
	PublishedAt interface{}           `json:"publishedAt"`
	CreatedAt   interface{}           `json:"createdAt"`
	UpdatedAt   interface{}           `json:"updatedAt"`
	DeletedAt   interface{}           `json:"deletedAt"`
}
