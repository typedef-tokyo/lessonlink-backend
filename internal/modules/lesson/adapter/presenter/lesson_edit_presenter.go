package presenter

type (
	ILessonEditPresenter interface {
		Present() *LessonEditResponse
	}

	LessonEditPresenter struct {
	}
)

func NewLessonEditPresenter() ILessonEditPresenter {
	return &LessonEditPresenter{}
}

type (
	LessonEditResponse struct {
		Msg string `json:"msg"`
	}
)

func (h *LessonEditPresenter) Present() *LessonEditResponse {

	return &LessonEditResponse{
		Msg: "更新しました",
	}
}
