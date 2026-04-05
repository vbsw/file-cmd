/*
 *          Copyright 2026 Vitali Baumtrok.
 * Distributed under the Boost Software License, Version 1.0.
 *     (See accompanying file LICENSE or copy at
 *        http://www.boost.org/LICENSE_1_0.txt)
 */

package main

import (
	"github.com/vbsw/go-lib/cl"
)

const (
	idxInfoHelp = iota
	idxInfoVersion
	idxInfoExample
	idxInfoCopyright
)

const (
	idxCmdSplit = iota
	idxCmdConcat
	idxCmdList
	idxCmdCount
	idxCmdCopy
	idxCmdMove
	idxCmdRemove
	idxCmdText
)

const (
	cmdNone = iota
	cmdInfo
	cmdInfoExample
	cmdInfoSplit
	cmdInfoConcat
	cmdInfoList
	cmdInfoCount
	cmdInfoCopy
	cmdInfoMove
	cmdInfoRemove
	cmdInfoText
	cmdVersion
	cmdExampleSplit
	cmdExampleConcat
	cmdExampleList
	cmdExampleCount
	cmdExampleCopy
	cmdExampleMove
	cmdExampleRemove
	cmdExampleText
	cmdCopyright
	cmdSplit
	cmdConcat
	cmdList
	cmdCount
	cmdCopy
	cmdMove
	cmdRemove
	cmdText
)

type tCommand struct {
	err error
	id  int
}

func getCommand(osArgs []string) *tCommand {
	var command *tCommand
	if len(osArgs) > 0 {
		command = new(tCommand)
		cmdLine := cl.New(osArgs)
		argsInfo := getInfoParams(cmdLine)
		argsCmd := getCmdParams(cmdLine)
		validateInfo(command, argsInfo, argsCmd, cmdLine)
	}
	return command
}

func getInfoParams(cmdLine *cl.CommandLine) []*cl.Arguments {
	args := make([]*cl.Arguments, 4, 4)
	args[idxInfoHelp] = cmdLine.Search("-h", "--help", "-help", "help")
	args[idxInfoVersion] = cmdLine.Search("-v", "--version", "-version", "version")
	args[idxInfoExample] = cmdLine.Search("-e", "--example", "-example", "example")
	args[idxInfoCopyright] = cmdLine.Search("--copyright", "-copyright", "copyright")
	return args
}

func getCmdParams(cmdLine *cl.CommandLine) []*cl.Arguments {
	args := make([]*cl.Arguments, 8, 8)
	args[idxCmdSplit] = cmdLine.Search("split")
	args[idxCmdConcat] = cmdLine.Search("concat")
	args[idxCmdList] = cmdLine.Search("ls")
	args[idxCmdCount] = cmdLine.Search("count")
	args[idxCmdCopy] = cmdLine.Search("cp")
	args[idxCmdMove] = cmdLine.Search("mv")
	args[idxCmdRemove] = cmdLine.Search("rm")
	args[idxCmdText] = cmdLine.Search("text")
	return args
}

func validateInfo(command *tCommand, argsInfo []*cl.Arguments, argsCmd []*cl.Arguments, cmdLine *cl.CommandLine) {
	if command.err == nil && anyAvailable(argsInfo) {
		if argsInfo[idxInfoHelp].Available() {
			validateInfoHelp(command, argsInfo, argsCmd, cmdLine)
		} else if argsInfo[idxInfoVersion].Available() {
			validateInfoVersion(command, argsInfo, argsCmd, cmdLine)
		} else if argsInfo[idxInfoExample].Available() {
			validateInfoExample(command, argsInfo, argsCmd, cmdLine)
		} else if argsInfo[idxInfoCopyright].Available() {
			validateInfoCopyright(command, argsInfo, argsCmd, cmdLine)
		}
	}
}

func validateInfoHelp(command *tCommand, argsInfo []*cl.Arguments, argsCmd []*cl.Arguments, cmdLine *cl.CommandLine) {
	if len(cmdLine.Arguments) == 1 {
		command.id = cmdInfo
	} else if anyMatches(argsInfo[idxInfoHelp].Keys, cmdLine.Arguments[0]) {
		if argsInfo[idxInfoHelp].Count() == 1 {
			if len(cmdLine.Arguments) == 2 {
				if anyAvailable(argsCmd) {
					command.id = getCmdInfoId(argsCmd)
				} else if anyMultiple(argsInfo) {
					command.err = errWrongArgumentUsage()
				} else {
					command.err = errUnknownCommand(cmdLine.Arguments[1])
				}
			} else {
				command.err = errTooManyArguments()
			}
		} else {
			command.err = errWrongArgumentUsage()
		}
	} // else: could be a command
}

