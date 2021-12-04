package watcher

import (
	"explorer/internal/server"
	"fmt"
	"net/http"
	"time"
)

func CheckBlocks(s *server.Server) {
	for {
		block, err := GetData(&http.Client{}, "head")
		if err != nil {
			fmt.Println(err)
		}

		if block.Hash == "" {
			continue
		}

		err = saveAllData(s, block)
		if err != nil {
			fmt.Println(err)
		}

		time.Sleep(time.Second * 30)
	}
}
