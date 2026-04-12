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
	"path/filepath"
	"strconv"
	"strings"
)

type tCommand struct {
	bytes          int64
	err            error
	inputPath      string
	outputPath     string
	str            string
	delimiter      string
	inputDir       string
	outputDir      string
	terminator     string
	fileNameFilter string
	contentFilter  []string
	parts, lines   int
	id             int
	threads        int
	overwrite      bool
	or             bool
	recursive      bool
	silent         bool
	lettersOnly    bool
	upperCase      bool
	lowerCase      bool
}

func getCommand(osArgs []string) *tCommand {
	command := new(tCommand)
	if len(osArgs) > 0 {
		var argsList [idxInfoTotal + idxCmdTotal]*cl.Arguments
		infoArgsList, cmdArgsList := argsList[:idxInfoTotal], argsList[idxInfoTotal:len(argsList)]
		cmdLine := cl.New(osArgs, cl.NewDelimiter("", " ", "="))
		command.threads = 1
		readCmdArgs(cmdArgsList, cmdLine)   // args having values
		readInfoArgs(infoArgsList, cmdLine) // args without values
		validateInfo(command, infoArgsList, cmdArgsList, cmdLine)
		validateCmd(command, infoArgsList, cmdArgsList, cmdLine)
		if command.err == nil && (command.id <= cmdNone || command.id >= cmdTotal) {
			command.err = errUnknownState()
		}
	} else {
		command.id = cmdInfo
	}
	return command
}

func readInfoArgs(argsList []*cl.Arguments, cmdLine *cl.CommandLine) {
	argsList[idxInfoHelp] = cmdLine.Match("-h", "--help", "-help", "help")
	argsList[idxInfoVersion] = cmdLine.Match("-v", "--version", "-version", "version")
	argsList[idxInfoExample] = cmdLine.Match("-e", "--example", "-example", "example")
	argsList[idxInfoCopyright] = cmdLine.Match("-c", "--copyright", "-copyright", "copyright")
}

func readCmdArgs(argsList []*cl.Arguments, cmdLine *cl.CommandLine) {
	argsList[idxCmdSplit] = cmdLine.Match("split")
	argsList[idxCmdConcat] = cmdLine.Match("concat")
	argsList[idxCmdList] = cmdLine.Match("ls", "list", "print")
	argsList[idxCmdCount] = cmdLine.Match("count")
	argsList[idxCmdCopy] = cmdLine.Match("cp", "copy")
	argsList[idxCmdMove] = cmdLine.Match("mv", "move")
	argsList[idxCmdRemove] = cmdLine.Match("rm", "remove")
	argsList[idxCmdText] = cmdLine.Match("text")
}

func validateInfo(command *tCommand, infoArgsList, cmdArgsList []*cl.Arguments, cmdLine *cl.CommandLine) {
	if infoArgsList[idxInfoHelp].HasIndex(0) {
		validateInfoHelp(command, infoArgsList, cmdArgsList, cmdLine)
	} else if infoArgsList[idxInfoVersion].HasIndex(0) {
		validateInfoVersion(command, infoArgsList, cmdArgsList, cmdLine)
	} else if infoArgsList[idxInfoExample].HasIndex(0) {
		validateInfoExample(command, infoArgsList, cmdArgsList, cmdLine)
	} else if infoArgsList[idxInfoCopyright].HasIndex(0) {
		validateInfoCopyright(command, infoArgsList, cmdArgsList, cmdLine)
	} // else: cound be a command
}

func validateCmd(command *tCommand, infoArgsList, cmdArgsList []*cl.Arguments, cmdLine *cl.CommandLine) {
	if command.err == nil && command.id == cmdNone {
		cmdLine.RevertMatched(infoArgsList...)
		if cmdArgsList[idxCmdSplit].HasIndex(0) {
			validateSplit(command, cmdArgsList[idxCmdSplit], cmdLine)
		} else if cmdArgsList[idxCmdConcat].HasIndex(0) {
			validateConcat(command, cmdArgsList[idxCmdConcat], cmdLine)
		} else if cmdArgsList[idxCmdList].HasIndex(0) {
			validateList(command, cmdArgsList[idxCmdList], cmdLine)
		} else if cmdArgsList[idxCmdCount].HasIndex(0) {
			validateCount(command, cmdArgsList[idxCmdCount], cmdLine)
		} else if cmdArgsList[idxCmdCopy].HasIndex(0) {
			validateCopy(command, cmdArgsList[idxCmdCopy], cmdLine)
		} else if cmdArgsList[idxCmdMove].HasIndex(0) {
			validateMove(command, cmdArgsList[idxCmdMove], cmdLine)
		} else if cmdArgsList[idxCmdRemove].HasIndex(0) {
			validateRemove(command, cmdArgsList[idxCmdRemove], cmdLine)
		} else if cmdArgsList[idxCmdText].HasIndex(0) {
			validateText(command, cmdArgsList[idxCmdText], cmdLine)
		} else {
			command.err = errUnknownCommand(cmdLine.Arguments[0])
		}
	}
}

