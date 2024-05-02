package progress

// Using pointers so we can know if they're assigned
type UpdateRequest struct {
	Token          string
	CardID         int64 `json:"card_id" binding:"required"`
	Priority       *int  `json:"priority"`
	DaysHidden     *int  `json:"days_hidden"`
	WatchCount     *int  `json:"watch_count"`
	PriorityExam   *int  `json:"priority_exam"`
	DaysHiddenExam *int  `json:"days_hidden_exam"`
	AnswerCount    *int  `json:"answer_count"`
	CorrectCount   *int  `json:"correct_count"`
	IsBuried       *bool `json:"is_buried"`
}

func (req *UpdateRequest) SetToken(token string) {
	req.Token = token
}
