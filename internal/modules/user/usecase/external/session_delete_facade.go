package external

import (
	"context"
)

type (
	ISessionDeleteFacade interface {
		Execute(ctx context.Context, deleteUserID int) error
	}
)
