package builder

import (
	"github.com/DevPulseLab/salat/internal/db/models"
	"github.com/DevPulseLab/salat/internal/dto"
)

type CloseIntervalDtoBuilder struct {
}

func NewCloseIntervalDtoBuilder() *CloseIntervalDtoBuilder {
	return &CloseIntervalDtoBuilder{}
}

func (builder *CloseIntervalDtoBuilder) BuildFromCloseIntervalModel(closeIntervalModels []models.CloseInterval) []dto.CloseInterval {
	closeIntervalDtos := []dto.CloseInterval{}

	for _, closeIntervalModel := range closeIntervalModels {
		closeIntervalDto := dto.CloseInterval{
			Id:        closeIntervalModel.ID,
			StartDate: closeIntervalModel.StartDate,
			EndDate:   closeIntervalModel.EndDate,
		}

		closeIntervalDtos = append(closeIntervalDtos, closeIntervalDto)
	}

	return closeIntervalDtos
}
