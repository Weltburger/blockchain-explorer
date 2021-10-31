package models

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
