package models

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
