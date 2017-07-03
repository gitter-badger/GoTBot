package structs

import "time"

type User struct {
	Name       string
	LastJoin   time.Time
	LastPart   time.Time
	LastActive time.Time
	FirstSeen  time.Time
}
