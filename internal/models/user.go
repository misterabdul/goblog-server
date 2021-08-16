package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type UserModel struct {
	UID                 primitive.ObjectID `bson:"_id" json:"id,omitempty"`
	FirstName           string             `json:"firstName"`
	LastName            string             `json:"lastName"`
	Username            string             `json:"username"`
	Email               string             `json:"email"`
	Password            string             `json:"password"`
	IssuedRefreshTokens []IssuedToken      `json:"issuedRefreshTokens"`
	IssuedAccessTokens  []IssuedToken      `json:"issuedAccessTokens"`
	RevokedAccessTokens []RevokedToken     `json:"revokedAccessTokens"`
	CreatedAt           interface{}        `json:"createdAt"`
	UpdatedAt           interface{}        `json:"updatedAt"`
	DeletedAt           interface{}        `json:"deletedAt"`
}

type UserCommonModel struct {
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
	Username  string `json:"username"`
	Email     string `json:"email"`
}

type IssuedToken struct {
	TokenUID  string             `json:"tokenUid"`
	Client    string             `json:"client"`
	IssuedAt  primitive.DateTime `json:"issuedAt"`
	ExpiredAt primitive.DateTime `json:"expiredAt"`
}

type RevokedToken struct {
	TokenUID string             `json:"tokenUid"`
	Until    primitive.DateTime `json:"until"`
}

func CreateUserCommonModel(user UserModel) UserCommonModel {
	return UserCommonModel{
		FirstName: user.FirstName,
		LastName:  user.LastName,
		Username:  user.Username,
		Email:     user.Email,
	}
}
