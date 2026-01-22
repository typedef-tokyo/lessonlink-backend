package campus

import (
	"context"
)

type QueryCampusDTO struct {
	Campus     string
	CampusName string
	OrderIndex int
}

type CampusQueryRepository interface {
	GetByCampus(ctx context.Context, campus string) (*QueryCampusDTO, error)
}
