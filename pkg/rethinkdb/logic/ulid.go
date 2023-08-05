package logic

import (
	"math/rand"
	"time"

	"github.com/oklog/ulid/v2"
)

func NewULID() ulid.ULID {

	entropy := rand.New(rand.NewSource(time.Now().UnixNano()))
	ms := ulid.Timestamp(time.Now())
	created := ulid.MustNew(ms, entropy)

	return created
}
