package models

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
