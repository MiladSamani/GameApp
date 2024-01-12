package userservice

import (
	"fmt"
	"gameAppProject/entity"
	"gameAppProject/pkg/phonenumber"
)

type Repository interface {
	// IsPhoneNumberUnique :: For unique phone number check
	IsPhoneNumberUnique(phoneNumber string) (bool, error)
	// Register :: Save new user in storage
	Register(u entity.User) (entity.User, error)
}

func New(repo Repository) Service {
	return Service{repo: repo}
}

type Service struct {
	repo Repository
}

// RegisterRequest :: Request for registration users
type RegisterRequest struct {
	Name        string
	PhoneNumber string
}

// RegisterResponse :: Response for registration users :: get phone and number and pass to register function
type RegisterResponse struct {
	User entity.User
}

func (s Service) Register(req RegisterRequest) (RegisterResponse, error) {
	//TODO :: We should verify phone number by verification code

	// Validate phone number
	if !phonenumber.IsValid(req.PhoneNumber) {
		return RegisterResponse{}, fmt.Errorf("phone number is not valid")
	}
	// Check uniqueness phone number :: Just one number on DB :: shorthand if for err scope!
	if isUnique, err := s.repo.IsPhoneNumberUnique(req.PhoneNumber); err != nil || !isUnique {
		if err != nil {
			return RegisterResponse{}, fmt.Errorf("unexpexted error %w", err)
		}
		if !isUnique {
			return RegisterResponse{}, fmt.Errorf("phone number is not unique")
		}
	}

	// Validate name
	if len(req.Name) < 3 {
		return RegisterResponse{}, fmt.Errorf("name length should be greater than 3")
	}
	// Create new user in storage
	user := entity.User{
		ID:          0,
		Name:        req.Name,
		PhoneNumber: req.PhoneNumber,
	}
	createdUser, err := s.repo.Register(user)
	if err != nil {
		return RegisterResponse{}, fmt.Errorf("unexpexted error %w", err)
	}

	// Return created user
	return RegisterResponse{User: createdUser}, nil
}
