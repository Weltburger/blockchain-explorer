package models

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
