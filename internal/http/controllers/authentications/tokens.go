package authentications

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"

	"github.com/misterabdul/goblog-server/internal/models"
	"github.com/misterabdul/goblog-server/internal/repositories"
	"github.com/misterabdul/goblog-server/pkg/jwt"
)

func noteRevokeToken(
	ctx context.Context,
	dbConn *mongo.Database,
	refreshClaims *jwt.CustomClaims,
	user *models.UserModel,
) (err error) {
	var revokeTokenData *models.RevokedTokenModel

	if revokeTokenData, err = createRevokeModelFromClaims(refreshClaims); err != nil {
		return err
	}
	revokeTokenData.Owner = user.ToCommonModel()

	return repositories.CreateRevokedToken(ctx, dbConn, revokeTokenData)
}

func createRevokeModelFromClaims(refreshClaims *jwt.CustomClaims) (
	model *models.RevokedTokenModel,
	err error,
) {
	var (
		revokeTokenUID primitive.ObjectID
		expiresAtTime  = time.Unix(refreshClaims.ExpiresAt, 0)
		_model         models.RevokedTokenModel
	)

	if revokeTokenUID, err = primitive.ObjectIDFromHex(refreshClaims.Id); err != nil {
		return nil, err
	}
	_model = models.RevokedTokenModel{
		UID:       revokeTokenUID,
		ExpiresAt: primitive.NewDateTimeFromTime(expiresAtTime),
	}

	return &_model, nil
}
