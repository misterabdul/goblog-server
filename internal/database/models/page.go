package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type PageModel struct {
	UID         primitive.ObjectID `bson:"_id" json:"id,omitempty"`
	Slug        string             `json:"slug"`
	Title       string             `json:"title"`
	Author      UserCommonModel    `json:"author"`
	PublishedAt interface{}        `json:"publishedAt"`
	CreatedAt   interface{}        `json:"createdAt"`
	UpdatedAt   interface{}        `json:"updatedAt"`
	DeletedAt   interface{}        `json:"deletedAt"`
}

type PageContentModel struct {
	UID     primitive.ObjectID `bson:"_id" json:"id,omitempty"`
	Content string             `json:"content"`
}
