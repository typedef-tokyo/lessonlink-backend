package presenter

import (
	"cmp"
	"slices"

	"github.com/samber/lo"
	"github.com/typedef-tokyo/lessonlink-backend/internal/configs"
	"github.com/typedef-tokyo/lessonlink-backend/internal/usecase/query/roomlist"
)

type IRoomListPresenter interface {
	Present(result *roomlist.RoomListQueryOutput) *RoomListResponse
}

type RoomListPresenter struct {
}

func NewRoomListPresenter(
	env configs.EnvConfig,
) IRoomListPresenter {
	return &RoomListPresenter{}
}

type (
	RoomListResponse struct {
		Rooms []*RoomListDTO `json:"rooms"`
	}

	RoomListDTO struct {
		RoomIndex int    `json:"room_index"`
		RoomName  string `json:"room_name"`
	}
)

func (h *RoomListPresenter) Present(result *roomlist.RoomListQueryOutput) *RoomListResponse {

	rooms := lo.Map(result.RoomList, func(item *roomlist.QueryRoomDTO, _ int) *RoomListDTO {
		return &RoomListDTO{
			RoomIndex: item.RoomIndex,
			RoomName:  item.RoomName,
		}
	})

	slices.SortFunc(rooms, func(a, b *RoomListDTO) int {
		return cmp.Compare(a.RoomIndex, b.RoomIndex)
	})

	return &RoomListResponse{
		Rooms: rooms,
	}
}
