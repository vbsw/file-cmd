/*
 *          Copyright 2026 Vitali Baumtrok.
 * Distributed under the Boost Software License, Version 1.0.
 *     (See accompanying file LICENSE or copy at
 *        http://www.boost.org/LICENSE_1_0.txt)
 */

package main

import (
	"errors"
	"github.com/vbsw/go-lib/cl"
	"os"
	"strconv"
)

const (
	idxInfoHelp = iota
	idxInfoVersion
	idxInfoExample
	idxInfoCopyright
	idxInfoTotal
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
	idxCmdTotal
)

const (
	idxOptSplitInput = iota
	idxOptSplitOutput
	idxOptSplitParts
	idxOptSplitBytes
	idxOptSplitLines
	idxOptSplitOvierwrite
	idxOptSplitTotal
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
	cmdTotal
)

type tCommand struct {
	err          error
	inputPath    string
	outputPath   string
	bytes        int64
	parts, lines int
	id           int
	overwrite    bool
}

func getCommand(osArgs []string) *tCommand {
	command := new(tCommand)
	if len(osArgs) > 0 {
		cmdLine := cl.New(osArgs, cl.NewDelimiter("", " ", "="))
		argsInfo := getInfoArgs(cmdLine)
		argsCmd := getCmdArgs(cmdLine)
		if argsInfo[idxInfoHelp].HasOrigin(0) {
			validateInfoHelp(command, argsInfo, argsCmd, cmdLine)
		} else if argsInfo[idxInfoVersion].HasOrigin(0) {
			validateInfoVersion(command, argsInfo, argsCmd, cmdLine)
		} else if argsInfo[idxInfoExample].HasOrigin(0) {
			validateInfoExample(command, argsInfo, argsCmd, cmdLine)
		} else if argsInfo[idxInfoCopyright].HasOrigin(0) {
			validateInfoCopyright(command, argsInfo, argsCmd, cmdLine)
		} else if argsCmd[idxCmdSplit].HasOrigin(0) {
			validateSplit(command, argsInfo, argsCmd, cmdLine)
		} else if argsCmd[idxCmdConcat].HasOrigin(0) {
			validateConcat(command, argsInfo, argsCmd, cmdLine)
		} else if argsCmd[idxCmdList].HasOrigin(0) {
			validateList(command, argsInfo, argsCmd, cmdLine)
		} else if argsCmd[idxCmdCount].HasOrigin(0) {
			validateCount(command, argsInfo, argsCmd, cmdLine)
		} else if argsCmd[idxCmdCopy].HasOrigin(0) {
			validateCopy(command, argsInfo, argsCmd, cmdLine)
		} else if argsCmd[idxCmdMove].HasOrigin(0) {
			validateMove(command, argsInfo, argsCmd, cmdLine)
		} else if argsCmd[idxCmdRemove].HasOrigin(0) {
			validateRemove(command, argsInfo, argsCmd, cmdLine)
		} else if argsCmd[idxCmdText].HasOrigin(0) {
			validateText(command, argsInfo, argsCmd, cmdLine)
		} else {
			command.err = errUnknownCommand(cmdLine.Arguments[0])
		}
		if command.err == nil && (command.id <= cmdNone || command.id >= cmdTotal) {
			command.err = errUnknownState()
		}
	} else {
		command.id = cmdInfo
	}
	return command
}

func getInfoArgs(cmdLine *cl.CommandLine) []*cl.Arguments {
	args := make([]*cl.Arguments, idxInfoTotal, idxInfoTotal)
	args[idxInfoHelp] = cmdLine.Search("-h", "--help", "-help", "help")
	args[idxInfoVersion] = cmdLine.Search("-v", "--version", "-version", "version")
	args[idxInfoExample] = cmdLine.Search("-e", "--example", "-example", "example")
	args[idxInfoCopyright] = cmdLine.Search("-c", "--copyright", "-copyright", "copyright")
	return args
}

