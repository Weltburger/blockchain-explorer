package main

import (
	"blockwatch.cc/tzgo/rpc"
	"context"
	"fmt"
	"log"
	"time"
)

/*type Block struct {
	Hash                        string    `json:"hash"`
	Baker                       string    `json:"baker"`
	Height                      int       `json:"height"`
	Cycle                       int       `json:"cycle"`
	IsCycleSnapshot             bool      `json:"is_cycle_snapshot"`
	Time                        time.Time `json:"time"`
	Solvetime                   int       `json:"solvetime"`
	Fitness                     int       `json:"fitness"`
	Priority                    int       `json:"priority"`
	NOps                        int       `json:"n_ops"`
	NEndorsement                int       `json:"n_endorsement"`
	Volume                      float64   `json:"volume"`
	Fee                         float64   `json:"fee"`
	GasUsed                     int       `json:"gas_used"`
	GasPrice                    float64   `json:"gas_price"`
	Protocol                    string    `json:"protocol"`
}*/

func main() {
	/*resp, err := http.Get("https://api.tzstats.com/explorer/block/1342853?meta=1&rights=1")
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}

	bodyString := string(body)
	fmt.Println(bodyString)*/

	//c, _ := rpc.NewClient("https://rpc.tzstats.com", nil)
	c, _ := rpc.NewClient("https://mainnet-tezos.giganode.io", nil)

	// create block header monitor
	mon := rpc.NewBlockHeaderMonitor()
	defer mon.Close()

	// all SDK functions take a context, here we just use a dummy
	ctx := context.TODO()

	// register the block monitor with our client
	if err := c.MonitorBlockHeader(ctx, mon); err != nil {
		log.Fatalln(err)
	}

	// wait for new block headers
	for {
		head, err := mon.Recv(ctx)
		if err != nil {
			log.Fatalln(err)
		}

		// do smth with the block header
		fmt.Printf("New block %s\n", head.Hash)

		time.Sleep(time.Second * 3)

		//block1 := new(Block)

		block, err := c.GetBlock(ctx, head.Hash)
		if err != nil {
			log.Fatal(err)
		}

		fmt.Println(block)
	}
}
