package repository

import (
	"database/sql"
	"fmt"
	"learn-swiping-api/erro"
	"learn-swiping-api/model"
	"log"
	"strings"
	"time"

	"github.com/go-sql-driver/mysql"
)

type DeckRepository interface {
	Create(model.Deck) (int64, error)
	ById(id int64, showHidden bool) (model.Deck, error)
	ByOwner(id int64, showHidden bool) ([]model.Deck, error)
	ByUserId(id int64, showHidden bool) ([]model.Deck, error) // ACC-DECK table
	Update(id int64, deck model.Deck) error
	Delete(id int64) error
}

type DeckRepositoryImpl struct {
	db           *sql.DB
	CreateStmt   *sql.Stmt
	ByIdStmt     *sql.Stmt
	ByOwnerStmt  *sql.Stmt
	ByUserIdStmt *sql.Stmt
	DeleteStmt   *sql.Stmt
}

func NewDeckRepository(db *sql.DB) *DeckRepositoryImpl {
	repo := &DeckRepositoryImpl{db: db}
	err := repo.InitStatements()
	if err != nil {
		log.Println(err)
	}
	return repo
}

func (repo *DeckRepositoryImpl) InitStatements() error {
	var err error
	repo.CreateStmt, err = repo.db.Prepare("INSERT INTO DECK (owner, title, description) VALUES (?,?,?)")
	if err != nil {
		return err
	}

	// visible = 1 OR visible = ?
	// To fetch only visible "?" needs to be 1
	// To fetch also hidden "?" needs to be 0
	repo.ByIdStmt, err = repo.db.Prepare("SELECT * FROM DECK WHERE deck_id = ? AND (visible = 1 OR visible = ?)")
	if err != nil {
		return err
	}

	repo.ByOwnerStmt, err = repo.db.Prepare("SELECT * FROM DECK WHERE acc_id = ? AND (visible = 1 OR visible = ?)")
	if err != nil {
		return err
	}

	repo.ByUserIdStmt, err = repo.db.Prepare("SELECT d.* FROM DECK as d, ACC_DECK as ad WHERE ad.acc_id = ? AND d.deck_id = ad.deck_id AND visible = ? AND (visible = 1 OR visible = ?)")
	if err != nil {
		return err
	}

	repo.DeleteStmt, err = repo.db.Prepare("DELETE FROM DECK WHERE deck_id = ?")
	if err != nil {
		return err
	}

	return nil
}

func (r *DeckRepositoryImpl) Create(deck model.Deck) (int64, error) {
	result, err := r.CreateStmt.Exec(deck.Owner, deck.Title, deck.Description)
	if err != nil {
		if err.(*mysql.MySQLError).Number == 1452 {
			return 0, erro.ErrUserNotFound
		}
		return 0, err
	}
	return result.LastInsertId()
}

func (r *DeckRepositoryImpl) ById(id int64, showHidden bool) (model.Deck, error) {
	row := r.ByIdStmt.QueryRow(id, !showHidden)
	return scanDeck(row)
}

func (r *DeckRepositoryImpl) ByOwner(id int64, showHidden bool) ([]model.Deck, error) {
	rows, err := r.ByOwnerStmt.Query(id, !showHidden)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	//TODO:
	return scanDecks(rows)
}

func (r *DeckRepositoryImpl) ByUserId(id int64, showHidden bool) ([]model.Deck, error) {
	rows, err := r.ByUserIdStmt.Query(id, !showHidden)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	return scanDecks(rows)
}

func (r *DeckRepositoryImpl) Update(id int64, deck model.Deck) error {
	var query strings.Builder
	var args []any
	query.WriteString("UPDATE DECK SET")

	updateDeckField(&query, &args, "title", deck.Title)
	updateDeckField(&query, &args, "description", deck.Description)
	updateDeckField(&query, &args, "updated_at", deck.UpdatedAt)
	updateDeckField(&query, &args, "visible", deck.Visible)

	args = append(args, id)
	query.WriteString(" WHERE deck_id = ?")

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
		return erro.ErrDeckNotFound
	}

	return nil
}

func (r *DeckRepositoryImpl) Delete(id int64) error {
	result, err := r.DeleteStmt.Exec(id)
	if err != nil {
		return err
	}

	affected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if affected == 0 {
		return erro.ErrDeckNotFound
	}

	return nil
}

func updateDeckField(query *strings.Builder, args *[]any, field string, value any) {
	// Just checking if it's a date and it isn't empty
	if _, ok := value.(time.Time); ok && value.(time.Time).IsZero() {
		return
	}

	if value == "" {
		return
	}

	if query.String() != "UPDATE DECK SET" {
		query.WriteString(",")
	}

	query.WriteString(fmt.Sprintf(" %s = ?", field))
	*args = append(*args, value)
}

func scanDeck(row *sql.Row) (model.Deck, error) {
	var deck model.Deck
	err := row.Scan(
		&deck.ID,
		&deck.Owner,
		&deck.Title,
		&deck.Description,
		&deck.UpdatedAt,
		&deck.CreatedAt,
		&deck.Visible,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return model.Deck{}, erro.ErrDeckNotFound
		}
		return model.Deck{}, err
	}

	return deck, nil
}

func scanDecks(rows *sql.Rows) ([]model.Deck, error) {
	var decks []model.Deck
	var deck model.Deck
	for rows.Next() {
		err := rows.Scan(
			&deck.ID,
			&deck.Owner,
			&deck.Title,
			&deck.Description,
			&deck.UpdatedAt,
			&deck.CreatedAt,
			&deck.Visible,
		)
		if err != nil {
			if err == sql.ErrNoRows {
				return decks, erro.ErrDeckNotFound
			}
			return decks, err
		}
		decks = append(decks, deck)
	}
	return decks, nil
}
