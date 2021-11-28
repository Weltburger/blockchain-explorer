package server

import (
	"explorer/models"
	"fmt"
	"github.com/spf13/viper"
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

		time.Sleep(time.Second * 30)
	}
}

func getData(index string) (*models.Block, error) {
	url := fmt.Sprintf(fmt.Sprintf("%s%s", viper.GetString("explorer.node"), index))
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