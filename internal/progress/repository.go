package progress

import (
	"database/sql"
	"fmt"
	"learn-swiping-api/erro"
	progress "learn-swiping-api/internal/progress/dto"
	"strings"

	"github.com/go-sql-driver/mysql"
)

type ProgressRepository interface {
	Create(progress.AccessRequest) (int64, error)
	Progress(progress.AccessRequest) (Progress, error)
	Update(progress.UpdateRequest) error
	Delete(progress.AccessRequest) error
}

type ProgressRepositoryImpl struct {
	db         *sql.DB
	CreateStmt *sql.Stmt
	ByCardID   *sql.Stmt
	DeleteStmt *sql.Stmt
}

func NewProgressRepository(db *sql.DB) ProgressRepository {
	return &ProgressRepositoryImpl{db: db}
}

func (r *ProgressRepositoryImpl) InitStatements() error {
	var err error
	r.CreateStmt, err = r.db.Prepare(`INSERT INTO PROGRESS (acc_id, card_id) 
								VALUES ((SELECT acc_id FROM ACCOUNT WHERE token = ?), ?)`)
	if err != nil {
		return err
	}

	r.ByCardID, err = r.db.Prepare(`SELECT p.* FROM PROGRESS p
										LEFT JOIN ACCOUNT a ON p.acc_id = a.acc_id
    									WHERE a.token = ? AND p.card_id = ?`)
	if err != nil {
		return err
	}

	r.DeleteStmt, err = r.db.Prepare("DELETE FROM PROGRESS WHERE acc_id = (SELECT acc_id FROM ACCOUNT WHERE token = ?) AND card_id = ?")
	if err != nil {
		return err
	}

	return nil
}

func (r *ProgressRepositoryImpl) Create(req progress.AccessRequest) (int64, error) {
	result, err := r.CreateStmt.Exec(req.Token, req.CardID)
	if err != nil {
		if err.(*mysql.MySQLError).Number == 1452 {
			return 0, erro.ErrCardNotFound
		}
		if err.(*mysql.MySQLError).Number == 1062 {
			return 0, erro.ErrProgressExists
		}
		if err.(*mysql.MySQLError).Number == 1048 {
			return 0, erro.ErrInvalidToken
		}
		return 0, err
	}

	return result.LastInsertId()
}

func (r *ProgressRepositoryImpl) Progress(req progress.AccessRequest) (Progress, error) {
	row := r.ByCardID.QueryRow(req.Token, req.CardID)

	var progress Progress
	err := row.Scan(
		&progress.ProgressID,
		&progress.AccID,
		&progress.CardID,
		&progress.Priority,
		&progress.DaysHidden,
		&progress.WatchCount,
		&progress.PriorityExam,
		&progress.DaysHiddenExam,
		&progress.AnswerCount,
		&progress.CorrectCount,
		&progress.IsBuried,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return Progress{}, erro.ErrProgressNotFound
		}
		return Progress{}, err
	}

	return progress, nil
}

func (r *ProgressRepositoryImpl) Update(req progress.UpdateRequest) error {
	var query strings.Builder
	var args []any
	query.WriteString("UPDATE PROGRESS SET")

	updateProgressField(&query, &args, "priority", req.Priority)
	updateProgressField(&query, &args, "days_hidden", req.DaysHidden)
	updateProgressField(&query, &args, "watch_count", req.WatchCount)
	updateProgressField(&query, &args, "priority_exam", req.PriorityExam)
	updateProgressField(&query, &args, "days_hidden_exam", req.DaysHiddenExam)
	updateProgressField(&query, &args, "answer_count", req.AnswerCount)
	updateProgressField(&query, &args, "correct_count", req.CorrectCount)
	updateProgressField(&query, &args, "is_buried", req.IsBuried)
	// TODO: last update, creation date

	args = append(args, req.Token, req.CardID)
	query.WriteString(" WHERE acc_id = (SELECT acc_id FROM ACCOUNT WHERE token = ?) AND card_id = ?")

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
		return erro.ErrProgressNotFound
	}

	return nil
}

func (r *ProgressRepositoryImpl) Delete(req progress.AccessRequest) error {
	result, err := r.DeleteStmt.Exec(req.Token, req.CardID)
	if err != nil {
		return err
	}

	affected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if affected == 0 {
		return erro.ErrProgressNotFound
	}

	return nil
}

func updateProgressField(query *strings.Builder, args *[]any, field string, value any) {
	if value == "" || value == nil {
		return
	}

	if query.String() != "UPDATE PROGRESS SET" {
		query.WriteString(",")
	}

	query.WriteString(fmt.Sprintf(" %s = ?", field))
	*args = append(*args, value)
}
