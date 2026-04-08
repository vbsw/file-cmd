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
	fileNameFilter []string
	filter         []string
	parts, lines   int
	id             int
	overwrite      bool
	or             bool
	recursive      bool
	silent         bool
}

func getCommand(osArgs []string) *tCommand {
	command := new(tCommand)
	if len(osArgs) > 0 {
		cmdLine := cl.New(osArgs, cl.NewDelimiter("", " ", "="))
		// first args with values
		cmdArgsList := getCmdArgs(cmdLine)
		// then just args
		infoArgsList := getInfoArgs(cmdLine)
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

func getInfoArgs(cmdLine *cl.CommandLine) []*cl.Arguments {
	argsList := make([]*cl.Arguments, idxInfoTotal, idxInfoTotal)
	argsList[idxInfoHelp] = cmdLine.Search("-h", "--help", "-help", "help")
	argsList[idxInfoVersion] = cmdLine.Search("-v", "--version", "-version", "version")
	argsList[idxInfoExample] = cmdLine.Search("-e", "--example", "-example", "example")
	argsList[idxInfoCopyright] = cmdLine.Search("-c", "--copyright", "-copyright", "copyright")
	return argsList
}

func getCmdArgs(cmdLine *cl.CommandLine) []*cl.Arguments {
	argsList := make([]*cl.Arguments, idxCmdTotal, idxCmdTotal)
	argsList[idxCmdSplit] = cmdLine.Search("split")
	argsList[idxCmdConcat] = cmdLine.Search("concat")
	argsList[idxCmdList] = cmdLine.Search("ls", "list", "print")
	argsList[idxCmdCount] = cmdLine.Search("count")
	argsList[idxCmdCopy] = cmdLine.Search("cp")
	argsList[idxCmdMove] = cmdLine.Search("mv")
	argsList[idxCmdRemove] = cmdLine.Search("rm")
	argsList[idxCmdText] = cmdLine.Search("text")
	return argsList
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
			validateSplit(command, cmdArgsList, cmdLine)
		} else if cmdArgsList[idxCmdConcat].HasIndex(0) {
			validateConcat(command, cmdArgsList, cmdLine)
		} else if cmdArgsList[idxCmdList].HasIndex(0) {
			validateList(command, cmdArgsList, cmdLine)
		} else if cmdArgsList[idxCmdCount].HasIndex(0) {
			validateCount(command, cmdArgsList, cmdLine)
		} else if cmdArgsList[idxCmdCopy].HasIndex(0) {
			validateCopy(command, cmdArgsList, cmdLine)
		} else if cmdArgsList[idxCmdMove].HasIndex(0) {
			validateMove(command, cmdArgsList, cmdLine)
		} else if cmdArgsList[idxCmdRemove].HasIndex(0) {
			validateRemove(command, cmdArgsList, cmdLine)
		} else if cmdArgsList[idxCmdText].HasIndex(0) {
			validateText(command, cmdArgsList, cmdLine)
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

func validateSplit(command *tCommand, cmdArgsList []*cl.Arguments, cmdLine *cl.CommandLine) {
	if len(cmdLine.Arguments) == 1 {
		command.err = errNotEnoughArguments("split")
	} else {
		optArgsList := getOptSplitArgs(cmdLine)
		argMultiple := getMultiple(optArgsList)
		if argMultiple == nil {
			if cmdArgsList[idxCmdSplit].Count() > 1 {
				cmdLine.RevertMatched(cmdArgsList[idxCmdSplit])
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
			command.bytes, command.err = argValToBytes(optArgsList[idxOptSplitBytes], 0, command.err)
			command.lines, command.err = argValToInt(optArgsList[idxOptSplitLines], 0, command.err)
			command.overwrite = optArgsList[idxOptSplitOvierwrite].Available()
			command.id = cmdSplit
			if command.parts <= 0 && command.bytes <= 0 && command.lines <= 0 {
				command.parts = 2
			}
		} else {
			command.err = errMultipleUsage(argMultiple.Keys)
		}
	}
}

func validateConcat(command *tCommand, cmdArgsList []*cl.Arguments, cmdLine *cl.CommandLine) {
	if len(cmdLine.Arguments) == 1 {
		command.err = errNotEnoughArguments("concat")
	} else {
		optArgsList := getOptConcatArgs(cmdLine)
		argMultiple := getMultiple(optArgsList)
		if argMultiple == nil {
			if cmdArgsList[idxCmdConcat].Count() > 1 {
				cmdLine.RevertMatched(cmdArgsList[idxCmdConcat])
				cmdLine.Matched[0] = true
				cmdLine.Matched[len(cmdLine.Arguments)] = false
			}
			unmatchedArgs := cmdLine.Unmatched()
			validateInputOutput(command, optArgsList[idxOptConcatInput], optArgsList[idxOptConcatOutput], unmatchedArgs)
			validateInputConcatFile(command)
			validateOutputConcatFile(command)
			command.overwrite = optArgsList[idxOptConcatOvierwrite].Available()
			command.id = cmdConcat
		} else {
			command.err = errMultipleUsage(argMultiple.Keys)
		}
	}
}

func validateList(command *tCommand, cmdArgsList []*cl.Arguments, cmdLine *cl.CommandLine) {
	if len(cmdLine.Arguments) == 1 {
		command.err = errNotEnoughArguments(cmdArgsList[idxCmdList].Keys[0])
	} else {
		optArgsList := getOptListArgs(cmdLine)
		argMultiple := getMultiple(optArgsList)
		if argMultiple == nil {
			if cmdArgsList[idxCmdList].Count() > 1 {
				cmdLine.RevertMatched(cmdArgsList[idxCmdList])
				cmdLine.Matched[0] = true
				cmdLine.Matched[len(cmdLine.Arguments)] = false
			}
			unmatchedArgs := cmdLine.Unmatched()
			validateInputListFile(command, optArgsList[idxOptConcatInput], unmatchedArgs, cmdLine)
			command.or = optArgsList[idxOptListOr].Available()
			command.recursive = optArgsList[idxOptListRecursive].Available()
			command.silent = optArgsList[idxOptListSilent].Available()
			command.id = cmdList
			if optArgsList[idxOptListFilter].Available() {
				filter := optArgsList[idxOptListFilter].Values[0]
				if filter > "" {
					if optArgsList[idxOptListDelimiter].Available() {
						delimiter := optArgsList[idxOptListDelimiter].Values[0]
						if delimiter > "" {
							command.filter = strings.Split(filter, delimiter)
						} else {
							command.filter = []string{filter}
						}
					} else if strings.IndexByte(filter, ',') >= 0 {
						command.filter = strings.Split(filter, ",")
					} else if strings.IndexByte(filter, ' ') >= 0 {
						command.filter = strings.Split(filter, " ")
					}
				}
			}
			for i, index := range unmatchedArgs.Indices {
				if !cmdLine.Matched[index] {
					command.filter = append(command.filter, unmatchedArgs.Keys[i])
				}
			}
		} else {
			command.err = errMultipleUsage(argMultiple.Keys)
		}
	}
}

func validateCount(command *tCommand, cmdArgsList []*cl.Arguments, cmdLine *cl.CommandLine) {
	// TODO
}

func validateCopy(command *tCommand, cmdArgsList []*cl.Arguments, cmdLine *cl.CommandLine) {
	// TODO
}

func validateMove(command *tCommand, cmdArgsList []*cl.Arguments, cmdLine *cl.CommandLine) {
	// TODO
}

func validateRemove(command *tCommand, cmdArgsList []*cl.Arguments, cmdLine *cl.CommandLine) {
	// TODO
}

func validateText(command *tCommand, cmdArgsList []*cl.Arguments, cmdLine *cl.CommandLine) {
	// TODO
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
	argsList[idxOptSplitParts] = cmdLine.SearchByDelimiter("-p")
	argsList[idxOptSplitBytes] = cmdLine.SearchByDelimiter("-b")
	argsList[idxOptSplitLines] = cmdLine.SearchByDelimiter("-l")
	argsList[idxOptSplitOvierwrite] = cmdLine.Search("-w")
	argsList[idxOptSplitInput] = cmdLine.SearchByDelimiter("-i", "--input")
	argsList[idxOptSplitOutput] = cmdLine.SearchByDelimiter("-o", "--output")
	return argsList
}

func getOptConcatArgs(cmdLine *cl.CommandLine) []*cl.Arguments {
	argsList := make([]*cl.Arguments, idxOptConcatTotal, idxOptConcatTotal)
	argsList[idxOptConcatOvierwrite] = cmdLine.Search("-w")
	argsList[idxOptConcatInput] = cmdLine.SearchByDelimiter("-i", "--input")
	argsList[idxOptConcatOutput] = cmdLine.SearchByDelimiter("-o", "--output")
	return argsList
}

func getOptListArgs(cmdLine *cl.CommandLine) []*cl.Arguments {
	argsList := make([]*cl.Arguments, idxOptListTotal, idxOptListTotal)
	argsList[idxOptListOr] = cmdLine.Search("--or")
	argsList[idxOptListRecursive] = cmdLine.Search("-r", "--recursive")
	argsList[idxOptListSilent] = cmdLine.Search("-s", "--silent")
	argsList[idxOptListInput] = cmdLine.SearchByDelimiter("-i", "--input")
	argsList[idxOptListFilter] = cmdLine.SearchByDelimiter("-f")
	argsList[idxOptListDelimiter] = cmdLine.SearchByDelimiter("-d")
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
		command.err = errInputEmpty()
	}
}

func validateInputFile(command *tCommand) {
	if command.err == nil {
		if command.inputPath > "" {
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
			command.err = errInputEmpty()
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

func validateInputListFile(command *tCommand, inputArgs, unmatchedArgs *cl.Arguments, cmdLine *cl.CommandLine) {
	if inputArgs.Available() {
		command.inputPath = inputArgs.Values[0]
	} else if unmatchedArgs.Count() > 0 {
		command.inputPath = unmatchedArgs.Keys[0]
		cmdLine.Matched[unmatchedArgs.Indices[0]] = true
	} else {
		command.err = errInputEmpty()
	}
	if command.err == nil {
		if command.inputPath > "" {
			command.inputDir = command.inputPath
			validateInputDir(command)
			if command.err != nil {
				command.err = nil
				command.inputDir = filepath.Dir(command.inputPath)
				command.fileNameFilter = strings.Split(filepath.Base(command.inputPath), "*")
				validateInputDir(command)
			}
		} else {
			command.err = errInputEmpty()
		}
	}
}

func validateInputDir(command *tCommand) {
	if command.err == nil {
		if command.inputDir > "" {
			info, err := os.Stat(command.inputDir)
			if err == nil {
				if info == nil {
					command.err = errInputDirWrongPathSyntax(command.inputDir)
				} else if !info.Mode().IsDir() {
					command.err = errInputFileNotADir(command.inputDir)
				}
			} else if errors.Is(err, os.ErrNotExist) {
				command.err = errInputDirNotExist(command.inputDir)
			} else if errors.Is(err, os.ErrInvalid) {
				command.err = errInputDirWrongPathSyntax(command.inputDir)
			} else {
				command.err = errInputDirCantRead(command.inputDir)
			}
		} else {
			command.err = errInputDirEmpty()
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

func argValToBytes(args *cl.Arguments, index int, err error) (int64, error) {
	if err == nil && args.Available() {
		value := args.Values[index]
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
			return bytes64, errArgNotInteger(args.Keys[index], args.Values[index])
		}
		return 0, errMissingArgValue(args.Keys[index])
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
