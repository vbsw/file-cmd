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

type tProcess struct {
	subDirs      map[string]struct{}
	absInputDir  string
	absOutputDir string
}

func (proc *tProcess) initInputOutputDir(command *tCommand) {
	if command.err == nil {
		proc.absInputDir, command.err = filepath.Abs(command.inputDir)
		if command.err == nil {
			proc.absOutputDir, command.err = filepath.Abs(command.outputDir)
		}
	}
}

func (proc *tProcess) ensureOutputDir(command *tCommand) {
	err := os.MkdirAll(command.outputDir, 0755)
	if err != nil {
		validateDir("output", command.outputDir, command)
	}
}

func (proc *tProcess) ensureOutputSubDir(command *tCommand, subDir string) {
	_, checked := proc.subDirs[subDir]
	if !checked {
		dir := filepath.Join(command.outputDir, subDir)
		err := os.MkdirAll(dir, 0755)
		if err != nil {
			validateDir("output", dir, command)
		}
		proc.subDirs[subDir] = struct{}{}
	}
}

func (proc *tProcess) fetchInputSubPaths(command *tCommand) {
	if command.err == nil {
		if command.recursive {
			proc.subDirs = make(map[string]struct{}, 64)
			proc.ensureOutputDir(command)
			proc.fetchInputSubPathsRecursive(command)
		} else {
			proc.ensureOutputDir(command)
			proc.fetchInputSubPathsFlat(command)
		}
	}
}

func (proc *tProcess) fetchInputSubPathsRecursive(command *tCommand) {
	inputDirLength := len(proc.absInputDir) + 1
	command.err = filepath.Walk(proc.absInputDir, func(path string, info os.FileInfo, err error) error {
		if err == nil {
			if info != nil && !info.IsDir() {
				fileName := info.Name()
				if match.WildcardMatch(command.fileNameFilter, fileName) {
					if inputDirLength == len(path)-len(fileName) {
						proc.addFile(fileName)
					} else {
						dir := filepath.Dir(path)
						subDir := dir[inputDirLength:]
						proc.ensureOutputSubDir(command, subDir)
						if command.err == nil {
							subDirPath := path[inputDirLength:]
							proc.addFile(subDirPath)
						} else {
							if !command.silent {
								fmt.Println("Warning:", command.err.Error())
							}
							command.err = nil
						}
					}
				}
				return nil
			}
		} else if !command.silent {
			fmt.Println("Warning:", err.Error())
		}
		return nil
	})
}

func (proc *tProcess) fetchInputSubPathsFlat(command *tCommand) {
	inputDirLength := len(proc.absInputDir) + 1
	command.err = filepath.Walk(proc.absInputDir, func(path string, info os.FileInfo, err error) error {
		if err == nil {
			fileName := info.Name()
			if info != nil && !info.IsDir() && inputDirLength == len(path)-len(fileName) {
				if match.WildcardMatch(command.fileNameFilter, fileName) {
					proc.addFile(fileName)
				}
				return nil
			}
		} else if !command.silent {
			fmt.Println("Warning:", err.Error())
		}
		return nil
	})
}

func (proc *tProcess) addFile(subFilePath string) {
	fmt.Println(subFilePath)
}
