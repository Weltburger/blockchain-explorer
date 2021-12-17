package watcher

import (
	"explorer/internal/server"
	"log"
	"net/http"
	"strconv"
	"time"
)

func CheckBlocks(s *server.Server) {
	var prevPos = int64(0)
	for {
		block, err := GetData(&http.Client{}, "head", prevPos)
		if err != nil {
			log.Println(err)
			continue
		}

		if prevPos != 0 && (block.Metadata.LevelInfo.Level - prevPos) > 1 {
			for i := prevPos + 1; i <= block.Metadata.LevelInfo.Level; {
				index := strconv.FormatInt(i, 10)
				block, err := GetData(&http.Client{}, index, prevPos)
				if err != nil {
					log.Println(err)
					continue
				}

				err = saveAllData(s, block)
				if err != nil {
					log.Println(err)
					continue
				}

				prevPos = block.Metadata.LevelInfo.Level
				i++
			}
			time.Sleep(time.Second * 5)
			continue
		}

		err = saveAllData(s, block)
		if err != nil {
			log.Println(err)
			continue
		}
		prevPos = block.Metadata.LevelInfo.Level

		time.Sleep(time.Second * 15)
	}
}
