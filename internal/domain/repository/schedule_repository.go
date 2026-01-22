package repository

import (
	"context"
	"database/sql"

	"github.com/typedef-tokyo/lessonlink-backend/internal/domain/model/schedule"
	"github.com/typedef-tokyo/lessonlink-backend/internal/domain/vo"
)

type ScheduleRepository interface {
	Save(ctx context.Context, tx *sql.Tx, rootModel *schedule.RootScheduleModel) (vo.ScheduleID, error)
	Delete(ctx context.Context, tx *sql.Tx, scheduleID vo.ScheduleID, deleteUserID vo.UserID) error
	FindByIDWithLock(ctx context.Context, tx *sql.Tx, scheduleID vo.ScheduleID) (*schedule.RootScheduleModel, error)
	FindByIDWithLockHistoryIndex(ctx context.Context, tx *sql.Tx, scheduleID vo.ScheduleID, historyIndex vo.HistoryIndex) (*schedule.RootScheduleModel, error)
	FindByIDWithHistoryIndex(ctx context.Context, scheduleID vo.ScheduleID, historyIndex vo.HistoryIndex) (*schedule.RootScheduleModel, error)
	FindByID(ctx context.Context, scheduleID vo.ScheduleID) (*schedule.RootScheduleModel, error)
}
