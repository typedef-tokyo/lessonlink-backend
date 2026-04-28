package presenter

import "github.com/typedef-tokyo/lessonlink-backend/internal/modules/schedule/usecase/interactor"

type IScheduleSavePresenter interface {
	Present(result *interactor.ScheduleSaveOutput) *ScheduleSaveResponse
}

type ScheduleSavePresenter struct {
}

func NewScheduleSavePresenter() IScheduleSavePresenter {
	return &ScheduleSavePresenter{}
}

type (
	ScheduleSaveResponse struct {
		Msg          string `json:"msg"`
		HistoryIndex int    `json:"history_index"`
	}
)

func (h *ScheduleSavePresenter) Present(result *interactor.ScheduleSaveOutput) *ScheduleSaveResponse {

	return &ScheduleSaveResponse{
		Msg:          "保存しました",
		HistoryIndex: result.HistoryIndex,
	}
}
