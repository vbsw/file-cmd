/*
 *          Copyright 2026 Vitali Baumtrok.
 * Distributed under the Boost Software License, Version 1.0.
 *     (See accompanying file LICENSE or copy at
 *        http://www.boost.org/LICENSE_1_0.txt)
 */

package main

import (
	"os"
	"path/filepath"
)

func processRemove(command *tCommand) {
	var proc tProcess
	proc.initInputOutputDir(command)
	proc.fetchInputSubPaths(command)
	if command.err == nil {
		if len(proc.subPaths) > 0 {
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
					inputPath := filepath.Join(command.inputDir, proc.subPaths[i])
					err := os.Remove(inputPath)
					if err != nil && command.verbose {
						printWarning(command, err)
					}
				} else if proc.resultsErr[i] != nil && command.verbose {
					printWarning(command, proc.resultsErr[i])
				}
			}
		}
	}
}
