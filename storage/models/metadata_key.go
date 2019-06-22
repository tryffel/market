package models

import "time"

type MetadataKey struct {
	Key       string
	Multiple  bool
	Comment   string
	Active    bool
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt time.Time
}

func (m *MetadataKey) BeforeCreate() error {
	m.CreatedAt = time.Now()
	return nil
}

func (m *MetadataKey) BeforeUpdate() error {
	m.UpdatedAt = time.Now()
	return nil
}
