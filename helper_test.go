package main

import (
	"os"
	"time"
)

func readFile(fn string) (string, error) {
	s, err := os.ReadFile(fn)
	if err != nil {
		return "", err
	}
	return string(s), nil
}

func NewDate(year int, month time.Month, day int, hour int, min int) time.Time {
	return time.Date(year, month, day, hour, min, 0, 0, time.UTC)
}
