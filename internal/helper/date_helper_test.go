package helper

import (
	"testing"
	"time"

	"github.com/DevPulseLab/salat/internal/dto"
	"github.com/uniplaces/carbon"
)

func TestIsDateInCurrentWeek(t *testing.T) {
	dateHelper := NewDateHelper()

	carbon.Freeze(time.Date(2025, time.March, 6, 10, 0, 0, 0, time.UTC))
	defer carbon.UnFreeze()

	now := carbon.Now()

	tests := []struct {
		inputDate time.Time
		expected  bool
	}{
		{now.AddDate(0, 0, -4), false}, // So, 2
		{now.AddDate(0, 0, -3), true},  // Mo, 3
		{now.AddDate(0, 0, -2), true},  // Di, 4
		{now.AddDate(0, 0, -1), true},  // Mi, 5
		{now.Time, true},               // Do, 6
		{now.AddDate(0, 0, 1), true},   // Fr, 7
		{now.AddDate(0, 0, 2), true},   // Sa, 8
		{now.AddDate(0, 0, 3), true},   // So, 9
		{now.AddDate(0, 0, 4), false},  // Mo, 10
	}

	for _, test := range tests {
		t.Run(test.inputDate.String(), func(t *testing.T) {
			result := dateHelper.IsDateInCurrentWeek(test.inputDate)
			if result != test.expected {
				t.Errorf("Expected %v, but got %v for date %v", test.expected, result, test.inputDate)
			}
		})
	}
}

func TestGetFridayOfWeek(t *testing.T) {
	dateHelper := NewDateHelper()

	tests := []struct {
		inputDate time.Time
		expected  time.Time
	}{
		{time.Date(2025, time.March, 1, 0, 0, 0, 0, time.UTC), time.Date(2025, time.February, 28, 0, 0, 0, 0, time.UTC)},
		{time.Date(2025, time.March, 2, 0, 0, 0, 0, time.UTC), time.Date(2025, time.February, 28, 0, 0, 0, 0, time.UTC)},
		{time.Date(2025, time.March, 3, 0, 0, 0, 0, time.UTC), time.Date(2025, time.March, 7, 0, 0, 0, 0, time.UTC)},
		{time.Date(2025, time.March, 4, 0, 0, 0, 0, time.UTC), time.Date(2025, time.March, 7, 0, 0, 0, 0, time.UTC)},
		{time.Date(2025, time.March, 5, 0, 0, 0, 0, time.UTC), time.Date(2025, time.March, 7, 0, 0, 0, 0, time.UTC)},
		{time.Date(2025, time.March, 6, 0, 0, 0, 0, time.UTC), time.Date(2025, time.March, 7, 0, 0, 0, 0, time.UTC)},
		{time.Date(2025, time.March, 7, 0, 0, 0, 0, time.UTC), time.Date(2025, time.March, 7, 0, 0, 0, 0, time.UTC)},
		{time.Date(2025, time.March, 8, 0, 0, 0, 0, time.UTC), time.Date(2025, time.March, 7, 0, 0, 0, 0, time.UTC)},
		{time.Date(2025, time.March, 9, 0, 0, 0, 0, time.UTC), time.Date(2025, time.March, 7, 0, 0, 0, 0, time.UTC)},
		{time.Date(2025, time.March, 10, 0, 0, 0, 0, time.UTC), time.Date(2025, time.March, 14, 0, 0, 0, 0, time.UTC)},
	}

	for _, test := range tests {
		t.Run(test.inputDate.String(), func(t *testing.T) {
			result := dateHelper.getFridayOfWeek(test.inputDate)
			if result != test.expected {
				t.Errorf("Expected %v, but got %v for date %v", test.expected, result, test.inputDate)
			}
		})
	}
}

