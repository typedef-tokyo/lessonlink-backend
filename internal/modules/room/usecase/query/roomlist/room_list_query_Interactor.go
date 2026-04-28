package roomlist

import (
	"context"

	"github.com/typedef-tokyo/lessonlink-backend/internal/pkg/log"
)

type (
	IRoomListQueryInputPort interface {
		Execute(ctx context.Context, campus string) (*RoomListQueryOutput, error)
	}
)

type (
	RoomListQueryOutput struct {
		RoomList []*QueryRoomDTO
	}
)

type RoomListQueryInteractor struct {
	repositoryQueryRoom RoomListQueryRepository
}

func NewRoomListQueryInteractor(
	repositoryQueryRoom RoomListQueryRepository,
) IRoomListQueryInputPort {
	return &RoomListQueryInteractor{
		repositoryQueryRoom: repositoryQueryRoom,
	}
}

func (r RoomListQueryInteractor) Execute(ctx context.Context, campus string) (*RoomListQueryOutput, error) {

	rooms, err := r.repositoryQueryRoom.GetListByCampus(ctx, campus)
	if err != nil {
		return nil, log.WrapErrorWithStackTrace(err)
	}

	return &RoomListQueryOutput{
		RoomList: rooms,
	}, nil
}