func validateInfoHelp(command *tCommand, infoArgsList, cmdArgsList []*cl.Arguments, cmdLine *cl.CommandLine) {
	if len(cmdLine.Arguments) == 1 {
		command.id = cmdInfo
	} else {
		if infoArgsList[idxInfoHelp].Count() == 1 {
			if len(cmdLine.Arguments) == 2 {
				cmdArgs := getAvailable(cmdArgsList)
				if cmdArgs == nil {
					otherInfoAvail := infoArgsList[idxInfoVersion].Available()
					otherInfoAvail = otherInfoAvail || infoArgsList[idxInfoExample].Available()
					otherInfoAvail = otherInfoAvail || infoArgsList[idxInfoCopyright].Available()
					if otherInfoAvail {
						command.err = errWrongArgumentUsage()
					} else {
						command.err = errUnknownCommand(cmdLine.Arguments[1])
					}
				} else {
					command.id = getCmdInfoId(cmdArgsList)
					command.str = cmdArgs.Keys[0]
				}
			} else {
				command.err = errTooManyArguments()
			}
		} else {
			command.err = errMultipleUsage(infoArgsList[idxInfoHelp].Keys)
		}
	}
}

func validateInfoVersion(command *tCommand, infoArgsList, cmdArgsList []*cl.Arguments, cmdLine *cl.CommandLine) {
	if len(cmdLine.Arguments) == 1 {
		command.id = cmdVersion
	} else {
		if infoArgsList[idxInfoVersion].Count() == 1 {
			command.err = errTooManyArguments()
		} else {
			command.err = errWrongArgumentUsage()
		}
	}
}

func validateInfoExample(command *tCommand, infoArgsList, cmdArgsList []*cl.Arguments, cmdLine *cl.CommandLine) {
	if len(cmdLine.Arguments) == 1 {
		command.id = cmdInfoExample
	} else {
		if infoArgsList[idxInfoExample].Count() == 1 {
			if len(cmdLine.Arguments) == 2 {
				if anyAvailable(cmdArgsList) {
					command.id = getCmdExampleId(cmdArgsList)
				} else {
					otherInfoAvail := infoArgsList[idxInfoHelp].Available()
					otherInfoAvail = otherInfoAvail || infoArgsList[idxInfoVersion].Available()
					otherInfoAvail = otherInfoAvail || infoArgsList[idxInfoCopyright].Available()
					if otherInfoAvail {
						command.err = errWrongArgumentUsage()
					} else {
						command.err = errUnknownCommand(cmdLine.Arguments[1])
					}
				}
			} else {
				command.err = errTooManyArguments()
			}
		} else {
			command.err = errMultipleUsage(infoArgsList[idxInfoExample].Keys)
		}
	}
}

func validateInfoCopyright(command *tCommand, infoArgsList, cmdArgsList []*cl.Arguments, cmdLine *cl.CommandLine) {
	if len(cmdLine.Arguments) == 1 {
		command.id = cmdCopyright
	} else {
		if infoArgsList[idxInfoCopyright].Count() == 1 {
			command.err = errTooManyArguments()
		} else {
			command.err = errMultipleUsage(infoArgsList[idxInfoCopyright].Keys)
		}
	}
}

func validateSplit(command *tCommand, cmdArgs *cl.Arguments, cmdLine *cl.CommandLine) {
	if len(cmdLine.Arguments) == 1 {
		command.err = errNotEnoughArguments("split")
	} else {
		optArgsList := getOptSplitArgs(cmdLine)
		argMultiple := getMultiple(optArgsList)
		if argMultiple == nil {
			if cmdArgs.Count() > 1 {
				cmdLine.RevertMatched(cmdArgs)
				cmdLine.Matched[0] = true
				cmdLine.Matched[len(cmdLine.Arguments)] = false
			}
			unmatchedArgs := cmdLine.Unmatched()
			validateInputOutput(command, optArgsList[idxOptSplitInput], optArgsList[idxOptSplitOutput], unmatchedArgs)
			validateInputFile(command)
			if command.outputPath == "" {
				command.outputPath = command.inputPath
			}
			validateOutputFile(command, ".0")
			command.parts, command.err = argValToInt(optArgsList[idxOptSplitParts], 0, command.err)
			command.bytes, command.err = argValToBytes(optArgsList[idxOptSplitBytes], command.err)
			command.lines, command.err = argValToInt(optArgsList[idxOptSplitLines], 0, command.err)
			command.overwrite = optArgsList[idxOptSplitOverwrite].Available()
			command.id = cmdSplit
			if command.parts <= 0 && command.bytes <= 0 && command.lines <= 0 {
				command.parts = 2
			}
		} else {
			command.err = errMultipleUsage(argMultiple.Keys)
		}
	}
}

