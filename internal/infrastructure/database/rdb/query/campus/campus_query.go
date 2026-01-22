package campus

import (
	"context"
	"database/sql"

	"github.com/typedef-tokyo/lessonlink-backend/internal/infrastructure/database/rdb"
	"github.com/typedef-tokyo/lessonlink-backend/internal/infrastructure/database/rdb/dto"
	"github.com/typedef-tokyo/lessonlink-backend/internal/pkg/log"
	campusRepository "github.com/typedef-tokyo/lessonlink-backend/internal/usecase/query/campus"
)

type CampusQuery struct {
	c *sql.DB
}

func NewCampusQueryRepository(c rdb.IMySQL) campusRepository.CampusQueryRepository {
	return &CampusQuery{c: c.GetConn()}
}

func (f *CampusQuery) GetByCampus(ctx context.Context, campus string) (*campusRepository.QueryCampusDTO, error) {

	record, err := dto.DataCampuses(
		dto.DataCampuseWhere.Campus.EQ(campus),
	).One(ctx, f.c)

	if err != nil && err != sql.ErrNoRows {
		return nil, log.WrapErrorWithStackTraceInternalServerError(err)
	}

	if record == nil {
		return nil, nil
	}

	return &campusRepository.QueryCampusDTO{
		Campus:     record.Campus,
		CampusName: record.CampusName,
		OrderIndex: record.OrderIndex,
	}, nil
}
