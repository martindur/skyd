package main

import (
	"log"
	"net"
	"net/http"
	"net/rpc"

	"github.com/google/uuid"
	"github.com/martindur/skyd/engine"
	"github.com/martindur/skyd/worknet"
)

type Server struct {
	game    *engine.Game
	players map[uuid.UUID]engine.Player
}

func (s *Server) Connect(p engine.Player, reply *worknet.PlayerStream) error {
	s.game.Players[p.ID] = &p

	for k, v := range s.game.Players {
		s.players[k] = *v
	}
	stream := worknet.PlayerStream{Players: s.players}
	*reply = stream

	return nil
}

func (s *Server) Sync(p engine.Player, reply *worknet.PlayerStream) error {
	s.game.Players[p.ID] = &p

	for k, v := range s.game.Players {
		s.players[k] = *v
	}
	stream := worknet.PlayerStream{Players: s.players}
	*reply = stream

	return nil
}

// func (s *Server) SyncPosition(stream worknet.ClientStream, reply *bool) error {
// 	s.game.Players[stream.ID].Pos = stream.Pos

// 	return nil
// }

func (s *Server) init() {
	player := s.game.SpawnPlayer()
	s.game.Players[player.ID] = player
	s.game.Local = player.ID
}

func RunServer(gs *Server) {
	rpc.Register(gs)
	rpc.HandleHTTP()

	log.Printf("Listening on port 8080...\n")

	l, err := net.Listen("tcp", ":8080")
	if err != nil {
		panic(err)
	}

	go http.Serve(l, nil)
}

func main() {
	game := engine.Game{Host: true}
	game.Players = make(map[uuid.UUID]*engine.Player)
	gs := Server{game: &game}
	gs.players = make(map[uuid.UUID]engine.Player)

	RunServer(&gs)

	gs.init()
	engine.GameLoop(&game)
}
