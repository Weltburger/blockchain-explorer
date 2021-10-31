package models

type ResultStorage struct {
	PurpleStorage   *PurpleStorage
	StorageArgArray []StorageArg
}

func (x *ResultStorage) UnmarshalJSON(data []byte) error {
	x.StorageArgArray = nil
	x.PurpleStorage = nil
	var c PurpleStorage
	object, err := unmarshalUnion(data, nil, nil, nil, nil, true, &x.StorageArgArray, true, &c, false, nil, false, nil, false)
	if err != nil {
		return err
	}
	if object {
		x.PurpleStorage = &c
	}
	return nil
}

func (x *ResultStorage) MarshalJSON() ([]byte, error) {
	return marshalUnion(nil, nil, nil, nil, x.StorageArgArray != nil, x.StorageArgArray, x.PurpleStorage != nil, x.PurpleStorage, false, nil, false, nil, false)
}
