package handlers

import (
	"net/http"

	"github.com/DevPulseLab/salat/internal/db/repositories"
	"github.com/DevPulseLab/salat/internal/forms"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type RealDayStatsHandler struct {
	RealDayStatsRepo *repositories.NewRealDayStatsRepository
}

func NewRealDayStatsHandler(db *gorm.DB) *RealDayStatsHandler {
	realDayStatsRepo := repositories.NewNewRealDayStatsRepository(db)
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
	var form forms.GetNumberOfPlatesForDayForm
	if err := ctx.ShouldBindJSON(&form); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusBadRequest, gin.H{"numberOfPlates": handler.RealDayStatsRepo.GetStatsForDay(form.StatsDay)})
}
