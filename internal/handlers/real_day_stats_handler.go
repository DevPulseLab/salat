package handlers

import (
	"net/http"
	"time"

	"github.com/DevPulseLab/salat/internal/db/repositories"
	"github.com/DevPulseLab/salat/internal/forms"
	"github.com/DevPulseLab/salat/internal/helper"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type RealDayStatsHandler struct {
	RealDayStatsRepo *repositories.RealDayStatsRepository
}

func NewRealDayStatsHandler(db *gorm.DB) *RealDayStatsHandler {
	realDayStatsRepo := repositories.NewRealDayStatsRepository(db)
	return &RealDayStatsHandler{realDayStatsRepo}
}

func (handler *RealDayStatsHandler) IncrementNumberOfPlatesForDay(ctx *gin.Context) {
	var form forms.IncrementNumberOfPlatesForm
	if err := ctx.ShouldBindJSON(&form); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if handler.RealDayStatsRepo.IncrementStatsForDay(form.StatsDay) {
		ctx.JSON(http.StatusOK, gin.H{"message": "Day stats data saved", "success": true})
		return
	}

	ctx.JSON(http.StatusBadRequest, gin.H{"message": "Day stats not saved", "success": false})
}

func (handler *RealDayStatsHandler) SaveNumberOfPlatesForDay(ctx *gin.Context) {
	var form forms.SaveNumberOfPlatesForDayForm
	if err := ctx.ShouldBindJSON(&form); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if handler.RealDayStatsRepo.SaveStatsForDay(form.StatsDay, form.NumberOfPlates) {
		ctx.JSON(http.StatusOK, gin.H{"message": "Day stats data saved", "success": true})
		return
	}

	ctx.JSON(http.StatusBadRequest, gin.H{"message": "Day stats not saved", "success": false})
}

func (handler *RealDayStatsHandler) GetNumberOfPlatesForDay(ctx *gin.Context) {
	statsDay, err := getStatsDateFromRequest(ctx)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid stats_date format. Use YYYY-MM-DD."})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"numberOfPlates": handler.RealDayStatsRepo.GetStatsForDay(statsDay)})
}

func getStatsDateFromRequest(ctx *gin.Context) (time.Time, error) {
	statsDate := ctx.DefaultQuery("stats_date", "")

	var startDate time.Time
	if statsDate == "" {
		startDate = time.Now()
	} else {
		var err error
		startDate, err = helper.ParseDate(statsDate)
		if err != nil {
			return time.Time{}, err
		}
	}

	return time.Date(startDate.Year(), startDate.Month(), startDate.Day(), 0, 0, 0, 0, time.Local), nil
}
