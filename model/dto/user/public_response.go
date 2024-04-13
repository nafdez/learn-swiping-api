package user

import "time"

type Public struct {
	ID       int64     `json:"acc_id"`
	Username string    `json:"username"`
	LastSeen time.Time `json:"last_seen"`
	Since    time.Time `json:"since"`
}
