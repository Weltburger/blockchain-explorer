package models

type ContentMetadata struct {
	BalanceUpdates           []MetadataBalanceUpdate   `json:"balance_updates"`
	Delegate                 string                    `json:"delegate,omitempty"`
	Slots                    []int64                   `json:"slots,omitempty"`
	OperationResult          *OperationResult          `json:"operation_result,omitempty"`
	InternalOperationResults []InternalOperationResult `json:"internal_operation_results,omitempty"`
}
