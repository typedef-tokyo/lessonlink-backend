package vo

import (
	"errors"
	"strings"

	"github.com/typedef-tokyo/lessonlink-backend/internal/pkg/log"
)

var ErrCampusNameEmpty = errors.New("キャンパス名が未設定です")

type CampusName string

const (
	CAMPUS_NAME_INVALID = CampusName("INVALID")
)

func NewCampusName(name string) (CampusName, error) {

	name = strings.TrimSpace(name)
	if name == "" {
		return CAMPUS_NAME_INVALID, log.WrapErrorWithStackTraceBadRequest(ErrCampusNameEmpty)
	}

	return CampusName(name), nil
}

func (r CampusName) Value() string {

	return string(r)
}
