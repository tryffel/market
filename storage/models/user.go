package models

import (
	"github.com/tryffel/market/modules/util"
	"strings"
	"time"
)

type User struct {
	Id        string
	Nid       int
	Name      string
	LowerName string
	Email     string
	Password  string
	LastSeen  time.Time
	IsActive  bool
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt time.Time
}

func (u *User) BeforeCreate() error {
	u.Id = util.NewUuid()
	u.LowerName = strings.ToLower(u.Name)
	u.CreatedAt = time.Now()
	return nil
}

func (u *User) BeforeUpdate() error {
	u.UpdatedAt = time.Now()
	return nil
}