func getCmdArgs(cmdLine *cl.CommandLine) []*cl.Arguments {
	args := make([]*cl.Arguments, idxCmdTotal, idxCmdTotal)
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

func validateInfoHelp(command *tCommand, argsInfo, argsCmd []*cl.Arguments, cmdLine *cl.CommandLine) {
	if len(cmdLine.Arguments) == 1 {
		command.id = cmdInfo
	} else {
		if argsInfo[idxInfoHelp].Count() == 1 {
			if len(cmdLine.Arguments) == 2 {
				if anyAvailable(argsCmd) {
					command.id = getCmdInfoId(argsCmd)
				} else {
					otherArgsAvail := argsInfo[idxInfoVersion].Available()
					otherArgsAvail = otherArgsAvail || argsInfo[idxInfoExample].Available()
					otherArgsAvail = otherArgsAvail || argsInfo[idxInfoCopyright].Available()
					if otherArgsAvail {
						command.err = errWrongArgumentUsage()
					} else {
						command.err = errUnknownCommand(cmdLine.Arguments[1])
					}
				}
			} else {
				command.err = errTooManyArguments()
			}
		} else {
			command.err = errMultipleUsage(argsInfo[idxInfoHelp].Keys)
		}
	}
}

func validateInfoVersion(command *tCommand, argsInfo, argsCmd []*cl.Arguments, cmdLine *cl.CommandLine) {
	if len(cmdLine.Arguments) == 1 {
		command.id = cmdVersion
	} else {
		if argsInfo[idxInfoVersion].Count() == 1 {
			command.err = errTooManyArguments()
		} else {
			command.err = errWrongArgumentUsage()
		}
	}
}

func validateInfoExample(command *tCommand, argsInfo, argsCmd []*cl.Arguments, cmdLine *cl.CommandLine) {
	if len(cmdLine.Arguments) == 1 {
		command.id = cmdInfoExample
	} else {
		if argsInfo[idxInfoExample].Count() == 1 {
			if len(cmdLine.Arguments) == 2 {
				if anyAvailable(argsCmd) {
					command.id = getCmdExampleId(argsCmd)
				} else {
					otherArgsAvail := argsInfo[idxInfoHelp].Available()
					otherArgsAvail = otherArgsAvail || argsInfo[idxInfoVersion].Available()
					otherArgsAvail = otherArgsAvail || argsInfo[idxInfoCopyright].Available()
					if otherArgsAvail {
						command.err = errWrongArgumentUsage()
					} else {
						command.err = errUnknownCommand(cmdLine.Arguments[1])
					}
				}
			} else {
				command.err = errTooManyArguments()
			}
		} else {
			command.err = errMultipleUsage(argsInfo[idxInfoExample].Keys)
		}
	}
}

func validateInfoCopyright(command *tCommand, argsInfo, argsCmd []*cl.Arguments, cmdLine *cl.CommandLine) {
	if len(cmdLine.Arguments) == 1 {
		command.id = cmdCopyright
	} else {
		if argsInfo[idxInfoCopyright].Count() == 1 {
			command.err = errTooManyArguments()
		} else {
			command.err = errMultipleUsage(argsInfo[idxInfoCopyright].Keys)
		}
	}
}

func validateSplit(command *tCommand, argsInfo, argsCmd []*cl.Arguments, cmdLine *cl.CommandLine) {
	if len(cmdLine.Arguments) == 1 {
		command.err = errNotEnoughArguments("split")
	} else {
		argsOpt := getOptSplitArgs(cmdLine)
		argMultiple := getMultiple(argsOpt)
		if argMultiple == nil {
			argUnm := cmdLine.Unmatched()
			validateInputOutput(command, argsOpt[idxOptSplitInput], argsOpt[idxOptSplitOutput], argUnm)
			validateInputFile(command)
			if command.outputPath == "" {
				command.outputPath = command.inputPath
			}
			validateOutputFile(command, ".0")
			command.parts, command.err = argValToInt(argsOpt[idxOptSplitParts], command.err)
			command.bytes, command.err = argValToBytes(argsOpt[idxOptSplitBytes], command.err)
			command.lines, command.err = argValToInt(argsOpt[idxOptSplitLines], command.err)
			command.overwrite = argsOpt[idxOptSplitOvierwrite].Available()
			command.id = cmdSplit
			if command.parts <= 0 && command.bytes <= 0 && command.lines <= 0 {
				command.parts = 2
			}
		} else {
			command.err = errMultipleUsage(argMultiple.Keys)
		}
	}
}

func argValToInt(arg *cl.Arguments, err error) (int, error) {
	if err == nil && arg.Available() {
		if arg.Values[0] > "" {
			n, errConv := strconv.Atoi(arg.Values[0])
			if errConv == nil {
				return n, nil
			}
			return 0, errArgNotInteger(arg.Keys[0], arg.Values[0])
		}
		return 0, errMissingArgValue(arg.Keys[0])
	}
	return 0, err
}

func argValToBytes(arg *cl.Arguments, err error) (int64, error) {
	if err == nil && arg.Available() {
		value := arg.Values[0]
		if value > "" {
			var bytes64 int64
			lastByte := value[len(value)-1]
			if lastByte == 'k' || lastByte == 'K' || lastByte == 'm' || lastByte == 'M' || lastByte == 'g' || lastByte == 'G' {
				value = value[:len(value)-1]
			} else {
				lastByte = 0
			}
			bytes, errConv := strconv.Atoi(value)
			if errConv == nil {
				switch lastByte {
				case 'k':
					bytes64 = int64(bytes) * 1000
				case 'K':
					bytes64 = int64(bytes) * 1024
				case 'm':
					bytes64 = int64(bytes) * 1000 * 1000
				case 'M':
					bytes64 = int64(bytes) * 1024 * 1024
				case 'g':
					bytes64 = int64(bytes) * 1000 * 1000 * 1000
				case 'G':
					bytes64 = int64(bytes) * 1024 * 1024 * 1024
				default:
					bytes64 = int64(bytes)
				}
				return bytes64, nil
			}
			return bytes64, errArgNotInteger(arg.Keys[0], arg.Values[0])
		}
		return 0, errMissingArgValue(arg.Keys[0])
	}
	return 0, err
}

func validateConcat(command *tCommand, argsInfo, argsCmd []*cl.Arguments, cmdLine *cl.CommandLine) {
}

func validateList(command *tCommand, argsInfo, argsCmd []*cl.Arguments, cmdLine *cl.CommandLine) {
}

func validateCount(command *tCommand, argsInfo, argsCmd []*cl.Arguments, cmdLine *cl.CommandLine) {
}

func validateCopy(command *tCommand, argsInfo, argsCmd []*cl.Arguments, cmdLine *cl.CommandLine) {
}

func validateMove(command *tCommand, argsInfo, argsCmd []*cl.Arguments, cmdLine *cl.CommandLine) {
}

func validateRemove(command *tCommand, argsInfo, argsCmd []*cl.Arguments, cmdLine *cl.CommandLine) {
}

func validateText(command *tCommand, argsInfo, argsCmd []*cl.Arguments, cmdLine *cl.CommandLine) {
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

func getOptSplitArgs(cmdLine *cl.CommandLine) []*cl.Arguments {
	args := make([]*cl.Arguments, idxOptSplitTotal, idxOptSplitTotal)
	args[idxOptSplitInput] = cmdLine.SearchByDelimiter("-i", "--input")
	args[idxOptSplitOutput] = cmdLine.SearchByDelimiter("-o", "--output")
	args[idxOptSplitParts] = cmdLine.SearchByDelimiter("-p")
	args[idxOptSplitBytes] = cmdLine.SearchByDelimiter("-b")
	args[idxOptSplitLines] = cmdLine.SearchByDelimiter("-l")
	args[idxOptSplitOvierwrite] = cmdLine.Search("-w")
	return args
}

func validateInputOutput(command *tCommand, argIn, argOut, argUnm *cl.Arguments) {
	if argIn.Available() {
		command.inputPath = argIn.Values[0]
		if argOut.Available() {
			if argUnm.Count() == 0 {
				command.outputPath = argOut.Values[0]
			} else {
				command.err = errUnknownArgument(argUnm.Keys[0])
			}
		} else if argUnm.Count() > 0 {
			if argUnm.Count() == 1 {
				command.outputPath = argUnm.Keys[0]
			} else {
				command.err = errUnknownArgument(argUnm.Keys[1])
			}
		}
	} else if argUnm.Count() > 0 {
		command.inputPath = argUnm.Keys[0]
		if argUnm.Count() > 1 {
			command.outputPath = argUnm.Keys[1]
			if argUnm.Count() > 2 {
				command.err = errUnknownArgument(argUnm.Keys[2])
			}
		}
	} else {
		command.err = errInputEmpty()
	}
}

func validateInputFile(command *tCommand) {
	if command.err == nil {
		if command.inputPath != "" {
			info, err := os.Stat(command.inputPath)
			if err == nil {
				if info == nil {
					command.err = errInputFileWrongPathSyntax(command.inputPath)
				} else if !info.Mode().IsRegular() {
					command.err = errInputFileNotAFile(command.inputPath)
				}
			} else if errors.Is(err, os.ErrNotExist) {
				command.err = errInputFileNotExist(command.inputPath)
			} else if errors.Is(err, os.ErrInvalid) {
				command.err = errInputFileWrongPathSyntax(command.inputPath)
			} else {
				command.err = errInputFileCantRead(command.inputPath)
			}
		} else {
			command.err = errInputEmpty()
		}
	}
}

func validateOutputFile(command *tCommand, suffix string) {
	if command.err == nil {
		outputPath := command.outputPath + suffix
		if outputPath > "" {
			_, err := os.Stat(outputPath)
			if err == nil || !errors.Is(err, os.ErrNotExist) {
				command.err = errOutputFileExists(command.outputPath)
			} else if errors.Is(err, os.ErrInvalid) {
				command.err = errOutputFileWrongPathSyntax(command.outputPath)
			} else if !errors.Is(err, os.ErrNotExist) {
				command.err = errOutputFileExists(command.outputPath)
			}
		} else {
			command.err = errOutputEmpty()
		}
	}
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

func getMultiple(argsList []*cl.Arguments) *cl.Arguments {
	for _, args := range argsList {
		if args.Count() > 1 {
			return args
		}
	}
	return nil
}
