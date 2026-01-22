package roomlist

import (
	"context"

	"github.com/typedef-tokyo/lessonlink-backend/internal/pkg/log"
	campusRepository "github.com/typedef-tokyo/lessonlink-backend/internal/usecase/query/campus"
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
	repositroryCampusQuery campusRepository.CampusQueryRepository
	repositoryQueryRoom    RoomListQueryRepository
}

func NewRoomListQueryInteractor(
	repositroryCampusQuery campusRepository.CampusQueryRepository,
	repositoryQueryRoom RoomListQueryRepository,
) IRoomListQueryInputPort {
	return &RoomListQueryInteractor{
		repositroryCampusQuery: repositroryCampusQuery,
		repositoryQueryRoom:    repositoryQueryRoom,
	}
}

func (r RoomListQueryInteractor) Execute(ctx context.Context, campus string) (*RoomListQueryOutput, error) {

	campusDTO, err := r.repositroryCampusQuery.GetByCampus(ctx, campus)
	if err != nil {
		return nil, log.WrapErrorWithStackTrace(err)
	}

	if campusDTO == nil {
		return nil, log.WrapErrorWithStackTraceNotFound(log.Errorf("指定したキャンパスはありません:%s", campus))
	}

	rooms, err := r.repositoryQueryRoom.GetListByCampus(ctx, campus)
	if err != nil {
		return nil, log.WrapErrorWithStackTrace(err)
	}

	return &RoomListQueryOutput{
		RoomList: rooms,
	}, nil
}
