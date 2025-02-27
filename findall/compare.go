//go:build exclude

package main

import (
	"encoding/json"
	"os"
)

func compare(left, right m2s) error {

	onlyLeft, onlyRight := m2s{}, m2s{}
	for lk, lv := range left {
		if _, found := right[lk]; !found {
			// only present in left
			onlyLeft[lk] = lv
		}
	}
	for rk, rv := range right {
		if _, found := left[rk]; !found {
			// only present in right
			onlyRight[rk] = rv
		}
	}
}

func compareFiles(leftName, rightName string) error {

	left, err := loadM2s(leftName)
	if err != nil {
		return err
	}
	right, err := loadM2s(rightName)
	if err != nil {
		return err
	}

	return compare(left, right)
}

func loadM2s(name string) (m2s, error) {
	raw, err := os.ReadFile(name)
	if err != nil {
		return nil, err
	}

	var data m2s
	if err := json.Unmarshal(raw, &data); err != nil {
		return nil, err
	}

	return data, nil
}
