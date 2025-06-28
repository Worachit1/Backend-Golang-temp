package model

import (
	"github.com/uptrace/bun"
)

type Student struct {
	bun.BaseModel `bun:"table:students"`

	ID             string `json:"id" bun:",pk,type:uuid,default:gen_random_uuid()"`
	First_name     string `bun:"first_name,notnull"`
	Last_name      string `bun:"last_name,notnull"`
	Student_number string `bun:"student_number,notnull"`

	Registrations []*Registration `bun:"rel:has-many,join:id=student_id"`

	CreateUpdateUnixTimestamp
	SoftDelete
}
