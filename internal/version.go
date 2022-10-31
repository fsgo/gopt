// Copyright(C) 2022 github.com/fsgo  All Rights Reserved.
// Author: hidu <duv123@gmail.com>
// Date: 2022/10/31

package internal

import (
	"context"
	"log"
	"os"
	"path/filepath"
	"time"

	"github.com/fsgo/fscache"
	"github.com/fsgo/fscache/filecache"
	"github.com/fsgo/gomodule"
)

var versionCache fscache.SCache

func init() {
	fc, err := filecache.NewSCache(&filecache.Option{
		Dir:        filepath.Join(os.TempDir(), "gopt", "latest_cache"),
		GCInterval: time.Hour,
	})
	if err != nil {
		log.Fatalln(err)
	}
	versionCache = fc
}

func latest(ctx context.Context, path string) (*gomodule.Info, error) {
	ret := versionCache.Get(ctx, path)
	if ret.Has() {
		var info *gomodule.Info
		if ok, e1 := ret.Value(&info); ok && e1 == nil {
			return info, nil
		}
	}
	pm := &gomodule.GoProxy{
		Module: path,
	}
	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()
	info, e2 := pm.Latest(ctx)
	if info != nil {
		versionCache.Set(ctx, path, info, time.Hour)
	}
	return info, e2
}
