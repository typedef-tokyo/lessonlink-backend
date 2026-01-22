package presenter

import (
	"github.com/typedef-tokyo/lessonlink-backend/internal/usecase"
)

type IScheduleSavePresenter interface {
	Present(result *usecase.ScheduleSaveOutput) *ScheduleSaveResponse
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

func (h *ScheduleSavePresenter) Present(result *usecase.ScheduleSaveOutput) *ScheduleSaveResponse {

	return &ScheduleSaveResponse{
		Msg:          "保存しました",
		HistoryIndex: result.HistoryIndex,
	}
}
