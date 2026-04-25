/*
 *          Copyright 2026 Vitali Baumtrok.
 * Distributed under the Boost Software License, Version 1.0.
 *     (See accompanying file LICENSE or copy at
 *        http://www.boost.org/LICENSE_1_0.txt)
 */

package main

import (
	"fmt"
	"github.com/vbsw/go-lib/fs"
)

func processConcat(command *tCommand) {
	var file fs.FileWriter
	defer func() {
		file.Close()
		command.err = file.Err
	}()
	if file.Open(command.outputPath) {
		inputPath := command.concatBase + fmt.Sprintf(".%0*d", command.concatNumLen, command.concatNumBegin)
		for i := command.concatNumBegin + 1; file.Err == nil && fs.IsExist(inputPath); i++ {
			file.CopyFrom(inputPath)
			inputPath = command.concatBase + fmt.Sprintf(".%0*d", command.concatNumLen, i)
		}
	}
}
