package util

import "github.com/mitchellh/mapstructure"

func MarshalToStruct[T, D any](dst D) (*T, error) {

	result := new(T)
	if err := mapstructure.Decode(dst, &result); err != nil {
		return nil, err
	}
	return result, nil
}
