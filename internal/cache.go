// Copyright(C) 2022 github.com/fsgo  All Rights Reserved.
// Author: hidu <duv123@gmail.com>
// Date: 2022/12/3

package internal

import (
	"bytes"
	"crypto/md5"
	"fmt"
	"net"
)

var keyInterfaces string

func init() {
	keyInterfaces, _ = netIfa()
}

func netIfa() (string, error) {
	ifas, err := net.Interfaces()
	if err != nil {
		return "", err
	}
	var b bytes.Buffer
	for _, ifa := range ifas {
		b.WriteString(ifa.Name)
		b.WriteString(ifa.Flags.String())

		addrs, _ := ifa.Addrs()
		for _, addr := range addrs {
			b.WriteString(addr.String())
		}
		b.WriteString("\n")
	}
	m5 := md5.New()
	m5.Write(b.Bytes())
	return fmt.Sprintf("%x", m5.Sum(nil)), nil
}
