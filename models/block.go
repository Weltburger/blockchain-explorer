package models

import (
	"bytes"
	"encoding/json"
	"errors"
)

func UnmarshalBlock(data []byte) (Block, error) {
	var r Block
	err := json.Unmarshal(data, &r)
	return r, err
}

func (r *Block) Marshal() ([]byte, error) {
	return json.Marshal(r)
}

type Block struct {
	Protocol   Protocol      `json:"protocol"`
	ChainID    ChainID       `json:"chain_id"`
	Hash       string        `json:"hash"`
	Header     Header        `json:"header"`
	Metadata   BlockMetadata `json:"metadata"`
	Operations [][]Operation `json:"operations"`
}

type Header struct {
	Level                     int64       `json:"level"`
	Proto                     int64       `json:"proto"`
	Predecessor               Predecessor `json:"predecessor"`
	Timestamp                 string      `json:"timestamp"`
	ValidationPass            int64       `json:"validation_pass"`
	OperationsHash            string      `json:"operations_hash"`
	Fitness                   []string    `json:"fitness"`
	Context                   string      `json:"context"`
	Priority                  int64       `json:"priority"`
	ProofOfWorkNonce          string      `json:"proof_of_work_nonce"`
	LiquidityBakingEscapeVote bool        `json:"liquidity_baking_escape_vote"`
	Signature                 string      `json:"signature"`
}

type BlockMetadata struct {
	Protocol                  Protocol                   `json:"protocol"`
	NextProtocol              Protocol                   `json:"next_protocol"`
	TestChainStatus           TestChainStatus            `json:"test_chain_status"`
	MaxOperationsTTL          int64                      `json:"max_operations_ttl"`
	MaxOperationDataLength    int64                      `json:"max_operation_data_length"`
	MaxBlockHeaderLength      int64                      `json:"max_block_header_length"`
	MaxOperationListLength    []MaxOperationListLength   `json:"max_operation_list_length"`
	Baker                     string                     `json:"baker"`
	LevelInfo                 LevelInfo                  `json:"level_info"`
	VotingPeriodInfo          VotingPeriodInfo           `json:"voting_period_info"`
	NonceHash                 interface{}                `json:"nonce_hash"`
	ConsumedGas               string                     `json:"consumed_gas"`
	Deactivated               []interface{}              `json:"deactivated"`
	BalanceUpdates            []MetadataBalanceUpdate    `json:"balance_updates"`
	LiquidityBakingEscapeEma  int64                      `json:"liquidity_baking_escape_ema"`
	ImplicitOperationsResults []ImplicitOperationsResult `json:"implicit_operations_results"`
}

type MetadataBalanceUpdate struct {
	Kind     BalanceUpdateKind `json:"kind"`
	Contract string           `json:"contract,omitempty"`
	Change   string            `json:"change"`
	Origin   Origin            `json:"origin"`
	Category Category         `json:"category,omitempty"`
	Delegate string           `json:"delegate,omitempty"`
	Cycle    int64            `json:"cycle,omitempty"`
}

type ImplicitOperationsResult struct {
	Kind             ImplicitOperationsResultKind            `json:"kind"`
	Storage          []ImplicitOperationsResultStorage       `json:"storage"`
	BalanceUpdates   []ImplicitOperationsResultBalanceUpdate `json:"balance_updates"`
	ConsumedGas      string                                  `json:"consumed_gas"`
	ConsumedMilligas string                                  `json:"consumed_milligas"`
	StorageSize      string                                  `json:"storage_size"`
}

type ImplicitOperationsResultBalanceUpdate struct {
	Kind     BalanceUpdateKind `json:"kind"`
	Contract string            `json:"contract"`
	Change   string            `json:"change"`
	Origin   Origin            `json:"origin"`
}

type ImplicitOperationsResultStorage struct {
	Int   string `json:"int,omitempty"`
	Bytes string `json:"bytes,omitempty"`
}

