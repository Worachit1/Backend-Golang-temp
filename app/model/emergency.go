package model

import (
	"github.com/uptrace/bun"
)

type Emergency struct {
	bun.BaseModel `bun:"table:emergencies"`

	ID          string `bun:",pk,type:uuid,default:gen_random_uuid()" json:"id"`
	UserID      string`bun:"user_id,type:uuid,notnull" json:"user_id"`
	OfficerID   string `bun:"officer_id,type:uuid,nullzero" json:"officer_id,omitempty"`
	Type        string    `json:"type"`
	Status      string    `json:"status"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	Location    string    `json:"location"`
	MapLink     string    `json:"map_link,omitempty"`
	ActionNote  string    `json:"action_note,omitempty"`
	CreateUpdateUnixTimestamp
	SoftDelete

	User    *User    `bun:"rel:belongs-to,join:user_id=id" json:"-"`
	Officer *Officer `bun:"rel:belongs-to,join:officer_id=id" json:"-"`
}