func validateInfoVersion(command *tCommand, argsInfo []*cl.Arguments, argsCmd []*cl.Arguments, cmdLine *cl.CommandLine) {
	if len(cmdLine.Arguments) == 1 {
		command.id = cmdVersion
	} else if anyMatches(argsInfo[idxInfoVersion].Keys, cmdLine.Arguments[0]) {
		if argsInfo[idxInfoVersion].Count() == 1 {
			command.err = errTooManyArguments()
		} else {
			command.err = errWrongArgumentUsage()
		}
	} // else: could be a command
}

func validateInfoExample(command *tCommand, argsInfo []*cl.Arguments, argsCmd []*cl.Arguments, cmdLine *cl.CommandLine) {
	if len(cmdLine.Arguments) == 1 {
		command.id = cmdInfoExample
	} else if anyMatches(argsInfo[idxInfoExample].Keys, cmdLine.Arguments[0]) {
		if argsInfo[idxInfoExample].Count() == 1 {
			if len(cmdLine.Arguments) == 2 {
				if anyAvailable(argsCmd) {
					command.id = getCmdExampleId(argsCmd)
				} else if anyMultiple(argsInfo) {
					command.err = errWrongArgumentUsage()
				} else {
					command.err = errUnknownCommand(cmdLine.Arguments[1])
				}
			} else {
				command.err = errTooManyArguments()
			}
		} else {
			command.err = errWrongArgumentUsage()
		}
	} // else: could be a command
}

func validateInfoCopyright(command *tCommand, argsInfo []*cl.Arguments, argsCmd []*cl.Arguments, cmdLine *cl.CommandLine) {
	if len(cmdLine.Arguments) == 1 {
		command.id = cmdCopyright
	} else if anyMatches(argsInfo[idxInfoCopyright].Keys, cmdLine.Arguments[0]) {
		if argsInfo[idxInfoCopyright].Count() == 1 {
			command.err = errTooManyArguments()
		} else {
			command.err = errWrongArgumentUsage()
		}
	} // else: could be a command
}

func getCmdInfoId(argsCmd []*cl.Arguments) int {
	cmdId := cmdNone
	if argsCmd[idxCmdSplit].Available() {
		cmdId = cmdInfoSplit
	} else if argsCmd[idxCmdConcat].Available() {
		cmdId = cmdInfoConcat
	} else if argsCmd[idxCmdList].Available() {
		cmdId = cmdInfoList
	} else if argsCmd[idxCmdCount].Available() {
		cmdId = cmdInfoCount
	} else if argsCmd[idxCmdCopy].Available() {
		cmdId = cmdInfoCopy
	} else if argsCmd[idxCmdMove].Available() {
		cmdId = cmdInfoMove
	} else if argsCmd[idxCmdRemove].Available() {
		cmdId = cmdInfoRemove
	} else if argsCmd[idxCmdText].Available() {
		cmdId = cmdInfoText
	}
	return cmdId
}

func getCmdExampleId(argsCmd []*cl.Arguments) int {
	cmdId := cmdNone
	if argsCmd[idxCmdSplit].Available() {
		cmdId = cmdExampleSplit
	} else if argsCmd[idxCmdConcat].Available() {
		cmdId = cmdExampleConcat
	} else if argsCmd[idxCmdList].Available() {
		cmdId = cmdExampleList
	} else if argsCmd[idxCmdCount].Available() {
		cmdId = cmdExampleCount
	} else if argsCmd[idxCmdCopy].Available() {
		cmdId = cmdExampleCopy
	} else if argsCmd[idxCmdMove].Available() {
		cmdId = cmdExampleMove
	} else if argsCmd[idxCmdRemove].Available() {
		cmdId = cmdExampleRemove
	} else if argsCmd[idxCmdText].Available() {
		cmdId = cmdExampleText
	}
	return cmdId
}

func anyAvailable(argsList []*cl.Arguments) bool {
	for _, args := range argsList {
		if args.Available() {
			return true
		}
	}
	return false
}

func anyMatches(args []string, arg string) bool {
	for _, a := range args {
		if a == arg {
			return true
		}
	}
	return false
}

func anyMultiple(argsList []*cl.Arguments) bool {
	for _, args := range argsList {
		if args.Count() > 1 {
			return true
		}
	}
	return false
}
