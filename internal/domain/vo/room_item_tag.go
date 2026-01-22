package vo

import (
	"errors"
	"strings"

	"github.com/typedef-tokyo/lessonlink-backend/internal/pkg/log"
)

var ErrRoomItemTagEmpty = errors.New("アイテムの種類が未設定です")
var ErrRoomItemTagInvalid = errors.New("アイテムの種類が不正です")

type RoomItemTag string

const (
	ROOM_ITEM_TAG_INVALID  = RoomItemTag("invalid")
	ROOM_ITEM_TAG_LESSON   = RoomItemTag("lesson")
	ROOM_ITEM_TAG_CLEANING = RoomItemTag("cleaning")
)

func NewRoomItemTag(itemTag string) (RoomItemTag, error) {

	itemTag = strings.TrimSpace(itemTag)
	if itemTag == "" {
		return ROOM_ITEM_TAG_INVALID, log.WrapErrorWithStackTrace(ErrRoomItemTagEmpty)
	}

	switch itemTag {
	case "lesson":
		return ROOM_ITEM_TAG_LESSON, nil
	case "cleaning":
		return ROOM_ITEM_TAG_CLEANING, nil
	default:
		return ROOM_ITEM_TAG_INVALID, log.WrapErrorWithStackTrace(ErrRoomItemTagInvalid)
	}
}

func (r RoomItemTag) Value() string {
	return string(r)
}

func (r RoomItemTag) IsLesson() bool {
	return r == ROOM_ITEM_TAG_LESSON
}

func (r RoomItemTag) IsCleaning() bool {
	return r == ROOM_ITEM_TAG_CLEANING
}
