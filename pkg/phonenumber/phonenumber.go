package phonenumber

import (
	"strconv"
)

// IsValid checks the validity of a phone number based on certain criteria.
// It currently performs the following checks:
// 1. The length of the phone number must be 11 digits.
// 2. The first two digits must be "09".
// 3. The remaining digits must be numeric.
func IsValid(phoneNumber string) bool {
	// TODO: We can use a regular expression to support +98 pattern.

	// Step 1: Check the length of the phone number
	if len(phoneNumber) != 11 {
		return false
	}

	// Step 2: Check the first two digits
	if phoneNumber[0:2] != "09" {
		return false
	}

	// Step 3: Check if the remaining digits are numeric
	if _, err := strconv.Atoi(phoneNumber[2:]); err != nil {
		return false
	}

	// If all checks pass, consider the phone number valid
	return true
}
