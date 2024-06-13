package progress

import (
	"database/sql"
	"fmt"
	"learn-swiping-api/erro"
	progress "learn-swiping-api/internal/progress/dto"
	"log"
	"reflect"

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
	repo := &ProgressRepositoryImpl{db: db}
	err := repo.InitStatements()
	if err != nil {
		log.Fatalln(err)
	}
	return repo
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
		&progress.Ease,
		&progress.Interval,
		&progress.DaysHidden,
		&progress.WatchCount,
		&progress.PriorityExam,
		&progress.DaysHiddenExam,
		&progress.AnswerCount,
		&progress.CorrectCount,
		&progress.IsRelearning,
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

func (r *ProgressRepositoryImpl) Delete(req progress.AccessRequest) error {
	result, err := r.DeleteStmt.Exec(req.Token, req.CardID)
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
		return erro.ErrProgressNotFound
	}

	return nil
}

func (r *ProgressRepositoryImpl) Update(req progress.UpdateRequest) error {
	query := make(map[string]string)
	var args []any

	query["insertColumns"] = "INSERT INTO PROGRESS ("
	query["insertValues"] = "SELECT"
	query["onDuplicate"] = "ON DUPLICATE KEY UPDATE"

	upsertProgressField(query, &args, "acc_id", req.Token)
	upsertProgressField(query, &args, "card_id", req.CardID)
	upsertProgressField(query, &args, "ease", req.Ease)
	upsertProgressField(query, &args, "`interval`", req.Interval)
	upsertProgressField(query, &args, "priority", req.Priority)
	upsertProgressField(query, &args, "days_hidden", req.DaysHidden)
	upsertProgressField(query, &args, "watch_count", req.WatchCount)
	upsertProgressField(query, &args, "priority_exam", req.PriorityExam)
	upsertProgressField(query, &args, "days_hidden_exam", req.DaysHiddenExam)
	upsertProgressField(query, &args, "answer_count", req.AnswerCount)
	upsertProgressField(query, &args, "correct_count", req.CorrectCount)
	upsertProgressField(query, &args, "is_relearning", req.IsRelearning)
	upsertProgressField(query, &args, "is_buried", req.IsBuried)
	// TODO: last update, creation date

	appendToStringMap(query, "insertColumns", ")")
	appendToStringMap(query, "insertValues", " FROM ACCOUNT WHERE token = ?")
	args = append(args, req.Token)

	strQuery := fmt.Sprintf("%s %s %s", query["insertColumns"], query["insertValues"], query["onDuplicate"])
	// log.Println(strQuery)

	// for _, arg := range args {
	// log.Print(arg)
	// }

	stmt, err := r.db.Prepare(strQuery)
	if err != nil {
		return err
	}

	_, err = stmt.Exec(args...)
	if err != nil {
		if err.(*mysql.MySQLError).Number == 1048 {
			return erro.ErrInvalidToken
		}
		return err
	}

	return nil
}

func upsertProgressField(query map[string]string, args *[]any, field string, value any) {
	if b, ok := value.(*bool); ok {
		if b == nil {
			return
		}
	} else if i, ok := value.(*int); ok {
		if i == nil {
			return
		}
	} else if value == "" || value == nil {
		return
	}

	if query["insertColumns"] != "INSERT INTO PROGRESS (" {
		appendToStringMap(query, "insertColumns", ",")
		appendToStringMap(query, "insertValues", ",")
		appendToStringMap(query, "onDuplicate", ",")
	}

	appendToStringMap(query, "insertColumns", fmt.Sprintf(" %s", field))
	appendToStringMap(query, "onDuplicate", fmt.Sprintf(" %s = VALUES(%s)", field, field))
	if field == "acc_id" {
		appendToStringMap(query, "insertValues", " acc_id")
		return
	} else {
		appendToStringMap(query, "insertValues", " ?")
	}

	if reflect.ValueOf(value).Kind() == reflect.Ptr {
		*args = append(*args, reflect.ValueOf(value).Elem().Interface())
		return
	}

	*args = append(*args, value)
}

func appendToStringMap(m map[string]string, key string, val string) {
	m[key] = m[key] + val
}

// func (r *ProgressRepositoryImpl) Update(req progress.UpdateRequest) error {
// 	var query strings.Builder
// 	var args []any
// 	query.WriteString("UPDATE PROGRESS SET")

// 	updateProgressField(&query, &args, "ease", req.Ease)
// 	updateProgressField(&query, &args, "`interval`", req.Interval)
// 	updateProgressField(&query, &args, "priority", req.Priority)
// 	updateProgressField(&query, &args, "days_hidden", req.DaysHidden)
// 	updateProgressField(&query, &args, "watch_count", req.WatchCount)
// 	updateProgressField(&query, &args, "priority_exam", req.PriorityExam)
// 	updateProgressField(&query, &args, "days_hidden_exam", req.DaysHiddenExam)
// 	updateProgressField(&query, &args, "answer_count", req.AnswerCount)
// 	updateProgressField(&query, &args, "correct_count", req.CorrectCount)
// 	updateProgressField(&query, &args, "is_relearning", req.IsRelearning)
// 	updateProgressField(&query, &args, "is_buried", req.IsBuried)
// 	// TODO: last update, creation date

// 	args = append(args, req.Token, req.CardID)
// 	query.WriteString(" WHERE acc_id = (SELECT acc_id FROM ACCOUNT WHERE token = ?) AND card_id = ?")

// 	stmt, err := r.db.Prepare(query.String())
// 	if err != nil {
// 		return err
// 	}

// 	result, err := stmt.Exec(args...)
// 	if err != nil {
// 		if err.(*mysql.MySQLError).Number == 1048 {
// 			return erro.ErrInvalidToken
// 		}
// 		return err
// 	}

// 	affected, err := result.RowsAffected()
// 	if err != nil {
// 		return err
// 	}

// 	if affected == 0 {
// 		return erro.ErrProgressNotFound
// 	}

// 	return nil
// }

// func updateProgressField(query *strings.Builder, args *[]any, field string, value any) {
// 	if b, ok := value.(*bool); ok {
// 		if b == nil {
// 			return
// 		}
// 	} else if i, ok := value.(*int); ok {
// 		if i == nil {
// 			return
// 		}
// 	} else if value == "" || value == nil {
// 		return
// 	}

// 	if query.String() != "UPDATE PROGRESS SET" {
// 		query.WriteString(",")
// 	}

// 	query.WriteString(fmt.Sprintf(" %s = ?", field))

// 	if reflect.ValueOf(value).Kind() == reflect.Ptr {
// 		*args = append(*args, reflect.ValueOf(value).Elem().Interface())
// 		return
// 	}

// 	*args = append(*args, value)
// }
