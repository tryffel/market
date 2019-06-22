package models

import "time"

type RoleLevels struct {
	Level     int
	Name      string
	Comment   string
	Group     int
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt time.Time
}

func (r *RoleLevels) BeforeCreate() error {
	r.CreatedAt = time.Now()
	return nil
}

func (r *RoleLevels) BeforeUpdate() error {
	r.UpdatedAt = time.Now()
	return nil
}
