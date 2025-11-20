package repositories

import (
	"fmt"
	"testing"
	"time"

	"github.com/DevPulseLab/salat/internal/db/models"
	"github.com/DevPulseLab/salat/internal/db/repositories/testutils"
	"github.com/DevPulseLab/salat/internal/enum"
	"github.com/uniplaces/carbon"
)

func TestFindByIdForUserIdSuccess(t *testing.T) {
	db := testutils.GetTestDb(t, &models.Calendar{})

	repo := NewCalendarRepository(db)

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
			result, err := repo.FindByIdForUserId(test.id, test.userId)
			if err != nil {
				t.Errorf("Expected model, return error: %v", err)
			}

			if result.UserId != test.userId {
				t.Errorf("Expected model with userId: %v, return %v", test.userId, result.UserId)
			}
		})
	}
}

func TestFindByIdForUserIdFailed(t *testing.T) {
	db := testutils.GetTestDb(t, &models.Calendar{})

	repo := NewCalendarRepository(db)

	t.Run("test record not found", func(t *testing.T) {
		_, err := repo.FindByIdForUserId(1, 10)
		if err == nil {
			t.Errorf("Expected error")
		}

		if err.Error() != "record not found" {
			t.Errorf("Expected error: error not found, return: %v", err)
		}
	})
}

func TestRemove(t *testing.T) {
	db := testutils.GetTestDb(t, &models.Calendar{})

	repo := NewCalendarRepository(db)

	db.Create(&models.Calendar{
		UserId: 10,
		Date:   time.Now(),
		Status: string(enum.Approved),
	})

	result, err := repo.FindByIdForUserId(1, 10)
	if err != nil {
		t.Errorf("Expected model, return error: %v", err)
	}

	repo.SoftDelete(&result)

	_, err = repo.FindByIdForUserId(1, 10)
	if err == nil {
		t.Errorf("Expected not found")
	}
}

func TestGetCalendarEntriesByUserId(t *testing.T) {
	db := testutils.GetTestDb(t, &models.Calendar{})

	repo := NewCalendarRepository(db)

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
		result, err := repo.FindByUserIdAndDateRange(test.UserId, test.StartDate, test.EndDate)
		if err != nil {
			t.Errorf("Result has an error %s", err.Error())
		}

		if len(result) != test.Expected {
			t.Errorf("Result len must be %v", test.Expected)
		}
	}
}

func TestFindByDateRange(t *testing.T) {
	db := testutils.GetTestDb(t, &models.Calendar{})

	repo := NewCalendarRepository(db)

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

	result, _ := repo.FindByDateRange(parseDate("2025-05-10"), parseDate("2025-05-14"))
	if len(result) != 3 {
		t.Errorf("Result len must be 3")
	}
}

func TestUpdateStatus(t *testing.T) {
	db := testutils.GetTestDb(t, &models.Calendar{})

	repo := NewCalendarRepository(db)

	model := models.Calendar{
		UserId: 20,
		Date:   parseDate("2025-05-13"),
		Status: string(enum.Rejected),
	}
	repo.Create(&model)

	if model.Status != string(enum.Rejected) {
		t.Errorf("Result status must be 'approved'")
	}

	repo.UpdateStatus(model.ID, string(enum.Approved))
	result, _ := repo.FindByIdForUserId(model.ID, model.UserId)

	if result.Status != string(enum.Approved) {
		t.Errorf("Result status must be 'approved'")
	}
}

func TestCountReservedByDate(t *testing.T) {
	db := testutils.GetTestDb(t, &models.Calendar{})

	repo := NewCalendarRepository(db)

	model := models.Calendar{
		UserId: 20,
		Date:   parseDate("2025-05-14"),
		Status: string(enum.Reserved),
	}
	repo.Create(&model)

	result, err := repo.CountReservedByDate(carbon.NewCarbon(parseDate("2025-05-14")))
	if err != nil {
		t.Errorf("Error count reserved date: %v", err.Error())
	}

	if result != 1 {
		t.Errorf("Result status must be 1")
	}
}

func TestFindDeletedByUserIdAndDateRange(t *testing.T) {
	db := testutils.GetTestDb(t, &models.Calendar{})
	repo := NewCalendarRepository(db)

	// Active entry (should not be found)
	db.Create(&models.Calendar{
		UserId: 10,
		Date:   parseDate("2025-05-10"),
		Status: string(enum.Approved),
	})

	// Soft-deleted entry within range (should be found)
	deletedInRange := models.Calendar{
		UserId: 10,
		Date:   parseDate("2025-05-12"),
		Status: string(enum.Approved),
	}
	db.Create(&deletedInRange)
	repo.SoftDelete(&deletedInRange)

	// Soft-deleted entry outside range (should not be found)
	deletedOutsideRange := models.Calendar{
		UserId: 10,
		Date:   parseDate("2025-05-20"),
		Status: string(enum.Approved),
	}
	db.Create(&deletedOutsideRange)
	repo.SoftDelete(&deletedOutsideRange)

	// Soft-deleted entry for another user (should not be found)
	deletedOtherUser := models.Calendar{
		UserId: 12,
		Date:   parseDate("2025-05-12"),
		Status: string(enum.Approved),
	}
	db.Create(&deletedOtherUser)
	repo.SoftDelete(&deletedOtherUser)

	results, err := repo.FindDeletedByUserIdAndDateRange(10, parseDate("2025-05-10"), parseDate("2025-05-15"))
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}

	if len(results) != 1 {
		t.Errorf("Expected 1 result, got %d", len(results))
	}

	if len(results) > 0 && results[0].ID != deletedInRange.ID {
		t.Errorf("Expected found entry ID to be %d, got %d", deletedInRange.ID, results[0].ID)
	}
}

func parseDate(dateString string) time.Time {
	res, _ := time.Parse("2006-01-02", dateString)
	return res
}
