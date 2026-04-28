package campus

import (
	"context"
	"database/sql"

	"github.com/typedef-tokyo/lessonlink-backend/internal/modules/campus/usecase/query/campus"
	"github.com/typedef-tokyo/lessonlink-backend/internal/pkg/log"
	"github.com/typedef-tokyo/lessonlink-backend/internal/platform/database/rdb"
	"github.com/typedef-tokyo/lessonlink-backend/internal/platform/database/rdb/dto"
)

type CampusQuery struct {
	c *sql.DB
}

func NewCampusQueryRepository(c rdb.IMySQL) campus.CampusQueryRepository {
	return &CampusQuery{c: c.GetConn()}
}

func (f *CampusQuery) GetByCampus(ctx context.Context, inputCampus string) (*campus.QueryCampusDTO, error) {

	record, err := dto.DataCampuses(
		dto.DataCampuseWhere.Campus.EQ(inputCampus),
	).One(ctx, f.c)

	if err != nil && err != sql.ErrNoRows {
		return nil, log.WrapErrorWithStackTraceInternalServerError(err)
	}

	if record == nil {
		return nil, nil
	}

	return &campus.QueryCampusDTO{
		Campus:     record.Campus,
		CampusName: record.CampusName,
		OrderIndex: record.OrderIndex,
	}, nil
}
