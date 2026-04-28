package vo

import (
	"errors"

	"github.com/typedef-tokyo/lessonlink-backend/internal/pkg/log"
)

var ErrItemDivideMinutesUnderMin = errors.New("分割時間は1分以上を設定する必要があります。")

type ItemDivideMinutes int

const (
	ITEM_DIVIDE_MINUTES_INVALID = ItemDivideMinutes(-1)
)

func NewItemDivideMinutes(minute int) (ItemDivideMinutes, error) {

	if minute < 1 {
		return ITEM_DIVIDE_MINUTES_INVALID, log.WrapErrorWithStackTrace(ErrItemDivideMinutesUnderMin)
	}

	return ItemDivideMinutes(minute), nil
}

func (r ItemDivideMinutes) Value() int {
	return int(r)
}
