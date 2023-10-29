package repo

type UserStatus int16

const (
	DRAFT UserStatus = iota + 1
	ACTIVE
)
