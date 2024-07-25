// Copyright(C) 2024 github.com/fsgo  All Rights Reserved.
// Author: hidu <duv123@gmail.com>
// Date: 2024/7/13

package internal

import (
	"context"
	_ "embed"
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/fatih/color"
	"github.com/fsgo/cmdutil"
)

func Install(args []string) error {
	u := &installer{}
	if err := u.setup(args); err != nil {
		return err
	}
	return u.Run()
}

type installer struct {
	flags   *flag.FlagSet
	timeout time.Duration
}

func (i *installer) setup(args []string) error {
	i.flags = flag.NewFlagSet("install", flag.ExitOnError)
	i.flags.DurationVar(&i.timeout, "t", 120*time.Second, `install timeout`)
	return i.flags.Parse(args)
}

func (i *installer) Run() error {
	args := i.flags.Args()
	if len(args) == 0 {
		return errors.New("what app name to install? e.g. 'dlv'")
	}
	var failedNames []string
	for _, name := range args {
		if err := i.installOne(name); err != nil {
			log.Println(color.RedString(name), err)
		}
	}
	if len(failedNames) == 0 {
		return nil
	}
	return fmt.Errorf("install %q failed", failedNames)
}

func (i *installer) getTimeout() time.Duration {
	if i.timeout > 0 {
		return i.timeout
	}
	return 120 * time.Second
}

func (i *installer) installOne(name string) error {
	info, ok := installInfos[name]
	if !ok {
		return fmt.Errorf("%q not found", name)
	}
	ctx, cancel := context.WithTimeout(context.Background(), i.getTimeout())
	defer cancel()
	args := []string{"install", info.Path + "@latest"}
	if info.Tags != "" {
		args = append(args, "--tags", info.Tags)
	}
	cmd := newGoCommand(ctx, args...)
	oe := &cmdutil.OSEnv{}
	oe.WithEnviron(cmd.Environ())
	if info.CGO {
		oe.Set("CGO_ENABLED", "1")
	}
	cmd.Env = oe.Environ()
	log.Println("will install:", cmd.String())
	return cmd.Run()
}

var installInfos = map[string]appInfo{}

type appInfo struct {
	Name string
	Path string
	Tags string
	CGO  bool
}

//go:embed db/app.jsonl
var appDB string

func init() {
	lines := strings.Split(appDB, "\n")
	for lineNo, content := range lines {
		content = strings.TrimSpace(content)
		if content == "" {
			continue
		}
		_, str, found := strings.Cut(content, " ")
		if !found {
			panic(fmt.Errorf("invalid line[%d] %q", lineNo, content))
		}
		info := &appInfo{}
		if err := json.Unmarshal([]byte(str), info); err != nil {
			panic(fmt.Errorf("invalid line[%d] %q, json decode failed: %w", lineNo, str, err))
		}
		installInfos[info.Name] = *info
	}
}
