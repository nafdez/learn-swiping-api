package account

import (
	"database/sql"
	"fmt"
	"learn-swiping-api/erro"
	"log"
	"strings"
	"time"

	"github.com/go-sql-driver/mysql"
)

type AccountRepository interface {
	Create(Account) (int64, error)
	ById(id int64) (Account, error)
	ByUsername(Username string) (Account, error)
	ByToken(token string) (Account, error)
	Update(id int64, account Account) error
	Delete(token string) error
}

type AccountRepositoryImpl struct {
	db              *sql.DB
	CreateStmt      *sql.Stmt
	ByIdStmt        *sql.Stmt
	ByUsernameStmt  *sql.Stmt
	DeleteStmt      *sql.Stmt
	UnlinkDecksStmt *sql.Stmt
}

func NewAccountRepository(db *sql.DB) *AccountRepositoryImpl {
	repo := &AccountRepositoryImpl{db: db}
	err := repo.InitStatements()
	if err != nil {
		log.Println(err)
	}
	return repo
}

func (r *AccountRepositoryImpl) InitStatements() error {
	var err error
	r.CreateStmt, err = r.db.Prepare("INSERT INTO ACCOUNT (Username, passwd, email, name, token, token_expire) VALUES (?, ?, ?, ?, ?, ?)")
	if err != nil {
		return err
	}

	r.ByIdStmt, err = r.db.Prepare("SELECT * FROM ACCOUNT WHERE acc_id = ?")
	if err != nil {
		return err
	}

	r.ByUsernameStmt, err = r.db.Prepare("SELECT * FROM ACCOUNT WHERE Username = ?")
	if err != nil {
		return err
	}

	r.DeleteStmt, err = r.db.Prepare("DELETE FROM ACCOUNT WHERE token = ?")
	if err != nil {
		return err
	}

	r.UnlinkDecksStmt, err = r.db.Prepare(`UPDATE DECK d 
											LEFT JOIN ACCOUNT a ON d.acc_id = a.acc_id
											SET d.acc_id = 1
											WHERE d.visible = 1 AND a.token = ?`)
	if err != nil {
		return err
	}

	return nil
}

func (r *AccountRepositoryImpl) Create(account Account) (int64, error) {
	// TODO: make this method take undetermined number of account parameters to create a new account
	result, err := r.CreateStmt.Exec(account.Username, account.Password, account.Email, account.Name, account.Token, account.TokenExpires)
	if err != nil {
		if err.(*mysql.MySQLError).Number == 1062 {
			return 0, erro.ErrAccountExists
		}
		return 0, err
	}
	return result.LastInsertId()
}

func (r *AccountRepositoryImpl) ById(id int64) (Account, error) {
	row := r.ByIdStmt.QueryRow(id)
	return scanaccount(row)
}

func (r *AccountRepositoryImpl) ByUsername(Username string) (Account, error) {
	row := r.ByUsernameStmt.QueryRow(Username)
	return scanaccount(row)
}

func (r *AccountRepositoryImpl) ByToken(token string) (Account, error) {
	// Checking token expire date on repository just for simplicity as is strange this is going to change
	// or cause problems
	stmt, err := r.db.Prepare("SELECT * FROM ACCOUNT WHERE token = ? AND token_expire >= NOW();")
	if err != nil {
		return Account{}, err
	}

	row := stmt.QueryRow(token)
	return scanaccount(row)
}

func (r *AccountRepositoryImpl) Update(id int64, account Account) error {
	var query strings.Builder
	var args []any
	query.WriteString("UPDATE ACCOUNT SET")

	// Don't know yet how to loop a struct and get its field json names
	// I feel ashamed
	updateField(&query, &args, "username", account.Username)
	updateField(&query, &args, "passwd", account.Password)
	updateField(&query, &args, "email", account.Email)
	updateField(&query, &args, "name", account.Name)
	updateField(&query, &args, "token", account.Token)
	updateField(&query, &args, "token_expire", account.TokenExpires)
	updateField(&query, &args, "last_seen", time.Now())

	args = append(args, id)
	query.WriteString(" WHERE acc_id = ?")

	stmt, err := r.db.Prepare(query.String())
	if err != nil {
		return err
	}

	result, err := stmt.Exec(args...)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return erro.ErrAccountNotFound
	}

	return nil
}

func (r *AccountRepositoryImpl) Delete(token string) error {
	// Necessary to not to delete decks when account is removed
	// deck's owner now is account 1 (deleted user)
	// Note that only public decks are saved into the
	// auxiliar account, the hidden ones are removed
	_, err := r.UnlinkDecksStmt.Exec(token)
	if err != nil {
		return err
	}

	result, err := r.DeleteStmt.Exec(token)
	if err != nil {
		return err
	}

	affected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if affected == 0 {
		return erro.ErrAccountNotFound
	}

	return nil
}

func updateField(query *strings.Builder, args *[]any, field string, value any) {
	// Just checking if it's a date and it isn't empty
	if _, ok := value.(time.Time); ok && value.(time.Time).IsZero() {
		return
	}

	if value == "" {
		return
	}

	if query.String() != "UPDATE ACCOUNT SET" {
		query.WriteString(",")
	}

	query.WriteString(fmt.Sprintf(" %s = ?", field))
	*args = append(*args, value)
}

func scanaccount(row *sql.Row) (Account, error) {
	var account Account
	err := row.Scan(
		&account.ID,
		&account.Username,
		&account.Email,
		&account.Password,
		&account.Name,
		&account.Token,
		&account.TokenExpires,
		&account.LastSeen,
		&account.Since,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return Account{}, erro.ErrAccountNotFound
		}
		return Account{}, err
	}

	return account, nil
}
