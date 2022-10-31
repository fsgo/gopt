// Copyright(C) 2022 github.com/hidu  All Rights Reserved.
// Author: hidu <duv123@gmail.com>
// Date: 2022/3/19

package internal

import (
	"fmt"
)

const consoleColorTag = 0x1B

// ConsoleColor 字符颜色
type ConsoleColor func(string) string

// ConsoleRed 控制台红色字符
func ConsoleRed(txt string) string {
	return fmt.Sprintf("%c[31m%s%c[0m", consoleColorTag, txt, consoleColorTag)
}

// ConsoleGreen 控制台绿色字符
func ConsoleGreen(txt string) string {
	return fmt.Sprintf("%c[32m%s%c[0m", consoleColorTag, txt, consoleColorTag)
}

// ConsoleGrey 控制台灰色字符
func ConsoleGrey(txt string) string {
	return fmt.Sprintf("%c[90m%s%c[0m", consoleColorTag, txt, consoleColorTag)
}