func validateConcat(command *tCommand, cmdArgs *cl.Arguments, cmdLine *cl.CommandLine) {
	if len(cmdLine.Arguments) == 1 {
		command.err = errNotEnoughArguments("concat")
	} else {
		optArgsList := getOptConcatArgs(cmdLine)
		argMultiple := getMultiple(optArgsList)
		if argMultiple == nil {
			if cmdArgs.Count() > 1 {
				cmdLine.RevertMatched(cmdArgs)
				cmdLine.Matched[0] = true
			}
			unmatchedArgs := cmdLine.Unmatched()
			validateInputOutput(command, optArgsList[idxOptConcatInput], optArgsList[idxOptConcatOutput], unmatchedArgs)
			validateInputConcatFile(command)
			validateOutputConcatFile(command)
			command.overwrite = optArgsList[idxOptConcatOverwrite].Available()
			command.id = cmdConcat
		} else {
			command.err = errMultipleUsage(argMultiple.Keys)
		}
	}
}

func validateList(command *tCommand, cmdArgs *cl.Arguments, cmdLine *cl.CommandLine) {
	if len(cmdLine.Arguments) == 1 {
		command.err = errNotEnoughArguments(cmdArgs.Keys[0])
	} else {
		optArgsList := getOptListArgs(cmdLine)
		argMultiple := getMultiple(optArgsList)
		if argMultiple == nil {
			if cmdArgs.Count() > 1 {
				cmdLine.RevertMatched(cmdArgs)
				cmdLine.Matched[0] = true
			}
			unmatchedArgs := cmdLine.Unmatched()
			validateInputDirFile(command, optArgsList[idxOptListInput], unmatchedArgs, cmdLine)
			command.or = optArgsList[idxOptListOr].Available()
			command.recursive = optArgsList[idxOptListRecursive].Available()
			command.silent = optArgsList[idxOptListSilent].Available()
			command.contentFilter = getFilter(optArgsList[idxOptListFilter], optArgsList[idxOptListDelimiter], unmatchedArgs, cmdLine)
			command.threads, command.err = getThreads(optArgsList[idxOptListThreads], command.err)
			command.id = cmdList
		} else {
			command.err = errMultipleUsage(argMultiple.Keys)
		}
	}
}

func validateCount(command *tCommand, cmdArgs *cl.Arguments, cmdLine *cl.CommandLine) {
	if len(cmdLine.Arguments) == 1 {
		command.err = errNotEnoughArguments(cmdArgs.Keys[0])
	} else {
		optArgsList := getOptCountArgs(cmdLine)
		argMultiple := getMultiple(optArgsList)
		if argMultiple == nil {
			if cmdArgs.Count() > 1 {
				cmdLine.RevertMatched(cmdArgs)
				cmdLine.Matched[0] = true
			}
			unmatchedArgs := cmdLine.Unmatched()
			validateInputDirFile(command, optArgsList[idxOptCountInput], unmatchedArgs, cmdLine)
			command.or = optArgsList[idxOptCountOr].Available()
			command.recursive = optArgsList[idxOptCountRecursive].Available()
			command.silent = optArgsList[idxOptCountSilent].Available()
			command.contentFilter = getFilter(optArgsList[idxOptCountFilter], optArgsList[idxOptCountDelimiter], unmatchedArgs, cmdLine)
			command.threads, command.err = getThreads(optArgsList[idxOptCountThreads], command.err)
			command.id = cmdCount
		} else {
			command.err = errMultipleUsage(argMultiple.Keys)
		}
	}
}

func validateCopy(command *tCommand, cmdArgs *cl.Arguments, cmdLine *cl.CommandLine) {
	if len(cmdLine.Arguments) == 1 {
		command.err = errNotEnoughArguments(cmdArgs.Keys[0])
	} else {
		optArgsList := getOptCopyArgs(cmdLine)
		argMultiple := getMultiple(optArgsList)
		if argMultiple == nil {
			if cmdArgs.Count() > 1 {
				cmdLine.RevertMatched(cmdArgs)
				cmdLine.Matched[0] = true
			}
			unmatchedArgs := cmdLine.Unmatched()
			validateInputDirFile(command, optArgsList[idxOptCopyInput], unmatchedArgs, cmdLine)
			validateOutputDir(command, optArgsList[idxOptCopyInput], unmatchedArgs, cmdLine)
			command.or = optArgsList[idxOptCopyOr].Available()
			command.overwrite = optArgsList[idxOptCopyOverwrite].Available()
			command.recursive = optArgsList[idxOptCopyRecursive].Available()
			command.silent = optArgsList[idxOptCopySilent].Available()
			command.contentFilter = getFilter(optArgsList[idxOptCopyFilter], optArgsList[idxOptCopyDelimiter], unmatchedArgs, cmdLine)
			command.threads, command.err = getThreads(optArgsList[idxOptCopyThreads], command.err)
			command.id = cmdCopy
		} else {
			command.err = errMultipleUsage(argMultiple.Keys)
		}
	}
}