type LevelInfo struct {
	Level              int64 `json:"level"`
	LevelPosition      int64 `json:"level_position"`
	Cycle              int64 `json:"cycle"`
	CyclePosition      int64 `json:"cycle_position"`
	ExpectedCommitment bool  `json:"expected_commitment"`
}

type MaxOperationListLength struct {
	MaxSize int64  `json:"max_size"`
	MaxOp   int64 `json:"max_op,omitempty"`
}

type TestChainStatus struct {
	Status string `json:"status"`
}

type VotingPeriodInfo struct {
	VotingPeriod VotingPeriod `json:"voting_period"`
	Position     int64        `json:"position"`
	Remaining    int64        `json:"remaining"`
}

type VotingPeriod struct {
	Index         int64  `json:"index"`
	Kind          string `json:"kind"`
	StartPosition int64  `json:"start_position"`
}

type Operation struct {
	Protocol  Protocol    `json:"protocol"`
	ChainID   ChainID     `json:"chain_id"`
	Hash      string      `json:"hash"`
	Branch    Predecessor `json:"branch"`
	Contents  []Content   `json:"contents"`
	Signature string     `json:"signature,omitempty"`
}

type Content struct {
	Kind         ImplicitOperationsResultKind `json:"kind"`
	Endorsement  EndorsementClass            `json:"endorsement,omitempty"`
	Slot         int64                       `json:"slot,omitempty"`
	Metadata     ContentMetadata              `json:"metadata"`
	Source       string                      `json:"source,omitempty"`
	Fee          string                      `json:"fee,omitempty"`
	Counter      string                      `json:"counter,omitempty"`
	GasLimit     string                      `json:"gas_limit,omitempty"`
	StorageLimit string                      `json:"storage_limit,omitempty"`
	Amount       string                      `json:"amount,omitempty"`
	Destination  string                      `json:"destination,omitempty"`
	Parameters   ContentParameters           `json:"parameters,omitempty"`
}

type EndorsementClass struct {
	Branch     Predecessor `json:"branch"`
	Operations Operations  `json:"operations"`
	Signature  string      `json:"signature"`
}

type Operations struct {
	Kind  OperationsKind `json:"kind"`
	Level int64          `json:"level"`
}

type ContentMetadata struct {
	BalanceUpdates           []MetadataBalanceUpdate   `json:"balance_updates"`
	Delegate                 string                   `json:"delegate,omitempty"`
	Slots                    []int64                   `json:"slots,omitempty"`
	OperationResult          OperationResult          `json:"operation_result,omitempty"`
	InternalOperationResults []InternalOperationResult `json:"internal_operation_results,omitempty"`
}

type InternalOperationResult struct {
	Kind        ImplicitOperationsResultKind       `json:"kind"`
	Source      Source                             `json:"source"`
	Nonce       int64                              `json:"nonce"`
	Amount      string                             `json:"amount"`
	Destination string                             `json:"destination"`
	Parameters  InternalOperationResultParameters `json:"parameters,omitempty"`
	Result      Result                             `json:"result"`
}

type InternalOperationResultParameters struct {
	Entrypoint Entrypoint     `json:"entrypoint"`
	Value      IndecentValue `json:"value"`
}

type PurpleValue struct {
	Prim PurplePrim `json:"prim"`
	Args []Arg8     `json:"args"`
}

type PurpleArg struct {
	Prim PurplePrim        `json:"prim"`
	Args []StorageArgClass `json:"args"`
}

type StorageArgClass struct {
	Bytes  string     `json:"bytes,omitempty"`
	Prim   PurplePrim `json:"prim,omitempty"`
	Args   []FluffyArg `json:"args,omitempty"`
	Int    string     `json:"int,omitempty"`
	String string     `json:"string,omitempty"`
}

type FluffyArg struct {
	Int string `json:"int"`
}

type TentacledArg struct {
	Bytes string `json:"bytes"`
}

