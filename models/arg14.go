package models

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
