package user

type Public struct {
	ID       string `json:"acc_id"`
	Username string `json:"username"`
	LastSeen string `json:"last_seen"`
	Since    string `json:"since"`
}