func validateMove(command *tCommand, cmdArgs *cl.Arguments, cmdLine *cl.CommandLine) {
	if len(cmdLine.Arguments) == 1 {
		command.err = errNotEnoughArguments(cmdArgs.Keys[0])
	} else {
		optArgsList := getOptMoveArgs(cmdLine)
		argMultiple := getMultiple(optArgsList)
		if argMultiple == nil {
			if cmdArgs.Count() > 1 {
				cmdLine.RevertMatched(cmdArgs)
				cmdLine.Matched[0] = true
			}
			unmatchedArgs := cmdLine.Unmatched()
			validateInputDirFile(command, optArgsList[idxOptMoveInput], unmatchedArgs, cmdLine)
			validateOutputDir(command, optArgsList[idxOptMoveInput], unmatchedArgs, cmdLine)
			command.or = optArgsList[idxOptMoveOr].Available()
			command.overwrite = optArgsList[idxOptMoveOverwrite].Available()
			command.recursive = optArgsList[idxOptMoveRecursive].Available()
			command.silent = optArgsList[idxOptMoveSilent].Available()
			command.contentFilter = getFilter(optArgsList[idxOptMoveFilter], optArgsList[idxOptMoveDelimiter], unmatchedArgs, cmdLine)
			command.threads, command.err = getThreads(optArgsList[idxOptMoveThreads], command.err)
			command.id = cmdMove
		} else {
			command.err = errMultipleUsage(argMultiple.Keys)
		}
	}
}

func validateRemove(command *tCommand, cmdArgs *cl.Arguments, cmdLine *cl.CommandLine) {
	if len(cmdLine.Arguments) == 1 {
		command.err = errNotEnoughArguments(cmdArgs.Keys[0])
	} else {
		optArgsList := getOptRemoveArgs(cmdLine)
		argMultiple := getMultiple(optArgsList)
		if argMultiple == nil {
			if cmdArgs.Count() > 1 {
				cmdLine.RevertMatched(cmdArgs)
				cmdLine.Matched[0] = true
			}
			unmatchedArgs := cmdLine.Unmatched()
			validateInputDirFile(command, optArgsList[idxOptRemoveInput], unmatchedArgs, cmdLine)
			command.or = optArgsList[idxOptRemoveOr].Available()
			command.recursive = optArgsList[idxOptRemoveRecursive].Available()
			command.silent = optArgsList[idxOptRemoveSilent].Available()
			command.contentFilter = getFilter(optArgsList[idxOptRemoveFilter], optArgsList[idxOptRemoveDelimiter], unmatchedArgs, cmdLine)
			command.threads, command.err = getThreads(optArgsList[idxOptRemoveThreads], command.err)
			command.id = cmdRemove
		} else {
			command.err = errMultipleUsage(argMultiple.Keys)
		}
	}
}

func validateText(command *tCommand, cmdArgs *cl.Arguments, cmdLine *cl.CommandLine) {
	if len(cmdLine.Arguments) == 1 {
		command.err = errNotEnoughArguments(cmdArgs.Keys[0])
	} else {
		optArgsList := getOptTextArgs(cmdLine)
		argMultiple := getMultiple(optArgsList)
		if argMultiple == nil {
			if cmdArgs.Count() > 1 {
				cmdLine.RevertMatched(cmdArgs)
				cmdLine.Matched[0] = true
			}
			unmatchedArgs := cmdLine.Unmatched()
			validateOutputTextFile(command, optArgsList[idxOptTextOutput], unmatchedArgs, cmdLine)
			validateFormat(command, optArgsList[idxOptTextFormat])
			command.overwrite = optArgsList[idxOptTextOverwrite].Available()
			command.bytes, command.err = argValToBytes(optArgsList[idxOptTextSize], command.err)
			command.terminator = optArgsList[idxOptTextTerminator].ValueAt(0, "")
			command.delimiter = optArgsList[idxOptTextDelimiter].ValueAt(0, "")
			command.id = cmdText
		} else {
			command.err = errMultipleUsage(argMultiple.Keys)
		}
	}
}

