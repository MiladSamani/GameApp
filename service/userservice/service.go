package userservice

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"gameAppProject/entity"
	"gameAppProject/pkg/phonenumber"
)

// Repository defines methods for user data storage.
type Repository interface {
	// IsPhoneNumberUnique checks if a phone number is unique in the storage.
	// It takes a phone number as a parameter and returns a boolean indicating
	// whether the phone number is unique. If an error occurs during the check,
	// it is returned as the second value.
	IsPhoneNumberUnique(phoneNumber string) (bool, error)

	// Register saves a new user in storage.
	// It takes a user entity as a parameter and returns the created user entity.
	// If an error occurs during the registration process, it is returned as the second value.
	Register(u entity.User) (entity.User, error)

	// GetUserByPhoneNumber retrieves a user by their phone number from the storage.
	// It takes a phone number as a parameter and returns the corresponding user entity.
	// If the user is not found, it returns nil as the first value and an error
	// indicating the absence of the user as the second value.
	// If an error occurs during the retrieval process, it is returned as the second value.
	GetUserByPhoneNumber(phoneNumber string) (entity.User, bool, error)
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
	Name        string `json:"name"`
	PhoneNumber string `json:"phone_number"`
	Password    string `json:"password"`
}

// RegisterResponse :: Response for registration users :: get phone and number and pass to register function
type RegisterResponse struct {
	User entity.User
}

// Register handles the user registration process.
// It takes a RegisterRequest containing user information and performs the following steps:
// 1. Verify the validity of the phone number using the IsValid function from the phonenumber package.
// 2. Check the uniqueness of the phone number by calling the IsPhoneNumberUnique method on the repository.
// 3. Validate the length of the user's name and ensure it is greater than 3 characters.
// 4. Validate the length of the user's password and ensure it is greater than 8 characters.
// 5. Create a new user entity with the provided information.
// 6. Call the repository's Register method to save the user to storage.
// 7. Return a RegisterResponse containing the created user or an error if any validation or registration step fails.
func (s Service) Register(req RegisterRequest) (RegisterResponse, error) {
	// TODO: Implement phone number verification using a verification code

	// Step 1: Validate the phone number
	if !phonenumber.IsValid(req.PhoneNumber) {
		return RegisterResponse{}, fmt.Errorf("invalid phone number format")
	}

	// Step 2: Check the uniqueness of the phone number
	if isUnique, err := s.repo.IsPhoneNumberUnique(req.PhoneNumber); err != nil || !isUnique {
		if err != nil {
			return RegisterResponse{}, fmt.Errorf("unexpected error during phone number uniqueness check: %w", err)
		}
		return RegisterResponse{}, fmt.Errorf("phone number is already registered")
	}

	// Step 3: Validate the length of the user's name
	if len(req.Name) < 3 {
		return RegisterResponse{}, fmt.Errorf("name should be at least 3 characters long")
	}

	// Step 4: Validate the length of the user's password
	if len(req.Password) < 8 {
		return RegisterResponse{}, fmt.Errorf("password should be at least 8 characters long")
	}

	//ToDo : replace md5 with bcrypt

	// Step 5: Create a new user entity
	user := entity.User{
		ID:          0,
		Name:        req.Name,
		PhoneNumber: req.PhoneNumber,
		Password:    getMD5Hash(req.Password),
	}

	// Step 6: Call the repository's Register method to save the user
	createdUser, err := s.repo.Register(user)
	if err != nil {
		return RegisterResponse{}, fmt.Errorf("unexpected error during user registration: %w", err)
	}

	// Step 7: Return the created user
	return RegisterResponse{User: createdUser}, nil
}

// LoginRequest represents the data structure for user login requests.
// It contains the user's phone number and password for authentication.
type LoginRequest struct {
	PhoneNumber string `json:"phone_number"` // PhoneNumber is the user's phone number used for login.
	Password    string `json:"password"`     // Password is the user's password used for authentication.
}

// LoginResponse represents the data structure for user login responses.
// It currently does not contain any specific fields, but it can be extended
// with relevant information about the user's login status or additional details.
type LoginResponse struct {
}

func (s Service) Login(req LoginRequest) (LoginResponse, error) {
	//ToDo : it would be better to use two separate method for existence check and getUserByPhonenumber
	user, exist, err := s.repo.GetUserByPhoneNumber(req.PhoneNumber)
	if err != nil {
		return LoginResponse{}, fmt.Errorf("unexpected error %w", err)
	}

	if !exist {
		return LoginResponse{}, fmt.Errorf("username or password isn't correct")
	}
	if user.Password != getMD5Hash(req.Password) {
		return LoginResponse{}, fmt.Errorf("username or password isn't correct")
	}
	return LoginResponse{}, nil
}

// getMD5Hash computes the MD5 hash of the input text and returns the hexadecimal representation.
// It uses the md5 package to calculate the hash, and the resulting hash is a fixed-size byte array.
// The hash is then converted to a hexadecimal string using hex.EncodeToString.
// Parameters:
//   - text: The input text for which the MD5 hash is to be computed.
//
// Returns:
//   - A string representing the hexadecimal MD5 hash of the input text.
func getMD5Hash(text string) string {
	hash := md5.Sum([]byte(text))
	return hex.EncodeToString(hash[:])
}
