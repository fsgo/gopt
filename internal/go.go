// Copyright(C) 2022 github.com/fsgo  All Rights Reserved.
// Author: hidu <duv123@gmail.com>
// Date: 2022/11/12

package internal

import (
	"bytes"
	"context"
	"os"
	"os/exec"
	"strings"
	"sync"
	"time"

	"github.com/fsgo/cmdutil/gosdk"
	"golang.org/x/mod/semver"
)

func newGoCommand(ctx context.Context, arg ...string) *exec.Cmd {
	goBin := gosdk.LatestOrDefault(ctx)
	cmd := exec.CommandContext(ctx, goBin, arg...)
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin
	cmd.Env = gosdk.GoCmdEnv(goBin, nil)
	return cmd
}

var goVersionStr string
var goVersionOnce sync.Once

func goVersion() string {
	goVersionOnce.Do(func() {
		ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
		defer cancel()
		goBin := gosdk.LatestOrDefault(ctx)
		cmd := exec.CommandContext(ctx, goBin, "version")
		out, _ := cmd.CombinedOutput()
		if bytes.HasPrefix(out, []byte("go version ")) {
			arr := strings.Fields(string(out))
			if len(arr) > 3 {
				goVersionStr = strings.ReplaceAll(arr[2], "go", "v")
			}
		}
	})
	return goVersionStr
}

// sameBinGoVersion 检查二进制使用的 go version 是否 >= 当前环境 go sdk 的版本
func sameBinGoVersion(binGoVersion string) bool {
	gv := goVersion()
	return gv == "" || semver.Compare(binGoVersion, gv) > 1
}
