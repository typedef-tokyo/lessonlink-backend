package vo

import (
	"errors"
	"fmt"
	"strings"
	"unicode/utf8"

	"github.com/typedef-tokyo/lessonlink-backend/internal/pkg/log"
)

var ErrUserDisplayNameEmpty = errors.New("ユーザー表示名が未設定です")
var ErrUserDisplayNameLengthOver = errors.New("ユーザー表示名に設定できる最大数を超えています")

type UserDisplayName string

const (
	USER_DISPLAY_NAME_INVALID = UserDisplayName("invalid")
)

func NewUserDisplayName(name string) (UserDisplayName, error) {

	name = strings.TrimSpace(name)
	if name == "" {
		return USER_DISPLAY_NAME_INVALID, log.WrapErrorWithStackTrace(ErrUserDisplayNameEmpty)
	}

	const MAX_USER_DISPLAY_NAME_LENGTH = 30
	if utf8.RuneCountInString(name) > MAX_USER_DISPLAY_NAME_LENGTH {
		return USER_DISPLAY_NAME_INVALID, log.WrapErrorWithStackTrace(fmt.Errorf("%w 最大:%d文字", ErrUserDisplayNameLengthOver, MAX_USER_DISPLAY_NAME_LENGTH))
	}

	return UserDisplayName(name), nil
}

func (r UserDisplayName) Value() string {

	return string(r)
}
