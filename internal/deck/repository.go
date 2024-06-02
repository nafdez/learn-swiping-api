package deck

import (
	"database/sql"
	"fmt"
	"learn-swiping-api/erro"
	deck "learn-swiping-api/internal/deck/dto"
	"log"
	"strings"
	"time"

	"github.com/go-sql-driver/mysql"
)

type DeckRepository interface {
	Create(deck.CreateRequest) (int64, error)
	ById(deckID int64, token string) (Deck, error)
	ByOwner(accID int64, username string, token string) ([]Deck, error)
	BySubsUsername(username string, token string) ([]Deck, error) // ACC-DECK table
	Update(id int64, deck Deck) error
	Delete(id int64) error
	// TODO: Only insert sub if token matches with the account's token of the id
	AddDeckSubscription(token string, deckId int64) error
	RemoveDeckSubscription(token string, deckId int64) error
	CheckOwnership(deckID int64, token string) bool

	DeckDetailsSubscription(deckID int64, token string) (deck.Details, error)
	DeckDetailsOwner(deckID int64, token string) (deck.Details, error)
	DeckDetailsShop(deckID int64) (deck.Details, error)
}

type DeckRepositoryImpl struct {
	db                         *sql.DB
	CreateStmt                 *sql.Stmt
	ByIdStmt                   *sql.Stmt
	ByOwnerStmt                *sql.Stmt
	BySubsUsernameStmt         *sql.Stmt
	DeleteStmt                 *sql.Stmt
	AddDeckSubscriptionStmt    *sql.Stmt
	RemoveDeckSubscriptionStmt *sql.Stmt
	CheckOwnerStmt             *sql.Stmt

	DeckDetailsOwnerStmt *sql.Stmt
	DeckDetailsShopStmt  *sql.Stmt
	DeckDetailsSubsStmt  *sql.Stmt
}

func NewDeckRepository(db *sql.DB) *DeckRepositoryImpl {
	repo := &DeckRepositoryImpl{db: db}
	err := repo.InitStatements()
	if err != nil {
		log.Fatalln(err)
	}
	return repo
}

