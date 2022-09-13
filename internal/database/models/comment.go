package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type CommentModel struct {
	UID              primitive.ObjectID `bson:"_id" json:"id,omitempty"`
	PostUid          primitive.ObjectID `json:"postUid,omitempty"`
	PostAuthorUid    primitive.ObjectID `json:"postAuthorUid,omitempty"`
	ParentCommentUid interface{}        `json:"parentCommentUid,omitempty"`
	Email            string             `json:"email"`
	Name             string             `json:"name"`
	Content          string             `json:"content"`
	ReplyCount       int16              `json:"replyCount"`
	CreatedAt        interface{}        `json:"createdAt"`
	DeletedAt        interface{}        `json:"deletedAt"`
}
