package models

type Arg10 struct {
	FluffyArg      *FluffyArg
	IndigoArgArray []IndigoArg
}

func (x *Arg10) UnmarshalJSON(data []byte) error {
	x.IndigoArgArray = nil
	x.FluffyArg = nil
	var c FluffyArg
	object, err := unmarshalUnion(data, nil, nil, nil, nil, true, &x.IndigoArgArray, true, &c, false, nil, false, nil, false)
	if err != nil {
		return err
	}
	if object {
		x.FluffyArg = &c
	}
	return nil
}

func (x *Arg10) MarshalJSON() ([]byte, error) {
	return marshalUnion(nil, nil, nil, nil, x.IndigoArgArray != nil, x.IndigoArgArray, x.FluffyArg != nil, x.FluffyArg, false, nil, false, nil, false)
}
