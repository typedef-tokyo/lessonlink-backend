package vo

import (
	"errors"

	"github.com/typedef-tokyo/lessonlink-backend/internal/pkg/log"
)

var ErrRoleIDUnderMin = errors.New("ロールIDは0以上を設定する必要があります。")

type RoleID int

const (
	ROLE_ID_INVALID = RoleID(-1)
	ROLE_ID_INITIAL = RoleID(0)
)

func NewRoleID(id int) (RoleID, error) {

	if id < 0 {
		return ROLE_ID_INVALID, log.WrapErrorWithStackTrace(ErrRoleIDUnderMin)
	}

	return RoleID(id), nil
}

func (r RoleID) IsValid() bool {
	return r > ROLE_ID_INITIAL
}
