package services

import (
	"time"

	"github.com/DevPulseLab/salat/internal/db/models"
	"github.com/DevPulseLab/salat/internal/db/repositories"
	"github.com/DevPulseLab/salat/internal/dto"
	"github.com/DevPulseLab/salat/internal/enum"
	"github.com/DevPulseLab/salat/internal/helper"
	"github.com/uniplaces/carbon"
)

type CalendarService struct {
	repo       *repositories.CalendarRepository
	dateHelper *helper.DateHelper
}

func NewCalendarService(repo *repositories.CalendarRepository, dh *helper.DateHelper) *CalendarService {
	return &CalendarService{repo: repo, dateHelper: dh}
}

func (s *CalendarService) AddCalendarEntries(user *models.User, startDate, endDate time.Time, closeIntervals []dto.CloseInterval) ([]models.Calendar, []error) {
	var errors []error
	var addedDays []models.Calendar
	var newEntries []models.Calendar

	currDate := time.Date(startDate.Year(), startDate.Month(), startDate.Day(), 0, 0, 0, 0, startDate.Location())
	endDate = time.Date(endDate.Year(), endDate.Month(), endDate.Day(), 0, 0, 0, 0, endDate.Location())

	deletedEntries, err := s.repo.FindDeletedByUserIdAndDateRange(user.ID, currDate, endDate)
	if err != nil {
		return nil, []error{err}
	}

	deletedEntriesMap := make(map[string]models.Calendar)
	for _, entry := range deletedEntries {
		deletedEntriesMap[entry.Date.Format("2006-01-02")] = entry
	}

	for endDate.Sub(currDate).Hours() >= 0 {
		if s.dateHelper.IsWeekend(currDate) || s.dateHelper.IsDateInCloseIntervals(currDate, closeIntervals) {
			currDate = currDate.AddDate(0, 0, 1)
			continue
		}

		status := s.determineStatus(user, currDate)
		dateStr := currDate.Format("2006-01-02")

		if deletedEntry, exists := deletedEntriesMap[dateStr]; exists {
			if err := s.repo.RestoreAndUpdate(&deletedEntry, status); err != nil {
				errors = append(errors, err)
			} else {
				addedDays = append(addedDays, deletedEntry)
			}
		} else {
			calendarEntry := models.Calendar{
				UserId: user.ID,
				Date:   currDate,
				Status: status,
			}
			newEntries = append(newEntries, calendarEntry)
		}

		currDate = currDate.AddDate(0, 0, 1)
	}

	if len(newEntries) > 0 {
		if err := s.repo.BatchCreate(newEntries); err != nil {
			errors = append(errors, err)
		} else {
			addedDays = append(addedDays, newEntries...)
		}
	}

	if len(errors) > 0 {
		return addedDays, errors
	}

	return addedDays, nil
}

func (s *CalendarService) determineStatus(user *models.User, date time.Time) string {
	now := carbon.Now().Time
	nowPlus30Days := carbon.Now().AddDate(0, 0, 30)

	if user.PenaltyCard == string(enum.Yellow) {
		return string(enum.Reserved)
	}
	if user.PenaltyCard == string(enum.Red) {
		return string(enum.Rejected)
	}
	if date.Equal(now) || date.Before(now) || date.After(nowPlus30Days) {
		return string(enum.Rejected)
	}
	if s.dateHelper.IsDateInCurrentWeek(date) {
		return string(enum.Reserved)
	}
	if s.dateHelper.IsDateNextWeekAndNowAfterFriday(date) {
		return string(enum.Reserved)
	}
	return string(enum.Approved)
}
