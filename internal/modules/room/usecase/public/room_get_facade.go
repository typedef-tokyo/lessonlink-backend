package public

import (
	"context"

	"github.com/samber/lo"
	"github.com/typedef-tokyo/lessonlink-backend/internal/modules/room/domain/model/room"
	"github.com/typedef-tokyo/lessonlink-backend/internal/modules/room/domain/repository"
	"github.com/typedef-tokyo/lessonlink-backend/internal/modules/room/domain/vo"
	"github.com/typedef-tokyo/lessonlink-backend/internal/pkg/log"
)

type (
	RoomsGetOutDTO struct {
		Rooms []RoomGetOutDTO
	}

	RoomGetOutDTO struct {
		Index int
		Name  string
	}
)

///////////////

type (
	RoomGetFacade struct {
		repositoryRoom repository.RoomRepository
	}
)

func NewRoomGetFacade(
	repositoryRoom repository.RoomRepository,
) *RoomGetFacade {
	return &RoomGetFacade{
		repositoryRoom: repositoryRoom,
	}
}

func (r RoomGetFacade) Execute(ctx context.Context, inputCampus string) (*RoomsGetOutDTO, error) {

	campus, err := vo.NewCampus(inputCampus)
	if err != nil {
		return nil, log.WrapErrorWithStackTraceBadRequest(err)
	}

	rooms, err := r.repositoryRoom.FindByCampus(ctx, campus)
	if err != nil {
		return nil, log.WrapErrorWithStackTrace(err)
	}

	return &RoomsGetOutDTO{
		Rooms: lo.Map(rooms, func(item *room.RootRoomModel, _ int) RoomGetOutDTO {
			return RoomGetOutDTO{
				Index: item.RoomIndex().Value(),
				Name:  item.RoomName().Value(),
			}
		}),
	}, nil
}
