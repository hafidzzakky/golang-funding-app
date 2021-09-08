package user

import (
	"errors"

	"golang.org/x/crypto/bcrypt"
)

type Service interface {
	RegisterUser(input RegisterUserInput) (User, error)
	LoginUser(input LoginUserInput) (User, error)
	IsEmailAvailable(input CheckEmailInput) (bool, error)
	SaveAvatar(id int, fileLocation string) (User, error)
}

type service struct {
	repository Repository
}

func NewService(repository Repository) *service {
	return &service{repository: repository}
}

func (s *service) RegisterUser(input RegisterUserInput) (User, error) {
	// Mapping Data
	user := User{}
	user.Name = input.Name
	user.Occupation = input.Occupation
	user.Email = input.Email
	passwordHash, err := bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.MinCost)
	if err != nil {
		return user, err
	}
	user.PasswordHash = string(passwordHash)
	user.Role = "user"

	// Input Via Repository
	newUser, err := s.repository.Save(user)
	if err != nil {
		return newUser, err
	}

	return newUser, nil
}

func (s *service) LoginUser(input LoginUserInput) (User, error) {
	email := input.Email
	password := input.Password
	findUser, err := s.repository.FindByEmail(email)

	if err != nil {
		return findUser, err
	}

	if findUser.ID == 0 {
		return findUser, errors.New("User not found")
	}

	// Check user password
	err = bcrypt.CompareHashAndPassword([]byte(findUser.PasswordHash), []byte(password))
	if err != nil {
		return findUser, errors.New("Email and password didn't match")
	}

	return findUser, nil
}

func (s *service) IsEmailAvailable(input CheckEmailInput) (bool, error) {
	email := input.Email
	user, err := s.repository.FindByEmail(email)
	if err != nil {
		return false, err
	}

	if user.ID == 0 {
		return true, nil
	}

	return false, nil
}

func (s *service) SaveAvatar(id int, fileLocation string) (User, error) {
	user, err := s.repository.FindByID(id)
	if err != nil {
		return user, err
	}

	user.AvatarFileName = fileLocation

	updated_user, err := s.repository.Update(user)
	if err != nil {
		return updated_user, err
	}

	return updated_user, nil
}
