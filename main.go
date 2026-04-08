/*
 *          Copyright 2026 Vitali Baumtrok.
 * Distributed under the Boost Software License, Version 1.0.
 *     (See accompanying file LICENSE or copy at
 *        http://www.boost.org/LICENSE_1_0.txt)
 */

// Package main executes various cmds on files.
package main

import (
	"fmt"
	"os"
)

func main() {
	command := getCommand(os.Args[1:])
	if command.err == nil {
		switch command.id {
		case cmdInfo:
			printInfo()
		case cmdInfoExample:
			printInfoExample()
		case cmdInfoSplit:
			printInfoSplit()
		case cmdInfoConcat:
			printInfoConcat()
		case cmdInfoList:
			printInfoList(command.str)
		case cmdInfoCount:
			printInfoCount()
		case cmdInfoCopy:
			printInfoCopy()
		case cmdInfoMove:
			printInfoMove()
		case cmdInfoRemove:
			printInfoRemove()
		case cmdInfoText:
			printInfoText()
		case cmdVersion:
			printVersion()
		case cmdExampleSplit:
			printExampleSplit()
		case cmdExampleConcat:
			printExampleConcat()
		case cmdExampleList:
			printExampleList()
		case cmdExampleCount:
			printExampleCount()
		case cmdExampleCopy:
			printExampleCopy()
		case cmdExampleMove:
			printExampleMove()
		case cmdExampleRemove:
			printExampleRemove()
		case cmdExampleText:
			printExampleText()
		case cmdCopyright:
			printCopyright()
		case cmdSplit:
			processSplit(command)
		case cmdConcat:
			processConcat(command)
		case cmdList:
			processList(command)
		case cmdCount:
			processCount(command)
		case cmdCopy:
			processCopy(command)
		case cmdMove:
			processMove(command)
		case cmdRemove:
			processRemove(command)
		case cmdText:
			processText(command)
		}
	}
	if command.err != nil {
		fmt.Println("error:", command.err.Error())
	}
}
