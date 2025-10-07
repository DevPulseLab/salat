package handlers

import (
	"net/http"
	"time"

	"github.com/DevPulseLab/salat/internal/builder"
	"github.com/DevPulseLab/salat/internal/db/models"
	"github.com/DevPulseLab/salat/internal/db/repositories"
	"github.com/DevPulseLab/salat/internal/forms"
	"github.com/DevPulseLab/salat/internal/helper"
	"github.com/DevPulseLab/salat/internal/services"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type UserCalendarHandler struct {
	CalendarRepo            *repositories.CalendarRepository
	CalendarService         *services.CalendarService
	UserRepo                *repositories.UserRepository
	CloseIntervalRepo       *repositories.CloseIntervalRepository
	DateHelper              *helper.DateHelper
	RequestHelper           *helper.RequestHelper
	CalendarDtoBuilder      *builder.CalendarDtoBuilder
	CloseIntervalDtoBuilder *builder.CloseIntervalDtoBuilder
}

func NewUserCalendarHandler(db *gorm.DB) *UserCalendarHandler {
	dateHelper := helper.NewDateHelper()
	calendarRepo := repositories.NewCalendarRepository(db)
	calendarService := services.NewCalendarService(calendarRepo, dateHelper)
	userRepo := repositories.NewUserRepository(db)
	closeIntervalsRepo := repositories.NewCloseIntervalsRepository(db)
	requestHelper := helper.NewRequestHelper()
	calendarDtoBuilder := builder.NewCalendarDtoBuilder()
	closeIntervalDtoBuilder := builder.NewCloseIntervalDtoBuilder()
	return &UserCalendarHandler{calendarRepo, calendarService, userRepo, closeIntervalsRepo, dateHelper, requestHelper, calendarDtoBuilder, closeIntervalDtoBuilder}
}

func (handler *UserCalendarHandler) Add(ctx *gin.Context) {
	var form forms.UserAddCalendarEntryForm
	if err := ctx.ShouldBindJSON(&form); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user := ctx.MustGet("user").(*models.User)

	closeIntervalModels := handler.CloseIntervalRepo.GetAllEntriesForInterval(form.StartDate, form.EndDate)

	addedCalendarModels, errors := handler.CalendarService.AddCalendarEntries(
		user,
		form.StartDate,
		form.EndDate,
		handler.CloseIntervalDtoBuilder.BuildFromCloseIntervalModel(closeIntervalModels))

	if len(errors) == 0 {
		calendarDtos := handler.CalendarDtoBuilder.BuildFromCalendarModels(addedCalendarModels)

		ctx.JSON(http.StatusOK, gin.H{"message": "Calendar data saved", "calendarEntries": calendarDtos})
		return
	}

	if len(errors) != 0 && len(addedCalendarModels) == 0 {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": "Calendar data was not saved"})
		return
	}

	calendarDtos := handler.CalendarDtoBuilder.BuildFromCalendarModels(addedCalendarModels)
	ctx.JSON(http.StatusOK, gin.H{"message": "Not all calendar data was saved", "calendarEntries": calendarDtos})
}

func (handler *UserCalendarHandler) CurrentUserList(ctx *gin.Context) {
	userId, err := handler.UserRepo.GetIdByUsername(ctx.GetString("username"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	startDate, err := handler.RequestHelper.GetStartDateFromRequest(ctx)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid start_date format. Use YYYY-MM-DD."})
		return
	}

	endDate, err := handler.RequestHelper.GetEndDateFromRequest(ctx)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid end_date format. Use YYYY-MM-DD."})
		return
	}

	calendars, err := handler.CalendarRepo.FindByUserIdAndDateRange(userId, startDate, endDate)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "An error occurred while find data for user"})
		return
	}

	calendarDtos := handler.CalendarDtoBuilder.BuildFromCalendarModels(calendars)

	ctx.JSON(http.StatusOK, gin.H{"calendarEntries": calendarDtos})
}

func (handler *UserCalendarHandler) RemoveEntryForCurrentUser(ctx *gin.Context) {
	var form forms.RemoveCalendarEntryForm
	if err := ctx.ShouldBindJSON(&form); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	userId, err := handler.UserRepo.GetIdByUsername(ctx.GetString("username"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	calendarEntry, err := handler.CalendarRepo.FindByIdForUserId(form.CalendarEntryId, userId)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Calendar entry not found"})
		return
	}

	if calendarEntry.Date.Before(time.Now()) ||
		calendarEntry.IsRejected() ||
		(calendarEntry.IsApproved() && handler.DateHelper.IsDateInCurrentWeek(calendarEntry.Date)) {

		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Can not remove past, rejected or approved this week entries"})
		return
	}

	handler.CalendarRepo.SoftDelete(&calendarEntry)

	ctx.JSON(http.StatusOK, gin.H{"success": true})
}

func (handler *UserCalendarHandler) GetCloseDateInterval(ctx *gin.Context) {
	startDate, err := handler.RequestHelper.GetStartDateFromRequest(ctx)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid start_date format. Use YYYY-MM-DD."})
		return
	}

	endDate, err := handler.RequestHelper.GetEndDateFromRequest(ctx)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid end_date format. Use YYYY-MM-DD."})
		return
	}

	closeDateIntervalModels := handler.CloseIntervalRepo.GetAllEntriesForInterval(startDate, endDate)

	closeDateIntervalsDto := handler.CloseIntervalDtoBuilder.BuildFromCloseIntervalModel(closeDateIntervalModels)

	ctx.JSON(http.StatusOK, gin.H{"closeDateIntervals": closeDateIntervalsDto})
}
