package models

type ResultLazyStorageDiff struct {
	Kind LazyStorageDiffKind `json:"kind"`
	ID   string              `json:"id"`
	Diff PurpleDiff          `json:"diff"`
}
