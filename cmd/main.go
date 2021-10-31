package main

import (
	"explorer/models"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"time"
)

func main() {
	for {
		resp, err := http.Get("https://mainnet-tezos.giganode.io/chains/main/blocks/head")
		if err != nil {
			log.Fatal(err)
		}

		body, _ := ioutil.ReadAll(resp.Body)
		resp.Body.Close()

		block, _ := models.UnmarshalBlock(body)
		fmt.Println(block.Operations)

		time.Sleep(time.Second * 15)
	}
}
