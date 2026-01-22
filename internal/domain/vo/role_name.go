package vo

import (
	"errors"
	"strings"

	"github.com/typedef-tokyo/lessonlink-backend/internal/pkg/log"
)

var ErrRoleNameEmpty = errors.New("権限名称が未設定です")

type RoleName string

const (
	ROLE_NAME_INVALID = RoleName("invalid_name")
	ROLE_NAME_UNKNOWN = RoleName("role_name_unknown")
)

func NewRoleName(name string) (RoleName, error) {

	name = strings.TrimSpace(name)
	if len(name) == 0 {
		return ROLE_NAME_INVALID, log.WrapErrorWithStackTrace(ErrRoleNameEmpty)
	}

	return RoleName(name), nil
}

func (r RoleName) Value() string {

	return string(r)
}
