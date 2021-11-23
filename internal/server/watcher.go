package server

import (
	"explorer/models"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"
)

func (s *Server) CheckBlocks() {
	for {
		block, err := getData("head")
		if err != nil {
			fmt.Println(err)
		}

		err = saveData(s, *block)
		if err != nil {
			fmt.Println(err)
		}

		fmt.Println(*block)
		time.Sleep(time.Second * 30)
	}
}

func getData(index string) (*models.Block, error) {
	url := fmt.Sprintf("https://testnet-tezos.giganode.io/chains/main/blocks/%s", index)
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	block, _ := models.UnmarshalBlock(body)

	return &block, nil
}