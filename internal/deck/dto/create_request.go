package deck

type CreateRequest struct {
	Token       string
	Owner       int64  // Obtained with token
	Title       string `json:"title" binding:"required"`
	Description string `json:"description"`
	Visible     bool   `json:"visible"` // Default hidden
}
