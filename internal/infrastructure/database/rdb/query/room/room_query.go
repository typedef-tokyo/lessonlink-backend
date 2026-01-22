package room

import (
	"context"
	"database/sql"

	"github.com/typedef-tokyo/lessonlink-backend/internal/infrastructure/database/rdb"
	"github.com/typedef-tokyo/lessonlink-backend/internal/infrastructure/database/rdb/dto"
	"github.com/typedef-tokyo/lessonlink-backend/internal/pkg/log"
	roomlist "github.com/typedef-tokyo/lessonlink-backend/internal/usecase/query/roomlist"
)

type RoomQuery struct {
	c *sql.DB
}

func NewRoomQueryRepository(c rdb.IMySQL) roomlist.RoomListQueryRepository {
	return &RoomQuery{c: c.GetConn()}
}

func (f *RoomQuery) GetListByCampus(ctx context.Context, campus string) ([]*roomlist.QueryRoomDTO, error) {

	roomRecords, err := dto.DataRooms(
		dto.DataRoomWhere.Campus.EQ(campus),
	).All(ctx, f.c)

	if err != nil {
		return nil, log.WrapErrorWithStackTraceInternalServerError(err)
	}

	return f.toList(roomRecords), nil
}

func (f *RoomQuery) toList(roomRecords []*dto.DataRoom) []*roomlist.QueryRoomDTO {

	roomList := make([]*roomlist.QueryRoomDTO, 0, len(roomRecords))

	for _, roomRecord := range roomRecords {

		room := &roomlist.QueryRoomDTO{
			ID:        roomRecord.ID,
			RoomIndex: roomRecord.RoomIndex,
			RoomName:  roomRecord.Name,
		}

		roomList = append(roomList, room)
	}

	return roomList
}
