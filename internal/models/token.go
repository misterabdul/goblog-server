package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type RevokedTokenModel struct {
	UID       primitive.ObjectID `bson:"_id" json:"id,omitempty"`
	ExpiresAt primitive.DateTime `json:"expiresAt"`
	Owner     UserCommonModel    `json:"owner"`
	CreatedAt interface{}        `json:"createdAt"`
	UpdatedAt interface{}        `json:"updatedAt"`
	DeletedAt interface{}        `json:"deletedAt"`
}
