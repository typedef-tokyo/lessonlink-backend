package presenter

import (
	"github.com/samber/lo"
	"github.com/typedef-tokyo/lessonlink-backend/internal/usecase"
)

type (
	ICampusListPresenter interface {
		Present(result *usecase.CampusListOutput) *CampusListResponse
	}

	CampusListPresenter struct {
	}
)

func NewCampusListPresenter() ICampusListPresenter {
	return &CampusListPresenter{}
}

type (
	CampusListResponse struct {
		Msg      string           `json:"msg"`
		Campuses []*CampusListDTO `json:"campuses"`
	}

	CampusListDTO struct {
		Campus     string `json:"campus"`
		CampusName string `json:"campus_name"`
	}
)

func (h *CampusListPresenter) Present(result *usecase.CampusListOutput) *CampusListResponse {

	return &CampusListResponse{

		Campuses: lo.Map(result.CampusList, func(item *usecase.CampusListOutputDTO, _ int) *CampusListDTO {
			return &CampusListDTO{
				Campus:     item.Campus,
				CampusName: item.CampusName,
			}
		}),
	}
}
