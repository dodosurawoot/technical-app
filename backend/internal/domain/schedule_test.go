package domain

import (
	"testing"
	"time"

	"airclean-tracker/backend/internal/models"
)

func d(v string) *time.Time {
	t, err := time.Parse("2006-01-02", v)
	if err != nil {
		panic(err)
	}
	return &t
}

func TestNextCleaningDateAddsSixMonths(t *testing.T) {
	got := NextCleaningDate(d("2026-06-10"), 6)
	if got.Format("2006-01-02") != "2026-12-10" {
		t.Fatalf("expected 2026-12-10, got %s", got.Format("2006-01-02"))
	}
}

func TestStatusNeverCleaned(t *testing.T) {
	if got := Status(nil, nil, nil, *d("2026-07-02")); got != models.StatusNeverCleaned {
		t.Fatalf("expected %s, got %s", models.StatusNeverCleaned, got)
	}
}

func TestStatusOverdue(t *testing.T) {
	if got := Status(d("2025-12-01"), d("2026-06-01"), nil, *d("2026-07-02")); got != models.StatusOverdue {
		t.Fatalf("expected %s, got %s", models.StatusOverdue, got)
	}
}

func TestStatusDueSoon(t *testing.T) {
	if got := Status(d("2026-01-20"), d("2026-07-20"), nil, *d("2026-07-02")); got != models.StatusDueSoon {
		t.Fatalf("expected %s, got %s", models.StatusDueSoon, got)
	}
}
