package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type NotificationModel struct {
	UID       primitive.ObjectID `bson:"_id" json:"id,omitempty"`
	Title     string             `json:"title"`
	Content   string             `json:"content"`
	Owner     UserCommonModel    `json:"owner"`
	ReadAt    interface{}        `json:"readAt"`
	CreatedAt interface{}        `json:"createdAt"`
	DeletedAt interface{}        `json:"deletedAt"`
}
