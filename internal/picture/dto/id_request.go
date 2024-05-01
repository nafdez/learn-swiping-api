package picture

type IDRequest struct {
	PicID string `json:"pic_id" binding:"required"`
}
