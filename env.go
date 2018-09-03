package main

import (
	"fmt"
	"os"
)

func setEnv(env map[string]string) error {
	for envKey, keySuffix := range env {
		key, err := getKey(keySuffix)
		if err != nil {
			fmt.Println(keySuffix)
			return err
		}
		os.Setenv(envKey, key)
	}
	return nil
}
