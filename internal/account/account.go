package account

import "time"

type Account struct {
	ID           int64     `json:"acc_id"`
	Username     string    `json:"username"`
	Password     string    `json:"-"`
	Email        string    `json:"email,omitempty"`
	Name         string    `json:"name"`
	Token        string    `json:"token"`
	TokenExpires time.Time `json:"token_expires"`
	LastSeen     time.Time `json:"last_seen"`
	Since        time.Time `json:"since"`
}
