/*
 *          Copyright 2026 Vitali Baumtrok.
 * Distributed under the Boost Software License, Version 1.0.
 *     (See accompanying file LICENSE or copy at
 *        http://www.boost.org/LICENSE_1_0.txt)
 */

package main

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
)

func processCopy(command *tCommand) {
	var proc tProcess
	proc.initInputOutputDir(command)
	proc.fetchInputSubPaths(command)
	if command.err == nil && len(proc.subPaths) > 0 {
		checkedDirs := make(map[string]bool, 64)
		proc.initOthers(command.threads, command.contentFilter)
		if command.or {
			for i := 0; i < proc.threads; i++ {
				from, to := proc.step(i)
				go proc.checkFileContentAny(command.inputDir, from, to)
			}
		} else {
			for i := 0; i < proc.threads; i++ {
				from, to := proc.step(i)
				go proc.checkFileContentAll(command.inputDir, from, to)
			}
		}
		for i := 0; i < len(proc.subPaths); i++ {
			proc.fetchResultsFromChannel(i)
			if proc.resultsIdx[i] == 1 {
				outputDirAvail := ensureOutputDir(command, proc.absOutputDir, proc.subPaths[i], &checkedDirs)
				if outputDirAvail {
					copyFile(command, proc.subPaths[i])
				}
			} else if proc.resultsErr[i] != nil && !command.silent {
				fmt.Println("Warning:", proc.resultsErr[i].Error())
			}
		}
	}
}

func copyFile(command *tCommand, subPath string) {
	inputPath := filepath.Join(command.inputDir, subPath)
	inputFile, inputErr := os.Open(inputPath)
	if inputErr == nil {
		outputPath := filepath.Join(command.outputDir, subPath)
		outputFile, outputErr := os.OpenFile(outputPath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0666)
		if outputErr == nil {
			_, copyErr := io.Copy(outputFile, inputFile)
			if copyErr == nil {
				outputFile.Sync()
			} else if !command.silent {
				fmt.Println("Warning:", copyErr.Error())
			}
			outputFile.Close()
		} else if !command.silent {
			fmt.Println("Warning:", outputErr.Error())
		}
		inputFile.Close()
	} else if !command.silent {
		fmt.Println("Warning:", inputErr.Error())
	}
}

func fileHasContent() (bool, error) {
	return true, nil
}
