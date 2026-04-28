package presenter

type IUpdateUserPresenter interface {
	Present() *UserUpdateResponse
}

type UpdateUserPresenter struct {
}

func NewUpdateUserPresenter() IUpdateUserPresenter {
	return &UpdateUserPresenter{}
}

type UserUpdateResponse struct {
	Msg string `json:"msg"`
}

func (h *UpdateUserPresenter) Present() *UserUpdateResponse {

	return &UserUpdateResponse{
		Msg: "ユーザー情報を更新しました",
	}
}
