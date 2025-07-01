package model

import (
	"github.com/uptrace/bun"
)

type Activity struct {
	bun.BaseModel `bun:"table:activities"`

	ID          string `json:"id" bun:",pk,type:uuid,default:gen_random_uuid()"`
	Name        string `bun:"name,notnull"`
	Description string `bun:"description,notnull"`
	ReleaseDate int64  `bun:"release_date,notnull"`

	Registrations []*Registration `bun:"rel:has-many,join:id=activity_id"`

	CreateUpdateUnixTimestamp
	SoftDelete
}
