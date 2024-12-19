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

	"github.com/fatih/color"
	"github.com/fsgo/gomodule"
)

func List(args []string) error {
	l := &list{
		id: 1,
	}
	if err := l.Setup(args); err != nil {
		return err
	}
	return l.Run()
}

type list struct {
	flags       *flag.FlagSet
	devel       string
	id          int
	timeout     int
	latest      bool
	onlyExpired bool
	printJSON   bool
}

func (l *list) Setup(args []string) error {
	l.flags = flag.NewFlagSet("list", flag.ExitOnError)
	l.flags.StringVar(&l.devel, "dev", "", `filter devel. 'yes': only devel; 'no': no devel; default is '': no filter`)
	l.flags.BoolVar(&l.latest, "l", true, "get latest version info")
	l.flags.BoolVar(&l.onlyExpired, "e", false, "filter only expired")
	l.flags.BoolVar(&l.printJSON, "json", false, "print JSON result")
	l.flags.IntVar(&l.timeout, "t", 5, `list timeout, seconds`)
	return l.flags.Parse(args)
}

func (l *list) Run() error {
	args := l.flags.Args()
	if len(args) > 0 {
		return fmt.Errorf("not support %q", args)
	}

	if l.onlyExpired {
		l.latest = true
		l.devel = "no"
	}

	sc := &Scanner{
		Call: l.onCall,
	}
	return sc.Run()
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

func (l *list) getTimeout() time.Duration {
	if l.timeout > 0 {
		return time.Duration(l.timeout) * time.Second
	}
	return 5 * time.Second
}

func (l *list) onCall(name string, bi *buildinfo.BuildInfo) error {
	if !l.develFilter(bi.Main) {
		return nil
	}
	si := &scanInfo{
		ID:        l.id,
		Name:      name,
		BuildInfo: bi,
	}
	si.FileInfo, _ = os.Stat(name)

	if l.latest {
		ctx, cancel := context.WithTimeout(context.Background(), l.getTimeout())
		defer cancel()
		si.Latest, si.Err = latest(ctx, bi.Main.Path)
	}

	if l.onlyExpired && !si.expired() {
		return nil
	}

	if l.printJSON {
		fmt.Println(si.JSON())
	} else {
		fmt.Println(si.String())
	}
	l.id++
	return nil
}

type scanInfo struct {
	FileInfo  os.FileInfo
	Err       error
	BuildInfo *buildinfo.BuildInfo
	Latest    *gomodule.Info
	Name      string
	ID        int
}

func (si *scanInfo) expired() bool {
	if si.Latest == nil || si.FileInfo == nil {
		return true
	}
	modTime := si.FileInfo.ModTime()
	return si.BuildInfo.Main.Version != develVersion && !modTime.IsZero() && si.Latest.Time.After(modTime)
}

const timeLayout = "2006-01-02 15:04:05"

func (si *scanInfo) String() string {
	const tpl = "%15s : %s\n"

	bs := &strings.Builder{}
	m := si.BuildInfo.Main
	fmt.Fprint(bs, color.GreenString("%3d %s\n", si.ID, si.Name))
	fmt.Fprintf(bs, tpl, "FilePath", si.BuildInfo.Path)
	fmt.Fprintf(bs, tpl, "Go", si.BuildInfo.GoVersion)

	var modTime time.Time
	if si.FileInfo != nil {
		modTime = si.FileInfo.ModTime()
		fmt.Fprintf(bs, tpl, "Install Time", modTime.Format(timeLayout))
	}

	fmt.Fprintf(bs, tpl, "Version", color.CyanString(m.Version))

	if si.Latest != nil {
		if si.Latest.Version == m.Version {
			fmt.Fprintf(bs, tpl, "Latest Version", color.CyanString(si.Latest.Version))
		} else {
			fmt.Fprintf(bs, tpl, "Latest Version", color.MagentaString(si.Latest.Version))
		}
		fmt.Fprintf(bs, tpl, "Latest Time", si.Latest.Time.Format(timeLayout))

		if m.Version != develVersion && !modTime.IsZero() && si.Latest.Time.After(modTime) {
			days := si.Latest.Time.Sub(modTime).Hours() / 24
			fmt.Fprint(bs, color.YellowString(tpl, "Expired", fmt.Sprintf("%.1f days", days)))
			fmt.Fprint(bs, color.YellowString(tpl, "Install", "go install "+si.BuildInfo.Path+"@latest"))
		}
	}
	if si.Err != nil {
		fmt.Fprint(bs, color.RedString(tpl, "Error", si.Err.Error()))
	}
	return bs.String()
}

func (si *scanInfo) JSON() string {
	info := map[string]any{
		"ID":       si.ID,
		"Ext":      si.Name,
		"FilePath": si.BuildInfo.Path,
		"Go":       si.BuildInfo.GoVersion,
		"Version":  si.BuildInfo.Main.Version,
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
