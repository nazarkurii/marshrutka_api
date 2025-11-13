package uuidutil

import (
	rfc7807 "maryan_api/pkg/problem"

	"github.com/d3code/uuid"
)

func Parse(idString string) (uuid.UUID, error) {
	id, err := uuid.Parse(idString)
	if err != nil {
		return uuid.Nil, rfc7807.UUID(err.Error())
	}
	return id, nil
}
