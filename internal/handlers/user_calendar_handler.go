package handlers

import (
	"log"
	"net/http"
	"time"

	"github.com/DevPulseLab/salat/internal/db/repositories"
	"github.com/DevPulseLab/salat/internal/dto"
	"github.com/DevPulseLab/salat/internal/enum"
	"github.com/DevPulseLab/salat/internal/forms"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type UserCalendarHandler struct {
	CalendarRepo *repositories.CalendarRepository
	UserRepo     *repositories.UserRepository
}

func NewUserCalendarHandler(db *gorm.DB) *UserCalendarHandler {
	calendarRepo := repositories.NewCalendarRepository(db)
	userRepo := repositories.NewUserRepository(db)
	return &UserCalendarHandler{calendarRepo, userRepo}
}

func (handler *UserCalendarHandler) Add(ctx *gin.Context) {
	var form forms.UserAddCalendarEntryForm
	if err := ctx.ShouldBindJSON(&form); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	userId, err := handler.UserRepo.GetIdByUsername(ctx.GetString("username"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	status := enum.Approved
	if isDateInCurrentWeek(form.StartDate) {
		status = enum.Reserved
	} else if form.StartDate.Before(time.Now()) {
		status = enum.Rejected
	}

	ok, errors := handler.CalendarRepo.AddCalendarEntry(userId, form.StartDate, form.EndDate, status)

	if ok {
		ctx.JSON(http.StatusOK, gin.H{"message": "Calendar data saved"})
	} else {
		log.Println(errors)
		ctx.JSON(http.StatusBadRequest, gin.H{"message": "Calendar data was not saved"})
	}
}

func isDateInCurrentWeek(t time.Time) bool {
	year, week := time.Now().ISOWeek()
	targetYear, targetWeek := t.ISOWeek()
	return year == targetYear && week == targetWeek
}

func (handler *UserCalendarHandler) AllUserList(ctx *gin.Context) {
	startDate, err := getStartDateFromRequest(ctx)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid start_date format. Use YYYY-MM-DD."})
		return
	}

	endDate, err := getEndDateFromRequest(ctx)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid end_date format. Use YYYY-MM-DD."})
		return
	}

	calendars := handler.CalendarRepo.GetCalendarEntriesForAllUsers(startDate, endDate)

	calendarDtos := []dto.Calendar{}
	for _, calendar := range calendars {
		calendarDto := dto.Calendar{
			Id:     calendar.ID,
			UserId: calendar.ID,
			Date:   calendar.Date,
			Status: calendar.Status,
		}
		calendarDtos = append(calendarDtos, calendarDto)
	}

	ctx.JSON(http.StatusOK, gin.H{"calendarEntries": calendarDtos})
}

func (handler *UserCalendarHandler) CurrentUserList(ctx *gin.Context) {
	userId, err := handler.UserRepo.GetIdByUsername(ctx.GetString("username"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	startDate, err := getStartDateFromRequest(ctx)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid start_date format. Use YYYY-MM-DD."})
		return
	}

	endDate, err := getEndDateFromRequest(ctx)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid end_date format. Use YYYY-MM-DD."})
		return
	}

	calendars := handler.CalendarRepo.GetCalendarEntriesByUserId(userId, startDate, endDate)

	calendarDtos := []dto.Calendar{}
	for _, calendar := range calendars {
		calendarDto := dto.Calendar{
			Id:     calendar.ID,
			UserId: userId,
			Date:   calendar.Date,
			Status: calendar.Status,
		}
		calendarDtos = append(calendarDtos, calendarDto)
	}

	ctx.JSON(http.StatusOK, gin.H{"calendarEntries": calendarDtos})
}

func getStartDateFromRequest(ctx *gin.Context) (time.Time, error) {
	startDateString := ctx.DefaultQuery("start_date", "")

	var startDate time.Time
	if startDateString == "" {
		startDate = time.Now().AddDate(0, 0, -int(time.Now().Weekday()-1))
	} else {
		var err error
		startDate, err = parseDate(startDateString)
		if err != nil {
			return time.Time{}, err
		}
	}

	return time.Date(startDate.Year(), startDate.Month(), startDate.Day(), 0, 0, 0, 0, startDate.Location()), nil
}

func getEndDateFromRequest(ctx *gin.Context) (time.Time, error) {
	endDateString := ctx.DefaultQuery("end_date", "")

	var endDate time.Time
	if endDateString == "" {
		endDate = time.Now().AddDate(0, 0, 7-int(time.Now().Weekday()))
	} else {
		var err error
		endDate, err = parseDate(endDateString)
		if err != nil {
			return time.Time{}, err
		}
	}

	return time.Date(endDate.Year(), endDate.Month(), endDate.Day(), 23, 59, 59, 999999999, endDate.Location()), nil
}

func parseDate(dateString string) (time.Time, error) {
	return time.Parse("2006-01-02", dateString)
}
