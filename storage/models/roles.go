package models

type Role struct {
	Id          int
	UserNid     int
	GroupNid    int
	RoleLevel   int
	MetadataKey int
}
