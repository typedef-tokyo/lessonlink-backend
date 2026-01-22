package utility

import (
	"time"

	"github.com/aarondl/null/v8"
)

func ConvertNullTimeToPointer(t null.Time) *time.Time {
	if t.Valid {
		return &t.Time
	}
	return nil
}

func ConvertStringToNullString(s string) null.String {
	if s == "" {
		return null.NewString("", false)
	}
	return null.NewString(s, true)
}

func ConvertNullStringToString(ns null.String) string {
	if !ns.Valid {
		return ""
	}
	return ns.String
}
