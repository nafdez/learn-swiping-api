package service

import (
	"crypto/rand"
	"database/sql"
	"encoding/hex"
	"errors"
	"learn-swiping-api/erro"
	"learn-swiping-api/model"
	"learn-swiping-api/model/dto/user"
	"learn-swiping-api/repository"
	"time"

	"golang.org/x/crypto/bcrypt"
)

type UserService interface {
	Register(user.RegisterRequest) (model.User, error)
	Login(user.LoginRequest) (model.User, error)
	Token(user.TokenRequest) (model.User, error) // Login with token
	Logout(user.TokenRequest) error
	Account(user.TokenRequest) (model.User, error)
	User(user.PublicRequest) (user.Public, error)
	Update(user.UpdateRequest) error
	Delete(user.TokenRequest) error
	updateToken(old string) (string, error)
	generateToken() (string, error)
	hashPassword(password string) (string, error)
	checkPasswordHash(password, hash string) bool
}

type UserServiceImpl struct {
	repository repository.UserRepository
}

func NewUserService(repository repository.UserRepository) UserService {
	return &UserServiceImpl{repository: repository}
}

func (s *UserServiceImpl) Register(request user.RegisterRequest) (model.User, error) {
	if request.Username == "" || request.Password == "" || request.Email == "" || request.Name == "" {
		return model.User{}, erro.ErrBadField
	}

	// TODO: check email is valid

	hash, err := s.hashPassword(request.Password)
	if err != nil {
		return model.User{}, err
	}

	token, err := s.generateToken()
	if err != nil {
		return model.User{}, err
	}

	user := model.User{
		Username:     request.Username,
		Password:     hash,
		Email:        request.Email,
		Name:         request.Name,
		Token:        token,
		TokenExpires: time.Now().AddDate(0, 0, 7),
	}

	id, err := s.repository.Create(user)
	if err != nil {
		return model.User{}, err
	}

	return s.repository.ById(id)
}

func (s *UserServiceImpl) Login(request user.LoginRequest) (model.User, error) {
	if request.Username == "" || request.Password == "" {
		return model.User{}, erro.ErrBadField
	}

	user, err := s.repository.ByUsername(request.Username)
	if err != nil {
		return model.User{}, err
	}

	if s.checkPasswordHash(request.Password, user.Password) {
		token, err := s.updateToken(user.Token)
		if err != nil {
			return model.User{}, err
		}

		user.Token = token
		return user, nil
	}

	return model.User{}, nil
}

// Same as login function but using a token instead of user and password
func (s *UserServiceImpl) Token(request user.TokenRequest) (model.User, error) {
	if request.Token == "" {
		return model.User{}, erro.ErrInvalidToken
	}

	user, err := s.repository.ByToken(request.Token)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return model.User{}, erro.ErrInvalidToken
		}
		return model.User{}, err
	}

	// Just updating last seen date
	user.LastSeen = time.Now()
	err = s.repository.Update(user.ID, model.User{LastSeen: user.LastSeen})
	if err != nil {
		return model.User{}, err
	}

	return user, nil
}

func (s *UserServiceImpl) Logout(request user.TokenRequest) error {
	if request.Token == "" {
		return erro.ErrInvalidToken
	}

	// Updating token and not returning to user to invalidate the previous token
	_, err := s.updateToken(request.Token)
	return err
}

func (s *UserServiceImpl) Account(request user.TokenRequest) (model.User, error) {
	if request.Token == "" {
		return model.User{}, erro.ErrInvalidToken
	}

	return s.repository.ByToken(request.Token)
}

// Returns an account details but with some fields hidden
func (s *UserServiceImpl) User(request user.PublicRequest) (user.Public, error) {
	// TODO: Also return a list of it's public decks
	if request.Username == "" {
		return user.Public{}, erro.ErrBadField
	}

	var account model.User
	account, err := s.repository.ByUsername(request.Username)
	if err != nil {
		return user.Public{}, err
	}

	user := user.Public{
		ID:       account.ID,
		Username: account.Username,
		LastSeen: account.LastSeen,
		Since:    account.Since,
	}

	return user, nil
}

func (s *UserServiceImpl) Update(request user.UpdateRequest) error {
	if request.Token == "" {
		return erro.ErrInvalidToken
	}

	// If all fields are empty, throw an error
	if request.Username == "" && request.Password == "" && request.Email == "" && request.Name == "" {
		return erro.ErrBadField
	}

	// TODO: if email isn't empty check it is valid

	user, err := s.repository.ByToken(request.Token)
	if err != nil {
		return erro.ErrInvalidToken
	}

	hash, err := s.hashPassword(request.Password)
	if err != nil {
		return err
	}

	update := model.User{
		Username: request.Username,
		Password: hash,
		Email:    request.Email,
		Name:     request.Name,
		Token:    request.Token,
	}

	return s.repository.Update(user.ID, update)
}

func (s *UserServiceImpl) Delete(request user.TokenRequest) error {
	if request.Token == "" {
		return erro.ErrInvalidToken
	}
	user, err := s.repository.ByToken(request.Token)
	if err != nil {
		return err
	}

	return s.repository.Delete(user.ID)
}

func (s *UserServiceImpl) updateToken(old string) (string, error) {
	user, err := s.repository.ByToken(old)
	if err != nil {
		return "", err
	}

	new, err := s.generateToken()
	if err != nil {
		return "", err
	}

	err = s.repository.Update(user.ID, model.User{Token: new, TokenExpires: time.Now().AddDate(0, 0, 7)})
	if err != nil {
		return "", err
	}

	return new, nil
}

// hashPassword hashes the password provided and returns it
func (s *UserServiceImpl) hashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(bytes), err
}

// checkPasswordHash takes a password and a hashed password and returns if the hashed
// password comes from the password
func (s *UserServiceImpl) checkPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

func (s *UserServiceImpl) generateToken() (string, error) {
	b := make([]byte, 16)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	return hex.EncodeToString(b), nil
}
