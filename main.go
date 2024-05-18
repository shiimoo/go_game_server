package main

import (
	"context"
	"runtime"
	"time"

	"github.com/shiimoo/go_game_server/blog"
	"github.com/shiimoo/go_game_server/net"
)

func main() {

	go func() {
		ticsker := time.NewTicker(time.Second * 1)
		for {
			<-ticsker.C
			runtime.GC()
		}
	}()

	ctx := context.Background()
	tl, err := net.NewTcpMgr(ctx, "", 8080)
	if err != nil {
		blog.Fatal(err)
	}
	tl.Start()

	<-ctx.Done()
}
