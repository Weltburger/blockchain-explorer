package models

type Arg9 struct {
	ImplicitOperationsResultStorage *ImplicitOperationsResultStorage
	IndigoArgArray                  []IndigoArg
}

func (x *Arg9) UnmarshalJSON(data []byte) error {
	x.IndigoArgArray = nil
	x.ImplicitOperationsResultStorage = nil
	var c ImplicitOperationsResultStorage
	object, err := unmarshalUnion(data, nil, nil, nil, nil, true, &x.IndigoArgArray, true, &c, false, nil, false, nil, false)
	if err != nil {
		return err
	}
	if object {
		x.ImplicitOperationsResultStorage = &c
	}
	return nil
}

func (x *Arg9) MarshalJSON() ([]byte, error) {
	return marshalUnion(nil, nil, nil, nil, x.IndigoArgArray != nil, x.IndigoArgArray, x.ImplicitOperationsResultStorage != nil, x.ImplicitOperationsResultStorage, false, nil, false, nil, false)
}
