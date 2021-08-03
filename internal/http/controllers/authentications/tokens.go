package authentications

import (
	"context"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"

	"github.com/misterabdul/goblog-server/internal/models"
	"github.com/misterabdul/goblog-server/internal/repositories"
	"github.com/misterabdul/goblog-server/pkg/hash"
	"github.com/misterabdul/goblog-server/pkg/jwt"
)

func saveToken(ctx context.Context, dbConn *mongo.Database, user *models.UserModel, accessToken *jwt.Claims, refreshToken *jwt.Claims) error {
	hashedRefreshToken, err := hash.Make(refreshToken.TokenUID)
	if err != nil {
		return err
	}

	user.IssuedRefreshTokens = append(user.IssuedRefreshTokens, models.IssuedToken{
		TokenUID:  hashedRefreshToken,
		Client:    "",
		IssuedAt:  primitive.NewDateTimeFromTime(refreshToken.IssuedAt),
		ExpiredAt: primitive.NewDateTimeFromTime(refreshToken.ExpiredAt),
	})
	user.IssuedAccessTokens = append(user.IssuedAccessTokens, models.IssuedToken{
		TokenUID:  accessToken.TokenUID,
		Client:    "",
		IssuedAt:  primitive.NewDateTimeFromTime(accessToken.IssuedAt),
		ExpiredAt: primitive.NewDateTimeFromTime(accessToken.ExpiredAt),
	})

	return repositories.UpdateUser(ctx, dbConn, user)
}
