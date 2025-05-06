package postgres

import (
	"strings"
	"unicode"
)

// ExtractPotentialNames splits an email address into potential first and last names
func ExtractPotentialNames(email string) (firstName, lastName string) {
	// Remove everything after @
	localPart := strings.Split(email, "@")[0]

	// Common separators in email local parts
	separators := []string{".", "_", "-"}

	// Try to split by each separator
	for _, sep := range separators {
		if strings.Contains(localPart, sep) {
			parts := strings.Split(localPart, sep)
			if len(parts) >= 2 {
				firstName = formatName(parts[0])
				lastName = formatName(parts[1])
				return
			}
		}
	}

	// If no separator found, try to split by case (camelCase)
	firstName, lastName = splitByCase(localPart)
	if firstName != "" && lastName != "" {
		return
	}

	// If all else fails, use the entire local part as first name
	firstName = formatName(localPart)
	return
}

// formatName capitalizes the first letter of a name
func formatName(name string) string {
	if len(name) == 0 {
		return ""
	}
	return strings.ToUpper(string(name[0])) + strings.ToLower(name[1:])
}

// splitByCase attempts to split a string by uppercase letters
func splitByCase(s string) (firstName, lastName string) {
	var parts []string
	lastCut := 0

	for i, r := range s {
		if i > 0 && unicode.IsUpper(r) {
			parts = append(parts, s[lastCut:i])
			lastCut = i
		}
	}

	if lastCut > 0 {
		parts = append(parts, s[lastCut:])
	}

	if len(parts) >= 2 {
		firstName = formatName(parts[0])
		lastName = formatName(parts[1])
	}

	return
}

