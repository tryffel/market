package models

import (
	"strings"
	"time"
)

type ActionType int

const (
	ActionAdd    ActionType = 1
	ACtionDel    ActionType = 2
	ActionUpdate ActionType = 3
)

type Action struct {
	Id           int
	UserNid      int
	Document     int
	Operation    int
	Comment      string
	CommentLower string
	Ts           time.Time
}

func (a *Action) BeforeCreate() error {
	a.Ts = time.Now()
	a.CommentLower = strings.ToLower(a.Comment)
	return nil
}
