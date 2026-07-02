package domain

import (
	"time"

	"airclean-tracker/backend/internal/models"
)

const (
	DefaultIntervalMonths = 6
	DueSoonDays           = 30
)

func NextCleaningDate(latest *time.Time, intervalMonths int) *time.Time {
	if latest == nil {
		return nil
	}
	if intervalMonths <= 0 {
		intervalMonths = DefaultIntervalMonths
	}
	next := latest.AddDate(0, intervalMonths, 0)
	return &next
}

func Status(latest, next, planned *time.Time, now time.Time) string {
	today := truncateDate(now)
	if latest == nil {
		return models.StatusNeverCleaned
	}
	if next == nil {
		next = NextCleaningDate(latest, DefaultIntervalMonths)
	}
	if next.Before(today) {
		return models.StatusOverdue
	}
	if !next.After(today.AddDate(0, 0, DueSoonDays)) {
		return models.StatusDueSoon
	}
	if planned != nil {
		return models.StatusPlanned
	}
	return models.StatusNormal
}

func Recalculate(latest, planned *time.Time, now time.Time) (*time.Time, string) {
	next := NextCleaningDate(latest, DefaultIntervalMonths)
	return next, Status(latest, next, planned, now)
}

func truncateDate(t time.Time) time.Time {
	y, m, d := t.Date()
	return time.Date(y, m, d, 0, 0, 0, 0, t.Location())
}
