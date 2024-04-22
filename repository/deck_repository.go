package repository

import (
	"database/sql"
	"fmt"
	"learn-swiping-api/erro"
	"learn-swiping-api/model"
	"learn-swiping-api/model/dto/deck"
	"log"
	"strings"
	"time"

	"github.com/go-sql-driver/mysql"
)

type DeckRepository interface {
	Create(model.Deck) (int64, error)
	ById(deck.ReadRequest) (model.Deck, error)
	ByOwner(deck.ReadRequest) ([]model.Deck, error)
	ByUserId(deck.ReadRequest) ([]model.Deck, error) // ACC-DECK table
	Update(id int64, deck model.Deck) error
	Delete(id int64) error
	CheckOwnership(deckID int64, token string) bool
}

type DeckRepositoryImpl struct {
	db             *sql.DB
	CreateStmt     *sql.Stmt
	ByIdStmt       *sql.Stmt
	ByOwnerStmt    *sql.Stmt
	ByUserIdStmt   *sql.Stmt
	DeleteStmt     *sql.Stmt
	CheckOwnerStmt *sql.Stmt
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
	// repo.ByIdStmt, err = repo.db.Prepare("SELECT * FROM DECK WHERE deck_id = ? AND (visible = 1 OR visible = ?)")
	repo.ByIdStmt, err = repo.db.Prepare(`SELECT d.* FROM DECK d
											LEFT JOIN ACC_DECK ad ON d.deck_id = ad.deck_id
											LEFT JOIN ACCOUNT a ON d.acc_id = a.acc_id
											HERE ad.acc_id = ?
											AND (d.visible = 1 OR (a.acc_id = ad.acc_id AND a.token = ?))`)
	if err != nil {
		return err
	}

	// repo.ByOwnerStmt, err = repo.db.Prepare("SELECT * FROM DECK WHERE acc_id = ? AND (visible = 1 OR visible = ?)")
	// Two birds in one shot. Only shows hidden when token of the account matches the deck's owner token
	repo.ByOwnerStmt, err = repo.db.Prepare("SELECT * FROM DECK d, ACCOUNT acc WHERE d.acc_id = ? AND (d.visible = 1 OR (acc.token = ?))")
	if err != nil {
		return err
	}

	// repo.ByUserIdStmt, err = repo.db.Prepare("SELECT d.* FROM DECK as d, ACC_DECK as ad WHERE ad.acc_id = ? AND d.deck_id = ad.deck_id AND (visible = 1 OR visible = ?)")
	// repo.ByUserIdStmt, err = repo.db.Prepare("SELECT d.* FROM DECK d JOIN ACC_DECK ad ON d.deck_id = ad.deck_id WHERE ad.acc_id = ? AND (d.visible = 1 OR EXISTS (SELECT 1 FROM ACCOUNT acc WHERE acc.acc_id = ? AND acc.token = ?))")
	repo.ByUserIdStmt, err = repo.db.Prepare("SELECT d.* FROM DECK d, ACC_DECK ad WHERE ad.acc_id = ? AND d.deck_id = ad.deck_id AND d.visible = 1")
	if err != nil {
		return err
	}

	repo.DeleteStmt, err = repo.db.Prepare("DELETE FROM DECK WHERE deck_id = ?")
	if err != nil {
		return err
	}

	repo.CheckOwnerStmt, err = repo.db.Prepare(`SELECT d.*
													FROM DECK d
													LEFT JOIN ACCOUNT a ON d.acc_id = a.acc_id
													WHERE d.deck_id = ? AND a.token = ?`)
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

func (r *DeckRepositoryImpl) ById(request deck.ReadRequest) (model.Deck, error) {
	row := r.ByIdStmt.QueryRow(request.Id, request.Token)
	return scanDeck(row)
}

func (r *DeckRepositoryImpl) ByOwner(request deck.ReadRequest) ([]model.Deck, error) {
	rows, err := r.ByOwnerStmt.Query(request.Id, request.Token)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	//TODO:
	return scanDecks(rows)
}

func (r *DeckRepositoryImpl) ByUserId(request deck.ReadRequest) ([]model.Deck, error) {
	rows, err := r.ByUserIdStmt.Query(request.Id, request.Token)
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

// This function should be in service
func (r *DeckRepositoryImpl) CheckOwnership(deckID int64, token string) bool {
	row := r.CheckOwnerStmt.QueryRow(deckID, token)
	_, err := scanDeck(row)
	return err == nil
}

func updateDeckField(query *strings.Builder, args *[]any, field string, value any) {
	// Just checking if it's a date and it isn't empty
	if _, ok := value.(time.Time); ok && value.(time.Time).IsZero() {
		return
	}

	if value == "" || value == nil {
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
		deck.Visible,
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
			deck.Visible,
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
