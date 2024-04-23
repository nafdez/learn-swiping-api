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

type UserRepository interface {
	Create(model.User) (int64, error)
	ById(id int64) (model.User, error)
	ByUsername(username string) (model.User, error)
	ByToken(token string) (model.User, error)
	Update(id int64, user model.User) error
	Delete(id int64) error
}

type UserRepositoryImpl struct {
	db             *sql.DB
	CreateStmt     *sql.Stmt
	ByIdStmt       *sql.Stmt
	ByUsernameStmt *sql.Stmt
	DeleteStmt     *sql.Stmt
}

func NewUserRepository(db *sql.DB) *UserRepositoryImpl {
	repo := &UserRepositoryImpl{db: db}
	err := repo.InitStatements()
	if err != nil {
		log.Println(err)
	}
	return repo
}

func (r *UserRepositoryImpl) InitStatements() error {
	var err error
	r.CreateStmt, err = r.db.Prepare("INSERT INTO ACCOUNT (username, passwd, email, name, token, token_expire) VALUES (?, ?, ?, ?, ?, ?)")
	if err != nil {
		return err
	}

	r.ByIdStmt, err = r.db.Prepare("SELECT * FROM ACCOUNT WHERE acc_id = ?")
	if err != nil {
		return err
	}

	r.ByUsernameStmt, err = r.db.Prepare("SELECT * FROM ACCOUNT WHERE username = ?")
	if err != nil {
		return err
	}

	r.DeleteStmt, err = r.db.Prepare("DELETE FROM ACCOUNT WHERE acc_id = ?")
	if err != nil {
		return err
	}

	return nil
}

func (r *UserRepositoryImpl) Create(user model.User) (int64, error) {
	// TODO: make this method take undetermined number of user parameters to create a new user
	result, err := r.CreateStmt.Exec(user.Username, user.Password, user.Email, user.Name, user.Token, user.TokenExpires)
	if err != nil {
		if err.(*mysql.MySQLError).Number == 1062 {
			return 0, erro.ErrUserExists
		}
		return 0, err
	}
	return result.LastInsertId()
}

func (r *UserRepositoryImpl) ById(id int64) (model.User, error) {
	row := r.ByIdStmt.QueryRow(id)
	return scanUser(row)
}

func (r *UserRepositoryImpl) ByUsername(username string) (model.User, error) {
	row := r.ByUsernameStmt.QueryRow(username)
	return scanUser(row)
}

func (r *UserRepositoryImpl) ByToken(token string) (model.User, error) {
	// Checking token expire date on repository just for simplicity as is strange this is going to change
	// or cause problems
	stmt, err := r.db.Prepare("SELECT * FROM ACCOUNT WHERE token = ? AND token_expire >= NOW();")
	if err != nil {
		return model.User{}, err
	}

	row := stmt.QueryRow(token)
	return scanUser(row)
}

func (r *UserRepositoryImpl) Update(id int64, user model.User) error {
	var query strings.Builder
	var args []any
	query.WriteString("UPDATE ACCOUNT SET")

	// Don't know yet how to loop a struct and get its field json names
	// I feel ashamed
	updateField(&query, &args, "username", user.Username)
	updateField(&query, &args, "passwd", user.Password)
	updateField(&query, &args, "email", user.Email)
	updateField(&query, &args, "name", user.Name)
	updateField(&query, &args, "token", user.Token)
	updateField(&query, &args, "token_expire", user.TokenExpires)
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
		return erro.ErrUserNotFound
	}

	return nil
}

func (r *UserRepositoryImpl) Delete(id int64) error {
	result, err := r.DeleteStmt.Exec(id)
	if err != nil {
		return err
	}

	affected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if affected == 0 {
		return erro.ErrUserNotFound
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

func scanUser(row *sql.Row) (model.User, error) {
	var user model.User
	err := row.Scan(
		&user.ID,
		&user.Username,
		&user.Email,
		&user.Password,
		&user.Name,
		&user.Token,
		&user.TokenExpires,
		&user.LastSeen,
		&user.Since,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return model.User{}, erro.ErrUserNotFound
		}
		return model.User{}, err
	}

	return user, nil
}
