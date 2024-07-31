//go:build !windows

package xos

import (
	"os"
	"strconv"
)

func MoveFile(from string, to string) error {
	to1 := to + "_" + strconv.Itoa(rand.Int())
	if e1 := os.Rename(to, to1); e1 != nil {
		return e1
	}
	e2 := os.Rename(from, to)
	if e2 == nil {
		_ = os.Remove(to1)
		return nil
	}
	_ = os.Rename(to1, to)
	return e2
}
