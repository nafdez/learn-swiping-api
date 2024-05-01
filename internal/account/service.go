package account

import (
	"bytes"
	"crypto/rand"
	"database/sql"
	"encoding/hex"
	"errors"
	"io"
	"learn-swiping-api/erro"
	account "learn-swiping-api/internal/account/dto"
	"learn-swiping-api/internal/picture"
	"net/mail"
	"path/filepath"
	"time"

	"golang.org/x/crypto/bcrypt"
)

type AccountService interface {
	Register(account.RegisterRequest) (Account, error)
	Login(account.LoginRequest) (Account, error)
	Token(account.TokenRequest) (Account, error) // Login with token
	Logout(account.TokenRequest) error
	Account(account.TokenRequest) (Account, error)
	account(account.PublicRequest) (account.Public, error)
	Update(account.UpdateRequest) error
	Delete(account.TokenRequest) error
	updateToken(old string) (string, error)
	generateToken() (string, error)
	hashPassword(password string) (string, error)
	checkPasswordHash(password, hash string) bool
}

type AccountServiceImpl struct {
	repository AccountRepository
}

func NewAccountService(repository AccountRepository) AccountService {
	return &AccountServiceImpl{repository: repository}
}

func (s *AccountServiceImpl) Register(request account.RegisterRequest) (Account, error) {
	if _, err := mail.ParseAddress(request.Email); err != nil {
		return Account{}, erro.ErrInvalidEmail
	}

	hash, err := s.hashPassword(request.Password)
	if err != nil {
		return Account{}, err
	}

	token, err := s.generateToken()
	if err != nil {
		return Account{}, err
	}

	account := Account{
		Username:     request.Username,
		Password:     hash,
		Email:        request.Email,
		Name:         request.Name,
		Token:        token,
		TokenExpires: time.Now().AddDate(0, 0, 7),
	}

	id, err := s.repository.Create(account)
	if err != nil {
		return Account{}, err
	}

	return s.repository.ById(id)
}

func (s *AccountServiceImpl) Login(request account.LoginRequest) (Account, error) {
	account, err := s.repository.ByUsername(request.Username)
	if err != nil {
		return Account{}, err
	}

	if s.checkPasswordHash(request.Password, account.Password) {
		token, err := s.updateToken(account.Token)
		if err != nil {
			return Account{}, err
		}

		account.Token = token
		return account, nil
	}

	return account, nil
}

// Same as login function but using a token instead of account and password
func (s *AccountServiceImpl) Token(request account.TokenRequest) (Account, error) {
	account, err := s.repository.ByToken(request.Token)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return Account{}, erro.ErrInvalidToken
		}
		return Account{}, err
	}

	// Just updating last seen date
	account.LastSeen = time.Now()
	err = s.repository.Update(account.ID, Account{LastSeen: account.LastSeen})
	if err != nil {
		return Account{}, err
	}

	return account, nil
}

func (s *AccountServiceImpl) Logout(request account.TokenRequest) error {
	// Updating token and not returning to account to invalidate the previous token
	_, err := s.updateToken(request.Token)
	return err
}

func (s *AccountServiceImpl) Account(request account.TokenRequest) (Account, error) {
	return s.repository.ByToken(request.Token)
}

// Returns an account details but with some fields hidden
func (s *AccountServiceImpl) account(request account.PublicRequest) (account.Public, error) {
	// TODO: Also return a list of it's public decks
	if request.Username == "" {
		return account.Public{}, erro.ErrBadField
	}

	var acc Account
	acc, err := s.repository.ByUsername(request.Username)
	if err != nil {
		return account.Public{}, err
	}

	accountPublic := account.Public{
		ID:       acc.ID,
		Username: acc.Username,
		PicID:    acc.PicID,
		LastSeen: acc.LastSeen,
		Since:    acc.Since,
	}

	return accountPublic, nil
}

func (s *AccountServiceImpl) Update(request account.UpdateRequest) error {
	// If all fields are empty, throw an error
	if request.Username == "" && request.Password == "" && request.Email == "" && request.Name == "" && request.Img == nil {
		return erro.ErrBadField
	}

	// Check if token is valid
	account, err := s.repository.ByToken(request.Token)
	if err != nil {
		return erro.ErrInvalidToken
	}

	// Check if email is valid
	if request.Email != "" {
		if _, err := mail.ParseAddress(request.Email); err != nil {
			return erro.ErrInvalidEmail
		}
	}

	var updateAcc Account

	// Check if password is not empty then hashing it
	if request.Password != "" {
		hash, err := s.hashPassword(request.Password)
		if err != nil {
			return err
		}
		updateAcc.Password = hash
	}

	if request.Img != nil {
		img, err := request.Img.Open()
		if err != nil {
			return erro.ErrBadField
		}
		defer img.Close()

		buf := bytes.NewBuffer(nil)
		if _, err := io.Copy(buf, img); err != nil {
			return erro.ErrBadField
		}

		picture.Remove(account.PicID)
		picID, err := picture.Store(filepath.Ext(request.Img.Filename), buf.Bytes())
		if err != nil {
			return err
		}
		updateAcc.PicID = picID
	}

	// If a value is empty it the repository won't update
	// it anyway
	updateAcc.Username = request.Username
	updateAcc.Email = request.Email
	updateAcc.Name = request.Name
	updateAcc.Token = request.Token

	return s.repository.Update(account.ID, updateAcc)
}

func (s *AccountServiceImpl) Delete(request account.TokenRequest) error {
	acc, err := s.repository.ByToken(request.Token)
	if err != nil {
		return erro.ErrInvalidToken
	}
	picture.Remove(acc.PicID)
	return s.repository.Delete(request.Token)
}

func (s *AccountServiceImpl) updateToken(old string) (string, error) {
	account, err := s.repository.ByToken(old)
	if err != nil {
		return "", err
	}

	new, err := s.generateToken()
	if err != nil {
		return "", err
	}

	err = s.repository.Update(account.ID, Account{Token: new, TokenExpires: time.Now().AddDate(0, 0, 7)})
	if err != nil {
		return "", err
	}

	return new, nil
}

// hashPassword hashes the password provided and returns it
func (s *AccountServiceImpl) hashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 8)
	return string(bytes), err
}

// checkPasswordHash takes a password and a hashed password and returns if the hashed
// password comes from the password
func (s *AccountServiceImpl) checkPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

func (s *AccountServiceImpl) generateToken() (string, error) {
	b := make([]byte, 16)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	return hex.EncodeToString(b), nil
}
