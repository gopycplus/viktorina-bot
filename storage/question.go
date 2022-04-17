package storage

import (
	model "github.com/shavkatjon/viktorina-bot/model"
	"github.com/shavkatjon/viktorina-bot/utils"
)

//game
func GameGetQuestion(subject string) model.GameQuestion {
	query := `
		SELECT 
			id,
			text,
			image,
			answer,
			status
		FROM question WHERE status = 1 and subject = $1 ORDER by RANDOM() LIMIT 1;
	`

	row := questionDb.QueryRow(query, subject)

	var question model.GameQuestion

	err := row.Scan(
		&question.Id,
		&question.Text,
		&question.Image,
		&question.Answer,
		&question.Status,
	)

	utils.Check(err)

	return question
}

func GameGetNumberOfQuestions(subject string) int64 {
	query := `
		SELECT
		count(*) 
		FROM question WHERE status = 1 and subject = $	`

	row := questionDb.QueryRow(query, subject)

	var count int64
	err := row.Scan(&count)

	if err != nil {
		return 0
	}

	return count
}

func GameGetVariant(answer string, subject string) (list []string) {
	rows, err := questionDb.Query("select answer from question where answer != $1 and subject=$2 order by random() limit 3", answer, subject)
	utils.Check(err)
	var word string
	for rows.Next() {
		err := rows.Scan(&word)
		utils.Check(err)
		list = append(list, word)
	}
	return list
}

//history
func InsertQuestion(question model.Question) int64 {
	query := `
		INSERT INTO question 
		(
			subject,
			text,
			image,
			answer,
			user_id,
			status
		)
		VALUES ($1, $2, $3, $4, $5, $6)
	`

	_, err := questionDb.Exec(
		query,
		question.Subject,
		question.Text,
		question.Image,
		question.Answer,
		question.UserId,
		question.Status,
	)
	utils.Check(err)

	var id int64

	query = `
		SELECT 
			id 
		FROM question 
		WHERE status = 0
		AND user_id = $1
	`
	row := questionDb.QueryRow(query, question.UserId)
	err = row.Scan(&id)
	utils.Check(err)

	return id
}

func UpdateQuestion(question model.Question) {
	query := `
		UPDATE question SET
			text = $1,
			image = $2,
			answer = $3,
			user_id = $4,
			status = $5,
			subject = $6
		WHERE id = $7
	`

	_, err := questionDb.Exec(
		query,
		question.Text,
		question.Image,
		question.Answer,
		question.UserId,
		question.Status,
		question.Subject,
		question.Id,
	)

	utils.Check(err)
}

func DeleteQuestion(id int64, userId int64) bool {
	query := `
		DELETE FROM question
		WHERE id = $1 and user_id = $2
	`

	_, err := questionDb.Exec(
		query,
		id,
		userId,
	)

	return err == nil
}

func GetQuestion(id int64) model.Question {
	query := `
		SELECT 
			id,
			subject,
			text,
			image,
			answer,
			user_id, 
			status
		FROM question WHERE id = $1
	`

	row := questionDb.QueryRow(query, id)

	var question model.Question

	err := row.Scan(
		&question.Id,
		&question.Subject,
		&question.Text,
		&question.Image,
		&question.Answer,
		&question.UserId,
		&question.Status,
	)

	utils.Check(err)

	return question
}

func IsQuestionExists(id int64) bool {
	query := `
		SELECT 
			count(*)
		FROM question WHERE id = $1
	`

	row := questionDb.QueryRow(query, id)

	var count int
	err := row.Scan(&count)

	if err != nil || count == 0 {
		return false
	}
	return true
}

func IsExists(id int64, user_id int64) bool {
	query := `
		SELECT 
			count(*)
		FROM question WHERE id = $1 and user_id = $2
	`

	row := questionDb.QueryRow(query, id, user_id)

	var count int
	err := row.Scan(&count)

	if err != nil || count == 0 {
		return false
	}
	return true
}

func GetQuestionList(Limit int64, Page int64) model.QuestionList {
	var qList model.QuestionList
	var question model.Question
	var count int64

	query := `
		SELECT 
			id,
			subject,
			text,
			answer,
			user_id, 
			status
		FROM question
		WHERE status = 1
		ORDER BY id
		LIMIT $1, $2
	`
	rows, err := questionDb.Query(query, (Page-1)*10, Limit)
	utils.Check(err)

	for rows.Next() {
		err := rows.Scan(&question.Id, &question.Subject, &question.Text, &question.Answer, &question.UserId, &question.Status)
		utils.Check(err)
		qList.List = append(qList.List, question)
	}

	query = `SELECT count(*) FROM question WHERE status = 1`
	err = questionDb.QueryRow(query).Scan(&count)
	utils.Check(err)

	qList.Limit = Limit
	qList.Page = Page
	qList.Count = count
	qList.PageCount = count / Limit
	if count%Limit > 0 {
		qList.PageCount++
	}

	return qList
}
