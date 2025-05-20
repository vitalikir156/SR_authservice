package types

import "time"

type Token struct {
	ID          int
	UID         int
	Token		string
	Created     time.Time
}

type User struct {
	ID        int
	Name      string
	Taskread  *bool
	Taskwrite *bool
	Userread  *bool
	Userwrite *bool
	Password string
}
