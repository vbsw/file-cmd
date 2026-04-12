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
	channel      chan tResult
	subPaths     []string
	resultsIdx   []int8
	resultsErr   []error
	absInputDir  string
	absOutputDir string
	threads      int
	offset       int
	delta        int
	rest         int
}

type tResult struct {
	err   error
	index int
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
	if command.err == nil {
		os.MkdirAll(command.outputDir, 0755)
		validateDir(command, command.outputDir, "output")
	}
}

func (proc *tProcess) fetchInputSubPaths(command *tCommand) {
	proc.ensureOutputDir(command)
	if command.err == nil {
		if command.recursive {
			proc.fetchInputSubPathsRecursive(command)
		} else {
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
						proc.subPaths = append(proc.subPaths, fileName)
					} else {
						subPath := path[inputDirLength:]
						proc.subPaths = append(proc.subPaths, subPath)
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
					proc.subPaths = append(proc.subPaths, fileName)
				}
				return nil
			}
		} else if !command.silent {
			fmt.Println("Warning:", err.Error())
		}
		return nil
	})
}

func (proc *tProcess) initOthers(threads int) {
	if len(proc.subPaths) < threads {
		proc.threads = len(proc.subPaths)
	} else {
		proc.threads = threads
	}
	proc.resultsIdx = make([]int8, len(proc.subPaths))
	proc.resultsErr = make([]error, len(proc.subPaths))
	proc.channel = make(chan tResult, max(min(len(proc.subPaths), 128), proc.threads*16))
	proc.delta = len(proc.subPaths) / proc.threads
	proc.rest = len(proc.subPaths) % proc.threads
}

func (proc *tProcess) step(i int) (int, int) {
	from := proc.offset
	if proc.rest == proc.threads-i {
		proc.delta++
		proc.rest = 0
	}
	proc.offset += proc.delta
	return from, proc.offset
}

func (proc *tProcess) fetchResultsFromChannel(i int) {
	for proc.resultsIdx[i] == 0 {
		rslt := <-proc.channel
		if rslt.index >= 0 {
			proc.resultsIdx[rslt.index] = 1
		} else {
			index := (rslt.index + 1) * -1
			proc.resultsIdx[index] = -1
			proc.resultsErr[index] = rslt.err
		}
	}
}

func (proc *tProcess) checkFileContentAll(inputDir string, from, to int, filter []string, maxFilterLength int) {
	//buffer := make([]byte, 8*1024*1024)
	for i := from; i < to && i < len(proc.subPaths); i++ {
		/*
			var hasContent bool
			var err error
			inputPath := filepath.Join(inputDir, subPaths[i])
			outputPath := filepath.Join(outputDir, subPaths[i])
			hasContent, err = fileHasContent(inputPath)
			if err == nil {
				if hasContent {
					outputParentDir := filepath.Dir(outputPath)
					if len(outputParentDir) != len(outputDir) {
						err = ensureDir(outputPath)
					}
					if err == nil {
						err = copyFile(inputPath, outputPath)
						if err == nil {
							channel <- tResult{nil, i}
						} else {
							channel <- tResult{err, i*-1-1}
						}
					} else {
						channel <- tResult{err, i*-1-1}
					}
				} else {
					channel <- tResult{nil, i*-1-1}
				}
			} else {
				channel <- tResult{err, i*-1-1}
			}
		*/
		proc.channel <- tResult{nil, i}
	}
}

func (proc *tProcess) checkFileContentAny(inputDir string, from, to int, filter []string, maxFilterLength int) {
	//buffer := make([]byte, 8*1024*1024)
	for i := from; i < to && i < len(proc.subPaths); i++ {
		/*
			var hasContent bool
			var err error
			inputPath := filepath.Join(inputDir, subPaths[i])
			outputPath := filepath.Join(outputDir, subPaths[i])
			hasContent, err = fileHasContent(inputPath)
			if err == nil {
				if hasContent {
					outputParentDir := filepath.Dir(outputPath)
					if len(outputParentDir) != len(outputDir) {
						err = ensureDir(outputPath)
					}
					if err == nil {
						err = copyFile(inputPath, outputPath)
						if err == nil {
							channel <- tResult{nil, i}
						} else {
							channel <- tResult{err, i*-1-1}
						}
					} else {
						channel <- tResult{err, i*-1-1}
					}
				} else {
					channel <- tResult{nil, i*-1-1}
				}
			} else {
				channel <- tResult{err, i*-1-1}
			}
		*/
		proc.channel <- tResult{nil, i}
	}
}

func ensureOutputDir(command *tCommand, absOutputDir, subPath string, checkedDirs *map[string]bool) bool {
	outputPath := filepath.Join(absOutputDir, subPath)
	outputDir := filepath.Dir(outputPath)
	outputDirAvail, inMap := true, true
	if len(outputDir) != len(absOutputDir) {
		outputDirAvail, inMap = (*checkedDirs)[outputDir]
		if !inMap {
			os.MkdirAll(outputDir, 0755)
			path := filepath.Join(command.outputDir, subPath)
			dir := filepath.Dir(path)
			validateDir(command, dir, "output")
			if command.err == nil {
				outputDirAvail = true
			} else if !command.silent {
				fmt.Println("Warning:", command.err.Error())
				command.err = nil
			}
			(*checkedDirs)[outputDir] = outputDirAvail
		}
	}
	return outputDirAvail
}

func maxStringLength(strings []string) int {
	var max int
	for _, str := range strings {
		if max < len(str) {
			max = len(str)
		}
	}
	return max
}