func (repo *DeckRepositoryImpl) InitStatements() error {
	var err error
	repo.CreateStmt, err = repo.db.Prepare(`INSERT INTO DECK (acc_id, title, description, visible) 
												VALUES ((SELECT acc_id FROM ACCOUNT WHERE token = ?), ?, ?, ?)`)
	if err != nil {
		return err
	}

	repo.ByIdStmt, err = repo.db.Prepare(`SELECT d.* FROM DECK d 
											LEFT JOIN ACCOUNT a ON d.acc_id = a.acc_id
											WHERE d.deck_id = ? 
												AND (d.visible = 1 OR a.token = ?)`)
	if err != nil {
		return err
	}

	repo.ByOwnerStmt, err = repo.db.Prepare(`SELECT d.*
												FROM DECK d
												LEFT JOIN ACCOUNT a ON d.acc_id = a.acc_id 
												WHERE (a.acc_id = ? OR a.username = ?) 
												AND (d.visible = 1 OR a.token = ?)`)
	if err != nil {
		return err
	}

	// revisar
	repo.BySubsUsernameStmt, err = repo.db.Prepare(`SELECT d.* FROM DECK d 
														LEFT JOIN ACC_DECK ad ON d.deck_id = ad.deck_id
														LEFT JOIN ACCOUNT a ON d.acc_id = a.acc_id
														LEFT JOIN ACCOUNT acc ON ad.acc_id = acc.acc_id
														WHERE acc.username = ?
														AND (d.visible = 1 OR (a.acc_id = ad.acc_id AND a.token = ?))`)
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

	repo.AddDeckSubscriptionStmt, err = repo.db.Prepare("INSERT INTO ACC_DECK(acc_id, deck_id) VALUES ((SELECT acc_id FROM ACCOUNT WHERE token = ?), ?)")
	if err != nil {
		return err
	}

	repo.RemoveDeckSubscriptionStmt, err = repo.db.Prepare("DELETE FROM ACC_DECK WHERE acc_id = (SELECT acc_id FROM ACCOUNT WHERE token = ?) AND deck_id = ?")
	if err != nil {
		return err
	}

	repo.DeckDetailsOwnerStmt, err = repo.db.Prepare(`SELECT 
														DECK.title,
														DECK.description,
														DECK.pic_id,
														COUNT(ACC_DECK.deck_id) AS subscriptions,
                                                        COUNT(adSubscribed.deck_id) As is_subscribed,
    													DECK.visible,
    													DECK.updated_at,
    													DECK.created_at
													FROM 
    													DECK 
													LEFT JOIN 
    													ACC_DECK ON DECK.deck_id = ACC_DECK.deck_id
													LEFT JOIN
														ACCOUNT ON DECK.acc_id = ACCOUNT.acc_id
													LEFT JOIN 
														ACC_DECK As adSubscribed ON adSubscribed.deck_id = DECK.deck_id AND adSubscribed.acc_id = ACCOUNT.acc_id
													WHERE 
    													DECK.deck_id = ?
    													AND ACCOUNT.token = ?
													GROUP BY 
    													DECK.deck_id`)
	if err != nil {
		return err
	}

	repo.DeckDetailsShopStmt, err = repo.db.Prepare(`SELECT 
														DECK.title,
														DECK.description,
														DECK.pic_id,
    													COUNT(ACC_DECK.deck_id) AS subscriptions,
                                                        COUNT(adSubscribed.deck_id) As is_subscribed,
														COUNT(CARD.card_id) AS cards,
														ACCOUNT.acc_id,
														ACCOUNT.username,
														DECK.updated_at,
														DECK.created_at
													FROM 
														DECK 
													LEFT JOIN 
														ACC_DECK ON DECK.deck_id = ACC_DECK.deck_id
													LEFT JOIN
														CARD ON DECK.deck_id = CARD.deck_id
													LEFT JOIN
														ACCOUNT ON DECK.acc_id = ACCOUNT.acc_id
													LEFT JOIN 
														ACC_DECK As adSubscribed ON adSubscribed.deck_id = DECK.deck_id AND adSubscribed.acc_id = ACCOUNT.acc_id
													WHERE 
														DECK.deck_id = ?
													GROUP BY
														DECK.deck_id`)
	if err != nil {
		return err
	}

	repo.DeckDetailsSubsStmt, err = repo.db.Prepare(`SELECT
														DECK.title,
														DECK.description,
														DECK.pic_id,
														coalesce((count(PROGRESS.progress_id) / COUNT(CARD.card_id) * 100), 0) AS total_progress,
														COUNT(PROGRESS.progress_id) AS cards_revised,
														coalesce((count(CARD.card_id) - count(PROGRESS.progress_id)), 0) AS cards_remaining,
														DECK.updated_at,
														DECK.created_at
													FROM 
														DECK
													LEFT JOIN 
														CARD ON DECK.deck_id = CARD.deck_id
													LEFT JOIN
														PROGRESS ON CARD.card_id = PROGRESS.card_id
													LEFT JOIN 
														ACCOUNT ON PROGRESS.acc_id = ACCOUNT.acc_id
													WHERE 
															DECK.deck_id = ?
                                                            AND ACCOUNT.token = ?
													GROUP BY
														DECK.deck_id`)
	if err != nil {
		return err
	}

	return nil
}

func (r *DeckRepositoryImpl) Create(deck deck.CreateRequest) (int64, error) {
	result, err := r.CreateStmt.Exec(deck.Token, deck.Title, deck.Description, deck.Visible)
	if err != nil {
		if err.(*mysql.MySQLError).Number == 1452 {
			return 0, erro.ErrAccountNotFound
		}
		if err.(*mysql.MySQLError).Number == 1062 {
			return 0, erro.ErrDeckExists
		}
		if err.(*mysql.MySQLError).Number == 1048 {
			return 0, erro.ErrInvalidToken
		}
		return 0, err
	}
	return result.LastInsertId()
}

func (r *DeckRepositoryImpl) ById(deckID int64, token string) (Deck, error) {
	row := r.ByIdStmt.QueryRow(deckID, token)
	return scanDeck(row)
}

func (r *DeckRepositoryImpl) ByOwner(accID int64, username string, token string) ([]Deck, error) {
	rows, err := r.ByOwnerStmt.Query(accID, username, token)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	return scanDecks(rows)
}

func (r *DeckRepositoryImpl) BySubsUsername(username string, token string) ([]Deck, error) {
	rows, err := r.BySubsUsernameStmt.Query(username, token)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	return scanDecks(rows)
}

