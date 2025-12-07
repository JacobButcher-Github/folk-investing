package util

import (
	"database/sql"
	"time"
)

// Helper to assign a NullInt64
func NullInt64(p *int64) sql.NullInt64 {
	if p == nil {
		return sql.NullInt64{Valid: false}
	}
	return sql.NullInt64{Int64: *p, Valid: true}
}

// Helper to assign a NullString
func NullString(p *string) sql.NullString {
	if p == nil {
		return sql.NullString{Valid: false}
	}
	return sql.NullString{String: *p, Valid: true}
}

func NullTime(p *time.Time) sql.NullTime {
	if p == nil {
		return sql.NullTime{Valid: false}
	}
	return sql.NullTime{Time: *p, Valid: true}
}
