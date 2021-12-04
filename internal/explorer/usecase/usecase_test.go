package usecase

import (
	"context"
	"explorer/internal/apperrors"
	"explorer/internal/explorer"
	"explorer/models"
	"reflect"
	"testing"
)

type testBlockRepo struct {}
type testTransRepo struct {}

var (
	block1 = &models.Block{
		Protocol:   "protocol",
		ChainID:    "chainID",
		Hash:       "BlockWithSomeHash",
	}
	block2 = &models.Block{
		Protocol:   "protocol",
		ChainID:    "chainID",
		Hash:       "someHash",
		Metadata: models.BlockMetadata(struct {
			Baker       string
			LevelInfo   models.LevelInfo
			ConsumedGas string
		}{
			Baker:"", LevelInfo: models.LevelInfo(struct {
				Level         int64
				Cycle         int64
				CyclePosition int64
			}{Level: 55, Cycle:415, CyclePosition:852}), ConsumedGas:""}),
	}
)

func (testBlockRepo) GetBlockByHash(ctx context.Context, blk string) (*models.Block, error) {
	if blk == "noBlock" {
		return nil, apperrors.NewNotFound("clickhouse", "such block was")
	}

	return &models.Block{
		Protocol:   "protocol",
		ChainID:    "chainID",
		Hash:       blk,
	}, nil
}

func (testBlockRepo) GetBlockByLevel(ctx context.Context, blk int64) (*models.Block, error) {
	if blk < 0 || blk > 150 {
		return nil, apperrors.NewNotFound("clickhouse", "such block was")
	}
	return &models.Block{
		Protocol:   "protocol",
		ChainID:    "chainID",
		Hash:       "someHash",
		Metadata: models.BlockMetadata(struct {
			Baker       string
			LevelInfo   models.LevelInfo
			ConsumedGas string
		}{
			Baker:"", LevelInfo: models.LevelInfo(struct {
			Level         int64
			Cycle         int64
			CyclePosition int64
		}{Level: blk, Cycle:415, CyclePosition:852}), ConsumedGas:""}),
	}, nil
}

func (testBlockRepo) GetBlocks(ctx context.Context, offset int, limit int) ([]models.Block, error) {
	var blocks []models.Block
	if offset < 0 {
		offset = 0
	}
	if limit < 1 {
		limit = 1
	}
	max := 150
	count := max - offset
	if count < 1 {
		return []models.Block{}, nil
	} else if limit > count {
		for i := 1; i <= count; i++ {
			blocks = append(blocks, models.Block{})
		}
	} else {
		for i := 1; i <= limit; i++ {
			blocks = append(blocks, models.Block{})
		}
	}

	return blocks, nil
}

func (testTransRepo) GetTransactions(ctx context.Context,
	offset int, limit int, blk string, hash string, acc string) ([]models.Transaction, error) {
	var transactions []models.Transaction
	max := 150
	count := max - offset
	if count < 1 {
		return []models.Transaction{}, nil
	} else if limit > count {
		for i := 0; i < count; i++ {
			transactions = append(transactions, models.Transaction{})
		}
	} else {
		for i := 0; i < limit; i++ {
			transactions = append(transactions, models.Transaction{})
		}
	}

	if blk != "" || hash != "" || acc != "" {
		if blk == "noBlock" || hash == "noHash" || acc == "noAcc" {
			return []models.Transaction{}, nil
		} else {
			transactions = []models.Transaction{
				{
					BlockHash: blk,
					Hash: hash, Branch: "",
					Source: acc,
					Destination: acc,
					Fee: "",
					Counter: "",
					GasLimit: "",
					Amount: "",
					ConsumedMilligas: "",
					StorageSize: "",
					Signature: "",
				},
			}
		}
	}

	return transactions, nil
}

