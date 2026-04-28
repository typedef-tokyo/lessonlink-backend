package vo

import (
	"errors"
	"regexp"
	"strings"

	"github.com/typedef-tokyo/lessonlink-backend/internal/pkg/log"
)

var ErrUserPasswordEmpty = errors.New("パスワードが未設定です")
var ErrUserPasswordIncludeFullWidthSpace = errors.New("パスワードに全角スペースは使えません")
var ErrUserPasswordIncludeHalfWidthSpace = errors.New("パスワードに半角スペースは使えません")
var ErrUserPasswordInvalidRule = errors.New("パスワードは6文字以上に設定してください。英数字記号が使用可能です")

type UserPassword string

const (
	USER_PASSWORD_INVALID = UserPassword("invalid")
	USER_PASSWORD_NONE    = UserPassword("")
)

func NewPasswordForCreation(value string) (UserPassword, error) {

	if value == "" {
		return USER_PASSWORD_INVALID, log.WrapErrorWithStackTrace(ErrUserPasswordEmpty)
	}

	if strings.ContainsRune(value, '　') {
		return USER_PASSWORD_INVALID, log.WrapErrorWithStackTrace(ErrUserPasswordIncludeFullWidthSpace)
	}

	if strings.ContainsRune(value, ' ') {
		return USER_PASSWORD_INVALID, log.WrapErrorWithStackTrace(ErrUserPasswordIncludeHalfWidthSpace)
	}

	var re = regexp.MustCompile(`^[\x21-\x7E]{6,}$`)
	if !re.MatchString(value) {
		return USER_PASSWORD_INVALID, log.WrapErrorWithStackTrace(ErrUserPasswordInvalidRule)
	}

	return UserPassword(value), nil
}

func ReconstructHashedPassword(_password string) (UserPassword, error) {

	if _password == "" {
		return USER_PASSWORD_INVALID, log.WrapErrorWithStackTrace(ErrUserPasswordEmpty)
	}

	return UserPassword(_password), nil
}

func (r UserPassword) Value() string {

	return string(r)
}

func (r UserPassword) IsNone() bool {

	return r == USER_PASSWORD_NONE
}