func (r *DeckRepositoryImpl) Update(id int64, deck Deck) error {
	var query strings.Builder
	var args []any
	query.WriteString("UPDATE DECK SET")

	updateDeckField(&query, &args, "title", deck.Title)
	updateDeckField(&query, &args, "description", deck.Description)
	updateDeckField(&query, &args, "pic_id", deck.PicID)
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

func (r *DeckRepositoryImpl) AddDeckSubscription(token string, deckId int64) error {
	_, err := r.AddDeckSubscriptionStmt.Exec(token, deckId)
	if err != nil {
		if err.(*mysql.MySQLError).Number == 1062 {
			return erro.ErrAlreadySuscribed
		}
		if err.(*mysql.MySQLError).Number == 1048 {
			return erro.ErrInvalidToken
		}
		return err
	}
	return nil
}

func (r *DeckRepositoryImpl) RemoveDeckSubscription(token string, deckId int64) error {
	result, err := r.RemoveDeckSubscriptionStmt.Exec(token, deckId)
	if err != nil {
		if err.(*mysql.MySQLError).Number == 1048 {
			return erro.ErrInvalidToken
		}
		return err
	}

	affected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if affected == 0 {
		return erro.ErrNotSuscribed
	}

	return nil
}

// This function should be in service
func (r *DeckRepositoryImpl) CheckOwnership(deckID int64, token string) bool {
	row := r.CheckOwnerStmt.QueryRow(deckID, token)
	_, err := scanDeck(row)
	return err == nil
}

func (r *DeckRepositoryImpl) DeckDetailsSubscription(deckID int64, token string) (deck.Details, error) {
	row := r.DeckDetailsOwnerStmt.QueryRow(deckID, token)

	var details deck.Details
	err := row.Scan(
		&details.Title,
		&details.Description,
		&details.PicID,
		&details.TotalProgress,
		&details.CardsRevised,
		&details.CardsRemaining,
		&details.UpdatedAt,
		&details.CreatedAt,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return deck.Details{}, erro.ErrDeckNotFound
		}
		return deck.Details{}, err
	}

	return details, nil
}

func (r *DeckRepositoryImpl) DeckDetailsOwner(deckID int64, token string) (deck.Details, error) {
	row := r.DeckDetailsOwnerStmt.QueryRow(deckID, token)

	var details deck.Details
	err := row.Scan(
		&details.Title,
		&details.Description,
		&details.PicID,
		&details.Subscriptions,
		&details.IsSubscribed,
		&details.IsVisible,
		&details.UpdatedAt,
		&details.CreatedAt,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return deck.Details{}, erro.ErrDeckNotFound
		}
		return deck.Details{}, err
	}

	return details, nil
}

func (r *DeckRepositoryImpl) DeckDetailsShop(deckID int64) (deck.Details, error) {
	row := r.DeckDetailsShopStmt.QueryRow(deckID)

	var details deck.Details
	err := row.Scan(
		&details.Title,
		&details.Description,
		&details.PicID,
		&details.Subscriptions,
		&details.IsSubscribed,
		&details.Cards,
		&details.OwnerID,
		&details.Owner,
		&details.UpdatedAt,
		&details.CreatedAt,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return deck.Details{}, erro.ErrDeckNotFound
		}
		return deck.Details{}, err
	}

	return details, nil
}

func updateDeckField(query *strings.Builder, args *[]any, field string, value any) {
	// Just checking if it's a date and it isn't empty
	if _, ok := value.(time.Time); ok && value.(time.Time).IsZero() {
		return
	}

	if b, ok := value.(*bool); ok {
		if b == nil {
			return
		}
	} else if value == "" || value == nil {
		return
	}

	if query.String() != "UPDATE DECK SET" {
		query.WriteString(",")
	}

	query.WriteString(fmt.Sprintf(" %s = ?", field))
	*args = append(*args, value)
}

func scanDeck(row *sql.Row) (Deck, error) {
	var deck Deck
	err := row.Scan(
		&deck.ID,
		&deck.Owner,
		&deck.Title,
		&deck.Description,
		&deck.PicID,
		&deck.Visible,
		&deck.UpdatedAt,
		&deck.CreatedAt,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return Deck{}, erro.ErrDeckNotFound
		}
		return Deck{}, err
	}

	return deck, nil
}

func scanDecks(rows *sql.Rows) ([]Deck, error) {
	var decks []Deck
	var deck Deck
	for rows.Next() {
		err := rows.Scan(
			&deck.ID,
			&deck.Owner,
			&deck.Title,
			&deck.Description,
			&deck.PicID,
			&deck.Visible,
			&deck.UpdatedAt,
			&deck.CreatedAt,
		)
		if err != nil {
			return []Deck{}, err
		}
		decks = append(decks, deck)
	}

	if len(decks) == 0 {
		return []Deck{}, erro.ErrDeckNotFound
	}

	return decks, nil
}
