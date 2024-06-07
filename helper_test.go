package main

import "os"

func readFile(fn string) (string, error) {
	s, err := os.ReadFile(fn)
	if err != nil {
		return "", err
	}
	return string(s), nil
}
