package presenter

type IUserDeletePresenter interface {
	Present() *UserDeleteResponse
}

type UserDeletePresenter struct {
}

func NewUserDeletePresenter() IUserDeletePresenter {
	return &UserDeletePresenter{}
}

type (
	UserDeleteResponse struct {
		Msg string `json:"msg"`
	}
)

func (h *UserDeletePresenter) Present() *UserDeleteResponse {

	return &UserDeleteResponse{
		Msg: "ユーザーを削除しました",
	}
}
