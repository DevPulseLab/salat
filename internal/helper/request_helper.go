package helper

import (
	"time"

	"github.com/gin-gonic/gin"
)

type RequestHelper struct {
}

func NewRequestHelper() *RequestHelper {
	return &RequestHelper{}
}

func (handler *RequestHelper) GetStartDateFromRequest(ctx *gin.Context) (time.Time, error) {
	startDateString := ctx.DefaultQuery("start_date", "")

	var startDate time.Time
	if startDateString == "" {
		startDate = time.Now().AddDate(0, 0, -int(time.Now().Weekday()-1))
	} else {
		var err error
		startDate, err = ParseDate(startDateString)
		if err != nil {
			return time.Time{}, err
		}
	}

	return time.Date(startDate.Year(), startDate.Month(), startDate.Day(), 0, 0, 0, 0, time.Local), nil
}

func (handler *RequestHelper) GetEndDateFromRequest(ctx *gin.Context) (time.Time, error) {
	endDateString := ctx.DefaultQuery("end_date", "")

	var endDate time.Time
	if endDateString == "" {
		endDate = time.Now().AddDate(0, 0, 7-int(time.Now().Weekday()))
	} else {
		var err error
		endDate, err = ParseDate(endDateString)
		if err != nil {
			return time.Time{}, err
		}
	}

	return time.Date(endDate.Year(), endDate.Month(), endDate.Day(), 23, 59, 59, 999999999, time.Local), nil
}
