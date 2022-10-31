// Copyright(C) 2022 github.com/fsgo  All Rights Reserved.
// Author: hidu <duv123@gmail.com>
// Date: 2022/10/31

package main

import (
	"flag"
	"log"
	"os"

	"github.com/fsgo/gopt/internal"
)

func main() {
	flag.Parse()
	args := stringSlice(os.Args)
	if len(args) < 2 || args.get(1) == "help" {
		flag.Usage()
		return
	}
	var err error
	switch args[1] {
	case "list":
		err = internal.List(args[2:])
	}
	if err != nil {
		log.Fatalf("error: %s failed, %v\n", args[1], err)
	}
}

type stringSlice []string

func (s stringSlice) get(index int) string {
	if index >= len(s) {
		return ""
	}
	return s[index]
}
