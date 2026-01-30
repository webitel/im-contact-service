package model

import "github.com/google/uuid"

type PeerKind int

const (
	PeerContact PeerKind = iota
	PeerBot
)

type Peer struct {
	Kind PeerKind 
	Id uuid.UUID
} 