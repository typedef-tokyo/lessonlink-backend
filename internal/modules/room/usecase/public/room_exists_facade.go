package public

import (
	"context"

	"github.com/typedef-tokyo/lessonlink-backend/internal/modules/room/domain/repository"
	"github.com/typedef-tokyo/lessonlink-backend/internal/modules/room/domain/vo"
	"github.com/typedef-tokyo/lessonlink-backend/internal/pkg/log"
)

type (
	RoomExistsFacade struct {
		repositoryRoom repository.RoomRepository
	}
)

func NewRoomExistsFacade(
	repositoryRoom repository.RoomRepository,
) *RoomExistsFacade {
	return &RoomExistsFacade{
		repositoryRoom: repositoryRoom,
	}
}

func (r RoomExistsFacade) Execute(ctx context.Context, inputCampus string, inputRoomIndexes []int) error {

	campus, err := vo.NewCampus(inputCampus)
	if err != nil {
		return log.WrapErrorWithStackTraceBadRequest(err)
	}

	rooms, err := r.repositoryRoom.FindByCampus(ctx, campus)
	if err != nil {
		return log.WrapErrorWithStackTrace(err)
	}

	for _, inputRoomIndex := range inputRoomIndexes {

		roomIndex, err := vo.NewRoomIndex(inputRoomIndex)
		if err != nil {
			return log.WrapErrorWithStackTrace(err)
		}

		if !rooms.IsExist(campus, roomIndex) {
			return log.WrapErrorWithStackTraceBadRequest(log.Errorf("指定した教室番号は存在しません 番号:%d", roomIndex.Value()))
		}
	}

	return nil
}
