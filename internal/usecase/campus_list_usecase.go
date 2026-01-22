package usecase

import (
	"context"

	"github.com/samber/lo"
	"github.com/typedef-tokyo/lessonlink-backend/internal/domain/model/campus"
	"github.com/typedef-tokyo/lessonlink-backend/internal/domain/repository"
	"github.com/typedef-tokyo/lessonlink-backend/internal/pkg/log"
)

type (
	ICampusListInputPort interface {
		Execute(ctx context.Context) (*CampusListOutput, error)
	}
)

type (
	CampusListOutput struct {
		CampusList []*CampusListOutputDTO
	}

	CampusListOutputDTO struct {
		Campus     string
		CampusName string
	}
)

type CampusListInteractor struct {
	repositoryCampus repository.CampusRepository
}

func NewCampusListInteractor(
	repositoryCampus repository.CampusRepository,
) ICampusListInputPort {
	return &CampusListInteractor{
		repositoryCampus: repositoryCampus,
	}
}

func (r *CampusListInteractor) Execute(ctx context.Context) (*CampusListOutput, error) {

	// リポジトリから工場を全て取得
	campuses, err := r.repositoryCampus.FindAll(ctx)
	if err != nil {
		return nil, log.WrapErrorWithStackTrace(err)
	}

	// 並び替えを実行
	campuses.SortByOrder()

	return &CampusListOutput{
		CampusList: lo.Map(campuses, func(item *campus.RootCumpusModel, _ int) *CampusListOutputDTO {
			return &CampusListOutputDTO{
				Campus:     item.Campus().Value(),
				CampusName: item.CampusName().Value(),
			}
		}),
	}, nil
}
