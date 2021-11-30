package server

import (
	"explorer/models"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
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

		time.Sleep(time.Second * 30)
	}
}

func getData(index string) (*models.Block, error) {
	url := fmt.Sprintf("%s%s", os.Getenv("NODE"), index)
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
