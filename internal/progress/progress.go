package progress

type Progress struct {
	ProgressID     int64   `json:"progress_id"`
	AccID          int64   `json:"acc_id"`
	CardID         int64   `json:"card_id"`
	Ease           float32 `json:"ease"`
	Interval       int     `json:"interval"`
	Priority       int     `json:"priority"`
	DaysHidden     int     `json:"days_hidden"`
	WatchCount     int     `json:"watch_count"`
	PriorityExam   int     `json:"priority_exam"`
	DaysHiddenExam int     `json:"days_hidden_exam"`
	AnswerCount    int     `json:"answer_count"`
	CorrectCount   int     `json:"correct_count"`
	IsRelearning   bool    `json:"is_relearning"`
	IsBuried       bool    `json:"is_buried"`
}
