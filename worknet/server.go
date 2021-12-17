package worknet

import (
	"sync"

	"github.com/google/uuid"
	"github.com/martindur/skyd/engine"
)

type client struct {
	pid uuid.UUID
	id  uuid.UUID
}

type Server struct {
	game      *engine.Game
	mu        sync.RWMutex
	clients   map[uuid.UUID]*client
	positions map[uuid.UUID]engine.Pos
}

func (s *Server) Connect(pid string, con *Connection) error {
	//use pid for player ID

	//Return token, players
	return nil
}

// func (s *Server) broadcast(msg *Message) {
// 	s.mu.Lock()
// 	for id, client := range s.clients {
// 		//Send message via client
// 	}
// }

// func (s *Server) Stream()

// func (s *Server) Action() {

// }

// A function that performs an action for an entity:
// This is a main sync function called from a client

// A function that broadcasts actions performed by other entities
// This is a function for updating news from the server

// A function that holds the "Source of truth" for e.g. positions
// This is for potential lerping, and keeping the game in sync
// This is called on a certain interval from the client.
