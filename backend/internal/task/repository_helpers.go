package task

import (
	"encoding/json"

	"github.com/jackc/pgx/v5/pgtype"
)

// jsonValue 实现当前函数行为。
func jsonValue(value json.RawMessage, fallback []byte) []byte {
	if len(value) == 0 {
		return fallback
	}
	return []byte(value)
}

// textValue 实现当前函数行为。
func textValue(value string) pgtype.Text {
	if value == "" {
		return pgtype.Text{}
	}
	return pgtype.Text{String: value, Valid: true}
}

// stringValue 实现当前函数行为。
func stringValue(value pgtype.Text) string {
	if !value.Valid {
		return ""
	}
	return value.String
}
