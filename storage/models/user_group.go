package models

// UserGroup has many-many binding for users and groups
type UserGroup struct {
	UserNid  int
	GroupNid int
}
