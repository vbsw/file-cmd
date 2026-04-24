/*
 *          Copyright 2026 Vitali Baumtrok.
 * Distributed under the Boost Software License, Version 1.0.
 *     (See accompanying file LICENSE or copy at
 *        http://www.boost.org/LICENSE_1_0.txt)
 */

package main

import (
	"fmt"
	"github.com/vbsw/go-lib/fs"
	"os"
	"path/filepath"
)

func processClean(command *tCommand) {
	var proc tProcess
	var file fs.File
	proc.initInputOutputDir(command)
	proc.fetchInputSubPaths(command)
	if command.err == nil && len(proc.subPaths) > 0 {
		for _, subPath := range proc.subPaths {
			inputPath := filepath.Join(command.inputDir, subPath)
			if file.IsEmpty(inputPath) {
				err := os.Remove(inputPath)
				if err == nil {
					if command.list {
						fmt.Println(subPath)
					}
					if command.all {
						dir := filepath.Dir(subPath)
						if dir != "." {
							err = os.Remove(filepath.Join(command.inputDir, dir))
							if err == nil && command.list {
								fmt.Println(dir)
							}
						}
					}
				} else if command.verbose {
					printWarningError(command, err)
				}
			}
		}
	}
}
