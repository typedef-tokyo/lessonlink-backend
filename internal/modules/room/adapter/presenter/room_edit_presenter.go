package presenter

type IRoomEditPresenter interface {
	Present() *RoomEditResponse
}

type RoomEditPresenter struct {
}

func NewRoomEditPresenter() IRoomEditPresenter {
	return &RoomEditPresenter{}
}

type (
	RoomEditResponse struct {
		Msg string `json:"msg"`
	}
)

func (h *RoomEditPresenter) Present() *RoomEditResponse {

	return &RoomEditResponse{
		Msg: "更新しました",
	}
}
