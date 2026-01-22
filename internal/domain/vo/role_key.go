package vo

import (
	"errors"
	"slices"
	"strings"

	"github.com/typedef-tokyo/lessonlink-backend/internal/pkg/log"
)

var ErrRoleKeyEmpty = errors.New("権限が未設定です")
var ErrRoleKeyInvalid = errors.New("ロール識別子が不適切です。")

type RoleKey string

const (
	ROLE_KEY_INVALID = RoleKey("invalid_key")
	ROLE_KEY_NONE    = RoleKey("none")
)

const (
	ROLE_KEY_OWNER  = RoleKey("owner")
	ROLE_KEY_EDITOR = RoleKey("editor")
	ROLE_KEY_VIEWER = RoleKey("viewer")
)

var validRoleKeys = []RoleKey{
	ROLE_KEY_OWNER,
	ROLE_KEY_EDITOR,
	ROLE_KEY_VIEWER,
}

func NewRoleKey(key string) (RoleKey, error) {

	key = strings.TrimSpace(key)
	if len(key) == 0 {
		return ROLE_KEY_INVALID, log.WrapErrorWithStackTrace(ErrRoleKeyEmpty)
	}

	roleKey := RoleKey(key)

	switch roleKey {
	case ROLE_KEY_OWNER:
		return roleKey, nil
	case ROLE_KEY_EDITOR:
		return roleKey, nil
	case ROLE_KEY_VIEWER:
		return roleKey, nil
	default:
		return ROLE_KEY_INVALID, log.WrapErrorWithStackTrace(ErrRoleKeyInvalid)
	}
}

func (r RoleKey) Value() string {

	return string(r)
}

func (r RoleKey) IsValid() bool {
	return slices.Contains(validRoleKeys, r)
}

func (r RoleKey) IsOwner() bool {
	return r == ROLE_KEY_OWNER
}

func (r RoleKey) IsEditor() bool {
	return r == ROLE_KEY_EDITOR
}

func (r RoleKey) IsViewer() bool {
	return r == ROLE_KEY_VIEWER
}
