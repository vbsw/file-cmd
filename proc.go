/*
 *          Copyright 2026 Vitali Baumtrok.
 * Distributed under the Boost Software License, Version 1.0.
 *     (See accompanying file LICENSE or copy at
 *        http://www.boost.org/LICENSE_1_0.txt)
 */

package main

import (
	"github.com/vbsw/go-lib/fs"
	"github.com/vbsw/go-lib/match"
	"os"
	"path/filepath"
)

const procInitialBufferSize = 8 * 1024 * 1024

type tProcess struct {
	channel      chan tResult
	subPaths     []string
	subDirs      []string
	resultsIdx   []int8
	resultsErr   []error
	filter       []string
	absInputDir  string
	absOutputDir string
	threads      int
	offset       int
	delta        int
	rest         int
	filterMin    int
	filterMax    int
	bufferSize   int
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
			if info != nil && info.Mode().IsRegular() {
				fileName := info.Name()
				if match.WildcardMatch(fileName, command.fileNameFilter) {
					if inputDirLength == len(path)-len(fileName) {
						proc.subPaths = append(proc.subPaths, fileName)
					} else {
						subPath := path[inputDirLength:]
						proc.subPaths = append(proc.subPaths, subPath)
					}
				}
				return nil
			}
		} else if command.verbose {
			printWarningError(command, err)
		}
		return nil
	})
}

func (proc *tProcess) fetchInputSubPathsFlat(command *tCommand) {
	inputDirLength := len(proc.absInputDir) + 1
	command.err = filepath.Walk(proc.absInputDir, func(path string, info os.FileInfo, err error) error {
		if err == nil {
			if info != nil && info.Mode().IsRegular() {
				fileName := info.Name()
				if inputDirLength == len(path)-len(fileName) {
					if match.WildcardMatch(fileName, command.fileNameFilter) {
						proc.subPaths = append(proc.subPaths, fileName)
					}
					return nil
				}
			}
		} else if command.verbose {
			printWarningError(command, err)
		}
		return nil
	})
}

func (proc *tProcess) initOthers(threads int, filter []string) {
	if len(proc.subPaths) < threads {
		proc.threads = len(proc.subPaths)
	} else {
		proc.threads = threads
	}
	proc.resultsIdx = make([]int8, len(proc.subPaths))
	proc.resultsErr = make([]error, len(proc.subPaths))
	proc.channel = make(chan tResult, min(len(proc.subPaths), max(128, proc.threads*16)))
	proc.delta = len(proc.subPaths) / proc.threads
	proc.rest = len(proc.subPaths) % proc.threads
	proc.filter = filter
	proc.filterMin, proc.filterMax = getBounds(filter)
	if proc.filterMax <= procInitialBufferSize {
		proc.bufferSize = procInitialBufferSize
	} else {
		proc.bufferSize = proc.filterMax + procInitialBufferSize
	}
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

func (proc *tProcess) checkContainsAll(inputDir string, from, to int) {
	var reader fs.FileReader
	reader.Buffer = make([]byte, proc.bufferSize)
	for i := from; i < to && i < len(proc.subPaths); i++ {
		checking := true
		inputPath := filepath.Join(inputDir, proc.subPaths[i])
		if reader.Open(inputPath) {
			checking = reader.Read(0) && reader.NRead >= proc.filterMax
			for checking {
				if match.Contains(reader.Buffer[:reader.NRead], proc.filter, match.And) {
					break
				} else {
					checking = reader.Read(proc.filterMax - 1)
				}
			}
			reader.Close()
		}
		if checking {
			proc.channel <- tResult{reader.Err, i}
		} else {
			proc.channel <- tResult{reader.Err, i*-1 - 1}
		}
	}
}

func (proc *tProcess) checkContainsAny(inputDir string, from, to int) {
	var reader fs.FileReader
	reader.Buffer = make([]byte, proc.bufferSize)
	for i := from; i < to && i < len(proc.subPaths); i++ {
		checking := true
		inputPath := filepath.Join(inputDir, proc.subPaths[i])
		if reader.Open(inputPath) {
			checking = reader.Read(0) && reader.NRead >= proc.filterMin
			for checking {
				if match.Contains(reader.Buffer[:reader.NRead], proc.filter, match.Or) {
					break
				} else {
					checking = reader.Read(proc.filterMax - 1)
				}
			}
			reader.Close()
		}
		if checking {
			proc.channel <- tResult{reader.Err, i}
		} else {
			proc.channel <- tResult{reader.Err, i*-1 - 1}
		}
	}
}

func ensureOutputSubDir(command *tCommand, absOutputDir, subPath string, checkedDirs *map[string]bool) (bool, bool) {
	outputPath := filepath.Join(absOutputDir, subPath)
	outputDir := filepath.Dir(outputPath)
	outputDirAvail, inMap, subDirRemovable := true, true, false
	if len(outputDir) != len(absOutputDir) {
		outputDirAvail, inMap = (*checkedDirs)[outputDir]
		if !inMap {
			os.MkdirAll(outputDir, 0755)
			path := filepath.Join(command.outputDir, subPath)
			dir := filepath.Dir(path)
			validateDir(command, dir, "output")
			if command.err == nil {
				outputDirAvail = true
			} else {
				if command.verbose {
					printWarningError(command, command.err)
				}
				command.err = nil
			}
			(*checkedDirs)[outputDir] = outputDirAvail
			subDirRemovable = outputDirAvail
		}
	}
	return outputDirAvail, subDirRemovable
}

func getBounds(strings []string) (int, int) {
	var min, max int
	for _, str := range strings {
		length := len(str)
		if max < length {
			max = length
		}
		if min > length {
			min = length
		}
	}
	return min, max
}
