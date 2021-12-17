package main

import (
	"log"
	"net/rpc"

	"github.com/google/uuid"
	"github.com/martindur/skyd/engine"
	"github.com/martindur/skyd/worknet"
)

type Nothing bool

type GameClient struct {
	Client *rpc.Client
	game   *engine.Game
	stream map[uuid.UUID]engine.Pos
}

func (c *GameClient) Connect() {
	player := c.game.SpawnPlayer()
	c.game.Local = player.ID
	var stream worknet.PlayerStream
	// stream.Players = make(map[uuid.UUID]engine.Player)
	// c.game.Players[player.ID] = player
	// stream := worknet.ClientStream{Pos: player.Pos, ID: c.game.Local}
	var err error

	if c.Client == nil {
		c.Client, err = rpc.DialHTTP("tcp", "localhost:8080")
		if err != nil {
			panic(err)
		}
	}

	err_connect := c.Client.Call("Server.Connect", player, &stream)
	if err_connect != nil {
		log.Printf("Error connecting: %q", err)
	} else {
		log.Print(stream.Players)
		for k, v := range stream.Players {
			c.game.Players[k] = &v
		}
		log.Printf("Connected with ID: %s\n", player.ID.String())
	}

}

func (c *GameClient) SyncPosition() {
	var stream worknet.PlayerStream
	// stream.Players = make(map[uuid.UUID]engine.Player)
	err := c.Client.Call(
		"Server.Sync",
		c.game.Players[c.game.Local],
		&stream,
	)
	if err != nil {
		log.Printf("Sync error: %q", err)
	} else {
		for k, v := range stream.Players {
			c.game.Players[k] = &v
		}
	}
}

func (c *GameClient) sync(pos chan engine.Pos) {
	for {
		<-pos
		c.SyncPosition()
	}
}

func main() {
	game := engine.Game{}
	game.Move = make(chan engine.Pos)
	game.Players = make(map[uuid.UUID]*engine.Player)

	// client := GameClient{game: &game}
	// client.stream = make(map[uuid.UUID]engine.Pos, 4)

	// client.Connect()

	// go client.sync(game.Move)

	engine.GameLoop(&game)
}
