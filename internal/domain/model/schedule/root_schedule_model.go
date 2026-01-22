package schedule

import (
	"errors"
	"fmt"
	"time"

	"github.com/samber/lo"
	"github.com/typedef-tokyo/lessonlink-backend/internal/domain/vo"
	"github.com/typedef-tokyo/lessonlink-backend/internal/pkg/log"
)

type RootScheduleModel struct {
	id             vo.ScheduleID
	campus         vo.Campus
	title          vo.ScheduleTitle
	historyIndex   vo.HistoryIndex
	createUser     vo.UserID
	lastUpdateUser vo.UserID
	items          ScheduleItemModelSlice
	roomItems      ScheduleRoomItemModelSlice
	scheduleTime   vo.ScheduleTime
	createdAt      time.Time
	updatedAt      time.Time
}

func NewRootScheduleModel(
	id vo.ScheduleID,
	campus vo.Campus,
	title vo.ScheduleTitle,
	historyIndex vo.HistoryIndex,
	createUser vo.UserID,
	lastUpdateUser vo.UserID,
	items ScheduleItemModelSlice,
	roomItems ScheduleRoomItemModelSlice,
	scheduleTime vo.ScheduleTime,
	createdAt time.Time,
	updatedAt time.Time,
) *RootScheduleModel {

	return &RootScheduleModel{
		id:             id,
		campus:         campus,
		title:          title,
		historyIndex:   historyIndex,
		createUser:     createUser,
		lastUpdateUser: lastUpdateUser,
		items:          items,
		roomItems:      roomItems,
		scheduleTime:   scheduleTime,
		createdAt:      createdAt,
		updatedAt:      updatedAt,
	}
}

func NewCreateRootScheduleModel(
	campus vo.Campus,
	createUser vo.UserID,
	scheduleTime vo.ScheduleTime,
) *RootScheduleModel {

	now := time.Now()

	return &RootScheduleModel{
		id:             vo.NewCreateInitialScheduleID(),
		campus:         campus,
		title:          vo.NewScheduleTitleInitialCreate(),
		historyIndex:   vo.HISTORY_INDEX_INITIAL,
		createUser:     createUser,
		lastUpdateUser: createUser,
		items:          []*ScheduleItemModel{},
		roomItems:      []*ScheduleRoomItemModel{},
		scheduleTime:   scheduleTime,
		createdAt:      now,
		updatedAt:      now,
	}
}

func (r RootScheduleModel) ID() vo.ScheduleID {
	return r.id
}

func (r RootScheduleModel) Campus() vo.Campus {
	return r.campus
}

func (r RootScheduleModel) Title() vo.ScheduleTitle {
	return r.title
}

func (r RootScheduleModel) HistoryIndex() vo.HistoryIndex {
	return r.historyIndex
}

func (r RootScheduleModel) CreateUser() vo.UserID {
	return r.createUser
}

func (r RootScheduleModel) LastUpdateUser() vo.UserID {
	return r.lastUpdateUser
}

func (r RootScheduleModel) Items() ScheduleItemModelSlice {
	return r.items
}

func (r RootScheduleModel) RoomItems() ScheduleRoomItemModelSlice {
	return r.roomItems
}

// func (r RootScheduleModel) ScheduleLesson() []*ScheduleLesson {
// 	return r.scheduleLesson
// }

func (r RootScheduleModel) ScheduleTime() vo.ScheduleTime {
	return r.scheduleTime
}

func (r RootScheduleModel) UpdatedAt() time.Time {
	return r.updatedAt
}

func (r RootScheduleModel) FilterScheduleItemByLessonID(lessonID vo.LessonID) ScheduleItemModelSlice {

	return r.items.filterByLessonID(lessonID)
}

func (r *RootScheduleModel) ChangeTitle(title vo.ScheduleTitle) {

	r.title = title
}

func (r *RootScheduleModel) RoomItemMove(item *ScheduleRoomItemModel) error {

	if !r.scheduleTime.IsWithinTimeRange(item.startTime) {
		return log.WrapErrorWithStackTrace(errors.New("講座開始時刻がスケジュール時刻以より前です"))
	}

	if !r.scheduleTime.IsWithinTimeRange(item.endTime) {
		return log.WrapErrorWithStackTrace(errors.New("講座終了時刻がスケジュール時刻より後です"))
	}

	removedItems := r.items.removeByIdentifier(item.identifier)
	replacedRoomItems := r.roomItems.replaceItem(item)

	err := r.validateUniqueIdentifiers(removedItems, replacedRoomItems)
	if err != nil {
		return log.WrapErrorWithStackTrace(err)
	}

	r.items = removedItems
	r.roomItems = replacedRoomItems

	return nil
}

