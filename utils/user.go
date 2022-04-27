package utils

import (
	"strconv"

	"github.com/shavkatjon/viktorina-bot/model"
)

func GameListToText(uList []model.GameUser) string {
	var text string
	text += "Top 10 talik o'yinchilar: \n"
	for i, user := range uList {
		text += "<b>" + strconv.Itoa(i+1) + " - o'rin:</b>" + "\n"
		text += " <b>O'yinchi: </b>" + user.FirstName + " " + user.LastName + "\n"
		text += " <b>O'yinchining bali: </b> " + strconv.Itoa(user.Score) + "\n"
	}

	return text
}
