package handlers

import (
	"fmt"
	"net/http"

	"github.com/DevPulseLab/salat/internal/builder"
	"github.com/DevPulseLab/salat/internal/config"
	"github.com/DevPulseLab/salat/internal/db/repositories"
	"github.com/DevPulseLab/salat/internal/enum"
	"github.com/DevPulseLab/salat/internal/forms"
	"github.com/DevPulseLab/salat/internal/helper"
	"github.com/DevPulseLab/salat/internal/services"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type AdminCalendarHandler struct {
	CalendarRepo         *repositories.CalendarRepository
	CloseIntervalRepo    *repositories.CloseIntervalRepository
	VisitStatsRepo       *repositories.VisitStatsRepository
	RequestHelper        *helper.RequestHelper
	CalendarDtoBuilder   *builder.CalendarDtoBuilder
	VisitStatsDtoBuilder *builder.VisitStatsDtoBuilder
	MessagingService     *services.MessagingService
	Logger               *logrus.Logger
}

func NewAdminCalendarHandler(db *gorm.DB, config *config.Config, log *logrus.Logger) *AdminCalendarHandler {
	closeIntervalRepo := repositories.NewCloseIntervalsRepository(db)
	calendarRepo := repositories.NewCalendarRepository(db)
	visitStatsRepo := repositories.NewVisitStatsRepository(db)
	requestHelper := helper.NewRequestHelper()
	calendarDtoBuilder := builder.NewCalendarDtoBuilder()
	visitStatsDtoBuilder := builder.NewVisitStatsDtoBuilder()
	ms := services.NewMessagingService(config.Slack.Token, db)
	return &AdminCalendarHandler{
		calendarRepo,
		closeIntervalRepo,
		visitStatsRepo,
		requestHelper,
		calendarDtoBuilder,
		visitStatsDtoBuilder,
		ms,
		log}
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

	calendars, err := handler.CalendarRepo.FindByDateRange(startDate, endDate)
	if err != nil {
		handler.Logger.Errorf("Errror while find data by date range: %s", err.Error())
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "An error occurred data by range"})
		return
	}

	calendarDtos := handler.CalendarDtoBuilder.BuildFromCalendarModels(calendars)

	ctx.JSON(http.StatusOK, gin.H{"calendarEntries": calendarDtos})
}

func (handler *AdminCalendarHandler) ChangeEntryStatus(ctx *gin.Context) {
	var form forms.ChangeCalendarEntryForm
	if err := ctx.ShouldBindJSON(&form); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	calendarModel, err := handler.CalendarRepo.UpdateStatus(form.CalendarEntryId, form.NewStatus)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Calendar entry not changed"})
		return
	}

	if enum.CalendarStatus(calendarModel.Status) == enum.Approved {
		err := handler.MessagingService.SendPrivateMessage(
			calendarModel.UserId,
			fmt.Sprintf(
				"Deine Anfrage wurde genehmigt, du darfst zur Salatbar am %s kommen",
				calendarModel.Date.Format("02.01.2006")))
		if err != nil {
			handler.Logger.Errorf("Errror while sending message to user: %s", err.Error())
		}
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

func (handler *AdminCalendarHandler) GetVisitStatsList(ctx *gin.Context) {
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

	visitStatsList := handler.VisitStatsRepo.GetVisitVisit(startDate, endDate)

	visitStatsDtos := handler.VisitStatsDtoBuilder.BuildFromVisitStatsModels(visitStatsList)

	ctx.JSON(http.StatusOK, gin.H{"calendarEntries": visitStatsDtos})
}

func (handler *AdminCalendarHandler) ToggleVisit(ctx *gin.Context) {
	var form forms.ToggleVisitForm
	if err := ctx.ShouldBindJSON(&form); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	visitStats, err := handler.VisitStatsRepo.ToggleVisit(form.UserId, form.VisitDate)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": "Visit stats not saved"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "Visit stats saved", "isVisit": visitStats.IsVisit})
}