func validateFormat(command *tCommand, args *cl.Arguments) {
	if command.err == nil && args.Available() {
		format := args.Values[0]
		if format > "" {
			for i := 0; i < len(format) && command.err == nil; i++ {
				b := format[i]
				if b == 'a' {
					command.lettersOnly = true
				} else if b == 'u' {
					command.upperCase = true
				} else if b == 'l' {
					command.lowerCase = true
				} else {
					command.err = errUnknownFormat(string(b))
				}
			}
			if command.upperCase && command.lowerCase {
				command.err = errIncompatibleFormats("u", "l")
			}
		}
	}
}

func getCmdInfoId(cmdArgsList []*cl.Arguments) int {
	cmdId := cmdNone
	if cmdArgsList[idxCmdSplit].Available() {
		cmdId = cmdInfoSplit
	} else if cmdArgsList[idxCmdConcat].Available() {
		cmdId = cmdInfoConcat
	} else if cmdArgsList[idxCmdList].Available() {
		cmdId = cmdInfoList
	} else if cmdArgsList[idxCmdCount].Available() {
		cmdId = cmdInfoCount
	} else if cmdArgsList[idxCmdCopy].Available() {
		cmdId = cmdInfoCopy
	} else if cmdArgsList[idxCmdMove].Available() {
		cmdId = cmdInfoMove
	} else if cmdArgsList[idxCmdRemove].Available() {
		cmdId = cmdInfoRemove
	} else if cmdArgsList[idxCmdText].Available() {
		cmdId = cmdInfoText
	}
	return cmdId
}

func getCmdExampleId(cmdArgsList []*cl.Arguments) int {
	cmdId := cmdNone
	if cmdArgsList[idxCmdSplit].Available() {
		cmdId = cmdExampleSplit
	} else if cmdArgsList[idxCmdConcat].Available() {
		cmdId = cmdExampleConcat
	} else if cmdArgsList[idxCmdList].Available() {
		cmdId = cmdExampleList
	} else if cmdArgsList[idxCmdCount].Available() {
		cmdId = cmdExampleCount
	} else if cmdArgsList[idxCmdCopy].Available() {
		cmdId = cmdExampleCopy
	} else if cmdArgsList[idxCmdMove].Available() {
		cmdId = cmdExampleMove
	} else if cmdArgsList[idxCmdRemove].Available() {
		cmdId = cmdExampleRemove
	} else if cmdArgsList[idxCmdText].Available() {
		cmdId = cmdExampleText
	}
	return cmdId
}

func getOptSplitArgs(cmdLine *cl.CommandLine) []*cl.Arguments {
	argsList := make([]*cl.Arguments, idxOptSplitTotal, idxOptSplitTotal)
	argsList[idxOptSplitParts] = cmdLine.MatchDelimited("-p")
	argsList[idxOptSplitBytes] = cmdLine.MatchDelimited("-b")
	argsList[idxOptSplitLines] = cmdLine.MatchDelimited("-l")
	argsList[idxOptSplitOverwrite] = cmdLine.Match("-w")
	argsList[idxOptSplitInput] = cmdLine.MatchDelimited("-i", "--input")
	argsList[idxOptSplitOutput] = cmdLine.MatchDelimited("-o", "--output")
	return argsList
}

func getOptConcatArgs(cmdLine *cl.CommandLine) []*cl.Arguments {
	argsList := make([]*cl.Arguments, idxOptConcatTotal, idxOptConcatTotal)
	argsList[idxOptConcatOverwrite] = cmdLine.Match("-w")
	argsList[idxOptConcatInput] = cmdLine.MatchDelimited("-i", "--input")
	argsList[idxOptConcatOutput] = cmdLine.MatchDelimited("-o", "--output")
	return argsList
}

func getOptListArgs(cmdLine *cl.CommandLine) []*cl.Arguments {
	argsList := make([]*cl.Arguments, idxOptListTotal, idxOptListTotal)
	argsList[idxOptListOr] = cmdLine.Match("--or")
	argsList[idxOptListRecursive] = cmdLine.Match("-r", "--recursive")
	argsList[idxOptListSilent] = cmdLine.Match("-s", "--silent")
	argsList[idxOptListInput] = cmdLine.MatchDelimited("-i", "--input")
	argsList[idxOptListFilter] = cmdLine.MatchDelimited("-f")
	argsList[idxOptListDelimiter] = cmdLine.MatchDelimited("-d")
	argsList[idxOptListThreads] = cmdLine.MatchDelimited("-t")
	return argsList
}

