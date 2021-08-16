package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type CommentModel struct {
	UID       primitive.ObjectID  `bson:"_id" json:"id,omitempty"`
	Slug      string              `json:"slug"`
	Content   string              `json:"content"`
	Email     string              `json:"email"`
	Name      string              `json:"name"`
	CreatedAt interface{}         `json:"createdAt"`
	DeletedAt interface{}         `json:"deletedAt"`
	Replies   []CommentReplyModel `json:"replies"`
}

type CommentReplyModel struct {
	Content   string      `json:"content"`
	Email     string      `json:"email"`
	Name      string      `json:"name"`
	CreatedAt interface{} `json:"createdAt"`
	DeletedAt interface{} `json:"deletedAt"`
}
