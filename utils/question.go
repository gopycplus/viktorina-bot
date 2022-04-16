package utils

import (
	"strconv"

	"github.com/shavkatjon/viktorina-bot/model"
)

func ListToText(list model.QuestionList) string {
	var text string

	for i, question := range list.List {
		text += "<b>" + strconv.Itoa(int(list.Page-1)*10+i+1) + "</b>." + " <b>Indeks:</b> " + strconv.Itoa(int(question.Id)) + "\n"
		text += " <b>Fan:</b> " + question.Subject + "\n"
		text += " <b>Savol:</b> " + question.Text + "\n"
		text += " <b>Javob:</b> " + question.Answer + "\n"
	}

	text += "\n<i>Umumiy savollar soni:</i> " + strconv.Itoa(int(list.Count))

	return text
}
