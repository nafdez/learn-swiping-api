package repository

import (
	"database/sql"
	"fmt"
	"learn-swiping-api/erro"
	"learn-swiping-api/model"
	"log"
	"strings"

	"github.com/go-sql-driver/mysql"
)

type CardRepository interface {
	Create(model.Card) (int64, error)
	ById(id int64) (model.Card, error)
	ByDeckId(id int64) ([]model.Card, error)
	Update(id int64, card model.Card) error
	Delete(id int64) error
	CreateWrong(wrong model.WrongAnswer) (int64, error)
	WrongByCardId(cardID int64) ([]model.WrongAnswer, error)
	UpdateWrong(id int64, wrong model.WrongAnswer) error
	DeleteWrong(id int64) error
}

type CardRepositoryImpl struct {
	db              *sql.DB
	ByIdStmt        *sql.Stmt
	ByDeckIdStmt    *sql.Stmt
	DeleteStmt      *sql.Stmt
	CreateWrongStmt *sql.Stmt
	WrongByIdStmt   *sql.Stmt
	DeleteWrongStmt *sql.Stmt
}

func NewCardRepository(db *sql.DB) *CardRepositoryImpl {
	repo := &CardRepositoryImpl{db: db}
	err := repo.InitStatements()
	if err != nil {
		log.Println(err)
	}
	return repo
}

func (repo *CardRepositoryImpl) InitStatements() error {
	var err error
	repo.ByIdStmt, err = repo.db.Prepare("SELECT * FROM CARD WHERE card_id = ?")
	if err != nil {
		return err
	}

	repo.ByDeckIdStmt, err = repo.db.Prepare("SELECT * FROM CARD WHERE deck_id = ?")
	if err != nil {
		return err
	}

	repo.DeleteStmt, err = repo.db.Prepare("DELETE FROM CARD WHERE card_id = ?")
	if err != nil {
		return err
	}

	repo.CreateWrongStmt, err = repo.db.Prepare("INSERT INTO WRONG_ANSWER (card_id, answer) VALUES (?, ?)")
	if err != nil {
		return err
	}

	repo.WrongByIdStmt, err = repo.db.Prepare("SELECT * FROM WRONG_ANSWER WHERE card_id = ?")
	if err != nil {
		return err
	}

	repo.DeleteWrongStmt, err = repo.db.Prepare("DELETE FROM WRONG_ANSWER WHERE wrong_id = ?")
	if err != nil {
		return err
	}

	return nil
}

func (r *CardRepositoryImpl) Create(card model.Card) (int64, error) {
	tx, err := r.db.Begin()
	if err != nil {
		return 0, err
	}

	// Cannot use globally prepared statements here because of the transaction
	cardStmt, err := tx.Prepare("INSERT INTO CARD (deck_id, to_study, question, answer) VALUES (?,?,?,?)")
	if err != nil {
		return 0, err
	}

	wrongStmt, err := tx.Prepare("INSERT INTO WRONG_ANSWER (card_id, answer) VALUES (?,?)")
	if err != nil {
		return 0, err
	}

	defer cardStmt.Close()
	defer wrongStmt.Close()

	// Insert Card
	result, err := cardStmt.Exec(card.DeckID, card.Study, card.Question, card.Answer)
	if err != nil {
		tx.Rollback()
		if err.(*mysql.MySQLError).Number == 1452 {
			return 0, erro.ErrDeckNotFound
		}
		return 0, err
	}

	id, err := result.LastInsertId()
	if err != nil {
		tx.Rollback()
		return 0, err
	}

	// Insert wrong answers
	for i := 0; i < len(card.Wrong); i++ {
		_, err := wrongStmt.Exec(id, card.Wrong[i].Answer)
		if err != nil {
			tx.Rollback()
			if err.(*mysql.MySQLError).Number == 1452 {
				return 0, erro.ErrCardNotFound
			}
			return 0, err
		}
	}

	err = tx.Commit()
	if err != nil {
		tx.Rollback()
		return 0, err
	}

	return id, nil
}

