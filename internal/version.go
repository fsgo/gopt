// Copyright(C) 2022 github.com/fsgo  All Rights Reserved.
// Author: hidu <duv123@gmail.com>
// Date: 2022/10/31

package internal

import (
	"context"
	"time"

	"github.com/fsgo/gomodule"
)

func latest(ctx context.Context, path string) (*gomodule.Info, error) {
	pm := &gomodule.GoProxy{
		Module: path,
	}
	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()
	return pm.Latest(ctx)
}
