package repositories

import (
	"testing"
	"time"

	"github.com/DevPulseLab/salat/internal/db/models"
	"github.com/DevPulseLab/salat/internal/db/repositories/testutils"
	"github.com/stretchr/testify/assert"
)

func TestSaveCloseInterval(t *testing.T) {
	db := testutils.GetTestDb(t, &models.CloseInterval{})
	repo := NewCloseIntervalsRepository(db)

	start := time.Now().Add(-time.Hour)
	end := time.Now().Add(time.Hour)

	err := repo.SaveCloseInterval(start, end)
	assert.NoError(t, err)

	var result models.CloseInterval
	res := db.First(&result)
	assert.NoError(t, res.Error)
	assert.WithinDuration(t, start, result.StartDate, time.Second)
	assert.WithinDuration(t, end, result.EndDate, time.Second)
}

func TestGetAllEntriesForInterval(t *testing.T) {
	db := testutils.GetTestDb(t, &models.CloseInterval{})
	repo := NewCloseIntervalsRepository(db)

	// Seed test data
	now := time.Now()
	intervals := []models.CloseInterval{
		{StartDate: now.Add(-3 * time.Hour), EndDate: now.Add(-2 * time.Hour)},      // before
		{StartDate: now.Add(-30 * time.Minute), EndDate: now.Add(30 * time.Minute)}, // overlap
		{StartDate: now.Add(2 * time.Hour), EndDate: now.Add(3 * time.Hour)},        // after
	}
	for _, i := range intervals {
		db.Create(&i)
	}

	results := repo.GetAllEntriesForInterval(now.Add(-time.Hour), now.Add(time.Hour))
	assert.Len(t, results, 1)
	assert.WithinDuration(t, intervals[1].StartDate, results[0].StartDate, time.Second)
}

func TestGetById(t *testing.T) {
	db := testutils.GetTestDb(t, &models.CloseInterval{})
	repo := NewCloseIntervalsRepository(db)

	entry := models.CloseInterval{
		StartDate: time.Now(),
		EndDate:   time.Now().Add(time.Hour),
	}
	db.Create(&entry)

	found, err := repo.GetById(entry.ID)
	assert.NoError(t, err)
	assert.Equal(t, entry.ID, found.ID)
}

func TestGetById_NotFound(t *testing.T) {
	db := testutils.GetTestDb(t, &models.CloseInterval{})
	repo := NewCloseIntervalsRepository(db)

	_, err := repo.GetById(999)
	assert.Error(t, err)
}

func TestRemoveInterval(t *testing.T) {
	db := testutils.GetTestDb(t, &models.CloseInterval{})
	repo := NewCloseIntervalsRepository(db)

	entry := models.CloseInterval{
		StartDate: time.Now(),
		EndDate:   time.Now().Add(time.Hour),
	}
	db.Create(&entry)

	repo.Remove(&entry)

	var count int64
	db.Model(&models.CloseInterval{}).Where("id = ?", entry.ID).Count(&count)
	assert.Equal(t, int64(0), count)
}
