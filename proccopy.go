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

func processCopy(command *tCommand) {
	var proc tProcess
	proc.initInputOutputDir(command)
	proc.ensureOutputDir(command)
	proc.fetchInputSubPaths(command)
	if command.err == nil && len(proc.subPaths) > 0 {
		checkedDirs := make(map[string]bool, 64)
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
				outputDirAvail, _ := ensureOutputSubDir(command, proc.absOutputDir, proc.subPaths[i], &checkedDirs)
				if outputDirAvail {
					copyFile(command, proc.subPaths[i])
				}
			} else if proc.resultsErr[i] != nil && command.verbose {
				printWarningError(command, proc.resultsErr[i])
			}
		}
	}
}

func copyFile(command *tCommand, subPath string) {
	inputPath := filepath.Join(command.inputDir, subPath)
	inputFile, err := os.Open(inputPath)
	if err == nil {
		var outputFile *os.File
		defer inputFile.Close()
		outputPath := filepath.Join(command.outputDir, subPath)
		outputFile, err = os.OpenFile(outputPath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0666)
		if err == nil {
			defer outputFile.Close()
			_, err = io.Copy(outputFile, inputFile)
			if err == nil {
				err = outputFile.Sync()
			}
		}
	}
	if err != nil {
		printWarningError(command, err)
	}
}
