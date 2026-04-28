package presenter

type IScheduleSaveTitlePresenter interface {
	Present() *ScheduleSaveTitleResponse
}

type ScheduleSaveTitlePresenter struct {
}

func NewScheduleSaveTitlePresenter() IScheduleSaveTitlePresenter {
	return &ScheduleSaveTitlePresenter{}
}

type (
	ScheduleSaveTitleResponse struct {
		Msg string `json:"msg"`
	}
)

func (h *ScheduleSaveTitlePresenter) Present() *ScheduleSaveTitleResponse {

	return &ScheduleSaveTitleResponse{
		Msg: "保存しました",
	}
}
