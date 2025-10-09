package repositories

import (
	"testing"
	"time"

	"github.com/DevPulseLab/salat/internal/db/models"
	"github.com/DevPulseLab/salat/internal/db/repositories/testutils"
	"github.com/stretchr/testify/assert"
)

func TestIncrementStatsForDay_NewEntry(t *testing.T) {
	db := testutils.GetTestDb(t, &models.RealDayStats{})
	repo := NewRealDayStatsRepository(db)

	day := time.Now().Truncate(24 * time.Hour)
	ok := repo.IncrementStatsForDay(day)
	assert.True(t, ok)

	var entry models.RealDayStats
	res := db.Where("DATE(date) = ?", day.Format("2006-01-02")).First(&entry)
	assert.NoError(t, res.Error)
	assert.Equal(t, uint(1), entry.NumberOfPlates)
}

func TestIncrementStatsForDay_ExistingEntry(t *testing.T) {
	db := testutils.GetTestDb(t, &models.RealDayStats{})
	repo := NewRealDayStatsRepository(db)

	day := time.Now().Truncate(24 * time.Hour)
	db.Create(&models.RealDayStats{Date: day, NumberOfPlates: 5})

	ok := repo.IncrementStatsForDay(day)
	assert.True(t, ok)

	var entry models.RealDayStats
	res := db.Where("DATE(date) = ?", day.Format("2006-01-02")).First(&entry)
	assert.NoError(t, res.Error)
	assert.Equal(t, uint(6), entry.NumberOfPlates)
}

func TestSaveStatsForDay_NewEntry(t *testing.T) {
	db := testutils.GetTestDb(t, &models.RealDayStats{})
	repo := NewRealDayStatsRepository(db)

	day := time.Now().Truncate(24 * time.Hour)
	ok := repo.SaveStatsForDay(day, 10)
	assert.True(t, ok)

	var entry models.RealDayStats
	res := db.Where("DATE(date) = ?", day.Format("2006-01-02")).First(&entry)
	assert.NoError(t, res.Error)
	assert.Equal(t, uint(10), entry.NumberOfPlates)
}

func TestSaveStatsForDay_UpdateExisting(t *testing.T) {
	db := testutils.GetTestDb(t, &models.RealDayStats{})
	repo := NewRealDayStatsRepository(db)

	day := time.Now().Truncate(24 * time.Hour)
	db.Create(&models.RealDayStats{Date: day, NumberOfPlates: 3})

	ok := repo.SaveStatsForDay(day, 12)
	assert.True(t, ok)

	var entry models.RealDayStats
	res := db.Where("DATE(date) = ?", day.Format("2006-01-02")).First(&entry)
	assert.NoError(t, res.Error)
	assert.Equal(t, uint(12), entry.NumberOfPlates)
}

func TestGetStatsForDay_Found(t *testing.T) {
	db := testutils.GetTestDb(t, &models.RealDayStats{})
	repo := NewRealDayStatsRepository(db)

	day := time.Now().Truncate(24 * time.Hour)
	db.Create(&models.RealDayStats{Date: day, NumberOfPlates: 8})

	val := repo.GetStatsForDay(day)
	assert.Equal(t, uint(8), val)
}

func TestGetStatsForDay_NotFound(t *testing.T) {
	db := testutils.GetTestDb(t, &models.RealDayStats{})
	repo := NewRealDayStatsRepository(db)

	day := time.Now().Truncate(24 * time.Hour)
	val := repo.GetStatsForDay(day)
	assert.Equal(t, uint(0), val)
}
