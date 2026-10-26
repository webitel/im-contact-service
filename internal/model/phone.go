package model

import (
	"regexp"
	"strings"
)

var e164 = regexp.MustCompile(`^\+[1-9]\d{6,14}$`)

// NormalizePhone converts a raw phone number into E.164 (+380671234567).
// Spaces, dashes, dots and parentheses are dropped; a leading "00" is treated
// as the international prefix and a bare digit string gets a "+" prefix.
// Returns "" when the result is not a valid E.164 number.
func NormalizePhone(raw string) string {
	raw = strings.TrimSpace(raw)
	if raw == "" {
		return ""
	}

	var b strings.Builder
	for i, r := range raw {
		switch {
		case r >= '0' && r <= '9':
			b.WriteRune(r)
		case r == '+' && i == 0:
			b.WriteRune(r)
		case r == ' ' || r == '-' || r == '.' || r == '(' || r == ')':
		default:
			return ""
		}
	}

	phone := b.String()
	switch {
	case strings.HasPrefix(phone, "+"):
	case strings.HasPrefix(phone, "00"):
		phone = "+" + phone[2:]
	default:
		phone = "+" + phone
	}

	if !e164.MatchString(phone) {
		return ""
	}

	return phone
}
