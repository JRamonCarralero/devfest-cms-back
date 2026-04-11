package utils

import (
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
)

// ParseIntOrDefault parses a string to an int and returns the default value if the parsing fails
func ParseIntOrDefault(value string, defaultValue int, min int) int {
	if value == "" {
		return defaultValue
	}

	i, err := strconv.Atoi(value)
	if err != nil || i < min {
		return defaultValue
	}

	return i
}

// GetPaginationParams returns the pagination parameters from the request
func GetPaginationParams(c *gin.Context) (page, pageSize int) {
	page = ParseIntOrDefault(c.Query("page"), 1, 1)
	pageSize = ParseIntOrDefault(c.Query("pageSize"), 10, 1)
	return
}

// uuidPtrToString converts pgtype.UUID to *string
func UuidPtrToString(u pgtype.UUID) *string {
	if !u.Valid {
		return nil
	}

	id, err := uuid.FromBytes(u.Bytes[:])
	if err != nil {
		return nil
	}

	s := id.String()
	return &s
}

// stringToNullUUID converts *string to pgtype.UUID
func StringToNullUUID(s *string) pgtype.UUID {
	if s == nil || *s == "" {
		return pgtype.UUID{Valid: false}
	}

	id, err := uuid.Parse(*s)
	if err != nil {
		return pgtype.UUID{Valid: false}
	}

	return pgtype.UUID{
		Bytes: [16]byte(id),
		Valid: true,
	}
}

// ToPgBool converts a *bool to pgtype.Bool
func ToPgBool(b *bool) pgtype.Bool {
	if b == nil {
		return pgtype.Bool{Valid: false}
	}
	return pgtype.Bool{Bool: *b, Valid: true}
}

// TextToPtr converts pgtype.Text to *string
func TextToPtr(t pgtype.Text) *string {
	if !t.Valid {
		return nil
	}
	s := t.String
	return &s
}

// PtrToText converts a *string to pgtype.Text
func PtrToText(s *string) pgtype.Text {
	if s == nil {
		return pgtype.Text{Valid: false}
	}
	return pgtype.Text{String: *s, Valid: true}
}

// PgStringToText converts pgtype.Text to string
func PgStringToText(t pgtype.Text) string {
	if !t.Valid {
		return ""
	}
	return t.String
}

// TextToPgString converts string to pgtype.Text
func TextToPgString(s string) pgtype.Text {
	return pgtype.Text{String: s, Valid: true}
}

// SafeString convets a *string to a string
func SafeString(s *string) string {
	if s == nil {
		return ""
	}
	return *s
}

// PtrToString converts a *string to a string
func PtrToString(s *string) string {
	if s == nil {
		return ""
	}
	return *s
}

func SafeInt(i *int) int {
	if i == nil {
		return 0
	}
	return *i
}

func Int4ToPtr(i pgtype.Int4) *int {
	if !i.Valid {
		return nil
	}
	val := int(i.Int32)
	return &val
}

func PtrToInt4(i *int) pgtype.Int4 {
	if i == nil {
		return pgtype.Int4{Valid: false}
	}
	return pgtype.Int4{Int32: int32(*i), Valid: true}
}

// ToPgDate
func ToPgDate(t time.Time) pgtype.Date {
	return pgtype.Date{
		Time:  t,
		Valid: !t.IsZero(),
	}
}

// ToPgTimestamp
func ToPgTimestamp(t time.Time) pgtype.Timestamp {
	return pgtype.Timestamp{
		Time:  t,
		Valid: !t.IsZero(),
	}
}

// PgTimeToTime
func PgTimeToTime(p pgtype.Time) time.Time {
	if !p.Valid {
		return time.Time{}
	}
	return time.Time{}.Add(time.Duration(p.Microseconds) * time.Microsecond)
}

// TimeToPgTime
func TimeToPgTime(t time.Time) pgtype.Time {
	if t.IsZero() {
		return pgtype.Time{Valid: false}
	}
	midnight := time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, t.Location())
	durationSinceMidnight := t.Sub(midnight)

	return pgtype.Time{
		Microseconds: durationSinceMidnight.Microseconds(),
		Valid:        true,
	}
}

// PgIntervalToDuration
func PgIntervalToDuration(p pgtype.Interval) time.Duration {
	if !p.Valid {
		return 0
	}

	const (
		hoursInDay = 24
	)
	totalMicroseconds := p.Microseconds

	if p.Days > 0 {
		totalMicroseconds += int64(p.Days) * hoursInDay * int64(time.Hour/time.Microsecond)
	}

	return time.Duration(totalMicroseconds) * time.Microsecond
}

// DurationToPgInterval
func DurationToPgInterval(d time.Duration) pgtype.Interval {
	if d == 0 {
		return pgtype.Interval{Valid: false}
	}

	return pgtype.Interval{
		Microseconds: d.Microseconds(),
		Valid:        true,
	}
}

// PgUUIDToUUID
func PgUUIDToUUID(p pgtype.UUID) uuid.UUID {
	if !p.Valid {
		return uuid.UUID{}
	}
	return uuid.UUID(p.Bytes)
}

// UUIDToPgUUID
func UUIDToPgUUID(u uuid.UUID) pgtype.UUID {
	return pgtype.UUID{
		Bytes: [16]byte(u),
		Valid: true,
	}
}
