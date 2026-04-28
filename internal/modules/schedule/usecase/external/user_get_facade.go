package external

import (
	"context"
)

type (
	IUserGetFacade interface {
		Execute(ctx context.Context, inputUserID int) (userName string, role string, err error)
	}
)
