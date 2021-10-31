package models

type IndecentValue struct {
	FluffyValue      *FluffyValue
	PurpleValueArray []PurpleValue
}

func (x *IndecentValue) UnmarshalJSON(data []byte) error {
	x.PurpleValueArray = nil
	x.FluffyValue = nil
	var c FluffyValue
	object, err := unmarshalUnion(data, nil, nil, nil, nil, true, &x.PurpleValueArray, true, &c, false, nil, false, nil, false)
	if err != nil {
		return err
	}
	if object {
		x.FluffyValue = &c
	}
	return nil
}

func (x *IndecentValue) MarshalJSON() ([]byte, error) {
	return marshalUnion(nil, nil, nil, nil, x.PurpleValueArray != nil, x.PurpleValueArray, x.FluffyValue != nil, x.FluffyValue, false, nil, false, nil, false)
}
