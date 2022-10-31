// Copyright(C) 2022 github.com/fsgo  All Rights Reserved.
// Author: hidu <duv123@gmail.com>
// Date: 2022/10/31

package internal

import (
	"context"
	"debug/buildinfo"
	"flag"
	"fmt"
	"path/filepath"
	"runtime/debug"
)

type list struct {
	devel    string
	fullPath bool
	latest   bool
}

func (l *list) setup(args []string) error {
	cf := flag.NewFlagSet("list", flag.ExitOnError)
	cf.StringVar(&l.devel, "dev", "", `filter devel. 'yes': only devel; 'no': no devel; default is '': no filter`)
	cf.BoolVar(&l.fullPath, "fp", false, "print full file path")
	cf.BoolVar(&l.latest, "l", true, "get latest version info")
	return cf.Parse(args)
}

func (l *list) develFilter(m debug.Module) bool {
	switch l.devel {
	case "yes", "only":
		return m.Version == "(devel)"
	case "no":
		return m.Version != "(devel)"
	default:
		return true
	}
}

func (l *list) getLatest(path string) (string, string) {
	if !l.latest {
		return "", ""
	}
	mp, err := latest(context.Background(), path)
	if err != nil {
		return "", ""
	}
	latestVersion := mp.Version
	if len(latestVersion) > 13 {
		latestVersion = latestVersion[:13] + "*"
	}
	latestTime := mp.Time.Format("20060201")
	return latestVersion, latestTime
}

func (l *list) onCall(name string, bi *buildinfo.BuildInfo) error {
	listTpl := "%-50s %-15s %-15s  %8s %s\n"
	bn := name
	if !l.fullPath {
		bn = filepath.Base(name)
	}
	if !l.develFilter(bi.Main) {
		return nil
	}
	m := bi.Main.Path + "@" + bi.Main.Version

	latestVersion, latestTime := l.getLatest(bi.Main.Path)
	fmt.Printf(listTpl, bn, bi.GoVersion, latestVersion, latestTime, m)
	return nil
}

func List(args []string) error {
	l := &list{}
	if err := l.setup(args); err != nil {
		return err
	}
	sc := &Scanner{
		Call: l.onCall,
	}
	return sc.Run()
}
