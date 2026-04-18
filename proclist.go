/*
 *          Copyright 2026 Vitali Baumtrok.
 * Distributed under the Boost Software License, Version 1.0.
 *     (See accompanying file LICENSE or copy at
 *        http://www.boost.org/LICENSE_1_0.txt)
 */

package main

import (
	"fmt"
)

func processList(command *tCommand) {
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
					fmt.Println(proc.subPaths[i])
				} else if proc.resultsErr[i] != nil && !command.silent {
					fmt.Println("Warning:", proc.resultsErr[i].Error())
				}
			}
		}
	}
}
