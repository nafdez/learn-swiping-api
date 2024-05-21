package deck

type CreateRequest struct {
	Token       string
	Owner       int64  // Obtained with token
	Title       string `json:"title" binding:"required"`
	Description string `json:"description"`
	PicID       string `json:"pic_id`
	Visible     bool   `json:"visible"` // Default hidden
}
