package worker

import (
	"time"

	"github.com/jackc/pgx/v5/pgtype"
)

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

// timestamptzValue 实现当前函数行为。
func timestamptzValue(value time.Time) pgtype.Timestamptz {
	if value.IsZero() {
		return pgtype.Timestamptz{}
	}
	return pgtype.Timestamptz{
		Time:  value,
		Valid: true,
	}
}
