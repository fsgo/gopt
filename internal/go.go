// Copyright(C) 2022 github.com/fsgo  All Rights Reserved.
// Author: hidu <duv123@gmail.com>
// Date: 2022/11/12

package internal

import (
	"context"
	"os"
	"os/exec"

	"github.com/fsgo/cmdutil/gosdk"
)

func newGoCommand(ctx context.Context, arg ...string) *exec.Cmd {
	goBin := gosdk.LatestOrDefault()
	cmd := exec.CommandContext(ctx, goBin, arg...)
	cmd.Stderr = os.Stderr
	cmd.Stdout = os.Stdout
	cmd.Stdin = os.Stdin
	cmd.Env = gosdk.GoCmdEnv(goBin, nil)
	return cmd
}
