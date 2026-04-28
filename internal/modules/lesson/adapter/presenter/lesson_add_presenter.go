package presenter

type (
	ILessonAddPresenter interface {
		Present() *LessonAddResponse
	}

	LessonAddPresenter struct {
	}
)

func NewLessonAddPresenter() ILessonAddPresenter {
	return &LessonAddPresenter{}
}

type (
	LessonAddResponse struct {
		Msg string `json:"msg"`
	}
)

func (h *LessonAddPresenter) Present() *LessonAddResponse {

	return &LessonAddResponse{
		Msg: "追加しました",
	}
}
