package models

type OperationResultLazyStorageDiff struct {
	Kind LazyStorageDiffKind `json:"kind"`
	ID   string              `json:"id"`
	Diff FluffyDiff          `json:"diff"`
}
