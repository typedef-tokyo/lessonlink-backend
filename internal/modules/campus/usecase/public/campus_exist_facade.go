package public

import (
	"context"

	"github.com/typedef-tokyo/lessonlink-backend/internal/modules/campus/domain/repository"
	"github.com/typedef-tokyo/lessonlink-backend/internal/modules/campus/domain/vo"
	"github.com/typedef-tokyo/lessonlink-backend/internal/pkg/log"
)

type CampusExistFacade struct {
	repositoryCampus repository.CampusRepository
}

func NewCampusExistFacade(
	repositoryCampus repository.CampusRepository,
) *CampusExistFacade {
	return &CampusExistFacade{
		repositoryCampus: repositoryCampus,
	}
}

func (r *CampusExistFacade) Execute(ctx context.Context, inputCampus string) (bool, error) {

	campus, err := vo.NewCampus(inputCampus)
	if err != nil {
		return false, log.WrapErrorWithStackTraceBadRequest(err)
	}

	campuses, err := r.repositoryCampus.FindAll(ctx)
	if err != nil {
		return false, log.WrapErrorWithStackTrace(err)
	}

	return campuses.IsExist(campus), nil
}
