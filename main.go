package main

import (
	"log"

	"github.com/shiimoo/go_game_server/pb"
)

func main() {
	if err := pb.LoadProto("./pb/testproto"); err != nil {
		log.Panicln(err)
	}
	bs, err := pb.Encode("test.AddFriendReq", pb.Message{
		"phone":   []string{"1", "111"},
		"keyword": "keyword",
		"tt": map[int32]string{
			1: "1",
			2: "2",
		},
		"f1": map[string]interface{}{
			"name": "name1",
			"age":  18,
		},
		// "fs": []pb.Message{
		// 	{
		// 		"name": "name1",
		// 		"age":  18,
		// 	},
		// 	{
		// 		"name": "name2",
		// 		"age":  180,
		// 	},
		// },
		// "fm": map[int32]pb.Message{
		// 	1: {
		// 		"name": "name1",
		// 		"age":  18,
		// 	},
		// },
	})
	if err != nil {
		log.Panicln(err)
	}
	log.Println("Encode", bs)
	msg, err := pb.Decode("test.AddFriendReq", bs)
	if err != nil {
		log.Panicln(err)
	}
	log.Println("Decode", msg)
}
