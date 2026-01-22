package presenter

type IUserAddPresenter interface {
	Present() *UserAddResponse
}

type UserAddPresenter struct {
}

func NewUserAddPresenter() IUserAddPresenter {
	return &UserAddPresenter{}
}

type (
	UserAddResponse struct {
		Msg string `json:"msg"`
	}
)

func (h *UserAddPresenter) Present() *UserAddResponse {

	return &UserAddResponse{
		Msg: "新規追加されました",
	}
}
