/*
 *          Copyright 2026 Vitali Baumtrok.
 * Distributed under the Boost Software License, Version 1.0.
 *     (See accompanying file LICENSE or copy at
 *        http://www.boost.org/LICENSE_1_0.txt)
 */

package main

import (
	"io"
	"os"
	"path/filepath"
)

func processMove(command *tCommand) {
	var proc tProcess
	proc.initInputOutputDir(command)
	proc.ensureOutputDir(command)
	proc.fetchInputSubPaths(command)
	if command.err == nil && len(proc.subPaths) > 0 {
		checkedDirs := make(map[string]bool, 64)
		//subDirs := make([]string, 64)
		proc.initOthers(command.threads, command.contentFilter)
		if command.or {
			for i := 0; i < proc.threads; i++ {
				from, to := proc.step(i)
				go proc.checkContainsAny(command.inputDir, from, to)
			}
		} else {
			for i := 0; i < proc.threads; i++ {
				from, to := proc.step(i)
				go proc.checkContainsAll(command.inputDir, from, to)
			}
		}
		for i := 0; i < len(proc.subPaths); i++ {
			proc.fetchResultsFromChannel(i)
			if proc.resultsIdx[i] == 1 {
				//outputDirAvail, subRem := ensureOutputSubDir(command, proc.absOutputDir, proc.subPaths[i], &checkedDirs)
				outputDirAvail, _ := ensureOutputSubDir(command, proc.absOutputDir, proc.subPaths[i], &checkedDirs)
				if outputDirAvail {
					moveFile(command, proc.subPaths[i])
					//if subRem {
					//	subDirs = append(subDirs, filepath.Dir(proc.subPaths[i]))
					//}
				}
			} else if proc.resultsErr[i] != nil && command.verbose {
				printWarningError(command, proc.resultsErr[i])
			}
		}
		//removeEmptyInputDir(proc.absInputDir, subDirs)
	}
}

func moveFile(command *tCommand, subPath string) {
	var inputFile, outputFile *os.File
	inputPath := filepath.Join(command.inputDir, subPath)
	outputPath := filepath.Join(command.outputDir, subPath)
	err := os.Rename(inputPath, outputPath)
	if err != nil {
		inputFile, err = os.Open(inputPath)
		if err == nil {
			defer inputFile.Close()
			outputFile, err = os.OpenFile(outputPath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0666)
			if err == nil {
				defer outputFile.Close()
				_, err = io.Copy(outputFile, inputFile)
				if err == nil {
					err = outputFile.Sync()
					if err == nil {
						err = os.Remove(inputPath)
					}
				}
			}
		}
		if err != nil && command.verbose {
			printWarningError(command, err)
		}
	}
}

func removeEmptyInputDir(inputDir string, subDirs []string) {
	for i := len(subDirs) - 1; i >= 0; i-- {
		inputDir := filepath.Join(inputDir, subDirs[i])
		_ = os.Remove(inputDir)
	}
}
