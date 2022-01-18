package forms

import "go.mongodb.org/mongo-driver/bson/primitive"

func toObjectIdArray(objectIdHexs []string) (
	objectIds []primitive.ObjectID,
	err error,
) {
	var objectId primitive.ObjectID

	for _, objectIdHex := range objectIdHexs {
		if objectId, err = primitive.ObjectIDFromHex(objectIdHex); err != nil {
			return nil, err
		}
		objectIds = append(objectIds, objectId)
	}

	return objectIds, nil
}
