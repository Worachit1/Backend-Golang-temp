package model

import (
	"github.com/uptrace/bun"
	"golang.org/x/crypto/bcrypt"
)

type Student struct {
	bun.BaseModel `bun:"table:students"`

	ID            string `json:"id" bun:",pk,type:uuid,default:gen_random_uuid()"`
	FirstName     string `bun:"first_name,notnull"`
	LastName      string `bun:"last_name,notnull"`
	StudentNumber string `bun:"student_number,notnull, unique"`
	Password	  string `bun:"password,notnull"`

	Registrations []*Registration `bun:"rel:has-many,join:id=student_id"`

	CreateUpdateUnixTimestamp
	SoftDelete
}

func (u *Student) CheckPassword(password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(password))
	return err == nil
}