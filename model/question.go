package model

type GameQuestion struct {
	Id       int64
	Subject  string
	Text     string
	Image    string
	Answer   string
	Status   int
	Variant1 string
	Variant2 string
	Variant3 string
}

type Question struct {
	Id       int64
	Subject  string
	Text     string
	Image    string
	Answer   string
	UserId   int64
	Status   int
	Variant1 string
	Variant2 string
	Variant3 string
}

type QuestionList struct {
	List      []Question
	Limit     int64
	Page      int64
	PageCount int64
	Count     int64
}