type FluffyValue struct {
	Prim PurplePrim  `json:"prim"`
	Args []StickyArg `json:"args"`
}

type StickyArg struct {
	Prim PurplePrim `json:"prim"`
	Args []Arg9     `json:"args"`
}

type IndigoArg struct {
	Prim string        `json:"prim"`
	Args []IndecentArg `json:"args"`
}

type IndecentArg struct {
	String string `json:"string,omitempty"`
	Bytes  string `json:"bytes,omitempty"`
}

type Result struct {
	Status              Status                                  `json:"status"`
	Storage             ResultStorage                          `json:"storage"`
	BigMapDiff          []ResultBigMapDiff                      `json:"big_map_diff,omitempty"`
	BalanceUpdates      []ImplicitOperationsResultBalanceUpdate `json:"balance_updates,omitempty"`
	ConsumedGas         string                                  `json:"consumed_gas"`
	ConsumedMilligas    string                                  `json:"consumed_milligas"`
	StorageSize         string                                 `json:"storage_size,omitempty"`
	PaidStorageSizeDiff string                                 `json:"paid_storage_size_diff,omitempty"`
	LazyStorageDiff     []ResultLazyStorageDiff                 `json:"lazy_storage_diff,omitempty"`
}

type ResultBigMapDiff struct {
	Action  Action         `json:"action"`
	BigMap  string         `json:"big_map"`
	KeyHash string         `json:"key_hash"`
	Key     KeyElement     `json:"key"`
	Value   TentacledValue `json:"value"`
}

type KeyElement struct {
	Prim  PurplePrim                       `json:"prim,omitempty"`
	Args  []ImplicitOperationsResultStorage `json:"args,omitempty"`
	Int   string                           `json:"int,omitempty"`
	Bytes string                           `json:"bytes,omitempty"`
}

type TentacledValue struct {
	Int  string     `json:"int,omitempty"`
	Prim PurplePrim `json:"prim,omitempty"`
	Args []Arg10     `json:"args,omitempty"`
}

type ResultLazyStorageDiff struct {
	Kind LazyStorageDiffKind `json:"kind"`
	ID   string              `json:"id"`
	Diff PurpleDiff          `json:"diff"`
}

type PurpleDiff struct {
	Action  Action         `json:"action"`
	Updates []PurpleUpdate `json:"updates"`
}

type PurpleUpdate struct {
	KeyHash string         `json:"key_hash"`
	Key     KeyElement     `json:"key"`
	Value   TentacledValue `json:"value"`
}

type StorageArg struct {
	Prim FluffyPrim       `json:"prim,omitempty"`
	Args []StorageArgClass `json:"args,omitempty"`
	Int  string           `json:"int,omitempty"`
}

type PurpleStorage struct {
	Prim PurplePrim `json:"prim"`
	Args []Arg11    `json:"args"`
}

type HilariousArg struct {
	Prim PurplePrim    `json:"prim,omitempty"`
	Args []AmbitiousArg `json:"args,omitempty"`
	Int  string        `json:"int,omitempty"`
}

type AmbitiousArg struct {
	Prim string      `json:"prim,omitempty"`
	Args []CunningArg `json:"args,omitempty"`
	Int  string      `json:"int,omitempty"`
}

type CunningArg struct {
	Bytes string     `json:"bytes,omitempty"`
	Prim  FluffyPrim `json:"prim,omitempty"`
}

type OperationResult struct {
	Status              Status                                  `json:"status"`
	Storage             OperationResultStorage                 `json:"storage"`
	BigMapDiff          []OperationResultBigMapDiff             `json:"big_map_diff,omitempty"`
	ConsumedGas         string                                  `json:"consumed_gas"`
	ConsumedMilligas    string                                  `json:"consumed_milligas"`
	StorageSize         string                                 `json:"storage_size,omitempty"`
	LazyStorageDiff     []OperationResultLazyStorageDiff        `json:"lazy_storage_diff,omitempty"`
	BalanceUpdates      []ImplicitOperationsResultBalanceUpdate `json:"balance_updates,omitempty"`
	PaidStorageSizeDiff string                                 `json:"paid_storage_size_diff,omitempty"`
}

