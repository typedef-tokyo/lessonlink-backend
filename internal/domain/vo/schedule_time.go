package vo

import (
	"errors"
	"fmt"

	"github.com/typedef-tokyo/lessonlink-backend/internal/pkg/log"
)

var ErrScheduleTimeInvalidStartHour = errors.New("スケジュール開始時間は0時～23時で設定してください")
var ErrScheduleTimeInvalidEndHour = errors.New("スケジュール終了時間は0時～23時で設定してください")
var ErrScheduleTimeInvalidIntegrity = errors.New("スケジュール開始時間は終了時間前である必要があります")

type ScheduleTime struct {
	scheduleStartTime int
	scheduleEndTime   int
}

func NewScheduleTime(startTime int, endTime int) (ScheduleTime, error) {

	if startTime < 0 || startTime > 23 {
		return ScheduleTime{}, log.WrapErrorWithStackTrace(ErrScheduleTimeInvalidStartHour)
	}

	if endTime < 0 || endTime > 23 {
		return ScheduleTime{}, log.WrapErrorWithStackTrace(ErrScheduleTimeInvalidEndHour)
	}

	if startTime >= endTime {
		return ScheduleTime{}, log.WrapErrorWithStackTrace(fmt.Errorf("%w 開始:%d時, 終了:%d時", ErrScheduleTimeInvalidIntegrity, startTime, endTime))

	}

	return ScheduleTime{scheduleStartTime: startTime, scheduleEndTime: endTime}, nil
}

func (r ScheduleTime) IsWithinTimeRange(time ScheduleLessonTime) bool {

	hour, minutes := time.Value()
	return r.scheduleStartTime <= time.scheduleItemTimeHour && (r.scheduleEndTime > hour || (r.scheduleEndTime == hour && minutes == 0))
}

func (r ScheduleTime) Value() (int, int) {

	return r.scheduleStartTime, r.scheduleEndTime
}

func (r ScheduleTime) EndTimeValueMinutes() int {
	return r.scheduleEndTime * 60
}
