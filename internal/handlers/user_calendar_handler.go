package handlers

import (
	"net/http"

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

	handler.CalendarRepo.AddCalendarEntry(userId, form.StartDate, form.EndDate, enum.Approved)

	ctx.JSON(http.StatusOK, gin.H{"message": "Calendar data saved"})
}

func (handler *UserCalendarHandler) CurrentUserList(ctx *gin.Context) {
	userId, err := handler.UserRepo.GetIdByUsername(ctx.GetString("username"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	calendars := handler.CalendarRepo.GetCalendarEntriesByUserId(userId)

	calendarDtos := []dto.Calendar{}
	for _, calendar := range calendars {
		calendarDto := dto.Calendar{
			Id:     calendar.ID,
			Date:   calendar.Date,
			Status: calendar.Status,
		}
		calendarDtos = append(calendarDtos, calendarDto)
	}

	ctx.JSON(http.StatusOK, gin.H{"users": calendarDtos})
}
