package task

import (
	"encoding/json"

	"github.com/jackc/pgx/v5/pgtype"
)

func jsonValue(value json.RawMessage, fallback []byte) []byte {
	if len(value) == 0 {
		return fallback
	}
	return []byte(value)
}

func textValue(value string) pgtype.Text {
	if value == "" {
		return pgtype.Text{}
	}
	return pgtype.Text{String: value, Valid: true}
}

func stringValue(value pgtype.Text) string {
	if !value.Valid {
		return ""
	}
	return value.String
}