type OperationResultBigMapDiff struct {
	Action  Action                `json:"action"`
	BigMap  string                `json:"big_map"`
	KeyHash string                `json:"key_hash"`
	Key     Key                   `json:"key"`
	Value   BigMapDiffValueUnion `json:"value"`
}

type Key struct {
	Int   string      `json:"int,omitempty"`
	Prim  FluffyPrim  `json:"prim,omitempty"`
	Args  []KeyElement `json:"args,omitempty"`
	Bytes Bytes       `json:"bytes,omitempty"`
}

type OperationResultLazyStorageDiff struct {
	Kind LazyStorageDiffKind `json:"kind"`
	ID   string              `json:"id"`
	Diff FluffyDiff          `json:"diff"`
}

type FluffyDiff struct {
	Action  Action         `json:"action"`
	Updates []FluffyUpdate `json:"updates"`
}

type FluffyUpdate struct {
	KeyHash string                `json:"key_hash"`
	Key     Key                   `json:"key"`
	Value   BigMapDiffValueUnion `json:"value"`
}

type FluffyStorage struct {
	Prim PurplePrim `json:"prim"`
	Args []Arg12    `json:"args"`
}

type MagentaArg struct {
	Prim  PurplePrim `json:"prim,omitempty"`
	Args  []Arg13     `json:"args,omitempty"`
	Int   string     `json:"int,omitempty"`
	Bytes string     `json:"bytes,omitempty"`
}

type FriskyArg struct {
	Prim PurplePrim      `json:"prim,omitempty"`
	Args []MischievousArg `json:"args,omitempty"`
	Int  string          `json:"int,omitempty"`
}

type MischievousArg struct {
	Prim PurplePrim        `json:"prim,omitempty"`
	Args []BraggadociousArg `json:"args,omitempty"`
	Int  string            `json:"int,omitempty"`
}

type BraggadociousArg struct {
	Bytes string        `json:"bytes,omitempty"`
	Prim  string        `json:"prim,omitempty"`
	Args  []TentacledArg `json:"args,omitempty"`
	Int   string        `json:"int,omitempty"`
}

type Arg1 struct {
	Prim PurplePrim `json:"prim,omitempty"`
	Args []Arg2      `json:"args,omitempty"`
	Int  string     `json:"int,omitempty"`
}

type Arg2 struct {
	Bytes string `json:"bytes,omitempty"`
	Prim  string `json:"prim,omitempty"`
	Args  []Arg3  `json:"args,omitempty"`
	Int   string `json:"int,omitempty"`
}

type Arg3 struct {
	Prim PurplePrim  `json:"prim"`
	Args []FluffyArg `json:"args"`
}

type ContentParameters struct {
	Entrypoint string          `json:"entrypoint"`
	Value      HilariousValue `json:"value"`
}

type StickyValue struct {
	Prim string  `json:"prim"`
	Args []Arg16 `json:"args"`
}

type Arg4 struct {
	String string     `json:"string,omitempty"`
	Prim   PurplePrim `json:"prim,omitempty"`
	Args   []Arg5      `json:"args,omitempty"`
}

type Arg5 struct {
	String string     `json:"string,omitempty"`
	Prim   PurplePrim `json:"prim,omitempty"`
	Args   []Arg6      `json:"args,omitempty"`
}

type Arg6 struct {
	String Source `json:"string,omitempty"`
	Int    string `json:"int,omitempty"`
}

type IndigoValue struct {
	Prim PurplePrim `json:"prim,omitempty"`
	Args []Arg7      `json:"args,omitempty"`
	Int  string     `json:"int,omitempty"`
}

