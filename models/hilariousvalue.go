package models

type HilariousValue struct {
	IndigoValue      *IndigoValue
	StickyValueArray []StickyValue
}

func (x *HilariousValue) UnmarshalJSON(data []byte) error {
	x.StickyValueArray = nil
	x.IndigoValue = nil
	var c IndigoValue
	object, err := unmarshalUnion(data, nil, nil, nil, nil, true, &x.StickyValueArray, true, &c, false, nil, false, nil, false)
	if err != nil {
		return err
	}
	if object {
		x.IndigoValue = &c
	}
	return nil
}

func (x *HilariousValue) MarshalJSON() ([]byte, error) {
	return marshalUnion(nil, nil, nil, nil, x.StickyValueArray != nil, x.StickyValueArray, x.IndigoValue != nil, x.IndigoValue, false, nil, false, nil, false)
}
