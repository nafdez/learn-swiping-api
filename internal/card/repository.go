package card

import (
	"database/sql"
	"fmt"
	"learn-swiping-api/erro"
	"log"
	"strings"

	"github.com/go-sql-driver/mysql"
)

type CardRepository interface {
	Create(Card) (int64, error)
	ById(cardID int64, deckID int64) (Card, error)
	ByDeckId(id int64) ([]Card, error)
	ByProgress(token string, deckID int64) ([]Card, error)
	Update(card Card) error
	Delete(cardID int64, deckID int64) error
	// CreateWrong(wrong WrongAnswer) (int64, error)
	WrongByCardId(cardID int64) ([]WrongAnswer, error)
	UpdateWrong(id int64, wrong WrongAnswer) error
	// DeleteWrong(id int64) error
}

type CardRepositoryImpl struct {
	db              *sql.DB
	ByIdStmt        *sql.Stmt
	ByDeckIdStmt    *sql.Stmt
	ByProgressStmt  *sql.Stmt
	DeleteStmt      *sql.Stmt
	CreateWrongStmt *sql.Stmt
	WrongByIdStmt   *sql.Stmt
}

func NewCardRepository(db *sql.DB) *CardRepositoryImpl {
	repo := &CardRepositoryImpl{db: db}
	err := repo.InitStatements()
	if err != nil {
		log.Fatalln(err)
	}
	return repo
}

func (repo *CardRepositoryImpl) InitStatements() error {
	var err error
	// TODO: Implement token check
	repo.ByIdStmt, err = repo.db.Prepare("SELECT * FROM CARD WHERE card_id = ? AND deck_id = ?")
	if err != nil {
		return err
	}

	repo.ByDeckIdStmt, err = repo.db.Prepare("SELECT * FROM CARD WHERE deck_id = ?")
	if err != nil {
		return err
	}

	repo.ByProgressStmt, err = repo.db.Prepare(`SELECT c.*
												FROM CARD c
												LEFT JOIN PROGRESS p ON c.card_id = p.card_id
												LEFT JOIN ACCOUNT a ON p.acc_id = a.acc_id
												WHERE ((p.days_hidden <= 0 AND p.is_buried = false)
   												OR p.card_id IS NULL) 
												AND c.deck_id = ? 
												AND a.token = ?`)
	if err != nil {
		return err
	}

	//
	repo.DeleteStmt, err = repo.db.Prepare("DELETE FROM CARD WHERE card_id = ? AND deck_id = ?")
	if err != nil {
		return err
	}

	repo.CreateWrongStmt, err = repo.db.Prepare("INSERT INTO WRONG_ANSWER (card_id, answer) VALUES (?, ?)")
	if err != nil {
		return err
	}

	repo.WrongByIdStmt, err = repo.db.Prepare("SELECT wrong_id, answer FROM WRONG_ANSWER WHERE card_id = ?")
	if err != nil {
		return err
	}

	return nil
}

func (r *CardRepositoryImpl) Create(card Card) (int64, error) {
	tx, err := r.db.Begin()
	if err != nil {
		return 0, err
	}

	// Cannot use globally prepared statements here because of the transaction
	cardStmt, err := tx.Prepare("INSERT INTO CARD (deck_id, title, front, back, question, answer) VALUES (?, ?, ?, ?, ?, ?)")
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
	result, err := cardStmt.Exec(card.DeckID, card.Title, card.Front, card.Back, card.Question, card.Answer)
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

func (r *CardRepositoryImpl) ById(cardID int64, deckID int64) (Card, error) {
	var card Card
	row := r.ByIdStmt.QueryRow(cardID, deckID)

	err := row.Scan(
		&card.CardID,
		&card.DeckID,
		&card.Title,
		&card.Front,
		&card.Back,
		&card.Question,
		&card.Answer,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return Card{}, erro.ErrCardNotFound
		}
		return Card{}, err
	}

	return card, nil
}

func (r *CardRepositoryImpl) ByDeckId(id int64) ([]Card, error) {
	var cards []Card
	rows, err := r.ByDeckIdStmt.Query(id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var card Card
	for rows.Next() {
		err := rows.Scan(
			&card.CardID,
			&card.DeckID,
			&card.Title,
			&card.Front,
			&card.Back,
			&card.Question,
			&card.Answer,
		)
		if err != nil {
			return cards, err
		}
		cards = append(cards, card)
	}

	if len(cards) == 0 {
		return nil, erro.ErrCardNotFound
	}

	return cards, nil
}

func (r *CardRepositoryImpl) ByProgress(token string, deckID int64) ([]Card, error) {
	var cards []Card
	rows, err := r.ByProgressStmt.Query(deckID, token)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var card Card
	for rows.Next() {
		err := rows.Scan(
			&card.CardID,
			&card.DeckID,
			&card.Title,
			&card.Front,
			&card.Back,
			&card.Question,
			&card.Answer,
		)
		if err != nil {
			return cards, err
		}
		cards = append(cards, card)
	}

	if len(cards) == 0 {
		return nil, erro.ErrCardNotFound
	}

	return cards, nil
}

func (r *CardRepositoryImpl) Update(card Card) error {
	var query strings.Builder
	var args []any
	query.WriteString("UPDATE CARD SET")

	updateCardField(&query, &args, "title", card.Title)
	updateCardField(&query, &args, "front", card.Front)
	updateCardField(&query, &args, "back", card.Back)
	updateCardField(&query, &args, "question", card.Question)
	updateCardField(&query, &args, "answer", card.Answer)

	args = append(args, card.CardID, card.DeckID)
	query.WriteString(" WHERE card_id = ? AND deck_id = ?")

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

func (r *CardRepositoryImpl) Delete(cardID int64, deckID int64) error {
	// No need to delete wrong answers too because ON DELETE CASCADE will delete them
	result, err := r.DeleteStmt.Exec(cardID, deckID)
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

func (r *CardRepositoryImpl) WrongByCardId(cardID int64) ([]WrongAnswer, error) {
	var wrong []WrongAnswer
	rows, err := r.WrongByIdStmt.Query(cardID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var wrongAnswer WrongAnswer
	for rows.Next() {
		err := rows.Scan(
			&wrongAnswer.WrongID,
			&wrongAnswer.Answer,
		)
		if err != nil {
			return wrong, err
		}
		wrong = append(wrong, wrongAnswer)
	}

	if len(wrong) == 0 {
		return nil, erro.ErrWrongNotFound
	}

	return wrong, nil
}

func (r *CardRepositoryImpl) UpdateWrong(id int64, wrong WrongAnswer) error {
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

func updateCardField(query *strings.Builder, args *[]any, field string, value any) {
	if value == "" {
		return
	}

	if query.String() != "UPDATE CARD SET" && query.String() != "UPDATE WRONG_ANSWER SET" {
		query.WriteString(",")
	}

	query.WriteString(fmt.Sprintf(" %s = ?", field))
	*args = append(*args, value)
}
