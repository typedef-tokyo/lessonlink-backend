package vo

import (
	"errors"
	"fmt"
	"strings"
	"unicode/utf8"

	"github.com/typedef-tokyo/lessonlink-backend/internal/pkg/log"
)

var ErrRoomNameEmpty = errors.New("教室名が未設定です")
var ErrRoomNameMaxOver = errors.New("教室名に設定できる最大数を超えています")

type RoomName string

const (
	ROOM_NAME_INVALID = RoomName("INVALID_ROOM_NAME")
)

func NewRoomName(name string) (RoomName, error) {

	name = strings.TrimSpace(name)
	if name == "" {
		return ROOM_NAME_INVALID, log.WrapErrorWithStackTrace(ErrRoomNameEmpty)
	}

	const NAME_MAX_LENGTH = 32
	if utf8.RuneCountInString(name) > NAME_MAX_LENGTH {
		return ROOM_NAME_INVALID, log.WrapErrorWithStackTrace(fmt.Errorf("%w 最大:%d文字", ErrRoomNameMaxOver, NAME_MAX_LENGTH))
	}

	return RoomName(name), nil
}

func (r RoomName) Value() string {

	return string(r)
}