func (r *RootScheduleModel) ItemReturnList(item *ScheduleItemModel) error {

	roomItem, found := r.roomItems.findByIdentifier(item.identifier)
	if !found {
		return log.WrapErrorWithStackTrace(log.Errorf("移動元アイテムが見つかりません"))
	}

	addItems := r.items
	if roomItem.itemTag.IsLesson() {
		addItems = r.items.addItem(item)
	}

	removedRoomItems := r.roomItems.removeByIdentifier(item.identifier)

	err := r.validateUniqueIdentifiers(addItems, removedRoomItems)
	if err != nil {
		return log.WrapErrorWithStackTrace(err)
	}

	r.items = addItems
	r.roomItems = removedRoomItems

	return nil
}

func (r RootScheduleModel) validateUniqueIdentifiers(items []*ScheduleItemModel, roomItems []*ScheduleRoomItemModel) error {

	allIdentifiers := append(
		lo.Map(items, func(i *ScheduleItemModel, _ int) vo.Identifier { return i.identifier }),
		lo.Map(roomItems, func(i *ScheduleRoomItemModel, _ int) vo.Identifier { return i.identifier })...,
	)

	dupes := lo.FindDuplicates(allIdentifiers)
	if len(dupes) > 0 {
		return log.WrapErrorWithStackTrace(fmt.Errorf("アイテムが重複しています: %v", dupes))
	}

	return nil
}

func (r *RootScheduleModel) ModifyEditing(historyIndex vo.HistoryIndex, lastUpdateUser vo.UserID) {

	r.historyIndex = historyIndex.Next()
	r.lastUpdateUser = lastUpdateUser
}

func (r *RootScheduleModel) ModifySaving(historyIndex vo.HistoryIndex, lastUpdateUser vo.UserID) {

	r.historyIndex = historyIndex
	r.lastUpdateUser = lastUpdateUser
}

func (r *RootScheduleModel) RoomItemShift(roomIndex vo.RoomIndex) error {

	shiftItems, err := r.roomItems.shiftedItems(roomIndex, r.scheduleTime)
	if err != nil {
		return log.WrapErrorWithStackTrace(err)
	}

	r.roomItems = r.roomItems.removeByRoomIndex(roomIndex)
	r.roomItems = append(r.roomItems, shiftItems...)

	return nil
}

func (r *RootScheduleModel) ItemDivide(lessonID vo.LessonID, initialLessonDuration vo.LessonDuration, identifier vo.Identifier, divideMinutes vo.ItemDivideMinutes) error {

	var divide = func(item *ScheduleItemModel) error {

		divideDurationFrom, divideDurationTo, err := item.duration.Divide(divideMinutes)
		if err != nil {
			return log.WrapErrorWithStackTrace(err)
		}

		divideFrom := NewScheduleItemModel(item.lessonID, identifier, divideDurationFrom)
		divideTo := NewScheduleItemModel(item.lessonID, vo.NewIdentifierGenerate(), divideDurationTo)

		items := r.items.removeByIdentifier(item.identifier)
		r.items = append(items, divideFrom, divideTo)

		return nil
	}

	var newDivideItem = func() error {

		newItem := NewScheduleItemModel(lessonID, identifier, initialLessonDuration)

		err := divide(newItem)
		if err != nil {
			return log.WrapErrorWithStackTrace(err)
		}

		return nil
	}

	item, found := r.items.findByIdentifier(identifier)
	if found {

		err := divide(item)
		if err != nil {
			return log.WrapErrorWithStackTrace(err)
		}

	} else {

		roomItem, found := r.roomItems.findByIdentifier(identifier)
		if !found {
			return newDivideItem()
		}

		divideDurationFrom, divideDurationTo, err := roomItem.duration.Divide(divideMinutes)
		if err != nil {
			return log.WrapErrorWithStackTrace(err)
		}

		newEndTime, err := vo.NewScheduleLessonTimeFromMinutes(roomItem.startTime.ValueMinutes() + divideDurationFrom.Value())
		if err != nil {
			return log.WrapErrorWithStackTrace(err)
		}
		divideFrom := NewScheduleRoomItemModel(
			roomItem.itemTag,
			roomItem.lessonID,
			roomItem.identifier,
			divideDurationFrom,
			roomItem.startTime,
			newEndTime,
			roomItem.roomIndex,
		)

		divideTo := NewScheduleRoomItemModel(
			roomItem.itemTag,
			roomItem.lessonID,
			vo.NewIdentifierGenerate(),
			divideDurationTo,
			newEndTime,
			roomItem.endTime,
			roomItem.roomIndex,
		)

		roomItems := r.roomItems.removeByIdentifier(roomItem.identifier)
		r.roomItems = append(roomItems, divideFrom, divideTo)
	}

	return nil
}

