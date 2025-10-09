package repositories

import (
	"testing"
	"time"

	"github.com/DevPulseLab/salat/internal/db/models"
	"github.com/DevPulseLab/salat/internal/db/repositories/testutils"
	"github.com/stretchr/testify/assert"
)

func TestToggleVisit_NewEntry(t *testing.T) {
	db := testutils.GetTestDb(t, &models.VisitStats{})
	repo := NewVisitStatsRepository(db)

	day := time.Now().Truncate(24 * time.Hour)

	entry, err := repo.ToggleVisit(1, day)
	assert.NoError(t, err)
	assert.Equal(t, uint(1), entry.UserId)
	assert.True(t, entry.IsVisit)

	var dbEntry models.VisitStats
	res := db.Where("user_id = ? AND DATE(date) = ?", 1, day.Format("2006-01-02")).First(&dbEntry)
	assert.NoError(t, res.Error)
	assert.True(t, dbEntry.IsVisit)
}

func TestToggleVisit_ExistingEntry(t *testing.T) {
	db := testutils.GetTestDb(t, &models.VisitStats{})
	repo := NewVisitStatsRepository(db)

	day := time.Now().Truncate(24 * time.Hour)
	db.Create(&models.VisitStats{UserId: 1, Date: day, IsVisit: true})

	entry, err := repo.ToggleVisit(1, day)
	assert.NoError(t, err)
	assert.False(t, entry.IsVisit)

	var updated models.VisitStats
	db.Where("user_id = ? AND DATE(date) = ?", 1, day.Format("2006-01-02")).First(&updated)
	assert.False(t, updated.IsVisit)
}

func TestGetVisitVisit(t *testing.T) {
	db := testutils.GetTestDb(t, &models.VisitStats{})
	repo := NewVisitStatsRepository(db)

	now := time.Now().Truncate(24 * time.Hour)
	data := []models.VisitStats{
		{UserId: 1, Date: now.AddDate(0, 0, -2), IsVisit: true}, // zu früh
		{UserId: 2, Date: now.AddDate(0, 0, -1), IsVisit: true}, // im Bereich
		{UserId: 3, Date: now, IsVisit: true},                   // im Bereich
		{UserId: 4, Date: now.AddDate(0, 0, 1), IsVisit: true},  // zu spät
	}
	for _, v := range data {
		db.Create(&v)
	}

	start := now.AddDate(0, 0, -1)
	end := now
	results := repo.GetVisitVisit(start, end)

	assert.Len(t, results, 2, "should only return visits within the date range")
	for _, v := range results {
		assert.True(t, v.Date.After(start.Add(-time.Second)) && v.Date.Before(end.Add(24*time.Hour)))
	}
}
