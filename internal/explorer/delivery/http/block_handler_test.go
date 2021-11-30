package http

import (
	"context"
	"errors"
	"explorer/internal/explorer"
	"explorer/models"
	"fmt"
	"github.com/labstack/echo/v4"
	"log"
	"net/http"
	"net/url"
	"strconv"
	"testing"
)

var block = models.Block{
	Protocol:   "123",
	ChainID:    "456",
	Hash:       "789",
	Header:     models.Header{},
	Metadata:   models.BlockMetadata{},
	Operations: nil,
}

type testBlockUseCase struct{}
type RWriter struct{}

func (RWriter) Header() http.Header {
	return http.Header{}
}

func (RWriter) Write(b []byte) (int, error) {
	log.Println(string(b))
	return len(b), nil
}

func (RWriter) WriteHeader(statusCode int) {
	log.Println(statusCode)
}

func (testBlockUseCase) GetBlock(ctx context.Context, blk string) (*models.Block, error) {
	if blk == "noBlock" {
		return nil, errors.New("there's no such block")
	}

	return &block, nil
}

func (testBlockUseCase) GetBlocks(ctx context.Context, offset int, limit int) ([]models.Block, error) {
	var blocks []models.Block
	max := 150
	count := max - offset
	if limit > count {
		for i := 1; i <= count; i++ {
			blocks = append(blocks, models.Block{Hash: strconv.Itoa(offset+i)})
		}
	} else if count < 1 {
		return []models.Block{}, nil
	} else {
		for i := 1; i <= limit; i++ {
			blocks = append(blocks, models.Block{Hash: strconv.Itoa(offset+i)})
		}
	}

	return blocks, nil
}

func TestBlockHandler_GetBlock(t *testing.T) {
	echoCtx1 := echo.New().NewContext(&http.Request{
		URL: &url.URL{},
	}, RWriter{})
	echoCtx1.SetParamNames("block")
	echoCtx1.SetParamValues("noBlock")

	echoCtx2 := echo.New().NewContext(&http.Request{
		URL: &url.URL{},
	}, RWriter{})
	echoCtx2.SetParamNames("block")
	echoCtx2.SetParamValues("Block")

	type fields struct {
		blockUseCase explorer.BlockUseCase
	}
	type args struct {
		c echo.Context
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			name:    "Test 1",
			fields:  struct{blockUseCase explorer.BlockUseCase }{blockUseCase: &testBlockUseCase{}},
			args:    struct{c echo.Context }{c: echoCtx1},
			wantErr: true,
		},
		{
			name:    "Test 2",
			fields:  struct{blockUseCase explorer.BlockUseCase }{blockUseCase: &testBlockUseCase{}},
			args:    struct{c echo.Context }{c: echoCtx2},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h := &BlockHandler{
				blockUseCase: tt.fields.blockUseCase,
			}
			if err := h.GetBlock(tt.args.c); (err != nil) != tt.wantErr {
				t.Errorf("GetBlock() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestBlockHandler_GetBlocks(t *testing.T) {
	URL1 := &url.URL{}
	URL1.RawQuery = fmt.Sprintf("limit=%s&offset=%s", "5", "50")
	echoCtx1 := echo.New().NewContext(&http.Request{
		URL: URL1,
	}, RWriter{})

	URL2 := &url.URL{}
	URL2.RawQuery = fmt.Sprintf("limit=%s&offset=%s", "50", "145")
	echoCtx2 := echo.New().NewContext(&http.Request{
		URL: URL2,
	}, RWriter{})

	URL3 := &url.URL{}
	URL3.RawQuery = fmt.Sprintf("limit=%s&offset=%s", "10", "150")
	echoCtx3 := echo.New().NewContext(&http.Request{
		URL: URL3,
	}, RWriter{})

	type fields struct {
		blockUseCase explorer.BlockUseCase
	}
	type args struct {
		c echo.Context
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			name: "Test 1",
			fields: struct{blockUseCase explorer.BlockUseCase }{blockUseCase: &testBlockUseCase{}},
			args: struct{c echo.Context}{c: echoCtx1},
			wantErr: false,
		},
		{
			name: "Test 2",
			fields: struct{blockUseCase explorer.BlockUseCase }{blockUseCase: &testBlockUseCase{}},
			args: struct{c echo.Context}{c: echoCtx2},
			wantErr: false,
		},
		{
			name: "Test 3",
			fields: struct{blockUseCase explorer.BlockUseCase }{blockUseCase: &testBlockUseCase{}},
			args: struct{c echo.Context}{c: echoCtx3},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h := &BlockHandler{
				blockUseCase: tt.fields.blockUseCase,
			}
			if err := h.GetBlocks(tt.args.c); (err != nil) != tt.wantErr {
				t.Errorf("GetBlocks() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestTransHandler_GetTransactions(t *testing.T) {
	type fields struct {
		transUseCase explorer.TransUseCase
	}
	type args struct {
		c echo.Context
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			name: ,
			fields: ,
			args: ,
			wantErr: ,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h := &TransHandler{
				transUseCase: tt.fields.transUseCase,
			}
			if err := h.GetTransactions(tt.args.c); (err != nil) != tt.wantErr {
				t.Errorf("GetTransactions() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