func TestBlockUseCase_GetBlock(t *testing.T) {
	type fields struct {
		blockRepo explorer.BlockRepo
	}
	type args struct {
		ctx context.Context
		blk string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *models.Block
		wantErr bool
	}{
		{
			name: "GetBlock Test 1",
			fields: struct{ blockRepo explorer.BlockRepo }{blockRepo: testBlockRepo{}},
			args: struct {
				ctx context.Context
				blk string
			}{ctx: context.Background(), blk: "BlockWithSomeHash"},
			want: block1,
			wantErr: false,
		},
		{
			name: "GetBlock Test 2",
			fields: struct{ blockRepo explorer.BlockRepo }{blockRepo: testBlockRepo{}},
			args: struct {
				ctx context.Context
				blk string
			}{ctx: context.Background(), blk: "55"},
			want: block2,
			wantErr: false,
		},
		{
			name: "GetBlock Test 3",
			fields: struct{ blockRepo explorer.BlockRepo }{blockRepo: testBlockRepo{}},
			args: struct {
				ctx context.Context
				blk string
			}{ctx: context.Background(), blk: "noBlock"},
			want: nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			b := &BlockUseCase{
				blockRepo: tt.fields.blockRepo,
			}
			got, err := b.GetBlock(tt.args.ctx, tt.args.blk)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetBlock() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetBlock() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestBlockUseCase_GetBlocks(t *testing.T) {
	type fields struct {
		blockRepo explorer.BlockRepo
	}
	type args struct {
		ctx    context.Context
		offset int
		limit  int
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    []models.Block
		wantErr bool
	}{
		{
			name: "GetBlocks Test 1",
			fields: struct{ blockRepo explorer.BlockRepo }{blockRepo: testBlockRepo{}},
			args: struct {
				ctx    context.Context
				offset int
				limit  int
			}{ctx: context.Background(), offset: 50, limit: 5},
			want: []models.Block{{},{},{},{},{}},
			wantErr: false,
		},
		{
			name: "GetBlocks Test 2",
			fields: struct{ blockRepo explorer.BlockRepo }{blockRepo: testBlockRepo{}},
			args: struct {
				ctx    context.Context
				offset int
				limit  int
			}{ctx: context.Background(), offset: 145, limit: 50},
			want: []models.Block{{},{},{},{},{}},
			wantErr: false,
		},
		{
			name: "GetBlocks Test 3",
			fields: struct{ blockRepo explorer.BlockRepo }{blockRepo: testBlockRepo{}},
			args: struct {
				ctx    context.Context
				offset int
				limit  int
			}{ctx: context.Background(), offset: 150, limit: 10},
			want: []models.Block{},
			wantErr: false,
		},
		{
			name: "GetBlocks Test 4",
			fields: struct{ blockRepo explorer.BlockRepo }{blockRepo: testBlockRepo{}},
			args: struct {
				ctx    context.Context
				offset int
				limit  int
			}{ctx: context.Background(), offset: -100, limit: -20},
			want: []models.Block{{}},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			b := &BlockUseCase{
				blockRepo: tt.fields.blockRepo,
			}
			got, err := b.GetBlocks(tt.args.ctx, tt.args.offset, tt.args.limit)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetBlocks() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetBlocks() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNewBlockUseCase(t *testing.T) {
	type args struct {
		blockRepo explorer.BlockRepo
	}
	tests := []struct {
		name string
		args args
		want *BlockUseCase
	}{
		{
			name: "NewBlockUseCase Test 1",
			args: struct{ blockRepo explorer.BlockRepo }{blockRepo: testBlockRepo{}},
			want: &BlockUseCase{blockRepo: testBlockRepo{}},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewBlockUseCase(tt.args.blockRepo); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewBlockUseCase() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNewTransUseCase(t *testing.T) {
	type args struct {
		transRepo explorer.TransRepo
	}
	tests := []struct {
		name string
		args args
		want *TransUseCase
	}{
		{
			name: "NewTransUseCase Test 1",
			args: struct{ transRepo explorer.TransRepo }{transRepo: testTransRepo{}},
			want: &TransUseCase{transRepo: testTransRepo{}},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewTransUseCase(tt.args.transRepo); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewTransUseCase() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestTransUseCase_GetTransactions(t1 *testing.T) {
	type fields struct {
		transRepo explorer.TransRepo
	}
	type args struct {
		ctx    context.Context
		offset int
		limit  int
		blk    string
		hash   string
		acc    string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    []models.Transaction
		wantErr bool
	}{
		{
			name: "GetTransactions Test 1",
			fields: struct{ transRepo explorer.TransRepo }{transRepo: testTransRepo{}},
			args: struct {
				ctx    context.Context
				offset int
				limit  int
				blk    string
				hash   string
				acc    string
			}{ctx: context.Background(), offset: 50, limit: 5, blk: "", hash: "", acc: ""},
			want: []models.Transaction{{},{},{},{},{}},
			wantErr: false,
		},
		{
			name: "GetTransactions Test 2",
			fields: struct{ transRepo explorer.TransRepo }{transRepo: testTransRepo{}},
			args: struct {
				ctx    context.Context
				offset int
				limit  int
				blk    string
				hash   string
				acc    string
			}{ctx: context.Background(), offset: 145, limit: 10, blk: "", hash: "", acc: ""},
			want: []models.Transaction{{},{},{},{},{}},
			wantErr: false,
		},
		{
			name: "GetTransactions Test 3",
			fields: struct{ transRepo explorer.TransRepo }{transRepo: testTransRepo{}},
			args: struct {
				ctx    context.Context
				offset int
				limit  int
				blk    string
				hash   string
				acc    string
			}{ctx: context.Background(), offset: 150, limit: 15, blk: "", hash: "", acc: ""},
			want: []models.Transaction{},
			wantErr: false,
		},
		{
			name: "GetTransactions Test 4",
			fields: struct{ transRepo explorer.TransRepo }{transRepo: testTransRepo{}},
			args: struct {
				ctx    context.Context
				offset int
				limit  int
				blk    string
				hash   string
				acc    string
			}{ctx: context.Background(), offset: 50, limit: 5, blk: "noBlock", hash: "someHash", acc: "someAcc"},
			want: []models.Transaction{},
			wantErr: false,
		},
		{
			name: "GetTransactions Test 5",
			fields: struct{ transRepo explorer.TransRepo }{transRepo: testTransRepo{}},
			args: struct {
				ctx    context.Context
				offset int
				limit  int
				blk    string
				hash   string
				acc    string
			}{ctx: context.Background(), offset: 50, limit: 5, blk: "someBlock", hash: "someHash", acc: "someAcc"},
			want: []models.Transaction{{BlockHash: "someBlock", Hash: "someHash", Source: "someAcc", Destination: "someAcc"}},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t1.Run(tt.name, func(t1 *testing.T) {
			t := &TransUseCase{
				transRepo: tt.fields.transRepo,
			}
			got, err := t.GetTransactions(tt.args.ctx, tt.args.offset, tt.args.limit, tt.args.blk, tt.args.hash, tt.args.acc)
			if (err != nil) != tt.wantErr {
				t1.Errorf("GetTransactions() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t1.Errorf("GetTransactions() got = %v, want %v", got, tt.want)
			}
		})
	}
}
