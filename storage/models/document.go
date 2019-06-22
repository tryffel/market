package models

import (
	"github.com/tryffel/market/modules/util"
	"strings"
	"time"
)

type Document struct {
	Id          string
	Nid         int
	Name        string
	LowerName   string
	Description string
	Mimetype    string
	Size        int
	CreatedBy   int
	CreatedAt   time.Time
	UpdatedAt   time.Time
	DeletedAt   time.Time
}

func (d *Document) BeforeCreate() error {
	d.Id = util.NewUuid()
	d.LowerName = strings.ToLower(d.Name)
	d.CreatedAt = time.Now()
	return nil
}

func (d *Document) BeforeUpdate() error {
	d.LowerName = strings.ToLower(d.Name)
	d.UpdatedAt = time.Now()
	return nil
}
