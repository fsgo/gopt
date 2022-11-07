// Copyright(C) 2022 github.com/fsgo  All Rights Reserved.
// Author: hidu <duv123@gmail.com>
// Date: 2022/10/31

package internal

import (
	"debug/buildinfo"
	"io/fs"
	"os"
	"path/filepath"
	"runtime"
	"strings"
)

type Scanner struct {
	Call func(name string, bi *buildinfo.BuildInfo) error
	Dirs []string
}

func (s *Scanner) getDirs() []string {
	if len(s.Dirs) > 0 {
		return s.Dirs
	}
	return strings.Split(os.Getenv("GOBIN"), ":")
}

func (s *Scanner) Run() error {
	dirs := s.getDirs()
	for _, dir := range dirs {
		err1 := filepath.Walk(dir, func(path string, info fs.FileInfo, err error) error {
			if err != nil {
				return err
			}
			if !info.IsDir() {
				bi, err2 := buildinfo.ReadFile(path)
				if err2 == nil {
					s.Call(path, bi)
				}
			}

			return nil
		})
		if err1 != nil {
			return err1
		}
	}
	return nil
}

func isWindows() bool {
	return runtime.GOOS == "windows"
}

func exe() string {
	if isWindows() {
		return ".exe"
	}
	return ""
}
