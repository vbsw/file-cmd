/*
 *          Copyright 2026 Vitali Baumtrok.
 * Distributed under the Boost Software License, Version 1.0.
 *     (See accompanying file LICENSE or copy at
 *        http://www.boost.org/LICENSE_1_0.txt)
 */

package main

import (
	"fmt"
	"github.com/vbsw/go-lib/match"
	"os"
	"path/filepath"
)

func processCopy(command *tCommand) {
	var proc tProcess
	proc.initInputOutputDir(command)
	proc.fetchInputSubPaths(command)
}

func iterateCopyRecursiveGo(command *tCommand, inputDir, outputDir string, threads int) {
}

func iterateCopyRecursive(command *tCommand, inputDir, outputDir string) {
	//buffer := make([]byte, 8*1024*1024)
	err := filepath.Walk(inputDir, func(path string, info os.FileInfo, err error) error {
		if err == nil {
			if info != nil && !info.IsDir() {
				if match.WildcardMatch(command.fileNameFilter, info.Name()) {
					fmt.Println(info.Name())
				}
				return nil
			}
		} else if !command.silent {
			fmt.Println("Warning:", err.Error())
		}
		return nil
	})
	command.err = err
}

func iterateCopyFlatGo(command *tCommand, inputDir, outputDir string, threads int) {
}

func iterateCopyFlat(command *tCommand, inputDir, outputDir string) {
	inputDirLength := len(inputDir) + 1
	//buffer := make([]byte, 8*1024*1024)
	err := filepath.Walk(inputDir, func(path string, info os.FileInfo, err error) error {
		if err == nil {
			if info != nil && !info.IsDir() && inputDirLength == len(path)-len(info.Name()) {
				if match.WildcardMatch(command.fileNameFilter, info.Name()) {
					fmt.Println(info.Name())
				}
				return nil
			}
		} else if !command.silent {
			fmt.Println("Warning:", err.Error())
		}
		return nil
	})
	command.err = err
}
