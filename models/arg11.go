package models

type Arg11 struct {
	FluffyArg         *FluffyArg
	HilariousArgArray []HilariousArg
}

func (x *Arg11) UnmarshalJSON(data []byte) error {
	x.HilariousArgArray = nil
	x.FluffyArg = nil
	var c FluffyArg
	object, err := unmarshalUnion(data, nil, nil, nil, nil, true, &x.HilariousArgArray, true, &c, false, nil, false, nil, false)
	if err != nil {
		return err
	}
	if object {
		x.FluffyArg = &c
	}
	return nil
}

func (x *Arg11) MarshalJSON() ([]byte, error) {
	return marshalUnion(nil, nil, nil, nil, x.HilariousArgArray != nil, x.HilariousArgArray, x.FluffyArg != nil, x.FluffyArg, false, nil, false, nil, false)
}
