package vo

import (
	"errors"

	"github.com/typedef-tokyo/lessonlink-backend/internal/pkg/log"
)

var ErrScheduleIDUnderMin = errors.New("スケジュールIDは0以上を設定してください")

type ScheduleID int

const (
	SCHEDULE_ID_INVALID = ScheduleID(-1)
	SCHEDULE_ID_INITIAL = ScheduleID(0)
)

func NewScheduleID(id int) (ScheduleID, error) {

	if id < 0 {
		return SCHEDULE_ID_INVALID, log.WrapErrorWithStackTraceBadRequest(ErrScheduleIDUnderMin)
	}

	return ScheduleID(id), nil
}

func NewCreateInitialScheduleID() ScheduleID {

	return SCHEDULE_ID_INITIAL
}

func (r ScheduleID) Value() int {

	return int(r)
}

func (r ScheduleID) IsInitial() bool {

	return r == SCHEDULE_ID_INITIAL
}
