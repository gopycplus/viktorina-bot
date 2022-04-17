package model

type GameQuestion struct {
	Id      int64
	Subject string
	Text    string
	Image   string
	Answer  string
	Status  int
}

type Question struct {
	Id      int64
	Subject string
	Text    string
	Image   string
	Answer  string
	UserId  int64
	Status  int
}

type QuestionList struct {
	List      []Question
	Limit     int64
	Page      int64
	PageCount int64
	Count     int64
}