func getOptCountArgs(cmdLine *cl.CommandLine) []*cl.Arguments {
	argsList := make([]*cl.Arguments, idxOptCountTotal, idxOptCountTotal)
	argsList[idxOptCountOr] = cmdLine.Match("--or")
	argsList[idxOptCountRecursive] = cmdLine.Match("-r", "--recursive")
	argsList[idxOptCountSilent] = cmdLine.Match("-s", "--silent")
	argsList[idxOptCountInput] = cmdLine.MatchDelimited("-i", "--input")
	argsList[idxOptCountFilter] = cmdLine.MatchDelimited("-f")
	argsList[idxOptCountDelimiter] = cmdLine.MatchDelimited("-d")
	argsList[idxOptCountThreads] = cmdLine.MatchDelimited("-t")
	return argsList
}

func getOptCopyArgs(cmdLine *cl.CommandLine) []*cl.Arguments {
	argsList := make([]*cl.Arguments, idxOptCopyTotal, idxOptCopyTotal)
	argsList[idxOptCopyOr] = cmdLine.Match("--or")
	argsList[idxOptCopyOverwrite] = cmdLine.Match("-w")
	argsList[idxOptCopyRecursive] = cmdLine.Match("-r", "--recursive")
	argsList[idxOptCopySilent] = cmdLine.Match("-s", "--silent")
	argsList[idxOptCopyInput] = cmdLine.MatchDelimited("-i", "--input")
	argsList[idxOptCopyOutput] = cmdLine.MatchDelimited("-o", "--output")
	argsList[idxOptCopyFilter] = cmdLine.MatchDelimited("-f")
	argsList[idxOptCopyDelimiter] = cmdLine.MatchDelimited("-d")
	argsList[idxOptCopyThreads] = cmdLine.MatchDelimited("-t")
	return argsList
}

func getOptMoveArgs(cmdLine *cl.CommandLine) []*cl.Arguments {
	argsList := make([]*cl.Arguments, idxOptMoveTotal, idxOptMoveTotal)
	argsList[idxOptMoveOr] = cmdLine.Match("--or")
	argsList[idxOptMoveOverwrite] = cmdLine.Match("-w")
	argsList[idxOptMoveRecursive] = cmdLine.Match("-r", "--recursive")
	argsList[idxOptMoveSilent] = cmdLine.Match("-s", "--silent")
	argsList[idxOptMoveInput] = cmdLine.MatchDelimited("-i", "--input")
	argsList[idxOptMoveOutput] = cmdLine.MatchDelimited("-o", "--output")
	argsList[idxOptMoveFilter] = cmdLine.MatchDelimited("-f")
	argsList[idxOptMoveDelimiter] = cmdLine.MatchDelimited("-d")
	argsList[idxOptMoveThreads] = cmdLine.MatchDelimited("-t")
	return argsList
}

func getOptRemoveArgs(cmdLine *cl.CommandLine) []*cl.Arguments {
	argsList := make([]*cl.Arguments, idxOptRemoveTotal, idxOptRemoveTotal)
	argsList[idxOptRemoveOr] = cmdLine.Match("--or")
	argsList[idxOptRemoveRecursive] = cmdLine.Match("-r", "--recursive")
	argsList[idxOptRemoveSilent] = cmdLine.Match("-s", "--silent")
	argsList[idxOptRemoveInput] = cmdLine.MatchDelimited("-i", "--input")
	argsList[idxOptRemoveFilter] = cmdLine.MatchDelimited("-f")
	argsList[idxOptRemoveDelimiter] = cmdLine.MatchDelimited("-d")
	argsList[idxOptRemoveThreads] = cmdLine.MatchDelimited("-t")
	return argsList
}

func getOptTextArgs(cmdLine *cl.CommandLine) []*cl.Arguments {
	argsList := make([]*cl.Arguments, idxOptTextTotal, idxOptTextTotal)
	argsList[idxOptTextOverwrite] = cmdLine.Match("-w")
	argsList[idxOptTextOutput] = cmdLine.MatchDelimited("-o", "--output")
	argsList[idxOptTextSize] = cmdLine.MatchDelimited("-s")
	argsList[idxOptTextTerminator] = cmdLine.MatchDelimited("-e")
	argsList[idxOptTextDelimiter] = cmdLine.MatchDelimited("-d")
	argsList[idxOptTextFormat] = cmdLine.MatchDelimited("-f")
	return argsList
}

func validateInputOutput(command *tCommand, inputArgs, outputArgs, unmatchedArgs *cl.Arguments) {
	if inputArgs.Available() {
		command.inputPath = inputArgs.Values[0]
		if outputArgs.Available() {
			if unmatchedArgs.Count() == 0 {
				command.outputPath = outputArgs.Values[0]
			} else {
				command.err = errUnknownArgument(unmatchedArgs.Keys[0])
			}
		} else if unmatchedArgs.Count() > 0 {
			if unmatchedArgs.Count() == 1 {
				command.outputPath = unmatchedArgs.Keys[0]
			} else {
				command.err = errUnknownArgument(unmatchedArgs.Keys[1])
			}
		}
	} else if unmatchedArgs.Count() > 0 {
		command.inputPath = unmatchedArgs.Keys[0]
		if unmatchedArgs.Count() > 1 {
			command.outputPath = unmatchedArgs.Keys[1]
			if unmatchedArgs.Count() > 2 {
				command.err = errUnknownArgument(unmatchedArgs.Keys[2])
			}
		}
	} else {
		command.err = errPathEmpty("input")
	}
}

