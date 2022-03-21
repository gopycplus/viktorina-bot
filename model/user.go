package model

type GameUser struct {
	Id        int64
	ChatId    int64
	FirstName string
	LastName  string
	UserName  string
	Step      int
	MessageId int64
	Score     int
	Subject   string
	Question  string
	Answer    string
	GroupName string
	GroupPass string
	IsInGroup int
}

type User struct {
	Id         int64
	ChatId     int64
	FirstName  string
	LastName   string
	UserName   string
	Step       int
	Subject    string
	MessageId  int64
	QuestionId int64
}
