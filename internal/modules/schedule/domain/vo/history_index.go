package vo

import (
	"errors"

	"github.com/typedef-tokyo/lessonlink-backend/internal/pkg/log"
)

var ErrHistoryIndexUnderMin = errors.New("履歴インデックスは0以上を設定してください")

type HistoryIndex int

const (
	HISTORY_INDEX_USE_LATEST = HistoryIndex(0)
	HISTORY_INDEX_INVALID    = HistoryIndex(-1)
	HISTORY_INDEX_INITIAL    = HistoryIndex(1)
)

func NewHistoryIndex(index int) (HistoryIndex, error) {

	if index < 0 {
		return HISTORY_INDEX_INVALID, log.WrapErrorWithStackTraceBadRequest(ErrHistoryIndexUnderMin)
	}

	if index == 0 {
		return HISTORY_INDEX_USE_LATEST, nil
	}

	return HistoryIndex(index), nil
}

func (r HistoryIndex) Value() int {

	return int(r)
}

func (r HistoryIndex) Next() HistoryIndex {

	return HistoryIndex(r + 1)
}

func (r HistoryIndex) IsUseLatest() bool {

	return r == HISTORY_INDEX_USE_LATEST
}