func validateInputFile(command *tCommand) {
	if command.err == nil {
		if command.inputPath > "" {
			info, err := os.Stat(command.inputPath)
			if err == nil {
				if info == nil {
					command.err = errFileWrongPathSyntax("input", command.inputPath)
				} else if !info.Mode().IsRegular() {
					command.err = errFileNotAFile("input", command.inputPath)
				}
			} else if errors.Is(err, os.ErrNotExist) {
				command.err = errFileNotExist("input", command.inputPath)
			} else if errors.Is(err, os.ErrInvalid) {
				command.err = errFileWrongPathSyntax("input", command.inputPath)
			} else {
				command.err = errFileCantRead("input", command.inputPath)
			}
		} else {
			command.err = errPathEmpty("input")
		}
	}
}

func validateOutputFile(command *tCommand, suffix string) {
	if command.err == nil {
		outputPath := command.outputPath + suffix
		if outputPath > "" {
			_, err := os.Stat(outputPath)
			if err == nil || !errors.Is(err, os.ErrNotExist) {
				command.err = errFileExists("output", command.outputPath)
			} else if errors.Is(err, os.ErrInvalid) {
				command.err = errFileWrongPathSyntax("output", command.outputPath)
			}
		} else {
			command.err = errPathEmpty("output")
		}
	}
}

func validateInputConcatFile(command *tCommand) {
	if command.err == nil {
		if command.inputPath > "" {
			command.inputPath = command.inputPath + ".0"
			validateInputFile(command)
			if command.err != nil {
				command.err = nil
				command.inputPath = command.inputPath[:len(command.inputPath)-2]
				validateInputFile(command)
			}
		} else {
			command.err = errPathEmpty("input")
		}
	}
}

func validateOutputConcatFile(command *tCommand) {
	if command.outputPath == "" {
		command.outputPath = command.inputPath
	}
	if command.outputPath > "" && command.outputPath == command.inputPath {
		command.outputPath = command.outputPath[:len(command.outputPath)-2]
	}
	validateOutputFile(command, "")
}

func validateOutputTextFile(command *tCommand, outputArgs, unmatchedArgs *cl.Arguments, cmdLine *cl.CommandLine) {
	if outputArgs.Available() {
		command.outputPath = outputArgs.Values[0]
	} else if unmatchedArgs.Count() > 0 {
		command.outputPath = unmatchedArgs.Keys[0]
		if unmatchedArgs.Count() == 1 {
			cmdLine.Matched[unmatchedArgs.Indices[0]] = true
			cmdLine.Matched[len(cmdLine.Arguments)] = true
		} else {
			command.err = errUnknownArgument(unmatchedArgs.Keys[2])
		}
	} else {
		command.err = errPathEmpty("output")
	}
	validateOutputFile(command, "")
}

func validateInputDirFile(command *tCommand, inputArgs, unmatchedArgs *cl.Arguments, cmdLine *cl.CommandLine) {
	if inputArgs.Available() {
		command.inputPath = inputArgs.Values[0]
	} else {
		unmatchedArgIdx := getUnmatchedIdx(unmatchedArgs, cmdLine.Matched)
		if unmatchedArgIdx >= 0 {
			command.inputPath = unmatchedArgs.Keys[unmatchedArgIdx]
			cmdLine.Matched[unmatchedArgs.Indices[unmatchedArgIdx]] = true
		}
	}
	if command.err == nil {
		if command.inputPath > "" {
			command.inputDir = command.inputPath
			validateDir(command, command.inputDir, "input")
			if command.err != nil {
				fileName := filepath.Base(command.inputPath)
				command.err = nil
				command.inputDir = filepath.Dir(command.inputPath)
				if fileName > "" && fileName != "*" {
					command.fileNameFilter = fileName
				}
				validateDir(command, command.inputDir, "input")
			}
		} else {
			command.err = errPathEmpty("input")
		}
	}
}

