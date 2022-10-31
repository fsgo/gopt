// Copyright(C) 2022 github.com/fsgo  All Rights Reserved.
// Author: hidu <duv123@gmail.com>
// Date: 2022/10/31

package internal

import (
	"context"
	"debug/buildinfo"
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"runtime/debug"
	"strings"
	"time"

	"github.com/fsgo/gomodule"
)

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

type list struct {
	devel     string
	latest    bool
	printJSON bool
	id        int
}

func (l *list) setup(args []string) error {
	cf := flag.NewFlagSet("list", flag.ExitOnError)
	cf.StringVar(&l.devel, "dev", "", `filter devel. 'yes': only devel; 'no': no devel; default is '': no filter`)
	cf.BoolVar(&l.latest, "l", true, "get latest version info")
	cf.BoolVar(&l.printJSON, "json", false, "print JSON result")
	return cf.Parse(args)
}

const develVersion = "(devel)"

func (l *list) develFilter(m debug.Module) bool {
	switch l.devel {
	case "yes", "only":
		return m.Version == develVersion
	case "no":
		return m.Version != develVersion
	default:
		return true
	}
}

func (l *list) onCall(name string, bi *buildinfo.BuildInfo) error {
	if !l.develFilter(bi.Main) {
		return nil
	}
	l.id++
	si := &scanInfo{
		ID:        l.id,
		Name:      name,
		BuildInfo: bi,
	}
	si.FileInfo, _ = os.Stat(name)

	if l.latest {
		mp, err := latest(context.Background(), bi.Main.Path)
		if err == nil {
			si.Latest = mp
		}
	}
	if l.printJSON {
		fmt.Println(si.JSON())
	} else {
		fmt.Println(si.String())
	}
	return nil
}

type scanInfo struct {
	FileInfo  os.FileInfo
	BuildInfo *buildinfo.BuildInfo
	Latest    *gomodule.Info
	Name      string
	ID        int
}

const timeLayout = "2006-01-02 15:04:05"

func (si *scanInfo) String() string {
	const tpl = "%15s : %s\n"

	bs := &strings.Builder{}
	m := si.BuildInfo.Main
	fmt.Fprint(bs, ConsoleGreen(fmt.Sprintf("%3d %s\n", si.ID, si.Name)))
	fmt.Fprintf(bs, tpl, "Path", si.BuildInfo.Path)
	fmt.Fprintf(bs, tpl, "Go", si.BuildInfo.GoVersion)
	fmt.Fprintf(bs, tpl, "Version", m.Version)

	var modTime time.Time
	if si.FileInfo != nil {
		modTime = si.FileInfo.ModTime()
		fmt.Fprintf(bs, tpl, "Install Time", modTime.Format(timeLayout))
	}

	if si.Latest != nil {
		fmt.Fprintf(bs, tpl, "Latest Version", si.Latest.Version)
		fmt.Fprintf(bs, tpl, "Latest Time", si.Latest.Time.Format(timeLayout))

		if m.Version != develVersion && !modTime.IsZero() && si.Latest.Time.After(modTime) {
			days := si.Latest.Time.Sub(modTime).Hours() / 24
			fmt.Fprint(bs, ConsoleRed(fmt.Sprintf(tpl, "Expired", fmt.Sprintf("%.1f days", days))))
			fmt.Fprint(bs, ConsoleRed(fmt.Sprintf(tpl, "Install", "go install "+si.BuildInfo.Path+"@latest")))
		}
	}
	return bs.String()
}

func (si *scanInfo) JSON() string {
	info := map[string]any{
		"ID":      si.ID,
		"Name":    si.Name,
		"Path":    si.BuildInfo.Path,
		"Go":      si.BuildInfo.GoVersion,
		"Version": si.BuildInfo.Main.Version,
	}
	if si.FileInfo != nil {
		info["InstallTime"] = si.FileInfo.ModTime().Format(timeLayout)
	}
	if si.Latest != nil {
		info["LatestVersion"] = si.Latest.Version
		info["LatestTime"] = si.Latest.Time.Format(timeLayout)
	}
	bs, _ := json.Marshal(info)
	return string(bs)
}