type Arg7 struct {
	String string           `json:"string,omitempty"`
	Int    string           `json:"int,omitempty"`
	Prim   PurplePrim       `json:"prim,omitempty"`
	Args   []StorageArgClass `json:"args,omitempty"`
}

type ChainID string
type Predecessor string
type Category string


type BalanceUpdateKind string

type Origin string

type ImplicitOperationsResultKind string

type Protocol string

type OperationsKind string

type Entrypoint string

type PurplePrim string

type Action string

type LazyStorageDiffKind string

type Status string

type FluffyPrim string
type Source string
type Bytes string

type IndecentValue struct {
	FluffyValue      *FluffyValue
	PurpleValueArray []PurpleValue
}

func (x *IndecentValue) UnmarshalJSON(data []byte) error {
	x.PurpleValueArray = nil
	x.FluffyValue = nil
	var c FluffyValue
	object, err := unmarshalUnion(data, nil, nil, nil, nil, true, &x.PurpleValueArray, true, &c, false, nil, false, nil, false)
	if err != nil {
		return err
	}
	if object {
		x.FluffyValue = &c
	}
	return nil
}

func (x *IndecentValue) MarshalJSON() ([]byte, error) {
	return marshalUnion(nil, nil, nil, nil, x.PurpleValueArray != nil, x.PurpleValueArray, x.FluffyValue != nil, x.FluffyValue, false, nil, false, nil, false)
}

type Arg8 struct {
	PurpleArgArray []PurpleArg
	TentacledArg   *TentacledArg
}

func (x *Arg8) UnmarshalJSON(data []byte) error {
	x.PurpleArgArray = nil
	x.TentacledArg = nil
	var c TentacledArg
	object, err := unmarshalUnion(data, nil, nil, nil, nil, true, &x.PurpleArgArray, true, &c, false, nil, false, nil, false)
	if err != nil {
		return err
	}
	if object {
		x.TentacledArg = &c
	}
	return nil
}

func (x *Arg8) MarshalJSON() ([]byte, error) {
	return marshalUnion(nil, nil, nil, nil, x.PurpleArgArray != nil, x.PurpleArgArray, x.TentacledArg != nil, x.TentacledArg, false, nil, false, nil, false)
}

type Arg9 struct {
	ImplicitOperationsResultStorage *ImplicitOperationsResultStorage
	IndigoArgArray                  []IndigoArg
}

func (x *Arg9) UnmarshalJSON(data []byte) error {
	x.IndigoArgArray = nil
	x.ImplicitOperationsResultStorage = nil
	var c ImplicitOperationsResultStorage
	object, err := unmarshalUnion(data, nil, nil, nil, nil, true, &x.IndigoArgArray, true, &c, false, nil, false, nil, false)
	if err != nil {
		return err
	}
	if object {
		x.ImplicitOperationsResultStorage = &c
	}
	return nil
}

func (x *Arg9) MarshalJSON() ([]byte, error) {
	return marshalUnion(nil, nil, nil, nil, x.IndigoArgArray != nil, x.IndigoArgArray, x.ImplicitOperationsResultStorage != nil, x.ImplicitOperationsResultStorage, false, nil, false, nil, false)
}

type Arg10 struct {
	FluffyArg      *FluffyArg
	IndigoArgArray []IndigoArg
}

func (x *Arg10) UnmarshalJSON(data []byte) error {
	x.IndigoArgArray = nil
	x.FluffyArg = nil
	var c FluffyArg
	object, err := unmarshalUnion(data, nil, nil, nil, nil, true, &x.IndigoArgArray, true, &c, false, nil, false, nil, false)
	if err != nil {
		return err
	}
	if object {
		x.FluffyArg = &c
	}
	return nil
}

func (x *Arg10) MarshalJSON() ([]byte, error) {
	return marshalUnion(nil, nil, nil, nil, x.IndigoArgArray != nil, x.IndigoArgArray, x.FluffyArg != nil, x.FluffyArg, false, nil, false, nil, false)
}

