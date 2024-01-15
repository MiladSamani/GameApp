package userservice

import (
	"fmt"
	"gameAppProject/entity"
	"gameAppProject/pkg/phonenumber"
)

// Repository defines methods for user data storage.
type Repository interface {
	// IsPhoneNumberUnique checks if a phone number is unique.
	IsPhoneNumberUnique(phoneNumber string) (bool, error)
	// Register saves a new user in storage.
	Register(u entity.User) (entity.User, error)
}

// Service provides user-related functionality.
type Service struct {
	repo Repository
}

// New creates a new Service instance with the provided repository.
func New(repo Repository) Service {
	return Service{repo: repo}
}

// RegisterRequest represents the request structure for user registration.
type RegisterRequest struct {
	Name        string
	PhoneNumber string
}

// RegisterResponse :: Response for registration users :: get phone and number and pass to register function
type RegisterResponse struct {
	User entity.User
}

// Register handles user registration.
// It takes a RegisterRequest containing user information and performs the following steps:
// 1. Validates the phone number using the IsValid function from the phonenumber package.
// 2. Checks the uniqueness of the phone number by calling the IsPhoneNumberUnique method on the repository.
// 3. Validates the length of the user's name.
// 4. Creates a new user entity and calls the repository's Register method to save the user.
// 5. Returns a RegisterResponse containing the created user or an error if any validation or registration step fails.
func (s Service) Register(req RegisterRequest) (RegisterResponse, error) {
	// TODO: Verify phone number by verification code

	// Step 1: Validate phone number
	if !phonenumber.IsValid(req.PhoneNumber) {
		return RegisterResponse{}, fmt.Errorf("phone number is not valid")
	}

	// Step 2: Check uniqueness of phone number
	if isUnique, err := s.repo.IsPhoneNumberUnique(req.PhoneNumber); err != nil || !isUnique {
		if err != nil {
			return RegisterResponse{}, fmt.Errorf("unexpected error %w", err)
		}
		return RegisterResponse{}, fmt.Errorf("phone number is not unique")
	}

	// Step 3: Validate name
	if len(req.Name) < 3 {
		return RegisterResponse{}, fmt.Errorf("name length should be greater than 3")
	}

	// Step 4: Create new user in storage
	user := entity.User{
		ID:          0,
		Name:        req.Name,
		PhoneNumber: req.PhoneNumber,
	}
	createdUser, err := s.repo.Register(user)
	if err != nil {
		return RegisterResponse{}, fmt.Errorf("unexpected error %w", err)
	}

	// Step 5: Return created user
	return RegisterResponse{User: createdUser}, nil
}
