package worknet

import (
	"github.com/google/uuid"
	"github.com/martindur/skyd/engine"
)

type Action interface {
	perform() // Not sure if this is the ideal approach
}

type ClientStream struct {
	ID  uuid.UUID
	Pos engine.Pos
}

type ServerStream struct {
	PosMap map[uuid.UUID]engine.Pos
}

type PlayerStream struct {
	Players map[uuid.UUID]engine.Player
}

type Connection struct {
	token   string
	players []*engine.Player
}

type Message struct {
	Action
}