func TestIsDateNextWeekAndNowAfterFriday(t *testing.T) {
	dateHelper := NewDateHelper()

	tests := []struct {
		inputDate time.Time
		nowDate   time.Time
		expected  bool
	}{
		// now 3 March
		{time.Date(2025, time.March, 3, 9, 0, 0, 0, time.UTC), time.Date(2025, time.March, 3, 10, 0, 0, 0, time.UTC), false},
		{time.Date(2025, time.March, 4, 9, 0, 0, 0, time.UTC), time.Date(2025, time.March, 3, 10, 0, 0, 0, time.UTC), false},
		{time.Date(2025, time.March, 5, 9, 0, 0, 0, time.UTC), time.Date(2025, time.March, 3, 10, 0, 0, 0, time.UTC), false},
		{time.Date(2025, time.March, 6, 9, 0, 0, 0, time.UTC), time.Date(2025, time.March, 3, 10, 0, 0, 0, time.UTC), false},
		{time.Date(2025, time.March, 7, 9, 0, 0, 0, time.UTC), time.Date(2025, time.March, 3, 10, 0, 0, 0, time.UTC), false},
		{time.Date(2025, time.March, 8, 9, 0, 0, 0, time.UTC), time.Date(2025, time.March, 3, 10, 0, 0, 0, time.UTC), false},
		{time.Date(2025, time.March, 9, 9, 0, 0, 0, time.UTC), time.Date(2025, time.March, 3, 10, 0, 0, 0, time.UTC), false},
		{time.Date(2025, time.March, 10, 9, 0, 0, 0, time.UTC), time.Date(2025, time.March, 3, 10, 0, 0, 0, time.UTC), false},
		{time.Date(2025, time.March, 11, 9, 0, 0, 0, time.UTC), time.Date(2025, time.March, 3, 10, 0, 0, 0, time.UTC), false},
		{time.Date(2025, time.March, 12, 9, 0, 0, 0, time.UTC), time.Date(2025, time.March, 3, 10, 0, 0, 0, time.UTC), false},

		// now 7 March before lunch
		{time.Date(2025, time.March, 3, 9, 0, 0, 0, time.UTC), time.Date(2025, time.March, 7, 10, 0, 0, 0, time.UTC), false},
		{time.Date(2025, time.March, 4, 9, 0, 0, 0, time.UTC), time.Date(2025, time.March, 7, 10, 0, 0, 0, time.UTC), false},
		{time.Date(2025, time.March, 5, 9, 0, 0, 0, time.UTC), time.Date(2025, time.March, 7, 10, 0, 0, 0, time.UTC), false},
		{time.Date(2025, time.March, 6, 9, 0, 0, 0, time.UTC), time.Date(2025, time.March, 7, 10, 0, 0, 0, time.UTC), false},
		{time.Date(2025, time.March, 7, 9, 0, 0, 0, time.UTC), time.Date(2025, time.March, 7, 10, 0, 0, 0, time.UTC), false},
		{time.Date(2025, time.March, 8, 9, 0, 0, 0, time.UTC), time.Date(2025, time.March, 7, 10, 0, 0, 0, time.UTC), false},
		{time.Date(2025, time.March, 9, 9, 0, 0, 0, time.UTC), time.Date(2025, time.March, 7, 10, 0, 0, 0, time.UTC), false},
		{time.Date(2025, time.March, 10, 9, 0, 0, 0, time.UTC), time.Date(2025, time.March, 7, 10, 0, 0, 0, time.UTC), false},
		{time.Date(2025, time.March, 11, 9, 0, 0, 0, time.UTC), time.Date(2025, time.March, 7, 10, 0, 0, 0, time.UTC), false},
		{time.Date(2025, time.March, 12, 9, 0, 0, 0, time.UTC), time.Date(2025, time.March, 7, 10, 0, 0, 0, time.UTC), false},

		// now 7 March after lunch
		{time.Date(2025, time.March, 3, 9, 0, 0, 0, time.UTC), time.Date(2025, time.March, 7, 22, 0, 0, 0, time.UTC), false},
		{time.Date(2025, time.March, 4, 9, 0, 0, 0, time.UTC), time.Date(2025, time.March, 7, 22, 0, 0, 0, time.UTC), false},
		{time.Date(2025, time.March, 5, 9, 0, 0, 0, time.UTC), time.Date(2025, time.March, 7, 22, 0, 0, 0, time.UTC), false},
		{time.Date(2025, time.March, 6, 9, 0, 0, 0, time.UTC), time.Date(2025, time.March, 7, 22, 0, 0, 0, time.UTC), false},
		{time.Date(2025, time.March, 7, 9, 0, 0, 0, time.UTC), time.Date(2025, time.March, 7, 22, 0, 0, 0, time.UTC), false},
		{time.Date(2025, time.March, 8, 9, 0, 0, 0, time.UTC), time.Date(2025, time.March, 7, 22, 0, 0, 0, time.UTC), false},
		{time.Date(2025, time.March, 9, 9, 0, 0, 0, time.UTC), time.Date(2025, time.March, 7, 22, 0, 0, 0, time.UTC), false},
		{time.Date(2025, time.March, 10, 9, 0, 0, 0, time.UTC), time.Date(2025, time.March, 7, 22, 0, 0, 0, time.UTC), true},
		{time.Date(2025, time.March, 11, 9, 0, 0, 0, time.UTC), time.Date(2025, time.March, 7, 22, 0, 0, 0, time.UTC), true},
		{time.Date(2025, time.March, 12, 9, 0, 0, 0, time.UTC), time.Date(2025, time.March, 7, 22, 0, 0, 0, time.UTC), true},

		// now 8 March
		{time.Date(2025, time.March, 3, 9, 0, 0, 0, time.UTC), time.Date(2025, time.March, 8, 10, 0, 0, 0, time.UTC), false},
		{time.Date(2025, time.March, 4, 9, 0, 0, 0, time.UTC), time.Date(2025, time.March, 8, 10, 0, 0, 0, time.UTC), false},
		{time.Date(2025, time.March, 5, 9, 0, 0, 0, time.UTC), time.Date(2025, time.March, 8, 10, 0, 0, 0, time.UTC), false},
		{time.Date(2025, time.March, 6, 9, 0, 0, 0, time.UTC), time.Date(2025, time.March, 8, 10, 0, 0, 0, time.UTC), false},
		{time.Date(2025, time.March, 7, 9, 0, 0, 0, time.UTC), time.Date(2025, time.March, 8, 10, 0, 0, 0, time.UTC), false},
		{time.Date(2025, time.March, 8, 9, 0, 0, 0, time.UTC), time.Date(2025, time.March, 8, 10, 0, 0, 0, time.UTC), false},
		{time.Date(2025, time.March, 9, 9, 0, 0, 0, time.UTC), time.Date(2025, time.March, 8, 10, 0, 0, 0, time.UTC), false},
		{time.Date(2025, time.March, 10, 9, 0, 0, 0, time.UTC), time.Date(2025, time.March, 8, 10, 0, 0, 0, time.UTC), true},
		{time.Date(2025, time.March, 11, 9, 0, 0, 0, time.UTC), time.Date(2025, time.March, 8, 10, 0, 0, 0, time.UTC), true},
		{time.Date(2025, time.March, 12, 9, 0, 0, 0, time.UTC), time.Date(2025, time.March, 8, 10, 0, 0, 0, time.UTC), true},
		{time.Date(2025, time.March, 13, 9, 0, 0, 0, time.UTC), time.Date(2025, time.March, 8, 10, 0, 0, 0, time.UTC), true},
		{time.Date(2025, time.March, 14, 9, 0, 0, 0, time.UTC), time.Date(2025, time.March, 8, 10, 0, 0, 0, time.UTC), true},
		{time.Date(2025, time.March, 15, 9, 0, 0, 0, time.UTC), time.Date(2025, time.March, 8, 10, 0, 0, 0, time.UTC), true},
		{time.Date(2025, time.March, 16, 9, 0, 0, 0, time.UTC), time.Date(2025, time.March, 8, 10, 0, 0, 0, time.UTC), true},
		{time.Date(2025, time.March, 17, 9, 0, 0, 0, time.UTC), time.Date(2025, time.March, 8, 10, 0, 0, 0, time.UTC), false},
		{time.Date(2025, time.March, 18, 9, 0, 0, 0, time.UTC), time.Date(2025, time.March, 8, 10, 0, 0, 0, time.UTC), false},
	}

	for _, test := range tests {
		t.Run(test.inputDate.String(), func(t *testing.T) {
			carbon.Freeze(test.nowDate)
			defer carbon.UnFreeze()
			result := dateHelper.IsDateNextWeekAndNowAfterFriday(test.inputDate)
			if result != test.expected {
				t.Errorf("Expected %v, but got %v for date %v", test.expected, result, test.inputDate)
			}
		})
	}
}

