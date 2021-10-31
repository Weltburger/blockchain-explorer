package models

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
