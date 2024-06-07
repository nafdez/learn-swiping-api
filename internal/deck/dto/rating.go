package deck

type Rating struct {
	Rating      int8    `json:"rating,omitempty"`
	RatingCount int64   `json:"rating_count,omitempty"`
	AvgRating   float32 `json:"avg_rating,omitempty"`
}
