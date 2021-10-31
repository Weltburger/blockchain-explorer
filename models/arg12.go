package models

type Arg12 struct {
	FluffyArgArray []FluffyArg
	MagentaArg     *MagentaArg
}

func (x *Arg12) UnmarshalJSON(data []byte) error {
	x.FluffyArgArray = nil
	x.MagentaArg = nil
	var c MagentaArg
	object, err := unmarshalUnion(data, nil, nil, nil, nil, true, &x.FluffyArgArray, true, &c, false, nil, false, nil, false)
	if err != nil {
		return err
	}
	if object {
		x.MagentaArg = &c
	}
	return nil
}

func (x *Arg12) MarshalJSON() ([]byte, error) {
	return marshalUnion(nil, nil, nil, nil, x.FluffyArgArray != nil, x.FluffyArgArray, x.MagentaArg != nil, x.MagentaArg, false, nil, false, nil, false)
}
