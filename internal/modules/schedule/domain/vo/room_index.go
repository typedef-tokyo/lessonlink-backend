package vo

import (
	"errors"

	"github.com/typedef-tokyo/lessonlink-backend/internal/pkg/log"
)

var ErrRoomIndexUnderMin = errors.New("教室番号は1以上を指定してください")

type RoomIndex int

const (
	ROOM_INDEX_INVALID = RoomIndex(-1)
)

func NewRoomIndex(index int) (RoomIndex, error) {

	if index <= 0 {
		return ROOM_INDEX_INVALID, log.WrapErrorWithStackTrace(ErrRoomIndexUnderMin)
	}

	return RoomIndex(index), nil
}

func (r RoomIndex) Value() int {

	return int(r)
}
