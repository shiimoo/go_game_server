package main

import (
	"log"

	"github.com/shiimoo/go_game_server/pb"
)

func main() {
	if err := pb.LoadProto("./proto/testproto"); err != nil {
		log.Panicln(err)
	}
}