func (r *RootScheduleModel) ItemJoin(joinFromID vo.Identifier, joinToID vo.Identifier) error {

	joinFromItem, foundItemFrom := r.items.findByIdentifier(joinFromID)
	joinToItem, foundItemTo := r.items.findByIdentifier(joinToID)

	joinFromRoomItem, foundRoomItemFrom := r.roomItems.findByIdentifier(joinFromID)
	joinToRoomItem, foundRoomItemTo := r.roomItems.findByIdentifier(joinToID)

	if foundItemFrom != foundItemTo || foundRoomItemFrom != foundRoomItemTo {
		return log.WrapErrorWithStackTrace(errors.New("講座の結合は一覧または教室に配置している同士のみ可能です"))
	}

	if !foundItemFrom && !foundRoomItemFrom {
		return log.WrapErrorWithStackTrace(errors.New("結合対象の講座がみつかりません"))
	}

	if foundItemFrom {

		if joinFromItem.lessonID != joinToItem.lessonID {
			return log.WrapErrorWithStackTrace(errors.New("結合は同じ講座のみできます"))
		}

		joinedItem := NewScheduleItemModel(
			joinToItem.lessonID,
			joinToItem.identifier,
			joinToItem.duration.Add(joinFromItem.duration),
		)

		removedItems := r.items.removeByIdentifier(joinFromItem.identifier).removeByIdentifier(joinToItem.identifier)
		r.items = append(removedItems, joinedItem)

	} else {

		if joinFromRoomItem.lessonID != joinToRoomItem.lessonID {
			return log.WrapErrorWithStackTrace(errors.New("結合は同じ講座のみできます"))
		}

		joinDuration := joinToRoomItem.duration.Add(joinFromRoomItem.duration)
		joinEndTime, err := vo.NewScheduleLessonTimeFromMinutes(joinToRoomItem.startTime.ValueMinutes() + joinDuration.Value())
		if err != nil {
			return log.WrapErrorWithStackTrace(err)
		}

		joinStartTime := joinToRoomItem.startTime
		if !r.scheduleTime.IsWithinTimeRange(joinEndTime) {

			_, scheduleEndTime := r.scheduleTime.Value()
			diffEndTimeMinute := joinEndTime.ValueMinutes() - scheduleEndTime
			joinStartTime, err = vo.NewScheduleLessonTimeFromMinutes(joinStartTime.ValueMinutes() - diffEndTimeMinute)
			if err != nil {
				return log.WrapErrorWithStackTrace(err)
			}

			joinEndTime, err = vo.NewScheduleLessonTimeFromMinutes(joinStartTime.ValueMinutes() - diffEndTimeMinute)
			if err != nil {
				return log.WrapErrorWithStackTrace(err)
			}
		}

		if !r.scheduleTime.IsWithinTimeRange(joinStartTime) {
			return log.WrapErrorWithStackTrace(errors.New("結合した講座は利用時間の範囲に収まりません"))
		}

		joinedItem := NewScheduleRoomItemModel(
			joinToRoomItem.itemTag,
			joinToRoomItem.lessonID,
			joinToRoomItem.identifier,
			joinDuration,
			joinStartTime,
			joinEndTime,
			joinToRoomItem.roomIndex,
		)

		removedRoomsItems := r.roomItems.removeByIdentifier(joinFromRoomItem.identifier).removeByIdentifier(joinToRoomItem.identifier)
		r.roomItems = append(removedRoomsItems, joinedItem)
	}

	return nil
}

func (r *RootScheduleModel) ChangeScheduleTime(newScheduleTime vo.ScheduleTime) error {

	if !lo.EveryBy(r.roomItems, func(item *ScheduleRoomItemModel) bool {
		return newScheduleTime.IsWithinTimeRange(item.startTime)
	}) {
		return log.WrapErrorWithStackTrace(errors.New("スケジュール開始時刻前に配置されている講座があります"))
	}

	if !lo.EveryBy(r.roomItems, func(item *ScheduleRoomItemModel) bool {
		return newScheduleTime.IsWithinTimeRange(item.endTime)
	}) {
		return log.WrapErrorWithStackTrace(errors.New("スケジュール開始時刻後に配置されている講座があります"))
	}

	r.scheduleTime = newScheduleTime
	r.historyIndex = vo.HISTORY_INDEX_INITIAL

	return nil
}

func (r RootScheduleModel) Duplicate(duplicateUser vo.UserID) *RootScheduleModel {

	duplicateSchedule := &RootScheduleModel{}
	*duplicateSchedule = r
	duplicateSchedule.id = vo.NewCreateInitialScheduleID()
	duplicateSchedule.title = r.title.Duplicate()
	duplicateSchedule.historyIndex = vo.HISTORY_INDEX_INITIAL
	duplicateSchedule.createUser = duplicateUser
	duplicateSchedule.lastUpdateUser = duplicateUser

	duplicateSchedule.items = lo.Map(r.items, func(item *ScheduleItemModel, _ int) *ScheduleItemModel {
		return item.duplicate()
	})

	duplicateSchedule.roomItems = lo.Map(r.roomItems, func(item *ScheduleRoomItemModel, _ int) *ScheduleRoomItemModel {
		return item.duplicate()
	})

	return duplicateSchedule
}
