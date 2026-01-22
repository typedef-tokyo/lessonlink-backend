package vo

import (
	"errors"
	"strings"

	"github.com/typedef-tokyo/lessonlink-backend/internal/pkg/log"
)

var ErrCampusEmpty = errors.New("キャンパス識別子が未設定です")

type Campus string

const (
	CAMPUS_INVALID = Campus("INVALID")
)

func NewCampus(key string) (Campus, error) {

	key = strings.TrimSpace(key)
	if key == "" {
		return CAMPUS_INVALID, log.WrapErrorWithStackTraceBadRequest(ErrCampusEmpty)
	}

	return Campus(key), nil
}

func (r Campus) Value() string {

	return string(r)
}
