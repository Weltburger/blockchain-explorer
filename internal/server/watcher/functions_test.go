package watcher

import (
	"bytes"
	"encoding/json"
	"errors"
	"explorer/models"
	"io/ioutil"
	"net/http"
	"reflect"
	"testing"
)

var (
	trans = []models.Transaction{
		{
			BlockHash: "hash",
			Hash: "ophash",
			Branch: "branch",
			Source: "src",
			Destination: "dst",
			Fee: "8415",
			Counter: "845",
			GasLimit: "6554253",
			Amount: "6845354",
			ConsumedMilligas: "68745355854",
			StorageSize: "645451",
			Signature:"sign",
		},
	}
	ops = [][]models.Operation{{},{},{},{models.Operation(struct {
		Protocol  string
		ChainID   string
		Hash      string
		Branch    string
		Contents  []models.Content
		Signature string
	}{
		Protocol:"protocol",
		ChainID:"id",
		Hash:"ophash",
		Branch:"branch",
		Contents: []models.Content{
			{
				Kind: "transaction",
				Source: "src",
				Fee: "8415",
				Counter: "845",
				GasLimit: "6554253",
				StorageLimit: "554",
				Amount: "6845354",
				Destination: "dst",
				Parameters: models.Parameters(struct {
					Entrypoint string
					Value      models.Value
				}{Entrypoint: "entr", Value: models.Value(struct{ Int string }{Int: "854"})}),
				Metadata: models.Metadata(struct {
					BalanceUpdates           []models.MetadataBalanceUpdate
					OperationResult          models.OperationResult
					InternalOperationResults []models.InternalOperationResult
				}{BalanceUpdates:nil, OperationResult: models.OperationResult(struct {
					Status           string
					Storage          []models.Storage
					BigMapDiff       []models.BigMapDiff
					ConsumedGas      string
					ConsumedMilligas string
					StorageSize      string
					LazyStorageDiff  []models.LazyStorageDiff
				}{
					Status:"",
					Storage:nil,
					BigMapDiff:nil,
					ConsumedGas:"",
					ConsumedMilligas:"68745355854",
					StorageSize:"645451",
					LazyStorageDiff:nil}), InternalOperationResults: nil})},
		}, Signature:"sign",
	})}}

	block0 = models.Block{
		Protocol:   "protocol",
		ChainID:    "chain id",
		Hash:       "hash",
		Header:     models.Header{},
		Metadata:   models.BlockMetadata{},
		Operations: nil,
	}
)

type testGetter struct {}

func (testGetter) Get(url string) (*http.Response, error) {
	if url == "error" {
		return &http.Response{
			Status:           "not found",
			StatusCode:       404,
			Body:             nil,
		}, errors.New("no data")
	} else if url == "empty" {
		resp := &http.Response{
			Body: ioutil.NopCloser(bytes.NewBufferString(string(""))),
		}
		return &http.Response{
			Status:           "200 OK",
			StatusCode:       200,
			Body:             resp.Body,
		}, nil
	} else {
		jsonBlock, _ := json.Marshal(block0)

		resp := &http.Response{
			Body: ioutil.NopCloser(bytes.NewBufferString(string(jsonBlock))),
		}
		return &http.Response{
			Status:           "200 OK",
			StatusCode:       200,
			Body:             resp.Body,
		}, nil
	}
}

func TestGetData(t *testing.T) {
	type args struct {
		getter Getter
		index  string
	}
	tests := []struct {
		name    string
		args    args
		want    models.Block
		wantErr bool
	}{
		{
			name: "Test 1",
			args: args{
				getter: testGetter{},
				index:  "error",
			},
			want: models.Block{},
			wantErr: true,
		},
		{
			name: "Test 2",
			args: args{
				getter: testGetter{},
				index:  "usual",
			},
			want: block0,
			wantErr: false,
		},
		{
			name: "Test 3",
			args: args{
				getter: testGetter{},
				index:  "empty",
			},
			want: models.Block{},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := GetData(tt.args.getter, tt.args.index)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetData() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetData() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGetTransactions(t *testing.T) {
	blk := block0
	blk.Operations = ops
	tests := []struct {
		name string
		block *models.Block
		want []models.Transaction
	}{
		{name: "Test 1", block: &block0, want: []models.Transaction{}},
		{name: "Test 2", block: &blk, want: trans},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := GetTransactions(tt.block); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetTransactions() = %v, want %v", got, tt.want)
			}
		})
	}
}