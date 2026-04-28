package vo

import (
	"errors"

	"github.com/typedef-tokyo/lessonlink-backend/internal/pkg/log"
)

var ErrUserIDMinUnder = errors.New("ユーザーIDは0以上を設定する必要があります")

type UserID int

const (
	USER_ID_INVALID = UserID(-1)
	USER_ID_INITIAL = UserID(0)
)

func NewUserID(id int) (UserID, error) {

	if id < 0 {
		return USER_ID_INVALID, log.WrapErrorWithStackTrace(ErrUserIDMinUnder)
	}

	return UserID(id), nil
}

func (r UserID) IsValid() bool {
	return r > USER_ID_INITIAL
}

func (r UserID) Value() int {

	return int(r)
}
