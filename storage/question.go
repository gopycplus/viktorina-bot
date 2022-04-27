package storage

import (
	"github.com/shavkatjon/viktorina-bot/model"
	"github.com/shavkatjon/viktorina-bot/utils"
)

// game.db

func GameGetQuestion(subject string) model.GameQuestion {
	query := `
		SELECT 
			"id",
			"text",
			"image",
			"answer",
			"status",
			"subject",
			"variant1",
			"variant2",
			"variant3"
		FROM question 
		WHERE "status" = 1 
		  AND "subject" = $1 
		ORDER BY RANDOM() LIMIT 1;`

	row := questionDb.QueryRow(query, subject)

	var question model.GameQuestion

	err := row.Scan(
		&question.Id,
		&question.Text,
		&question.Image,
		&question.Answer,
		&question.Status,
		&question.Subject,
		&question.Variant1,
		&question.Variant2,
		&question.Variant3,
	)

	utils.Check(err)

	return question
}

func GameGetNumberOfQuestions(subject string) int64 {
	query := `
		SELECT
		count(*) 
		FROM question 
		WHERE "status" = 1 
		  AND "subject" = $1`

	row := questionDb.QueryRow(query, subject)

	var count int64
	err := row.Scan(&count)

	if err != nil {
		return 0
	}

	return count
}

// qa.db

func InsertQuestion(question model.Question) int64 {
	query := `
		INSERT INTO question 
		(
			"subject",
			"text",
			"image",
			"answer",
			"user_id",
			"status", 
		 	"variant1",
		 	"variant2",
		 	"variant3"
		)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)`

	_, err := questionDb.Exec(
		query,
		question.Subject,
		question.Text,
		question.Image,
		question.Answer,
		question.UserId,
		question.Status,
		question.Variant1,
		question.Variant2,
		question.Variant3,
	)
	utils.Check(err)

	var id int64

	query = `
		SELECT 
			"id" 
		FROM question 
		WHERE "status" = 0
		AND "user_id" = $1`

	row := questionDb.QueryRow(query, question.UserId)
	err = row.Scan(&id)
	utils.Check(err)

	return id
}

func UpdateQuestion(question model.Question) {
	query := `
		UPDATE question SET
			"text" = $1,
			"image" = $2,
			"answer" = $3,
			"user_id" = $4,
			"status" = $5,
			"subject" = $6,
			"variant1" = $7,
			"variant2" = $8,
			"variant3" = $9
		WHERE "id" = $10`

	_, err := questionDb.Exec(
		query,
		question.Text,
		question.Image,
		question.Answer,
		question.UserId,
		question.Status,
		question.Subject,
		question.Variant1,
		question.Variant2,
		question.Variant3,
		question.Id,
	)

	utils.Check(err)
}

func DeleteQuestion(id int64, userId int64) bool {
	query := `
		DELETE 
		FROM question
		WHERE "id" = $1 
		  AND "user_id" = $2`

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
			"id",
			"subject",
			"text",
			"image",
			"answer",
			"user_id", 
			"status",
			"variant1",
			"variant2",
			"variant3"
		FROM question WHERE "id" = $1
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
		&question.Variant1,
		&question.Variant2,
		&question.Variant3,
	)

	utils.Check(err)

	return question
}

func IsQuestionExists(id int64) bool {
	query := `
		SELECT 
			count(*)
		FROM question 
		WHERE "id" = $1`

	row := questionDb.QueryRow(query, id)

	var count int
	err := row.Scan(&count)

	if err != nil || count == 0 {
		return false
	}
	return true
}

func IsExists(id int64, userId int64) bool {
	query := `
		SELECT
			count(*)
		FROM question 
		WHERE "id" = $1 
		  AND "user_id" = $2`

	row := questionDb.QueryRow(query, id, userId)

	var count int
	err := row.Scan(&count)

	if err != nil || count == 0 {
		return false
	}
	return true
}

func GetQuestionList(subject string, limit int64, page int64) model.QuestionList {
	var (
		qList    model.QuestionList
		question model.Question
		count    int64
	)

	query := `
		SELECT 
			"id",
			"subject",
			"text",
			"answer"
		FROM question
		WHERE "status" = 1 AND "subject" = $1
		ORDER BY "id"
		LIMIT $2, $3`

	rows, err := questionDb.Query(query, subject, (page-1)*10, limit)
	utils.Check(err)

	for rows.Next() {
		err := rows.Scan(&question.Id, &question.Subject, &question.Text, &question.Answer)
		utils.Check(err)
		qList.List = append(qList.List, question)
	}

	query = `
		SELECT 
		    count(*) 
		FROM question 
		WHERE "status" = 1`

	err = questionDb.QueryRow(query).Scan(&count)
	utils.Check(err)

	qList.Limit = limit
	qList.Page = page
	qList.Count = count
	qList.PageCount = count / limit
	if count%limit > 0 {
		qList.PageCount++
	}

	return qList
}
