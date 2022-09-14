package payloads

import (
	"encoding/json"

	"github.com/misterabdul/goblog-server/internal/database/models"
)

type UpdateMePayload struct {
	models.UserModel
}

func (p *UpdateMePayload) Marshall() (
	data []byte,
	err error,
) {
	return json.Marshal(p)
}

func NewUpdateMePayload(updatedUser models.UserModel) (
	payload *UpdateMePayload,
) {
	return &UpdateMePayload{
		updatedUser}
}

func UnmarshallUpdateMePayload(data []byte) (
	payload *UpdateMePayload,
	err error,
) {
	var _payload UpdateMePayload

	if err = json.Unmarshal(data, &_payload); err != nil {
		return nil, err
	}

	return &_payload, nil
}
