package identity

import (
	"time"
)

func parseTime(s string) time.Time {
	t, err := time.Parse(time.RFC3339Nano, s)
	if err != nil {
		return time.Time{}
	}
	return t
}

func parseTimePtr(s *string) *time.Time {
	if s == nil {
		return nil
	}
	t := parseTime(*s)
	return &t
}

func timePtrToStringPtr(t *time.Time) *string {
	if t == nil {
		return nil
	}
	s := t.Format(time.RFC3339Nano)
	return &s
}

func timeToString(t time.Time) string {
	return t.Format(time.RFC3339Nano)
}

func int64PtrFromUint8Ptr(v *uint8) any {
	if v == nil {
		return nil
	}
	i := int64(*v)
	return i
}

func stringPtrFromString(s string) *string {
	if s == "" {
		return nil
	}
	return &s
}

func stringFromStringPtr(s *string) string {
	if s == nil {
		return ""
	}
	return *s
}
