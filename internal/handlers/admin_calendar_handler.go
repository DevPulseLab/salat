package handlers

import (
	"net/http"

	"github.com/DevPulseLab/salat/internal/builder"
	"github.com/DevPulseLab/salat/internal/db/repositories"
	"github.com/DevPulseLab/salat/internal/forms"
	"github.com/DevPulseLab/salat/internal/helper"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type AdminCalendarHandler struct {
	CalendarRepo       *repositories.CalendarRepository
	CloseIntervalRepo  *repositories.CloseIntervalRepository
	RequestHelper      *helper.RequestHelper
	CalendarDtoBuilder *builder.CalendarDtoBuilder
}

func NewAdminCalendarHandler(db *gorm.DB) *AdminCalendarHandler {
	dateHelper := helper.NewDateHelper()
	closeIntervalRepo := repositories.NewCloseIntervalsRepository(db)
	calendarRepo := repositories.NewCalendarRepository(db, dateHelper)
	requestHelper := helper.NewRequestHelper()
	calendarDtoBuilder := builder.NewCalendarDtoBuilder()
	return &AdminCalendarHandler{calendarRepo, closeIntervalRepo, requestHelper, calendarDtoBuilder}
}

func (handler *AdminCalendarHandler) AllUserList(ctx *gin.Context) {
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

	calendars := handler.CalendarRepo.GetCalendarEntriesForAllUsers(startDate, endDate)

	calendarDtos := handler.CalendarDtoBuilder.BuildFromCalendarModels(calendars)

	ctx.JSON(http.StatusOK, gin.H{"calendarEntries": calendarDtos})
}

func (handler *AdminCalendarHandler) ChangeEntryStatus(ctx *gin.Context) {
	var form forms.ChangeCalendarEntryForm
	if err := ctx.ShouldBindJSON(&form); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err := handler.CalendarRepo.ChangeEntryStatus(form.CalendarEntryId, form.NewStatus)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Calendar entry not changed"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"success": true})
}

func (handler *AdminCalendarHandler) AddCloseDateInterval(ctx *gin.Context) {
	var form forms.AddCloseCalendarEntryForm
	if err := ctx.ShouldBindJSON(&form); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err := handler.CloseIntervalRepo.SaveCloseInterval(form.StartDate, form.EndDate)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Close interval not saved"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"success": true})
}

func (handler *AdminCalendarHandler) RemoveCloseDateInterval(ctx *gin.Context) {
	var form forms.RemoveCloseIntervalEntryForm
	if err := ctx.ShouldBindJSON(&form); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	closeIntervalEntry, err := handler.CloseIntervalRepo.GetById(form.CloseIntervalEntryId)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Close interval not found"})
		return
	}

	handler.CloseIntervalRepo.Remove(&closeIntervalEntry)

	ctx.JSON(http.StatusOK, gin.H{"success": true})
}
