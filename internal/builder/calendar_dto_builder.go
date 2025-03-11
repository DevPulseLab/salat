package builder

import (
	"github.com/DevPulseLab/salat/internal/db/models"
	"github.com/DevPulseLab/salat/internal/dto"
)

type CalendarDtoBuilder struct {
}

func NewCalendarDtoBuilder() *CalendarDtoBuilder {
	return &CalendarDtoBuilder{}
}

func (builder *CalendarDtoBuilder) BuildFromCalendarModels(calendarModels []models.Calendar) []dto.Calendar {
	calendarDtos := []dto.Calendar{}
	for _, calendar := range calendarModels {
		calendarDto := dto.Calendar{
			Id:     calendar.ID,
			UserId: calendar.UserId,
			Date:   calendar.Date,
			Status: calendar.Status,
		}
		calendarDtos = append(calendarDtos, calendarDto)
	}
	return calendarDtos
}
