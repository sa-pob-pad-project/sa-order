package utils

import (
	"time"
)

func ParseNullableTime(value *string) *time.Time {
	if value == nil {
		return nil
	}
	t, err := time.Parse(time.RFC3339, *value)
	if err != nil {
		return nil
	}
	return &t
}
