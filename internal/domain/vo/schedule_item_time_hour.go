package vo

import (
	"errors"
	"fmt"

	"github.com/typedef-tokyo/lessonlink-backend/internal/pkg/log"
)

var ErrScheduleLessonTimeInvalidHour = errors.New("時間は0時～23時で設定してください")
var ErrScheduleLessonTimeInvalidMinutes = errors.New("分は0分～59分で設定してください")
var ErrScheduleLessonTimeMinUnderMinutes = errors.New("分は0以上である必要があります")
var ErrScheduleLessonTimeMaxOverMinutes = errors.New("分に設定できる最大値を超えています")

type ScheduleLessonTime struct {
	scheduleItemTimeHour    int
	scheduleItemTimeMinutes int
}

func NewScheduleLessonTime(hour int, minutes int) (ScheduleLessonTime, error) {

	if hour < 0 || hour > 23 {
		return ScheduleLessonTime{}, log.WrapErrorWithStackTrace(ErrScheduleLessonTimeInvalidHour)
	}

	if minutes < 0 || minutes > 59 {
		return ScheduleLessonTime{}, log.WrapErrorWithStackTrace(ErrScheduleLessonTimeInvalidMinutes)
	}

	return ScheduleLessonTime{scheduleItemTimeHour: hour, scheduleItemTimeMinutes: minutes}, nil
}

func NewScheduleLessonTimeFromMinutes(minutes int) (ScheduleLessonTime, error) {

	const ONE_HOUR = 60

	if minutes < 0 {
		return ScheduleLessonTime{}, log.WrapErrorWithStackTrace(ErrScheduleLessonTimeMinUnderMinutes)
	}

	oneDayMinute := ONE_HOUR * 24
	if minutes > (oneDayMinute) {
		return ScheduleLessonTime{}, log.WrapErrorWithStackTrace(fmt.Errorf("%w 最大:%d", ErrScheduleLessonTimeMaxOverMinutes, oneDayMinute))
	}

	return ScheduleLessonTime{
		scheduleItemTimeHour:    minutes / ONE_HOUR,
		scheduleItemTimeMinutes: minutes % ONE_HOUR,
	}, nil
}

func (r ScheduleLessonTime) Value() (int, int) {

	return r.scheduleItemTimeHour, r.scheduleItemTimeMinutes
}

func (r ScheduleLessonTime) ValueMinutes() int {
	return r.scheduleItemTimeHour*60 + r.scheduleItemTimeMinutes
}
