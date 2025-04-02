package builder

import (
	"github.com/DevPulseLab/salat/internal/db/models"
	"github.com/DevPulseLab/salat/internal/dto"
)

type VisitStatsDtoBuilder struct {
}

func NewVisitStatsDtoBuilder() *VisitStatsDtoBuilder {
	return &VisitStatsDtoBuilder{}
}

func (builder *VisitStatsDtoBuilder) BuildFromVisitStatsModels(visitStatsModels []models.VisitStats) []dto.VisitStats {
	visitStatsDtos := []dto.VisitStats{}
	for _, visitStats := range visitStatsModels {
		visitStatsDto := dto.VisitStats{
			Id:      visitStats.ID,
			UserId:  visitStats.UserId,
			Date:    visitStats.Date,
			IsVisit: visitStats.IsVisit,
		}
		visitStatsDtos = append(visitStatsDtos, visitStatsDto)
	}
	return visitStatsDtos
}
