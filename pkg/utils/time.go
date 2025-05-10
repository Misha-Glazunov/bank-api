package utils

import "time"

func ParseRFC3339(dateStr string) (time.Time, error) {
    return time.Parse(time.RFC3339, dateStr)
}

func BeginningOfMonth(t time.Time) time.Time {
    return time.Date(t.Year(), t.Month(), 1, 0, 0, 0, 0, t.Location())
}

func EndOfMonth(t time.Time) time.Time {
    return BeginningOfMonth(t).AddDate(0, 1, -1)
}
