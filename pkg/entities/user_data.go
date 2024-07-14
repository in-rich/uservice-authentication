package entities

import (
	"github.com/google/uuid"
	"github.com/uptrace/bun"
)

type User struct {
	bun.BaseModel `bun:"table:users"`

	ID *uuid.UUID `bun:"id,pk,type:uuid"`

	PublicIdentifier string `bun:"public_identifier,unique,notnull"`
	FirebaseUID      string `bun:"firebase_uid,unique,notnull"`
}
