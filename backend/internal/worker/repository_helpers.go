package worker

import "github.com/jackc/pgx/v5/pgtype"

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
