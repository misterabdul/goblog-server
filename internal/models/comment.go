package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type CommentModel struct {
	UID       primitive.ObjectID  `bson:"_id" json:"id,omitempty"`
	PostUid   primitive.ObjectID  `json:"postUid"`
	Email     string              `json:"email"`
	Name      string              `json:"name"`
	Content   string              `json:"content"`
	Replies   []CommentReplyModel `json:"replies"`
	CreatedAt interface{}         `json:"createdAt"`
	DeletedAt interface{}         `json:"deletedAt"`
}

type CommentReplyModel struct {
	Email     string      `json:"email"`
	Name      string      `json:"name"`
	Content   string      `json:"content"`
	CreatedAt interface{} `json:"createdAt"`
	DeletedAt interface{} `json:"deletedAt"`
}
