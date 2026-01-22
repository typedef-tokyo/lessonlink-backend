package schedule

import (
	"context"
	"database/sql"

	"github.com/aarondl/sqlboiler/v4/queries/qm"
	"github.com/typedef-tokyo/lessonlink-backend/internal/infrastructure/database/rdb"
	"github.com/typedef-tokyo/lessonlink-backend/internal/infrastructure/database/rdb/dto"
	"github.com/typedef-tokyo/lessonlink-backend/internal/pkg/log"
	schedulelist "github.com/typedef-tokyo/lessonlink-backend/internal/usecase/query/schedule_list"
)

type ScheduleQuery struct {
	c *sql.DB
}

func NewScheduleQueryRepository(c rdb.IMySQL) schedulelist.ScheduleListQueryRepository {
	return &ScheduleQuery{c: c.GetConn()}
}

func (f *ScheduleQuery) GetListByCampus(ctx context.Context, campus string) ([]*schedulelist.QueryScheduleDTO, error) {

	scheduleDTOs, err := dto.TBLSchedules(
		dto.TBLScheduleWhere.Campus.EQ(campus),
		qm.Load(dto.TBLScheduleRels.CreateUserTBLUser),
		qm.Load(dto.TBLScheduleRels.LastUpdateUserTBLUser),
	).All(ctx, f.c)

	if err != nil {
		return nil, log.WrapErrorWithStackTraceInternalServerError(err)
	}

	return f.toList(scheduleDTOs), nil
}

func (f *ScheduleQuery) toList(scheduleDTOs []*dto.TBLSchedule) []*schedulelist.QueryScheduleDTO {

	scheduleList := make([]*schedulelist.QueryScheduleDTO, 0, len(scheduleDTOs))

	for _, scheduleDTO := range scheduleDTOs {

		schedule := &schedulelist.QueryScheduleDTO{
			ID:                 scheduleDTO.ID,
			Campus:             scheduleDTO.Campus,
			Title:              scheduleDTO.Title,
			CreateUserName:     "不明なユーザー",
			LastUpdateUserName: "不明なユーザー",
			CreateUser:         scheduleDTO.CreateUser,
			UpdatedAt:          scheduleDTO.UpdatedAt,
			CreatedAt:          scheduleDTO.CreatedAt,
		}

		if scheduleDTO.R != nil {

			rel := scheduleDTO.R

			if rel.CreateUserTBLUser != nil {
				schedule.CreateUserName = rel.CreateUserTBLUser.Name
			}

			if rel.LastUpdateUserTBLUser != nil {
				schedule.LastUpdateUserName = rel.LastUpdateUserTBLUser.Name
			}
		}

		scheduleList = append(scheduleList, schedule)
	}

	return scheduleList
}
