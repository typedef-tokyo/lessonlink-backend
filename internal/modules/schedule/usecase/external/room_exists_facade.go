package external

import (
	"context"
)

type (
	IRoomExistsFacade interface {
		Execute(ctx context.Context, inputCampus string, inputRoomIndexes []int) error
	}
)
