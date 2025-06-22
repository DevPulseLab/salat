package repositories

import (
	"fmt"
	"testing"
	"time"

	"github.com/DevPulseLab/salat/internal/db/models"
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

func TestGetByIdForUserIdSuccess(t *testing.T) {
	db := getTestDb(t)

	repo := NewCalendarRepository(db, helper.NewDateHelper())

	db.Create(&models.Calendar{
		UserId: 10,
		Date:   time.Now(),
		Status: string(enum.Approved),
	})
	db.Create(&models.Calendar{
		UserId: 12,
		Date:   time.Now(),
		Status: string(enum.Rejected),
	})

	tests := []struct {
		id     uint
		userId uint
		status enum.CalendarStatus
	}{
		{
			id:     1,
			userId: 10,
			status: enum.Approved,
		},
		{
			id:     2,
			userId: 12,
			status: enum.Rejected,
		},
	}

	for _, test := range tests {
		t.Run(fmt.Sprintf("test userId: %v", test.userId), func(t *testing.T) {
			result, err := repo.GetByIdForUserId(test.id, test.userId)
			if err != nil {
				t.Errorf("Expected model, return error: %v", err)
			}

			if result.UserId != test.userId {
				t.Errorf("Expected model with userId: %v, return %v", test.userId, result.UserId)
			}
		})
	}
}

func TestGetByIdForUserIdFailed(t *testing.T) {
	db := getTestDb(t)

	repo := NewCalendarRepository(db, helper.NewDateHelper())

	t.Run("test record not found", func(t *testing.T) {
		_, err := repo.GetByIdForUserId(1, 10)
		if err == nil {
			t.Errorf("Expected error")
		}

		if err.Error() != "record not found" {
			t.Errorf("Expected error: error not found, return: %v", err)
		}
	})
}

func TestRemove(t *testing.T) {
	db := getTestDb(t)

	repo := NewCalendarRepository(db, helper.NewDateHelper())

	db.Create(&models.Calendar{
		UserId: 10,
		Date:   time.Now(),
		Status: string(enum.Approved),
	})

	result, err := repo.GetByIdForUserId(1, 10)
	if err != nil {
		t.Errorf("Expected model, return error: %v", err)
	}

	repo.Remove(&result)

	_, err = repo.GetByIdForUserId(1, 10)
	if err == nil {
		t.Errorf("Expected not found")
	}
}

func TestAddCalendarEntry(t *testing.T) {
	carbon.Freeze(time.Date(2025, time.May, 1, 0, 0, 0, 0, time.Local))
	defer carbon.UnFreeze()

	db := getTestDb(t)

	repo := NewCalendarRepository(db, helper.NewDateHelper())

	tests := []struct {
		userId         uint
		startDate      time.Time
		endDate        time.Time
		closeIntervals []dto.CloseInterval
		expected       []models.Calendar
	}{
		{
			userId:         1,
			startDate:      parseDate("2025-05-01"),
			endDate:        parseDate("2025-05-01"),
			closeIntervals: []dto.CloseInterval{},
			expected: []models.Calendar{
				{Date: parseDate("2025-05-01"), Status: string(enum.Rejected)},
			},
		},
		{
			userId:         1,
			startDate:      parseDate("2025-05-02"),
			endDate:        parseDate("2025-05-03"),
			closeIntervals: []dto.CloseInterval{},
			expected: []models.Calendar{
				{Date: parseDate("2025-05-02"), Status: string(enum.Reserved)},
			},
		},
		{
			userId:         1,
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
			userId:         1,
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
			userId:         1,
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
			userId:         1,
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
			userId:         1,
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
			userId:    1,
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
	}

	for testNum, test := range tests {
		t.Run(fmt.Sprintf("test userId: %v startDate: %v endDate: %v", test.userId, test.startDate, test.endDate), func(t *testing.T) {
			result, errors := repo.AddCalendarEntry(test.userId, test.startDate, test.endDate, test.closeIntervals)
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

func TestGetCalendarEntriesByUserId(t *testing.T) {
	db := getTestDb(t)

	repo := NewCalendarRepository(db, helper.NewDateHelper())

	db.Create(&models.Calendar{
		UserId: 10,
		Date:   parseDate("2025-05-10"),
		Status: string(enum.Approved),
	})
	db.Create(&models.Calendar{
		UserId: 10,
		Date:   parseDate("2025-05-12"),
		Status: string(enum.Approved),
	})
	db.Create(&models.Calendar{
		UserId: 10,
		Date:   parseDate("2025-05-13"),
		Status: string(enum.Approved),
	})
	db.Create(&models.Calendar{
		UserId: 10,
		Date:   parseDate("2025-05-14"),
		Status: string(enum.Approved),
	})

	db.Create(&models.Calendar{
		UserId: 12,
		Date:   parseDate("2025-05-13"),
		Status: string(enum.Rejected),
	})

	tests := []struct {
		UserId    uint
		StartDate time.Time
		EndDate   time.Time
		Expected  int
	}{
		{10, parseDate("2025-05-10"), parseDate("2025-05-14"), 4},
		{12, parseDate("2025-05-10"), parseDate("2025-05-20"), 1},
		{1, parseDate("2025-05-10"), parseDate("2025-05-20"), 0},
	}

	for _, test := range tests {
		result := repo.GetCalendarEntriesByUserId(test.UserId, test.StartDate, test.EndDate)
		if len(result) != test.Expected {
			t.Errorf("Result len must be %v", test.Expected)
		}
	}
}

func TestGetCalendarEntriesForAllUsers(t *testing.T) {
	db := getTestDb(t)

	repo := NewCalendarRepository(db, helper.NewDateHelper())

	db.Create(&models.Calendar{
		UserId: 10,
		Date:   parseDate("2025-05-10"),
		Status: string(enum.Approved),
	})
	db.Create(&models.Calendar{
		UserId: 10,
		Date:   parseDate("2025-05-12"),
		Status: string(enum.Approved),
	})
	db.Create(&models.Calendar{
		UserId: 12,
		Date:   parseDate("2025-05-13"),
		Status: string(enum.Rejected),
	})

	result := repo.GetCalendarEntriesForAllUsers(parseDate("2025-05-10"), parseDate("2025-05-14"))
	if len(result) != 3 {
		t.Errorf("Result len must be 3")
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
