// Copyright(C) 2022 github.com/fsgo  All Rights Reserved.
// Author: hidu <duv123@gmail.com>
// Date: 2022/10/31

package main

import (
	"context"

	"github.com/fsgo/gopt/internal"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	internal.Run(ctx)
}
