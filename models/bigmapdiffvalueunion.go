package models

type BigMapDiffValueUnion struct {
	KeyArray   []Key
	KeyElement *KeyElement
}

func (x *BigMapDiffValueUnion) UnmarshalJSON(data []byte) error {
	x.KeyArray = nil
	x.KeyElement = nil
	var c KeyElement
	object, err := unmarshalUnion(data, nil, nil, nil, nil, true, &x.KeyArray, true, &c, false, nil, false, nil, false)
	if err != nil {
		return err
	}
	if object {
		x.KeyElement = &c
	}
	return nil
}

func (x *BigMapDiffValueUnion) MarshalJSON() ([]byte, error) {
	return marshalUnion(nil, nil, nil, nil, x.KeyArray != nil, x.KeyArray, x.KeyElement != nil, x.KeyElement, false, nil, false, nil, false)
}
