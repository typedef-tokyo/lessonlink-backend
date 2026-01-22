package repository

import (
	"context"
	"database/sql"

	"github.com/typedef-tokyo/lessonlink-backend/internal/domain/model/room"
	"github.com/typedef-tokyo/lessonlink-backend/internal/domain/vo"
)

type RoomRepository interface {
	Save(ctx context.Context, tx *sql.Tx, campus vo.Campus, slice room.RootRoomModelSlice) error
	FindByCampus(ctx context.Context, campus vo.Campus) (room.RootRoomModelSlice, error)
}
