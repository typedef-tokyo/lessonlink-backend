package infrastructure

import (
	"context"
	"database/sql"
	"errors"

	"github.com/typedef-tokyo/lessonlink-backend/internal/modules/campus/domain/model/campus"
	"github.com/typedef-tokyo/lessonlink-backend/internal/modules/campus/domain/repository"
	"github.com/typedef-tokyo/lessonlink-backend/internal/modules/campus/domain/vo"
	"github.com/typedef-tokyo/lessonlink-backend/internal/pkg/log"
	"github.com/typedef-tokyo/lessonlink-backend/internal/platform/database/rdb"
	"github.com/typedef-tokyo/lessonlink-backend/internal/platform/database/rdb/dto"
	"github.com/typedef-tokyo/lessonlink-backend/internal/platform/utility"
)

type Campus struct {
	c *sql.DB
}

func NewCampusRepository(c rdb.IMySQL) repository.CampusRepository {
	return &Campus{c: c.GetConn()}
}

func (f *Campus) FindAll(ctx context.Context) (campus.CampusModelSlice, error) {

	campusDTOList, err := dto.DataCampuses().All(ctx, f.c)

	if err != nil {
		return nil, log.WrapErrorWithStackTraceInternalServerError(err)
	}

	campuseModels := make([]*campus.RootCumpusModel, 0, len(campusDTOList))
	for _, campusDTO := range campusDTOList {

		model, err := f.toModel(campusDTO)
		if err != nil {
			return nil, log.WrapErrorWithStackTrace(err)
		}

		campuseModels = append(campuseModels, model)
	}

	return campuseModels, nil
}

func (f *Campus) toModel(dto *dto.DataCampuse) (*campus.RootCumpusModel, error) {

	var campusVO vo.Campus
	var campusName vo.CampusName

	var errs error

	errs = errors.Join(errs, utility.SetVOConstructor(&campusVO, vo.NewCampus, dto.Campus))
	errs = errors.Join(errs, utility.SetVOConstructor(&campusName, vo.NewCampusName, dto.CampusName))

	if errs != nil {
		return nil, log.WrapErrorWithStackTraceInternalServerError(log.Errorf("%s", errs.Error()))
	}

	return campus.NewRootCampus(
		campusVO,
		campusName,
		dto.OrderIndex,
	), nil

}
