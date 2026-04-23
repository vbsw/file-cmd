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
	"github.com/vbsw/go-lib/fs"
	"os"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"
)

type tCommand struct {
	err            error
	inputPath      string
	outputPath     string
	str            string
	delimiter      string
	inputDir       string
	outputDir      string
	lineTerminator string
	fileNameFilter string
	contentFilter  []string
	bytes          int64
	parts, lines   int
	id             int
	threads        int
	overwrite      bool
	or             bool
	recursive      bool
	verbose        bool
	lettersOnly    bool
	upperCase      bool
	lowerCase      bool
	list           bool
	all            bool
}

func getCommand(osArgs []string) *tCommand {
	command := new(tCommand)
	if len(osArgs) > 0 {
		var args [idxInfoTotal + idxCmdTotal]*cl.Arguments
		infoArgs, cmdArgs := args[:idxInfoTotal], args[idxInfoTotal:len(args)]
		cmdLine := cl.New(osArgs, cl.NewDelimiter("", " ", "="))
		command.threads = 1
		readCmdArgs(cmdArgs, cmdLine)   // args having values
		readInfoArgs(infoArgs, cmdLine) // args without values
		validateInfo(command, infoArgs, cmdArgs, cmdLine)
		validateCmd(command, infoArgs, cmdArgs, cmdLine)
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
	argsList[idxCmdClean] = cmdLine.Match("cl", "clean")
	argsList[idxCmdText] = cmdLine.Match("text")
}

func validateInfo(command *tCommand, infoArgs, cmdArgs []*cl.Arguments, cmdLine *cl.CommandLine) {
	if infoArgs[idxInfoHelp].HasIndex(0) {
		validateInfoHelp(command, infoArgs, cmdArgs, cmdLine)
	} else if infoArgs[idxInfoVersion].HasIndex(0) {
		validateInfoVersion(command, infoArgs, cmdArgs, cmdLine)
	} else if infoArgs[idxInfoExample].HasIndex(0) {
		validateInfoExample(command, infoArgs, cmdArgs, cmdLine)
	} else if infoArgs[idxInfoCopyright].HasIndex(0) {
		validateInfoCopyright(command, infoArgs, cmdArgs, cmdLine)
	} // else: cound be a command
}

func validateCmd(command *tCommand, infoArgs, cmdArgs []*cl.Arguments, cmdLine *cl.CommandLine) {
	if command.err == nil && command.id == cmdNone {
		cmdLine.RevertMatched(infoArgs...)
		if cmdArgs[idxCmdSplit].HasIndex(0) {
			validateSplit(command, cmdArgs[idxCmdSplit], cmdLine)
		} else if cmdArgs[idxCmdConcat].HasIndex(0) {
			validateConcat(command, cmdArgs[idxCmdConcat], cmdLine)
		} else if cmdArgs[idxCmdList].HasIndex(0) {
			validateList(command, cmdArgs[idxCmdList], cmdLine)
		} else if cmdArgs[idxCmdCount].HasIndex(0) {
			validateCount(command, cmdArgs[idxCmdCount], cmdLine)
		} else if cmdArgs[idxCmdCopy].HasIndex(0) {
			validateCopy(command, cmdArgs[idxCmdCopy], cmdLine)
		} else if cmdArgs[idxCmdMove].HasIndex(0) {
			validateMove(command, cmdArgs[idxCmdMove], cmdLine)
		} else if cmdArgs[idxCmdRemove].HasIndex(0) {
			validateRemove(command, cmdArgs[idxCmdRemove], cmdLine)
		} else if cmdArgs[idxCmdClean].HasIndex(0) {
			validateClean(command, cmdArgs[idxCmdClean], cmdLine)
		} else if cmdArgs[idxCmdText].HasIndex(0) {
			validateText(command, cmdArgs[idxCmdText], cmdLine)
		} else {
			command.err = errUnknownCommand(cmdLine.Arguments[0])
		}
	}
}

func validateInfoHelp(command *tCommand, infoArgs, cmdArgs []*cl.Arguments, cmdLine *cl.CommandLine) {
	if len(cmdLine.Arguments) == 1 {
		command.id = cmdInfo
	} else {
		if infoArgs[idxInfoHelp].Count() == 1 {
			if len(cmdLine.Arguments) == 2 {
				cmdAArgs := getAvailable(cmdArgs)
				if cmdAArgs == nil {
					otherInfoAvail := infoArgs[idxInfoVersion].Available()
					otherInfoAvail = otherInfoAvail || infoArgs[idxInfoExample].Available()
					otherInfoAvail = otherInfoAvail || infoArgs[idxInfoCopyright].Available()
					if otherInfoAvail {
						command.err = errWrongArgumentUsage()
					} else {
						command.err = errUnknownCommand(cmdLine.Arguments[1])
					}
				} else {
					command.id = getCmdInfoId(cmdArgs)
					command.str = cmdAArgs.Keys[0]
				}
			} else {
				command.err = errTooManyArguments()
			}
		} else {
			command.err = errMultipleUsage(infoArgs[idxInfoHelp].Keys)
		}
	}
}

func validateInfoVersion(command *tCommand, infoArgs, cmdArgs []*cl.Arguments, cmdLine *cl.CommandLine) {
	if len(cmdLine.Arguments) == 1 {
		command.id = cmdVersion
	} else {
		if infoArgs[idxInfoVersion].Count() == 1 {
			command.err = errTooManyArguments()
		} else {
			command.err = errWrongArgumentUsage()
		}
	}
}

func validateInfoExample(command *tCommand, infoArgs, cmdArgs []*cl.Arguments, cmdLine *cl.CommandLine) {
	if len(cmdLine.Arguments) == 1 {
		command.id = cmdInfoExample
	} else {
		if infoArgs[idxInfoExample].Count() == 1 {
			if len(cmdLine.Arguments) == 2 {
				if anyAvailable(cmdArgs) {
					command.id = getCmdExampleId(cmdArgs)
				} else {
					otherInfoAvail := infoArgs[idxInfoHelp].Available()
					otherInfoAvail = otherInfoAvail || infoArgs[idxInfoVersion].Available()
					otherInfoAvail = otherInfoAvail || infoArgs[idxInfoCopyright].Available()
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
			command.err = errMultipleUsage(infoArgs[idxInfoExample].Keys)
		}
	}
}

func validateInfoCopyright(command *tCommand, infoArgs, cmdArgs []*cl.Arguments, cmdLine *cl.CommandLine) {
	if len(cmdLine.Arguments) == 1 {
		command.id = cmdCopyright
	} else {
		if infoArgs[idxInfoCopyright].Count() == 1 {
			command.err = errTooManyArguments()
		} else {
			command.err = errMultipleUsage(infoArgs[idxInfoCopyright].Keys)
		}
	}
}

func validateSplit(command *tCommand, cmdArgs *cl.Arguments, cmdLine *cl.CommandLine) {
	if len(cmdLine.Arguments) == 1 {
		command.err = errNotEnoughArguments("split")
	} else {
		optArgs := getOptSplitArgs(cmdLine)
		argMultiple := getMultiple(optArgs)
		if argMultiple == nil {
			if cmdArgs.Count() > 1 {
				cmdLine.RevertMatched(cmdArgs)
				cmdLine.Matched[0] = true
				cmdLine.Matched[len(cmdLine.Arguments)] = false
			}
			unmatchedArgs := cmdLine.Unmatched()
			validateInputOutput(command, optArgs[idxOptSplitInput], optArgs[idxOptSplitOutput], unmatchedArgs)
			validateInputFile(command)
			if command.outputPath == "" {
				command.outputPath = command.inputPath
			}
			validateOutputFile(command, ".0")
			command.parts, command.err = argValToInt(optArgs[idxOptSplitParts], 0, command.err)
			command.bytes, command.err = argValToBytes(optArgs[idxOptSplitBytes], command.err)
			command.lines, command.err = argValToInt(optArgs[idxOptSplitLines], 0, command.err)
			command.overwrite = optArgs[idxOptSplitOverwrite].Available()
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
		optArgs := getOptConcatArgs(cmdLine)
		argMultiple := getMultiple(optArgs)
		if argMultiple == nil {
			if cmdArgs.Count() > 1 {
				cmdLine.RevertMatched(cmdArgs)
				cmdLine.Matched[0] = true
			}
			unmatchedArgs := cmdLine.Unmatched()
			validateInputOutput(command, optArgs[idxOptConcatInput], optArgs[idxOptConcatOutput], unmatchedArgs)
			validateInputConcatFile(command)
			validateOutputConcatFile(command)
			command.overwrite = optArgs[idxOptConcatOverwrite].Available()
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
		optArgs := getOptListArgs(cmdLine)
		argMultiple := getMultiple(optArgs)
		if argMultiple == nil {
			if cmdArgs.Count() > 1 {
				cmdLine.RevertMatched(cmdArgs)
				cmdLine.Matched[0] = true
			}
			unmatchedArgs := cmdLine.Unmatched()
			validateInputDirFile(command, optArgs[idxOptListInput], unmatchedArgs, cmdLine)
			command.or = optArgs[idxOptListOr].Available()
			command.recursive = optArgs[idxOptListRecursive].Available()
			command.verbose = optArgs[idxOptListVerbose].Available()
			command.contentFilter = getFilter(optArgs[idxOptListFilter], optArgs[idxOptListDelimiter], unmatchedArgs, cmdLine)
			command.threads, command.err = getThreads(optArgs[idxOptListThreads], command.err)
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
		optArgs := getOptCountArgs(cmdLine)
		argMultiple := getMultiple(optArgs)
		if argMultiple == nil {
			if cmdArgs.Count() > 1 {
				cmdLine.RevertMatched(cmdArgs)
				cmdLine.Matched[0] = true
			}
			unmatchedArgs := cmdLine.Unmatched()
			validateInputDirFile(command, optArgs[idxOptCountInput], unmatchedArgs, cmdLine)
			command.or = optArgs[idxOptCountOr].Available()
			command.recursive = optArgs[idxOptCountRecursive].Available()
			command.verbose = optArgs[idxOptCountVerbose].Available()
			command.contentFilter = getFilter(optArgs[idxOptCountFilter], optArgs[idxOptCountDelimiter], unmatchedArgs, cmdLine)
			command.threads, command.err = getThreads(optArgs[idxOptCountThreads], command.err)
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
		optArgs := getOptCopyArgs(cmdLine)
		argMultiple := getMultiple(optArgs)
		if argMultiple == nil {
			if cmdArgs.Count() > 1 {
				cmdLine.RevertMatched(cmdArgs)
				cmdLine.Matched[0] = true
			}
			unmatchedArgs := cmdLine.Unmatched()
			validateInputDirFile(command, optArgs[idxOptCopyInput], unmatchedArgs, cmdLine)
			validateOutputDir(command, optArgs[idxOptCopyOutput], unmatchedArgs, cmdLine)
			command.or = optArgs[idxOptCopyOr].Available()
			command.overwrite = optArgs[idxOptCopyOverwrite].Available()
			command.recursive = optArgs[idxOptCopyRecursive].Available()
			command.verbose = optArgs[idxOptCopyVerbose].Available()
			command.contentFilter = getFilter(optArgs[idxOptCopyFilter], optArgs[idxOptCopyDelimiter], unmatchedArgs, cmdLine)
			command.threads, command.err = getThreads(optArgs[idxOptCopyThreads], command.err)
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
		optArgs := getOptMoveArgs(cmdLine)
		argMultiple := getMultiple(optArgs)
		if argMultiple == nil {
			if cmdArgs.Count() > 1 {
				cmdLine.RevertMatched(cmdArgs)
				cmdLine.Matched[0] = true
			}
			unmatchedArgs := cmdLine.Unmatched()
			validateInputDirFile(command, optArgs[idxOptMoveInput], unmatchedArgs, cmdLine)
			validateOutputDir(command, optArgs[idxOptMoveOutput], unmatchedArgs, cmdLine)
			command.or = optArgs[idxOptMoveOr].Available()
			command.overwrite = optArgs[idxOptMoveOverwrite].Available()
			command.recursive = optArgs[idxOptMoveRecursive].Available()
			command.verbose = optArgs[idxOptMoveVerbose].Available()
			command.contentFilter = getFilter(optArgs[idxOptMoveFilter], optArgs[idxOptMoveDelimiter], unmatchedArgs, cmdLine)
			command.threads, command.err = getThreads(optArgs[idxOptMoveThreads], command.err)
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
		optArgs := getOptRemoveArgs(cmdLine)
		argMultiple := getMultiple(optArgs)
		if argMultiple == nil {
			if cmdArgs.Count() > 1 {
				cmdLine.RevertMatched(cmdArgs)
				cmdLine.Matched[0] = true
			}
			unmatchedArgs := cmdLine.Unmatched()
			validateInputDirFile(command, optArgs[idxOptRemoveInput], unmatchedArgs, cmdLine)
			command.or = optArgs[idxOptRemoveOr].Available()
			command.recursive = optArgs[idxOptRemoveRecursive].Available()
			command.verbose = optArgs[idxOptRemoveVerbose].Available()
			command.contentFilter = getFilter(optArgs[idxOptRemoveFilter], optArgs[idxOptRemoveDelimiter], unmatchedArgs, cmdLine)
			command.threads, command.err = getThreads(optArgs[idxOptRemoveThreads], command.err)
			command.id = cmdRemove
		} else {
			command.err = errMultipleUsage(argMultiple.Keys)
		}
	}
}

func validateClean(command *tCommand, cmdArgs *cl.Arguments, cmdLine *cl.CommandLine) {
	if len(cmdLine.Arguments) == 1 {
		command.err = errNotEnoughArguments(cmdArgs.Keys[0])
	} else {
		optArgs := getOptCleanArgs(cmdLine)
		argMultiple := getMultiple(optArgs)
		if argMultiple == nil {
			if cmdArgs.Count() > 1 {
				cmdLine.RevertMatched(cmdArgs)
				cmdLine.Matched[0] = true
			}
			unmatchedArgs := cmdLine.Unmatched()
			validateInputDirFile(command, optArgs[idxOptCleanInput], unmatchedArgs, cmdLine)
			validateUnmatchedOption(command, unmatchedArgs, cmdLine)
			command.all = optArgs[idxOptCleanAll].Available()
			command.recursive = optArgs[idxOptCleanRecursive].Available()
			command.verbose = optArgs[idxOptCleanVerbose].Available()
			command.list = optArgs[idxOptCleanList].Available()
			command.id = cmdClean
		} else {
			command.err = errMultipleUsage(argMultiple.Keys)
		}
	}
}

func validateText(command *tCommand, cmdArgs *cl.Arguments, cmdLine *cl.CommandLine) {
	if len(cmdLine.Arguments) == 1 {
		command.err = errNotEnoughArguments(cmdArgs.Keys[0])
	} else {
		optArgs := getOptTextArgs(cmdLine)
		argMultiple := getMultiple(optArgs)
		if argMultiple == nil {
			if cmdArgs.Count() > 1 {
				cmdLine.RevertMatched(cmdArgs)
				cmdLine.Matched[0] = true
			}
			unmatchedArgs := cmdLine.Unmatched()
			validateOutputTextFile(command, optArgs[idxOptTextOutput], unmatchedArgs, cmdLine)
			if command.err != nil && command.outputPath == "" {
				command.err = nil
			}
			validateFormat(command, optArgs[idxOptTextFormat])
			command.threads, command.err = getThreads(optArgs[idxOptTextThreads], command.err)
			command.overwrite = optArgs[idxOptTextOverwrite].Available()
			command.bytes, command.err = argValToBytes(optArgs[idxOptTextSize], command.err)
			command.lineTerminator = getLineTerminator(optArgs[idxOptTextTerminator])
			command.delimiter = optArgs[idxOptTextDelimiter].ValueAt(0, " ")
			command.id = cmdText
			validateOverwrite(command)
		} else {
			command.err = errMultipleUsage(argMultiple.Keys)
		}
	}
}

func getLineTerminator(args *cl.Arguments) string {
	lineTerminator := args.ValueAt(0, "")
	if len(lineTerminator) > 0 {
		lineTerminator = strings.ReplaceAll(lineTerminator, "CR", "\r")
		lineTerminator = strings.ReplaceAll(lineTerminator, "LF", "\n")
		return lineTerminator
	}
	if runtime.GOOS == "windows" {
		return "\r\n"
	}
	return "\n"
}

func getCmdInfoId(cmdArgs []*cl.Arguments) int {
	cmdId := cmdNone
	if cmdArgs[idxCmdSplit].Available() {
		cmdId = cmdInfoSplit
	} else if cmdArgs[idxCmdConcat].Available() {
		cmdId = cmdInfoConcat
	} else if cmdArgs[idxCmdList].Available() {
		cmdId = cmdInfoList
	} else if cmdArgs[idxCmdCount].Available() {
		cmdId = cmdInfoCount
	} else if cmdArgs[idxCmdCopy].Available() {
		cmdId = cmdInfoCopy
	} else if cmdArgs[idxCmdMove].Available() {
		cmdId = cmdInfoMove
	} else if cmdArgs[idxCmdRemove].Available() {
		cmdId = cmdInfoRemove
	} else if cmdArgs[idxCmdClean].Available() {
		cmdId = cmdInfoClean
	} else if cmdArgs[idxCmdText].Available() {
		cmdId = cmdInfoText
	}
	return cmdId
}

func getCmdExampleId(cmdArgs []*cl.Arguments) int {
	cmdId := cmdNone
	if cmdArgs[idxCmdSplit].Available() {
		cmdId = cmdExampleSplit
	} else if cmdArgs[idxCmdConcat].Available() {
		cmdId = cmdExampleConcat
	} else if cmdArgs[idxCmdList].Available() {
		cmdId = cmdExampleList
	} else if cmdArgs[idxCmdCount].Available() {
		cmdId = cmdExampleCount
	} else if cmdArgs[idxCmdCopy].Available() {
		cmdId = cmdExampleCopy
	} else if cmdArgs[idxCmdMove].Available() {
		cmdId = cmdExampleMove
	} else if cmdArgs[idxCmdRemove].Available() {
		cmdId = cmdExampleRemove
	} else if cmdArgs[idxCmdClean].Available() {
		cmdId = cmdExampleClean
	} else if cmdArgs[idxCmdText].Available() {
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
	argsList[idxOptListVerbose] = cmdLine.Match("-v", "--verbose")
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
	argsList[idxOptCountVerbose] = cmdLine.Match("-v", "--verbose")
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
	argsList[idxOptCopyVerbose] = cmdLine.Match("-v", "--verbose")
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
	argsList[idxOptMoveVerbose] = cmdLine.Match("-v", "--verbose")
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
	argsList[idxOptRemoveVerbose] = cmdLine.Match("-v", "--verbose")
	argsList[idxOptRemoveInput] = cmdLine.MatchDelimited("-i", "--input")
	argsList[idxOptRemoveFilter] = cmdLine.MatchDelimited("-f")
	argsList[idxOptRemoveDelimiter] = cmdLine.MatchDelimited("-d")
	argsList[idxOptRemoveThreads] = cmdLine.MatchDelimited("-t")
	return argsList
}

func getOptCleanArgs(cmdLine *cl.CommandLine) []*cl.Arguments {
	argsList := make([]*cl.Arguments, idxOptRemoveTotal, idxOptRemoveTotal)
	argsList[idxOptCleanAll] = cmdLine.Match("-a", "--all")
	argsList[idxOptCleanRecursive] = cmdLine.Match("-r", "--recursive")
	argsList[idxOptCleanVerbose] = cmdLine.Match("-v", "--verbose")
	argsList[idxOptCleanList] = cmdLine.Match("-l", "--list")
	argsList[idxOptCleanInput] = cmdLine.MatchDelimited("-i", "--input")
	return argsList
}

func getOptTextArgs(cmdLine *cl.CommandLine) []*cl.Arguments {
	argsList := make([]*cl.Arguments, idxOptTextTotal, idxOptTextTotal)
	argsList[idxOptTextOutput] = cmdLine.MatchDelimited("-o", "--output")
	argsList[idxOptTextSize] = cmdLine.MatchDelimited("-s")
	argsList[idxOptTextTerminator] = cmdLine.MatchDelimited("-e")
	argsList[idxOptTextDelimiter] = cmdLine.MatchDelimited("-d")
	argsList[idxOptTextFormat] = cmdLine.MatchDelimited("-f")
	argsList[idxOptTextOverwrite] = cmdLine.Match("-w")
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
}

func validateOverwrite(command *tCommand) {
	if command.err == nil {
		var file fs.File
		if file.Stat(command.outputPath) {
			if command.overwrite {
				if !file.Info.Mode().IsRegular() {
					command.err = errCantOverwriteNotRegular(command.outputPath)
				}
			} else {
				command.err = errFileExists("output", command.outputPath)
			}
		}
	}
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
			var file fs.File
			if file.IsDir(command.inputPath) {
				command.inputDir = command.inputPath
				command.fileNameFilter = "*"
			} else {
				command.inputDir = filepath.Dir(command.inputPath)
				command.fileNameFilter = filepath.Base(command.inputPath)
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

func validateUnmatchedOption(command *tCommand, unmatchedArgs *cl.Arguments, cmdLine *cl.CommandLine) {
	if command.err == nil && unmatchedArgs.Available() {
		for i, index := range unmatchedArgs.Indices {
			if !cmdLine.Matched[index] {
				command.err = errUnknownOption(unmatchedArgs.Keys[i])
			}
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

func anyMatches(argsList []string, arg string) bool {
	for _, args := range argsList {
		if args == arg {
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
