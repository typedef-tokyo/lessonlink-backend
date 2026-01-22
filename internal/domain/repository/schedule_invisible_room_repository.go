package repository

import (
	"context"
	"database/sql"

	"github.com/typedef-tokyo/lessonlink-backend/internal/domain/model/invisible"
	"github.com/typedef-tokyo/lessonlink-backend/internal/domain/vo"
)

type ScheduleInvisibleRoomRepository interface {
	Save(ctx context.Context, tx *sql.Tx, sheduleID vo.ScheduleID, models []*invisible.RootScheduleInvisibleRoomModel) error
	FindBySheduleID(ctx context.Context, sheduleID vo.ScheduleID) (invisible.RootScheduleInvisibleRoomModelSlice, error)
}
