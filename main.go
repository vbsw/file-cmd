/*
 *          Copyright 2026 Vitali Baumtrok.
 * Distributed under the Boost Software License, Version 1.0.
 *     (See accompanying file LICENSE or copy at
 *        http://www.boost.org/LICENSE_1_0.txt)
 */

// Package main executes various cmds on files.
package main

import (
	"fmt"
	"os"
)

func main() {
	command := getCommand(os.Args[1:])
	if command.err == nil {
		switch command.id {
		case cmdNone:
			fmt.Println("error: unknown state")
		case cmdInfo:
			fmt.Println("use --help")
		}
	} else {
		fmt.Println("error:", command.err.Error())
	}
}