type ResultStorage struct {
	PurpleStorage   *PurpleStorage
	StorageArgArray []StorageArg
}

func (x *ResultStorage) UnmarshalJSON(data []byte) error {
	x.StorageArgArray = nil
	x.PurpleStorage = nil
	var c PurpleStorage
	object, err := unmarshalUnion(data, nil, nil, nil, nil, true, &x.StorageArgArray, true, &c, false, nil, false, nil, false)
	if err != nil {
		return err
	}
	if object {
		x.PurpleStorage = &c
	}
	return nil
}

func (x *ResultStorage) MarshalJSON() ([]byte, error) {
	return marshalUnion(nil, nil, nil, nil, x.StorageArgArray != nil, x.StorageArgArray, x.PurpleStorage != nil, x.PurpleStorage, false, nil, false, nil, false)
}

type Arg11 struct {
	FluffyArg         *FluffyArg
	HilariousArgArray []HilariousArg
}

func (x *Arg11) UnmarshalJSON(data []byte) error {
	x.HilariousArgArray = nil
	x.FluffyArg = nil
	var c FluffyArg
	object, err := unmarshalUnion(data, nil, nil, nil, nil, true, &x.HilariousArgArray, true, &c, false, nil, false, nil, false)
	if err != nil {
		return err
	}
	if object {
		x.FluffyArg = &c
	}
	return nil
}

func (x *Arg11) MarshalJSON() ([]byte, error) {
	return marshalUnion(nil, nil, nil, nil, x.HilariousArgArray != nil, x.HilariousArgArray, x.FluffyArg != nil, x.FluffyArg, false, nil, false, nil, false)
}

type BigMapDiffValueUnion struct {
	KeyArray   []Key
	KeyElement *KeyElement
}

func (x *BigMapDiffValueUnion) UnmarshalJSON(data []byte) error {
	x.KeyArray = nil
	x.KeyElement = nil
	var c KeyElement
	object, err := unmarshalUnion(data, nil, nil, nil, nil, true, &x.KeyArray, true, &c, false, nil, false, nil, false)
	if err != nil {
		return err
	}
	if object {
		x.KeyElement = &c
	}
	return nil
}

func (x *BigMapDiffValueUnion) MarshalJSON() ([]byte, error) {
	return marshalUnion(nil, nil, nil, nil, x.KeyArray != nil, x.KeyArray, x.KeyElement != nil, x.KeyElement, false, nil, false, nil, false)
}

type OperationResultStorage struct {
	FluffyStorage *FluffyStorage
	UnionArray    []StorageStorageUnion
}

func (x *OperationResultStorage) UnmarshalJSON(data []byte) error {
	x.UnionArray = nil
	x.FluffyStorage = nil
	var c FluffyStorage
	object, err := unmarshalUnion(data, nil, nil, nil, nil, true, &x.UnionArray, true, &c, false, nil, false, nil, false)
	if err != nil {
		return err
	}
	if object {
		x.FluffyStorage = &c
	}
	return nil
}

func (x *OperationResultStorage) MarshalJSON() ([]byte, error) {
	return marshalUnion(nil, nil, nil, nil, x.UnionArray != nil, x.UnionArray, x.FluffyStorage != nil, x.FluffyStorage, false, nil, false, nil, false)
}

type StorageStorageUnion struct {
	Key      *Key
	KeyArray []Key
}

func (x *StorageStorageUnion) UnmarshalJSON(data []byte) error {
	x.KeyArray = nil
	x.Key = nil
	var c Key
	object, err := unmarshalUnion(data, nil, nil, nil, nil, true, &x.KeyArray, true, &c, false, nil, false, nil, false)
	if err != nil {
		return err
	}
	if object {
		x.Key = &c
	}
	return nil
}

