package models

import (
	"github.com/tryffel/market/modules/util"
	"time"
)

type Session struct {
	Id        int
	UserNid   int
	Nonce     string
	UserAgent string
	IssuedAt  time.Time
}

// BeforeCreate creates nonce and sets issued_at = now
func (s *Session) BeforeCreate() error {
	s.Nonce = util.RandomKey(20)
	s.IssuedAt = time.Now()
	return nil
}
