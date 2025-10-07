package services

import (
	"fmt"
	"testing"
	"time"

	"github.com/DevPulseLab/salat/internal/db/models"
	"github.com/DevPulseLab/salat/internal/db/repositories"
	"github.com/DevPulseLab/salat/internal/dto"
	"github.com/DevPulseLab/salat/internal/enum"
	"github.com/DevPulseLab/salat/internal/helper"
	"github.com/google/go-cmp/cmp"
	"github.com/uniplaces/carbon"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func getTestDb(t *testing.T) *gorm.DB {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		t.Fatalf("failed to connect database: %v", err)
	}
	db.AutoMigrate(&models.Calendar{})
	return db
}

func TestAddCalendarEntry(t *testing.T) {
	carbon.Freeze(time.Date(2025, time.May, 1, 0, 0, 0, 0, time.Local))
	defer carbon.UnFreeze()

	db := getTestDb(t)

	calendarRepo := repositories.NewCalendarRepository(db)
	sut := NewCalendarService(calendarRepo, helper.NewDateHelper())

	tests := []struct {
		user           models.User
		startDate      time.Time
		endDate        time.Time
		closeIntervals []dto.CloseInterval
		expected       []models.Calendar
	}{
		{
			user:           models.User{Model: gorm.Model{ID: 1}},
			startDate:      parseDate("2025-05-01"),
			endDate:        parseDate("2025-05-01"),
			closeIntervals: []dto.CloseInterval{},
			expected: []models.Calendar{
				{Date: parseDate("2025-05-01"), Status: string(enum.Rejected)},
			},
		},
		{
			user:           models.User{Model: gorm.Model{ID: 1}},
			startDate:      parseDate("2025-05-02"),
			endDate:        parseDate("2025-05-03"),
			closeIntervals: []dto.CloseInterval{},
			expected: []models.Calendar{
				{Date: parseDate("2025-05-02"), Status: string(enum.Reserved)},
			},
		},
		{
			user:           models.User{Model: gorm.Model{ID: 1}},
			startDate:      parseDate("2025-05-05"),
			endDate:        parseDate("2025-05-08"),
			closeIntervals: []dto.CloseInterval{},
			expected: []models.Calendar{
				{Date: parseDate("2025-05-05"), Status: string(enum.Approved)},
				{Date: parseDate("2025-05-06"), Status: string(enum.Approved)},
				{Date: parseDate("2025-05-07"), Status: string(enum.Approved)},
				{Date: parseDate("2025-05-08"), Status: string(enum.Approved)},
			},
		},
		{
			user:           models.User{Model: gorm.Model{ID: 1}},
			startDate:      parseDate("2025-05-12"),
			endDate:        parseDateTime("2025-05-15 10:00:00"),
			closeIntervals: []dto.CloseInterval{},
			expected: []models.Calendar{
				{Date: parseDate("2025-05-12"), Status: string(enum.Approved)},
				{Date: parseDate("2025-05-13"), Status: string(enum.Approved)},
				{Date: parseDate("2025-05-14"), Status: string(enum.Approved)},
				{Date: parseDate("2025-05-15"), Status: string(enum.Approved)},
			},
		},
		{
			user:           models.User{Model: gorm.Model{ID: 1}},
			startDate:      parseDateTime("2025-05-19 00:00:00"),
			endDate:        parseDateTime("2025-05-21 00:00:00"),
			closeIntervals: []dto.CloseInterval{},
			expected: []models.Calendar{
				{Date: parseDate("2025-05-19"), Status: string(enum.Approved)},
				{Date: parseDate("2025-05-20"), Status: string(enum.Approved)},
				{Date: parseDate("2025-05-21"), Status: string(enum.Approved)},
			},
		},
		{
			user:           models.User{Model: gorm.Model{ID: 1}},
			startDate:      parseDateTime("2025-05-27 12:00:00"),
			endDate:        parseDateTime("2025-05-30 10:00:00"),
			closeIntervals: []dto.CloseInterval{},
			expected: []models.Calendar{
				{Date: parseDate("2025-05-27"), Status: string(enum.Approved)},
				{Date: parseDate("2025-05-28"), Status: string(enum.Approved)},
				{Date: parseDate("2025-05-29"), Status: string(enum.Approved)},
				{Date: parseDate("2025-05-30"), Status: string(enum.Approved)},
			},
		},
		{
			user:           models.User{Model: gorm.Model{ID: 1}},
			startDate:      parseDateTime("2025-05-22 00:00:00"),
			endDate:        parseDateTime("2025-05-26 00:00:00"),
			closeIntervals: []dto.CloseInterval{},
			expected: []models.Calendar{
				{Date: parseDate("2025-05-22"), Status: string(enum.Approved)},
				{Date: parseDate("2025-05-23"), Status: string(enum.Approved)},
				{Date: parseDate("2025-05-26"), Status: string(enum.Approved)},
			},
		},
		{
			user:      models.User{Model: gorm.Model{ID: 1}},
			startDate: parseDateTime("2025-06-02 00:00:00"),
			endDate:   parseDateTime("2025-06-05 00:00:00"),
			closeIntervals: []dto.CloseInterval{
				{Id: 1, StartDate: parseDate("2025-06-03"), EndDate: parseDate("2025-06-04")},
			},
			expected: []models.Calendar{
				{Date: parseDate("2025-06-02"), Status: string(enum.Rejected)},
				{Date: parseDate("2025-06-05"), Status: string(enum.Rejected)},
			},
		},
		{
			user:           models.User{Model: gorm.Model{ID: 2}, PenaltyCard: string(enum.Yellow)},
			startDate:      parseDateTime("2025-05-22 00:00:00"),
			endDate:        parseDateTime("2025-05-23 00:00:00"),
			closeIntervals: []dto.CloseInterval{},
			expected: []models.Calendar{
				{Date: parseDate("2025-05-22"), Status: string(enum.Reserved)},
				{Date: parseDate("2025-05-23"), Status: string(enum.Reserved)},
			},
		},
		{
			user:           models.User{Model: gorm.Model{ID: 2}, PenaltyCard: string(enum.Red)},
			startDate:      parseDateTime("2025-05-26 00:00:00"),
			endDate:        parseDateTime("2025-05-27 00:00:00"),
			closeIntervals: []dto.CloseInterval{},
			expected: []models.Calendar{
				{Date: parseDate("2025-05-26"), Status: string(enum.Rejected)},
				{Date: parseDate("2025-05-27"), Status: string(enum.Rejected)},
			},
		},
	}

	for testNum, test := range tests {
		t.Run(fmt.Sprintf("test userId: %v startDate: %v endDate: %v", test.user.ID, test.startDate, test.endDate), func(t *testing.T) {
			result, errors := sut.AddCalendarEntries(&test.user, test.startDate, test.endDate, test.closeIntervals)
			if len(errors) != 0 {
				t.Errorf("Return errors: %v", errors)
			}

			if len(result) != len(test.expected) {
				t.Errorf("The result and the expected length are not equal %v != %v", len(test.expected), len(result))
			}

			prepared := []models.Calendar{}
			for _, item := range result {
				prepared = append(prepared, models.Calendar{Date: item.Date, Status: item.Status})
			}

			if !cmp.Equal(prepared, test.expected) {
				t.Errorf("Wrong result %v, actual: %v", testNum, prepared)
			}
		})
	}
}

func parseDate(dateString string) time.Time {
	res, _ := time.ParseInLocation("2006-01-02", dateString, time.Local)
	return res
}

func parseDateTime(dateString string) time.Time {
	res, _ := time.ParseInLocation("2006-01-02 15:04:05", dateString, time.Local)
	return res
}