func TestIsWeekend(t *testing.T) {
	dateHelper := NewDateHelper()

	tests := []struct {
		inputDate time.Time
		expected  bool
	}{
		{time.Date(2025, time.May, 2, 0, 0, 0, 0, time.UTC), false},
		{time.Date(2025, time.May, 3, 0, 0, 0, 0, time.UTC), true},
		{time.Date(2025, time.May, 4, 0, 0, 0, 0, time.UTC), true},
		{time.Date(2025, time.May, 5, 0, 0, 0, 0, time.UTC), false},
	}

	for _, test := range tests {
		result := dateHelper.IsWeekend(test.inputDate)
		if result != test.expected {
			t.Errorf("Wrong result for: %v", test.inputDate)
		}
	}
}

func TestIsDateInCloseIntervals(t *testing.T) {
	dataHelper := NewDateHelper()

	tests := []struct {
		inputDate      time.Time
		closeIntervals []dto.CloseInterval
		expected       bool
	}{
		{
			parseDate("2025-04-29"),
			[]dto.CloseInterval{{Id: 1, StartDate: parseDate("2025-05-01"), EndDate: parseDate("2025-05-04")}},
			false,
		},
		{
			parseDate("2025-05-01"),
			[]dto.CloseInterval{{Id: 1, StartDate: parseDate("2025-05-01"), EndDate: parseDate("2025-05-04")}},
			true,
		},
		{
			parseDate("2025-05-02"),
			[]dto.CloseInterval{{Id: 1, StartDate: parseDate("2025-05-01"), EndDate: parseDate("2025-05-04")}},
			true,
		},
		{
			parseDate("2025-05-04"),
			[]dto.CloseInterval{{Id: 1, StartDate: parseDate("2025-05-01"), EndDate: parseDate("2025-05-04")}},
			true,
		},
		{
			parseDate("2025-05-05"),
			[]dto.CloseInterval{{Id: 1, StartDate: parseDate("2025-05-01"), EndDate: parseDate("2025-05-04")}},
			false,
		},
	}

	for _, test := range tests {
		result := dataHelper.IsDateInCloseIntervals(test.inputDate, test.closeIntervals)
		if result != test.expected {
			t.Errorf("Wrong result")
		}
	}
}

func parseDate(dateString string) time.Time {
	res, _ := time.Parse("2006-01-02", dateString)
	return res
}
