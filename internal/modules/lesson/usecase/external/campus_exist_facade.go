package external

import (
	"context"
)

type (
	ICampusExistFacade interface {
		Execute(ctx context.Context, inputCampus string) (bool, error)
	}
)
