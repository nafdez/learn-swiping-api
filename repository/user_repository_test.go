package repository

import (
	"database/sql"
	"errors"
	"fmt"
	"learn-swiping-api/model"
	"log"
	"math/rand/v2"
	"os"
	"testing"

	"github.com/go-sql-driver/mysql"
	"github.com/joho/godotenv"
)

var (
	userRepository UserRepository
)

func TestMain(m *testing.M) {
	err := godotenv.Load("./../.env")
	if err != nil {
		log.Fatalf("Error loading .env file: %s", err)
	}

	user := os.Getenv("DB_USER")
	pass := os.Getenv("DB_PASS")
	url := os.Getenv("DB_URL")
	dbName := os.Getenv("DB_NAME")

	_ = mysql.Config{}

	db, err := sql.Open("mysql", fmt.Sprintf("%v:%v@tcp(%v)/%v?parseTime=true", user, pass, url, dbName))
	if err != nil {
		log.Fatal("Failed to open connection to database. Error: ", err)
	}
	defer db.Close()

	if err = db.Ping(); err != nil {
		log.Fatal("Failed to connect to database. Error: ", err)
	}

	userRepository = NewUserRepository(db)

	m.Run()
}

func TestCreate(t *testing.T) {
	var randInt int64 = rand.Int64()
	var randInt2 int64 = rand.Int64()
	testcases := []model.User{
		{
			Username: fmt.Sprintf("user%v", randInt),
			Password: "1jfk534ljl23$^'/asdfio",
			Email:    fmt.Sprintf("email@%v.com", randInt),
			Name:     fmt.Sprintf("name%v", randInt),
		},
		{
			Username: fmt.Sprintf("user%v", randInt2),
			Password: "1jfk534ljl23$^'/asdfio",
			Email:    fmt.Sprintf("email@%v.com", randInt2),
			Name:     fmt.Sprintf("name%v", randInt2),
		},
	}

	for _, testcase := range testcases {
		lastId, err := userRepository.Create(testcase)
		if err != nil {
			if err == ErrUserAlreadyExists {
				t.Error(testcase, "already exists")
			}
			t.Error(err)
		}

		if lastId == 0 {
			t.Error(errors.New("user not inserted"))
		}
	}
}

func TestById(t *testing.T) {
	var randInt int64 = rand.Int64()
	testcase := model.User{
		Username: fmt.Sprintf("user%v", randInt),
		Password: "1jfk534ljl23$^'/asdfio",
		Email:    fmt.Sprintf("email@%v.com", randInt),
		Name:     fmt.Sprintf("name%v", randInt),
	}

	id, err := userRepository.Create(testcase)
	if err != nil {
		t.Error(err)
	}

	user, err := userRepository.ById(id)
	if err != nil {
		t.Error(err)
	}

	if user.Username != testcase.Username {
		t.Error(errors.New("wrong user"))
	}
}

func TestByUsername(t *testing.T) {
	var randInt int64 = rand.Int64()
	testcase := model.User{
		Username: fmt.Sprintf("user%v", randInt),
		Password: "1jfk534ljl23$^'/asdfio",
		Email:    fmt.Sprintf("email@%v.com", randInt),
		Name:     fmt.Sprintf("name%v", randInt),
	}

	_, err := userRepository.Create(testcase)
	if err != nil {
		t.Error(err)
	}

	user, err := userRepository.ByUsername(testcase.Username)
	if err != nil {
		t.Error(err)
	}

	if user.Username != testcase.Username {
		t.Error(errors.New("wrong user"))
	}
}

func TestUpdate(t *testing.T) {
	var randInt int64 = rand.Int64()
	testcase := model.User{
		Username: fmt.Sprintf("user%v", randInt),
		Password: "1jfk534ljl23$^'/asdfio",
		Email:    fmt.Sprintf("email@%v.com", randInt),
		Name:     fmt.Sprintf("name%v", randInt),
	}

	id, err := userRepository.Create(testcase)
	if err != nil {
		t.Error(err)
	}

	newUsername := fmt.Sprintf("updatedUser%v", randInt)
	testcase.Username = newUsername

	err = userRepository.Update(id, testcase)
	if err != nil {
		t.Error(err)
	}

	user, err := userRepository.ById(id)
	if err != nil {
		t.Error(err)
	}

	if user.Username != newUsername {
		t.Error(errors.New("not updated"))
	}
}

func TestDelete(t *testing.T) {
	var randInt int64 = rand.Int64()
	testcase := model.User{
		Username: fmt.Sprintf("user%v", randInt),
		Password: "1jfk534ljl23$^'/asdfio",
		Email:    fmt.Sprintf("email@%v.com", randInt),
		Name:     fmt.Sprintf("name%v", randInt),
	}

	id, err := userRepository.Create(testcase)
	if err != nil {
		t.Error(err)
	}

	err = userRepository.Delete(id)
	if err != nil {
		t.Error(err)
	}

	_, err = userRepository.ById(id)
	if err == nil {
		t.Error(errors.New("not deleted"))
	}
}

func BenchmarkCreate(b *testing.B) {
	var randInt int64 = rand.Int64()
	testcase := model.User{
		Username: fmt.Sprintf("user%v", randInt),
		Password: "1jfk534ljl23$^'/asdfio",
		Email:    fmt.Sprintf("email@%v.com", randInt),
		Name:     fmt.Sprintf("name%v", randInt),
	}

	_, err := userRepository.Create(testcase)
	if err != nil {
		b.Error(err)
	}
}

func BenchmarkById(b *testing.B) {}

func BenchmarkByUsername(b *testing.B) {}

func BenchmarkUpdate(b *testing.B) {}

func BenchmarkDelete(b *testing.B) {}
