package timezone

import (
	"fmt"
	"strings"
)

// ParseOffset parses a string representing a time zone offset
// and returns the offset in seconds east of UTC.
//
// The function supports numerical offset formats (e.g., "+05", "-07:00", "+0530")
// as well as the "Z" (UTC) designator. It also supports "UTC" and "GMT"
// as names for the zero offset, and "UTC" or "GMT" prefixed
// offsets (e.g., "UTC+5", "GMT-07:00").
//
// This function does not parse time zone location names (e.g., "PST", "EST").
// Parsing location names is ambiguous without a full time zone database context.
// Use time.LoadLocation for location-based parsing.
//
// If the string does not match a supported numerical offset or "Z", "UTC", "GMT",
// it returns an error.
func ParseOffset(s string) (int, error) {
	if len(s) == 0 {
		return 0, fmt.Errorf("timezone: invalid time zone %q", s)
	}

	if s == "Z" || s == "UTC" || s == "GMT" {
		return 0, nil
	}

	// Handle "UTC" or "GMT" prefixes (e.g., "UTC+5", "GMT-07:00")
	if strings.HasPrefix(s, "UTC") {
		return parseNumericalOffset(s[3:])
	}
	if strings.HasPrefix(s, "GMT") {
		return parseNumericalOffset(s[3:])
	}

	// Try to parse the entire string as a numerical offset
	if offset, err := parseNumericalOffset(s); err == nil {
		return offset, nil
	}

	// If all parsing fails
	return 0, fmt.Errorf("timezone: invalid time zone %q", s)
}

// parseNumericalOffset parses numerical time zone offsets like
// "+05:00", "-07", or "+0530".
func parseNumericalOffset(s string) (int, error) {
	sOrig := s

	sign := 1
	switch s[0] {
	case '+':
		s = s[1:]
	case '-':
		sign = -1
		s = s[1:]
	default:
		return 0, fmt.Errorf("timezone: invalid time zone offset %q", sOrig)
	}

	if len(s) == 0 {
		return 0, fmt.Errorf("timezone: invalid time zone offset %q", sOrig)
	}

	var h, m int

	switch len(s) {
	case 1: // ±H
		if !isDigits(s) {
			return 0, fmt.Errorf("timezone: invalid time zone offset %q", sOrig)
		}
		h = parseDigits(s)
		m = 0

	case 2: // ±HH
		if !isDigits(s) {
			return 0, fmt.Errorf("timezone: invalid time zone offset %q", sOrig)
		}
		h = parseDigits(s)
		m = 0

	case 3: // ±HMM
		if !isDigits(s) {
			return 0, fmt.Errorf("timezone: invalid time zone offset %q", sOrig)
		}
		h = parseDigits(s[:1])
		m = parseDigits(s[1:])

	case 4: // ±HHMM
		if !isDigits(s) {
			return 0, fmt.Errorf("timezone: invalid time zone offset %q", sOrig)
		}
		h = parseDigits(s[:2])
		m = parseDigits(s[2:])

	case 5: // ±HH:MM
		if s[2] != ':' || !isDigits(s[:2]) || !isDigits(s[3:]) {
			return 0, fmt.Errorf("timezone: invalid time zone offset %q", sOrig)
		}
		h = parseDigits(s[:2])
		m = parseDigits(s[3:])

	default:
		return 0, fmt.Errorf("timezone: invalid time zone offset %q", sOrig)
	}

	if h > 14 || m > 59 {
		return 0, fmt.Errorf("timezone: invalid time zone offset %q", sOrig)
	}

	return sign * (h*3600 + m*60), nil
}

// isDigits checks if string s contains only ASCII digits.
func isDigits(s string) bool {
	for i := 0; i < len(s); i++ {
		if s[i] < '0' || s[i] > '9' {
			return false
		}
	}
	return true
}

// parseDigits converts ASCII digits s to integer.
// Assumes s contains only digits.
func parseDigits(s string) int {
	result := 0
	for i := 0; i < len(s); i++ {
		result = result*10 + int(s[i]-'0')
	}
	return result
}
