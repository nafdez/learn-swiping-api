package account

import "time"

type Public struct {
	ID       int64     `json:"acc_id"`
	Username string    `json:"username"`
	PicID    string    `json:"pic_id"`
	LastSeen time.Time `json:"last_seen"`
	Since    time.Time `json:"since"`
}
