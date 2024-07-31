package xos

import (
	"syscall"
)

func MoveFile(src string, dst string) error {
	return syscall.Rename(src, dst)
}