func (x *StorageStorageUnion) MarshalJSON() ([]byte, error) {
	return marshalUnion(nil, nil, nil, nil, x.KeyArray != nil, x.KeyArray, x.Key != nil, x.Key, false, nil, false, nil, false)
}

type Arg12 struct {
	FluffyArgArray []FluffyArg
	MagentaArg     *MagentaArg
}

func (x *Arg12) UnmarshalJSON(data []byte) error {
	x.FluffyArgArray = nil
	x.MagentaArg = nil
	var c MagentaArg
	object, err := unmarshalUnion(data, nil, nil, nil, nil, true, &x.FluffyArgArray, true, &c, false, nil, false, nil, false)
	if err != nil {
		return err
	}
	if object {
		x.MagentaArg = &c
	}
	return nil
}

func (x *Arg12) MarshalJSON() ([]byte, error) {
	return marshalUnion(nil, nil, nil, nil, x.FluffyArgArray != nil, x.FluffyArgArray, x.MagentaArg != nil, x.MagentaArg, false, nil, false, nil, false)
}

type Arg13 struct {
	Arg1       *Arg1
	UnionArray []Arg14
}

func (x *Arg13) UnmarshalJSON(data []byte) error {
	x.UnionArray = nil
	x.Arg1 = nil
	var c Arg1
	object, err := unmarshalUnion(data, nil, nil, nil, nil, true, &x.UnionArray, true, &c, false, nil, false, nil, false)
	if err != nil {
		return err
	}
	if object {
		x.Arg1 = &c
	}
	return nil
}

func (x *Arg13) MarshalJSON() ([]byte, error) {
	return marshalUnion(nil, nil, nil, nil, x.UnionArray != nil, x.UnionArray, x.Arg1 != nil, x.Arg1, false, nil, false, nil, false)
}

type Arg14 struct {
	StorageArg *StorageArg
	UnionArray []Arg15
}

func (x *Arg14) UnmarshalJSON(data []byte) error {
	x.UnionArray = nil
	x.StorageArg = nil
	var c StorageArg
	object, err := unmarshalUnion(data, nil, nil, nil, nil, true, &x.UnionArray, true, &c, false, nil, false, nil, false)
	if err != nil {
		return err
	}
	if object {
		x.StorageArg = &c
	}
	return nil
}

func (x *Arg14) MarshalJSON() ([]byte, error) {
	return marshalUnion(nil, nil, nil, nil, x.UnionArray != nil, x.UnionArray, x.StorageArg != nil, x.StorageArg, false, nil, false, nil, false)
}

type Arg15 struct {
	FriskyArgArray []FriskyArg
	StorageArg     *StorageArg
}

func (x *Arg15) UnmarshalJSON(data []byte) error {
	x.FriskyArgArray = nil
	x.StorageArg = nil
	var c StorageArg
	object, err := unmarshalUnion(data, nil, nil, nil, nil, true, &x.FriskyArgArray, true, &c, false, nil, false, nil, false)
	if err != nil {
		return err
	}
	if object {
		x.StorageArg = &c
	}
	return nil
}

func (x *Arg15) MarshalJSON() ([]byte, error) {
	return marshalUnion(nil, nil, nil, nil, x.FriskyArgArray != nil, x.FriskyArgArray, x.StorageArg != nil, x.StorageArg, false, nil, false, nil, false)
}

type HilariousValue struct {
	IndigoValue      *IndigoValue
	StickyValueArray []StickyValue
}

func (x *HilariousValue) UnmarshalJSON(data []byte) error {
	x.StickyValueArray = nil
	x.IndigoValue = nil
	var c IndigoValue
	object, err := unmarshalUnion(data, nil, nil, nil, nil, true, &x.StickyValueArray, true, &c, false, nil, false, nil, false)
	if err != nil {
		return err
	}
	if object {
		x.IndigoValue = &c
	}
	return nil
}

