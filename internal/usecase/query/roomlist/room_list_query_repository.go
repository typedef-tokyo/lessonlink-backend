package roomlist

import (
	"context"
)

type QueryRoomDTO struct {
	ID        int
	RoomIndex int
	RoomName  string
}

type RoomListQueryRepository interface {
	GetListByCampus(ctx context.Context, campus string) ([]*QueryRoomDTO, error)
}