func validateOutputDir(command *tCommand, outputArgs, unmatchedArgs *cl.Arguments, cmdLine *cl.CommandLine) {
	if outputArgs.Available() {
		command.outputPath = outputArgs.Values[0]
	} else {
		unmatchedArgIdx := getUnmatchedIdx(unmatchedArgs, cmdLine.Matched)
		if unmatchedArgIdx >= 0 {
			command.outputPath = unmatchedArgs.Keys[unmatchedArgIdx]
			cmdLine.Matched[unmatchedArgs.Indices[unmatchedArgIdx]] = true
		}
	}
	if command.err == nil {
		if command.outputPath > "" {
			command.outputDir = command.outputPath
			//validateDir(command, command.outputDir, "output")
		} else {
			command.err = errPathEmpty("output")
		}
	}
}

func validateDir(command *tCommand, dir, io string) {
	if command.err == nil {
		if dir > "" {
			info, err := os.Stat(dir)
			if err == nil {
				if info == nil {
					command.err = errDirWrongPathSyntax(io, dir)
				} else if !info.Mode().IsDir() {
					command.err = errFileNotADir(io, dir)
				}
			} else if errors.Is(err, os.ErrNotExist) {
				command.err = errDirNotExist(io, dir)
			} else if errors.Is(err, os.ErrInvalid) {
				command.err = errDirWrongPathSyntax(io, dir)
			} else {
				command.err = errDirCantRead(io, dir)
			}
		} else {
			command.err = errDirEmpty(io)
		}
	}
}

func getFilter(filterArgs, delimiterArgs, unmatchedArgs *cl.Arguments, cmdLine *cl.CommandLine) []string {
	var result []string
	if filterArgs.Available() {
		contentFilter := filterArgs.Values[0]
		if contentFilter > "" {
			if delimiterArgs.Available() {
				delimiter := delimiterArgs.Values[0]
				if delimiter > "" {
					result = strings.Split(contentFilter, delimiter)
				} else {
					result = []string{contentFilter}
				}
			} else if strings.IndexByte(contentFilter, ',') >= 0 {
				result = strings.Split(contentFilter, ",")
			} else if strings.IndexByte(contentFilter, ' ') >= 0 {
				result = strings.Split(contentFilter, " ")
			}
		}
	}
	if unmatchedArgs != nil {
		for i, index := range unmatchedArgs.Indices {
			if !cmdLine.Matched[index] {
				result = append(result, unmatchedArgs.Keys[i])
			}
		}
	}
	return removeEmpty(result)
}

func removeEmpty(result []string) []string {
	var j int
	for i := 0; i < len(result); i++ {
		if len(result[i]) > 0 {
			result[j] = result[i]
			j++
		}
	}
	return result[:j]
}

func getThreads(args *cl.Arguments, err error) (int, error) {
	if err == nil && args.Available() {
		value := args.Values[0]
		if value > "" {
			threads, errConv := strconv.Atoi(value)
			if threads < 1 {
				threads = 1
			}
			if errConv == nil {
				return threads, nil
			}
			return threads, errArgNotInteger(args.Keys[0], args.Values[0])
		}
		return 1, errMissingArgValue(args.Keys[0])
	}
	return 1, err
}

func anyAvailable(argsList []*cl.Arguments) bool {
	for _, args := range argsList {
		if args.Available() {
			return true
		}
	}
	return false
}

func getAvailable(argsList []*cl.Arguments) *cl.Arguments {
	for _, args := range argsList {
		if args.Available() {
			return args
		}
	}
	return nil
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

func argValToInt(args *cl.Arguments, index int, err error) (int, error) {
	if err == nil && args.Available() {
		if args.Values[index] > "" {
			n, errConv := strconv.Atoi(args.Values[index])
			if errConv == nil {
				return n, nil
			}
			return 0, errArgNotInteger(args.Keys[index], args.Values[index])
		}
		return 0, errMissingArgValue(args.Keys[index])
	}
	return 0, err
}

func argValToBytes(args *cl.Arguments, err error) (int64, error) {
	if err == nil && args.Available() {
		value := args.Values[0]
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
			return bytes64, errArgNotInteger(args.Keys[0], args.Values[0])
		}
		return 0, errMissingArgValue(args.Keys[0])
	}
	return 0, err
}

func byteIndex(str string, b byte) int {
	for i := 0; i < len(str); i++ {
		if str[i] == b {
			return i
		}
	}
	return -1
}

func rIndex(str string, b byte) int {
	for i := len(str) - 1; i >= 0; i-- {
		if str[i] == b {
			return i
		}
	}
	return -1
}

func getUnmatchedIdx(unmatchedArgs *cl.Arguments, matched []bool) int {
	if unmatchedArgs != nil {
		for i, index := range unmatchedArgs.Indices {
			if !matched[index] {
				return i
			}
		}
	}
	return -1
}

func isEmpty(str string) bool {
	for i := 0; i < len(str); i++ {
		b := str[i]
		if b >= 0 && b <= 32 || b != '*' {
			return false
		}
	}
	return true
}
