package presenter

import "github.com/typedef-tokyo/lessonlink-backend/internal/modules/schedule/usecase/interactor"

type IScheduleCreatePresenter interface {
	Present(result *interactor.ScheduleCreateOutput) *ScheduleCreateResponse
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

func (h *ScheduleCreatePresenter) Present(result *interactor.ScheduleCreateOutput) *ScheduleCreateResponse {

	return &ScheduleCreateResponse{
		ScheduleID: result.ScheduleID,
	}
}
