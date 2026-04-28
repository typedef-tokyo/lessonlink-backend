package vo

import (
	"errors"
	"strings"

	"github.com/google/uuid"
	"github.com/typedef-tokyo/lessonlink-backend/internal/pkg/log"
)

var ErrIdentifierEmpty = errors.New("識別子が未設定です")

type Identifier string

const (
	IDENTIFIER_INITIAL = Identifier("")
	IDENTIFIER_INVALID = Identifier("INVALID")
)

func NewIdentifier(id string) (Identifier, error) {

	key := strings.TrimSpace(id)
	if key == "" {
		return IDENTIFIER_INVALID, log.WrapErrorWithStackTraceBadRequest(ErrIdentifierEmpty)
	}

	return Identifier(key), nil
}

func NewIdentifierGenerate() Identifier {

	return Identifier(uuid.New().String())
}

func (r Identifier) Value() string {

	return string(r)
}
