package vo

import (
	"errors"
	"fmt"
	"strings"
	"unicode/utf8"

	"github.com/typedef-tokyo/lessonlink-backend/internal/pkg/log"
)

var ErrUserNameEmpty = errors.New("ユーザーネームが未設定です")
var ErrUserNameLengthOver = errors.New("ユーザーネームに設定できる最大文字数を超えています")

type UserName string

const (
	USER_NAME_INVALID = UserName("invalid")
)

func NewUserName(name string) (UserName, error) {

	name = strings.TrimSpace(name)
	if name == "" {
		return USER_NAME_INVALID, log.WrapErrorWithStackTrace(ErrUserNameEmpty)
	}

	const MAX_USER_NAME_LENGTH = 30
	if utf8.RuneCountInString(name) > MAX_USER_NAME_LENGTH {
		return USER_NAME_INVALID, log.WrapErrorWithStackTrace(fmt.Errorf("%w 最大:%d文字", ErrUserNameLengthOver, MAX_USER_NAME_LENGTH))
	}

	return UserName(name), nil
}

func (r UserName) Value() string {

	return string(r)
}
