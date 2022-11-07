// Copyright(C) 2022 github.com/fsgo  All Rights Reserved.
// Author: hidu <duv123@gmail.com>
// Date: 2022/10/31

package internal

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"time"

	"github.com/fsgo/fscache"
	"github.com/fsgo/fscache/filecache"
	"github.com/fsgo/gomodule"
)

var versionCache fscache.SCache

func init() {
	cacheDir := filepath.Join(os.TempDir(), "gopt", "latest_cache")
	fc, err := filecache.NewSCache(&filecache.Option{
		Dir:        cacheDir,
		GCInterval: time.Hour,
	})
	if err != nil {
		log.Fatalln(err)
	}
	versionCache = fc
}

func latest(ctx context.Context, path string) (*gomodule.Info, error) {
	ret := versionCache.Get(ctx, path)
	var info *gomodule.Info
	if ok, _ := ret.Value(&info); ok {
		return info, nil
	}
	pm := &gomodule.GoProxy{
		Module: path,
	}
	info, e2 := pm.Latest(ctx)
	if info != nil {
		versionCache.Set(ctx, path, info, time.Hour)
	}
	if e2 == nil {
		return info, nil
	}

	cmd1 := exec.CommandContext(ctx, "go", "list", "-m", "-json", path+"@latest")
	cmd1.Stdin = os.Stdin
	buf1 := &bytes.Buffer{}
	cmd1.Stderr = buf1
	out, err3 := cmd1.Output()
	if err3 != nil {
		return nil, fmt.Errorf("exec %q: %w\n %s", cmd1.String(), err3, buf1.String())
	}
	if err4 := json.Unmarshal(out, &info); err4 != nil {
		return nil, err4
	}
	if info == nil || len(info.Version) == 0 {
		return nil, fmt.Errorf("invald go list response: %q", out)
	}
	return info, nil
}
