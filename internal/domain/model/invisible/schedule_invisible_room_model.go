package invisible

import (
	"github.com/samber/lo"
	"github.com/typedef-tokyo/lessonlink-backend/internal/domain/vo"
)

type RootScheduleInvisibleRoomModelSlice []*RootScheduleInvisibleRoomModel

func (r RootScheduleInvisibleRoomModelSlice) IsInvisible(roomIndex vo.RoomIndex) bool {

	_, found := lo.Find(r, func(item *RootScheduleInvisibleRoomModel) bool {
		return item.roomIndex == roomIndex
	})

	return found
}

type RootScheduleInvisibleRoomModel struct {
	scheduleID vo.ScheduleID
	roomIndex  vo.RoomIndex
}

func NewRootScheduleInvisibleRoomModel(
	scheduleID vo.ScheduleID,
	roomIndex vo.RoomIndex,
) *RootScheduleInvisibleRoomModel {

	return &RootScheduleInvisibleRoomModel{
		scheduleID: scheduleID,
		roomIndex:  roomIndex,
	}
}

func (r RootScheduleInvisibleRoomModel) ScheduleID() vo.ScheduleID {
	return r.scheduleID
}

func (r RootScheduleInvisibleRoomModel) RoomIndex() vo.RoomIndex {
	return r.roomIndex
}
