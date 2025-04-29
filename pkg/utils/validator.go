package utils

func ValidateEmail(email string) bool {
	// Simple email validation
	return len(email) > 3 && len(email) < 254
}

func ValidatePassword(password string) bool {
	return len(password) >= 8
}
