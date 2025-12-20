package model

type ContactType int

//go:generate stringer -type=ContactType -output=contact_type_string.go
const (
	Webitel ContactType = iota
	User
	Bot
)
