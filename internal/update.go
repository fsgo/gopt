// Copyright(C) 2022 github.com/fsgo  All Rights Reserved.
// Author: hidu <duv123@gmail.com>
// Date: 2022/10/31

package internal

import (
	"context"
	"debug/buildinfo"
	"flag"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"time"

	"golang.org/x/mod/semver"
)

func Update(args []string) error {
	u := &updater{}
	if err := u.setup(args); err != nil {
		return err
	}
	sc := &Scanner{
		Call: u.onCall,
	}
	return sc.Run()
}

type updater struct {
	timeout   int
	numFailed int
}

func (u *updater) setup(args []string) error {
	cf := flag.NewFlagSet("list", flag.ExitOnError)
	cf.IntVar(&u.timeout, "timeout", 60, `update timeout, seconds`)
	return cf.Parse(args)
}

func (u *updater) getTimeout() time.Duration {
	if u.timeout > 0 {
		return time.Duration(u.timeout) * time.Second
	}
	return 60 * time.Second
}

func (u *updater) onCall(name string, bi *buildinfo.BuildInfo) error {
	bn := filepath.Base(name)
	log.SetPrefix(ConsoleGreen("[" + bn + "] "))
	log.Println(ConsoleGreen(name), bi.Path+"@"+bi.Main.Version)

	if bi.Main.Version == develVersion {
		log.Println("skipped update by version:", develVersion)
		return nil
	}

	ctx, cancel := context.WithTimeout(context.Background(), u.getTimeout())
	defer cancel()
	mp, err := latest(ctx, bi.Path)
	if err != nil {
		u.numFailed++
		log.Println("get latest info failed:", err.Error())
		return nil
	}
	log.Println("found latest:", mp.Version, mp.Time.String())
	if semver.Compare(mp.Version, bi.Main.Version) < 1 {
		return nil
	}
	// 二进制文件名已经改名，暂时不能直接使用 go install 安装替换
	if filepath.Base(bi.Path) != bn {
		log.Println("filename not match, skipped")
		return nil
	}
	cmd := exec.CommandContext(ctx, "go", "install", bi.Path+"@latest")
	log.Println("will update:", cmd.String())
	cmd.Stderr = os.Stderr
	cmd.Stdout = os.Stdout
	cmd.Stdin = os.Stdin
	if err = cmd.Run(); err != nil {
		log.Println(ConsoleRed("install failed: " + err.Error()))
	} else {
		log.Println(ConsoleGreen("install success"))
	}
	return nil
}
