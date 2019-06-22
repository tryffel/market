package models

import (
	"github.com/tryffel/market/modules/util"
	"strings"
	"time"
)

type Group struct {
	Id        string
	Nid       int
	Name      string
	LowerName string
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt time.Time
}

func (g *Group) BeforeCreate() error {
	g.Id = util.NewUuid()
	g.LowerName = strings.ToLower(g.Name)
	g.CreatedAt = time.Now()
	return nil
}

func (g *Group) BeforeUpdate() error {
	g.UpdatedAt = time.Now()
	return nil
}
