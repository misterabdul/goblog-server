package authentications

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"

	"github.com/misterabdul/goblog-server/internal/database/models"
	"github.com/misterabdul/goblog-server/internal/service"
	"github.com/misterabdul/goblog-server/pkg/jwt"
)

func noteRevokeToken(
	service *service.RevokedTokenService,
	refreshClaims *jwt.CustomClaims,
	user *models.UserModel,
) (err error) {
	var revokeTokenData *models.RevokedTokenModel

	if revokeTokenData, err = createRevokeModelFromClaims(refreshClaims); err != nil {
		return err
	}
	revokeTokenData.Owner = user.ToCommonModel()

	return service.CreateRevokedToken(revokeTokenData)
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
