package models

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
