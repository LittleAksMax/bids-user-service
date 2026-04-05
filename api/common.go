package api

import (
	"errors"
	"net/http"

	"github.com/google/uuid"
)

const profileIDPath = "profileID"
const regionPath = "region"
const apiKeyHeader = "X-Api-Key"
const serviceUserIDHeader = "X-User-ID"
const uuidSubjectKey = "uuidSubject"

func subjectUUIDFromContext(r *http.Request) (uuid.UUID, error) {
	userID, ok := r.Context().Value(uuidSubjectKey).(uuid.UUID)
	if !ok {
		return uuid.Nil, errors.New("missing user subject")
	}

	return userID, nil
}