func (x *HilariousValue) MarshalJSON() ([]byte, error) {
	return marshalUnion(nil, nil, nil, nil, x.StickyValueArray != nil, x.StickyValueArray, x.IndigoValue != nil, x.IndigoValue, false, nil, false, nil, false)
}

type Arg16 struct {
	Arg4           *Arg4
	PurpleArgArray []PurpleArg
}

func (x *Arg16) UnmarshalJSON(data []byte) error {
	x.PurpleArgArray = nil
	x.Arg4 = nil
	var c Arg4
	object, err := unmarshalUnion(data, nil, nil, nil, nil, true, &x.PurpleArgArray, true, &c, false, nil, false, nil, false)
	if err != nil {
		return err
	}
	if object {
		x.Arg4 = &c
	}
	return nil
}

func (x *Arg16) MarshalJSON() ([]byte, error) {
	return marshalUnion(nil, nil, nil, nil, x.PurpleArgArray != nil, x.PurpleArgArray, x.Arg4 != nil, x.Arg4, false, nil, false, nil, false)
}

func unmarshalUnion(data []byte, pi **int64, pf **float64, pb **bool, ps **string, haveArray bool, pa interface{}, haveObject bool, pc interface{}, haveMap bool, pm interface{}, haveEnum bool, pe interface{}, nullable bool) (bool, error) {
	if pi != nil {
		*pi = nil
	}
	if pf != nil {
		*pf = nil
	}
	if pb != nil {
		*pb = nil
	}
	if ps != nil {
		*ps = nil
	}

	dec := json.NewDecoder(bytes.NewReader(data))
	dec.UseNumber()
	tok, err := dec.Token()
	if err != nil {
		return false, err
	}

	switch v := tok.(type) {
	case json.Number:
		if pi != nil {
			i, err := v.Int64()
			if err == nil {
				*pi = &i
				return false, nil
			}
		}
		if pf != nil {
			f, err := v.Float64()
			if err == nil {
				*pf = &f
				return false, nil
			}
			return false, errors.New("Unparsable number")
		}
		return false, errors.New("Union does not contain number")
	case float64:
		return false, errors.New("Decoder should not return float64")
	case bool:
		if pb != nil {
			*pb = &v
			return false, nil
		}
		return false, errors.New("Union does not contain bool")
	case string:
		if haveEnum {
			return false, json.Unmarshal(data, pe)
		}
		if ps != nil {
			*ps = &v
			return false, nil
		}
		return false, errors.New("Union does not contain string")
	case nil:
		if nullable {
			return false, nil
		}
		return false, errors.New("Union does not contain null")
	case json.Delim:
		if v == '{' {
			if haveObject {
				return true, json.Unmarshal(data, pc)
			}
			if haveMap {
				return false, json.Unmarshal(data, pm)
			}
			return false, errors.New("Union does not contain object")
		}
		if v == '[' {
			if haveArray {
				return false, json.Unmarshal(data, pa)
			}
			return false, errors.New("Union does not contain array")
		}
		return false, errors.New("Cannot handle delimiter")
	}
	return false, errors.New("Cannot unmarshal union")

}

func marshalUnion(pi *int64, pf *float64, pb *bool, ps *string, haveArray bool, pa interface{}, haveObject bool, pc interface{}, haveMap bool, pm interface{}, haveEnum bool, pe interface{}, nullable bool) ([]byte, error) {
	if pi != nil {
		return json.Marshal(*pi)
	}
	if pf != nil {
		return json.Marshal(*pf)
	}
	if pb != nil {
		return json.Marshal(*pb)
	}
	if ps != nil {
		return json.Marshal(*ps)
	}
	if haveArray {
		return json.Marshal(pa)
	}
	if haveObject {
		return json.Marshal(pc)
	}
	if haveMap {
		return json.Marshal(pm)
	}
	if haveEnum {
		return json.Marshal(pe)
	}
	if nullable {
		return json.Marshal(nil)
	}
	return nil, errors.New("Union must not be null")
}

