package presenter

import (
	"github.com/typedef-tokyo/lessonlink-backend/internal/usecase"
)

type IScheduleCreatePresenter interface {
	Present(result *usecase.ScheduleCreateOutput) *ScheduleCreateResponse
}

type ScheduleCreatePresenter struct {
}

func NewScheduleCreatePresenter() IScheduleCreatePresenter {
	return &ScheduleCreatePresenter{}
}

type (
	ScheduleCreateResponse struct {
		ScheduleID int `json:"schedule_id"`
	}
)

func (h *ScheduleCreatePresenter) Present(result *usecase.ScheduleCreateOutput) *ScheduleCreateResponse {

	return &ScheduleCreateResponse{
		ScheduleID: result.ScheduleID,
	}
}
