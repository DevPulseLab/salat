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

func (handler *RealDayStatsHandler) IncrementCountForDay(ctx *gin.Context) {
	var form forms.IncrementNumberOfPlatesForm
	if err := ctx.ShouldBindJSON(&form); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	handler.RealDayStatsRepo.IncrementStatsForDay(form.StatsDay)

	ctx.JSON(http.StatusOK, gin.H{"message": "Day stats data saved", "success": true})
}
