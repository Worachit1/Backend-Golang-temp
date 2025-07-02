package model

import (
	"github.com/uptrace/bun"
	"golang.org/x/crypto/bcrypt"
)

type User struct {
	bun.BaseModel `bun:"table:users"`

	ID            string `bun:",pk,type:uuid,default:gen_random_uuid()"`
	StudentNumber string `bun:"student_number,notnull,unique"`
	FirstName     string `bun:"first_name,notnull"`
	LastName      string `bun:"last_name,notnull"`
	Phone         string `bun:"phone,notnull"`
	Email         string `bun:"email,notnull,unique"`
	Address       string `bun:"address,notnull"`
	Role		  string `bun:"role,notnull,default:'user'"` // user, officer, admin
	Password      string `bun:"password,notnull"`

	CreateUpdateUnixTimestamp
	SoftDelete
}

func (u *User) CheckPassword(password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(password))
	return err == nil
}
