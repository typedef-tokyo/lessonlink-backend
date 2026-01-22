package room

import (
	"fmt"

	"github.com/samber/lo"
	"github.com/typedef-tokyo/lessonlink-backend/internal/domain/vo"
)

type RootRoomModelSlice []*RootRoomModel

func (r RootRoomModelSlice) IsUniq() bool {

	uniqueRoomIndex := lo.UniqBy(r, func(r *RootRoomModel) string {
		return fmt.Sprintf("%v|%v", r.campus, r.roomIndex)
	})

	uniqueRoomName := lo.UniqBy(r, func(r *RootRoomModel) string {
		return fmt.Sprintf("%v|%v", r.campus, r.roomName)
	})

	return len(uniqueRoomIndex) == len(r) && len(uniqueRoomName) == len(r)
}

func (r RootRoomModelSlice) IsExist(campus vo.Campus, roomIndex vo.RoomIndex) bool {

	_, found := lo.Find(r, func(item *RootRoomModel) bool {
		return item.campus == campus && item.roomIndex == roomIndex
	})

	return found
}

type RootRoomModel struct {
	campus    vo.Campus
	roomIndex vo.RoomIndex
	roomName  vo.RoomName
}

func NewRootRoomModel(
	campus vo.Campus,
	roomIndex vo.RoomIndex,
	roomName vo.RoomName,
) *RootRoomModel {

	return &RootRoomModel{
		campus:    campus,
		roomIndex: roomIndex,
		roomName:  roomName,
	}
}

func (r RootRoomModel) Campus() vo.Campus {
	return r.campus
}

func (r RootRoomModel) RoomIndex() vo.RoomIndex {
	return r.roomIndex
}

func (r RootRoomModel) RoomName() vo.RoomName {
	return r.roomName
}
