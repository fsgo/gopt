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
	"strings"
	"time"

	"github.com/fatih/color"
	"github.com/fsgo/cmdutil"
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
	timeout   time.Duration
	numFailed int
}

func (u *updater) setup(args []string) error {
	u.flags = flag.NewFlagSet("update", flag.ExitOnError)
	u.flags.DurationVar(&u.timeout, "t", 120*time.Second, `update timeout`)
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
		return u.timeout
	}
	return 120 * time.Second
}

func (u *updater) onCall(name string, bi *buildinfo.BuildInfo) error {
	bn := filepath.Base(name)
	log.SetPrefix(color.GreenString("[" + bn + "] "))
	log.Println(color.CyanString(name), bi.Path+"@"+bi.Main.Version)

	if !strings.Contains(bi.Path, ".") && bi.Main.Version == develVersion {
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
	log.Println("found latest:", mp.Version, mp.Time.Local().String())
	if semver.Compare(mp.Version, bi.Main.Version) < 1 {
		return nil
	}
	u.install(ctx, bi, name)
	return nil
}

func (u *updater) install(ctx context.Context, bi *buildinfo.BuildInfo, rawName string) {
	// 二进制文件名已经改名，不能直接使用 go install 安装替换
	// 先安装到临时目录，然后再替换
	useRawDir := filepath.Base(bi.Path)+exe() == filepath.Base(rawName)

	biEnv := make(map[string]string)
	args := []string{"install"}
	for _, tm := range bi.Settings {
		switch tm.Key {
		case "-tags":
			args = append(args, "-tags", tm.Value)
		case "CGO_ENABLED":
			biEnv[tm.Key] = tm.Value
		}
	}
	args = append(args, bi.Path+"@latest")

	cmd := newGoCommand(ctx, args...)
	oe := &cmdutil.OSEnv{}
	oe.WithEnviron(cmd.Environ())
	for k, v := range biEnv {
		_ = oe.Set(k, v)
	}

	goBinTMP := goBinTMPDir()
	if useRawDir {
		_ = oe.Set("GOBIN", filepath.Dir(rawName))
	} else {
		_ = oe.Set("GOBIN", goBinTMP)
		log.Println("TMP GOBIN=", goBinTMP)
	}
	cmd.Env = oe.Environ()
	log.Println("will update:", cmd.String())
	err := cmd.Run()
	if err != nil {
		log.Println(color.RedString("install failed: " + err.Error()))
		return
	}

	if !useRawDir {
		distPath := filepath.Join(goBinTMP, filepath.Base(bi.Path)) + exe()
		if e1 := mv(distPath, rawName); e1 != nil {
			log.Println(color.RedString("mv(%q,%q) failed:%v", distPath, rawName, e1))
			return
		}
	}
	log.Println(color.GreenString("install success"))
}
