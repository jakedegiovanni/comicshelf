package main

import (
	"encoding/json"
	"fmt"
)

func prettyPrint(i interface{}) error {
	b, err := json.MarshalIndent(i, "", "  ")
	if err != nil {
		return err
	}

	fmt.Println(string(b))
	return nil
}
