/*
 *          Copyright 2026 Vitali Baumtrok.
 * Distributed under the Boost Software License, Version 1.0.
 *     (See accompanying file LICENSE or copy at
 *        http://www.boost.org/LICENSE_1_0.txt)
 */

package main

import (
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
		checkedDirs := make(map[string]bool, 64)
		for _, subPath := range proc.subPaths {
			inputPath := filepath.Join(command.inputDir, subPath)
			inputDir := filepath.Dir(inputPath)
			if _, inMap := checkedDirs[inputDir]; !inMap {
				checkedDirs[inputDir] = false
			}
			if file.IsEmpty(inputPath) {
				err := os.Remove(inputPath)
				if err != nil && !command.silent {
					printWarning(command, err)
				}
			}
		}
		for key, _ := range checkedDirs {
			os.Remove(key)
		}
	}
}
