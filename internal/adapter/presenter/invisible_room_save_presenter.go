package presenter

type (
	IInvisibleRoomPresenter interface {
		Present() *InvisibleRoomSaveResponse
	}

	InvisibleRoomPresenter struct {
	}
)

func NewInvisibleRoom() IInvisibleRoomPresenter {
	return &InvisibleRoomPresenter{}
}

type (
	InvisibleRoomSaveResponse struct {
		Msg string `json:"msg"`
	}
)

func (h *InvisibleRoomPresenter) Present() *InvisibleRoomSaveResponse {

	return &InvisibleRoomSaveResponse{
		Msg: "更新しました",
	}
}
