package vo

import (
	"errors"
	"fmt"
	"strings"
	"time"
	"unicode/utf8"

	"github.com/typedef-tokyo/lessonlink-backend/internal/pkg/log"
)

var ErrScheduleTitleEmpty = errors.New("スケジュールタイトルが未設定です")
var ErrScheduleTitleLengthOver = errors.New("スケジュールタイトルに設定できる最大数を超えています")

type ScheduleTitle string

const (
	SCHEDULE_TITLE_INVALID = ScheduleTitle("invalid")
)

func NewScheduleTitle(title string) (ScheduleTitle, error) {

	title = strings.TrimSpace(title)
	if len(title) == 0 {
		return SCHEDULE_TITLE_INVALID, log.WrapErrorWithStackTraceBadRequest(ErrScheduleTitleEmpty)
	}

	const NAME_MAX_LENGTH = 64
	if utf8.RuneCountInString(title) > NAME_MAX_LENGTH {
		return SCHEDULE_TITLE_INVALID, log.WrapErrorWithStackTrace(fmt.Errorf("%w 最大:%d文字", ErrScheduleTitleLengthOver, NAME_MAX_LENGTH))
	}

	return ScheduleTitle(title), nil
}

func NewScheduleTitleInitialCreate() ScheduleTitle {

	return ScheduleTitle(fmt.Sprintf(`%s_スケジュール`, time.Now().Format("20060102_1504")))
}

func (r ScheduleTitle) Value() string {

	return string(r)
}

func (r ScheduleTitle) Duplicate() ScheduleTitle {

	return ScheduleTitle(fmt.Sprintf("%s_コピー", r))
}
