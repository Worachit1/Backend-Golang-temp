package model

import (
	"github.com/uptrace/bun"
)

type Officer struct {
	bun.BaseModel `bun:"table:officers"`

	ID        string `bun:",pk,type:uuid,default:gen_random_uuid()"`
	FirstName string `bun:"first_name,notnull"`
	LastName  string `bun:"last_name,notnull"`
	Phone     string `bun:"phone,notnull"`
	Email     string `bun:"email,notnull,unique"`
	Password      string `bun:"password,notnull"`

	CreateUpdateUnixTimestamp
	SoftDelete
}