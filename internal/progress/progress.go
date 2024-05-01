package progress

type Progress struct {
	ProgressID     int64 `json:"progress_id"`
	AccID          int64 `json:"acc_id"`
	CardID         int64 `json:"card_id"`
	Priority       int   `json:"priority"`
	DaysHidden     int   `json:"days_hidden"`
	PriorityExam   int   `json:"priority_exam"`
	DaysHiddenExam int   `json:"days_hidden_exam"`
	AnswerCount    int   `json:"answer_count"`
	CorrectCount   int   `json:"correct_count"`
	IsBuried       bool  `json:"is_buried"`
}
