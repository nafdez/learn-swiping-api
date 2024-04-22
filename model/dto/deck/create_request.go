package deck

type CreateRequest struct {
	Owner       int64  `json:"owner" binding:"required"`
	Title       string `json:"title" binding:"required"`
	Description string `json:"description"`
	Visible     bool   `json:"visible"` // Default hidden
}
