// Copyright(C) 2022 github.com/fsgo  All Rights Reserved.
// Author: hidu <duv123@gmail.com>
// Date: 2022/10/31

package internal

import (
	"context"
	"debug/buildinfo"
	"flag"
	"fmt"
	"log"
	"path/filepath"
	"time"

	"github.com/fatih/color"
	"golang.org/x/mod/semver"
)

func Update(args []string) error {
	u := &updater{}
	if err := u.setup(args); err != nil {
		return err
	}
	return u.Run()
}

type updater struct {
	flags     *flag.FlagSet
	timeout   int
	numFailed int
}

func (u *updater) setup(args []string) error {
	u.flags = flag.NewFlagSet("update", flag.ExitOnError)
	u.flags.IntVar(&u.timeout, "T", 60, `update timeout, seconds`)
	return u.flags.Parse(args)
}

func (u *updater) Run() error {
	args := u.flags.Args()
	if len(args) > 1 {
		return fmt.Errorf("not support %q", args)
	}
	if len(args) == 0 {
		sc := &Scanner{
			Call: u.onCall,
		}
		return sc.Run()
	}
	wantName := args[0] + exe()
	var found bool
	sc := &Scanner{
		Call: func(name string, bi *buildinfo.BuildInfo) error {
			bn := filepath.Base(name)
			if bn != wantName {
				return nil
			}
			found = true
			return u.onCall(name, bi)
		},
	}
	err := sc.Run()
	if err != nil {
		return err
	}
	if !found {
		return fmt.Errorf("not found %q", wantName)
	}
	return nil
}

func (u *updater) getTimeout() time.Duration {
	if u.timeout > 0 {
		return time.Duration(u.timeout) * time.Second
	}
	return 60 * time.Second
}

func (u *updater) onCall(name string, bi *buildinfo.BuildInfo) error {
	bn := filepath.Base(name)
	log.SetPrefix(color.GreenString("[" + bn + "] "))
	log.Println(color.CyanString(name), bi.Path+"@"+bi.Main.Version)

	if bi.Main.Version == develVersion {
		log.Println("skipped update by version:", develVersion)
		return nil
	}

	ctx, cancel := context.WithTimeout(context.Background(), u.getTimeout())
	defer cancel()
	mp, err := latest(ctx, bi.Main.Path)
	if err != nil {
		u.numFailed++
		log.Println(color.RedString("get latest info failed: " + err.Error()))
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
	cmd := newGoCommand(ctx, "install", bi.Path+"@latest")
	log.Println("will update:", cmd.String())
	if err = cmd.Run(); err != nil {
		log.Println(color.RedString("install failed: " + err.Error()))
	} else {
		log.Println(color.GreenString("install success"))
	}
	return nil
}
