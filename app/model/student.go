package model

import (
	"github.com/uptrace/bun"
)

type Student struct {
	bun.BaseModel `bun:"table:students"`

	ID            string `json:"id" bun:",pk,type:uuid,default:gen_random_uuid()"`
	FirstName     string `bun:"first_name,notnull"`
	LastName      string `bun:"last_name,notnull"`
	StudentNumber string `bun:"student_number,notnull"`

	Registrations []*Registration `bun:"rel:has-many,join:id=student_id"`

	CreateUpdateUnixTimestamp
	SoftDelete
}