func (r *CardRepositoryImpl) ById(id int64) (model.Card, error) {
	var card model.Card
	row := r.ByIdStmt.QueryRow(id)

	err := row.Scan(
		&card.ID,
		&card.DeckID,
		&card.Study,
		&card.Question,
		&card.Answer,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return model.Card{}, erro.ErrCardNotFound
		}
		return model.Card{}, err
	}

	return card, nil
}

func (r *CardRepositoryImpl) ByDeckId(id int64) ([]model.Card, error) {
	var cards []model.Card
	rows, err := r.ByDeckIdStmt.Query(id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var card model.Card
	for rows.Next() {
		err := rows.Scan(
			&card.ID,
			&card.DeckID,
			&card.Study,
			&card.Question,
			&card.Answer,
		)
		if err != nil {
			if err == sql.ErrNoRows {
				return nil, erro.ErrCardNotFound
			}
			return cards, err
		}
		cards = append(cards, card)
	}
	return cards, nil
}

func (r *CardRepositoryImpl) Update(id int64, card model.Card) error {
	var query strings.Builder
	var args []any
	query.WriteString("UPDATE CARD SET")

	updateCardField(&query, &args, "to_study", card.Study)
	updateCardField(&query, &args, "question", card.Question)
	updateCardField(&query, &args, "answer", card.Answer)

	args = append(args, id)
	query.WriteString(" WHERE card_id = ?")

	stmt, err := r.db.Prepare(query.String())
	if err != nil {
		return err
	}

	result, err := stmt.Exec(args...)
	if err != nil {
		return err
	}

	affected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if affected == 0 {
		return erro.ErrCardNotFound
	}

	return nil
}

func (r *CardRepositoryImpl) Delete(id int64) error {
	// No need to delete wrong answers too because ON DELETE CASCADE will delete them
	result, err := r.DeleteStmt.Exec(id)
	if err != nil {
		return err
	}

	affected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if affected == 0 {
		return erro.ErrCardNotFound
	}

	return nil
}

func (r *CardRepositoryImpl) CreateWrong(wrong model.WrongAnswer) (int64, error) {
	result, err := r.CreateWrongStmt.Exec(wrong.CardID, wrong.Answer)
	if err != nil {
		return 0, err
	}
	return result.LastInsertId()
}

func (r *CardRepositoryImpl) WrongByCardId(cardID int64) ([]model.WrongAnswer, error) {
	var wrong []model.WrongAnswer
	rows, err := r.WrongByIdStmt.Query(cardID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var wrongAnswer model.WrongAnswer
	for rows.Next() {
		err := rows.Scan(
			&wrongAnswer.CardID,
			&wrongAnswer.Answer,
		)
		if err != nil {
			if err == sql.ErrNoRows {
				return nil, erro.ErrWrongNotFound
			}
			return wrong, err
		}
		wrong = append(wrong, wrongAnswer)
	}
	return wrong, nil
}

func (r *CardRepositoryImpl) UpdateWrong(id int64, wrong model.WrongAnswer) error {
	var query strings.Builder
	var args []any
	query.WriteString("UPDATE WRONG_ANSWER SET")

	updateCardField(&query, &args, "answer", wrong.Answer)

	args = append(args, id)
	query.WriteString(" WHERE wrong_id = ?")

	stmt, err := r.db.Prepare(query.String())
	if err != nil {
		return err
	}

	result, err := stmt.Exec(args...)
	if err != nil {
		return err
	}

	affected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if affected == 0 {
		return erro.ErrCardNotFound
	}

	return nil
}

func (r *CardRepositoryImpl) DeleteWrong(id int64) error {
	result, err := r.DeleteWrongStmt.Exec(id)
	if err != nil {
		return err
	}

	affected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if affected == 0 {
		return erro.ErrCardNotFound
	}

	return nil
}

func updateCardField(query *strings.Builder, args *[]any, field string, value any) {
	if value == "" {
		return
	}

	if query.String() != "UPDATE DECK SET" {
		query.WriteString(",")
	}

	query.WriteString(fmt.Sprintf(" %s = ?", field))
	*args = append(*args, value)
}
