package deck

import "mime/multipart"

type UpdateRequest struct {
	DeckID      int64  `json:"deck_id" binding:"required"`
	Title       string `json:"title"`
	Description string `json:"description"`
	Img         *multipart.FileHeader
	Visible     *bool `json:"visible"` // pointer to check if empty or not
}
