package structs

import "time"

type User struct {
	Name       string `json:"name"`
	LastJoin   time.Time `json:"last_join"`
	LastPart   time.Time `json:"last_part"`
	LastActive time.Time`json:"last_active"`
	FirstSeen  time.Time`json:"first_seen"`
}
