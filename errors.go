/*
 *          Copyright 2026 Vitali Baumtrok.
 * Distributed under the Boost Software License, Version 1.0.
 *     (See accompanying file LICENSE or copy at
 *        http://www.boost.org/LICENSE_1_0.txt)
 */

package main

import "errors"

func errMultipleUsage(args []string) error {
	errStr := "same parameter multiple"
	for i, arg := range args {
		if i == 0 {
			errStr = errStr + " \"" + arg + "\""
		} else if !inList(args[:i], arg) {
			errStr = errStr + ", \"" + arg + "\""
		}
	}
	return errors.New(errStr)
}

func errWrongArgumentUsage() error {
	return errors.New("wrong argument usage")
}

func errUnknownCommand(arg string) error {
	return errors.New("unknown command \"" + arg + "\"")
}

func errTooManyArguments() error {
	return errors.New("too many arguments")
}

func errUnknownArgument(arg string) error {
	return errors.New("unknown argument \"" + arg + "\"")
}

func errNotEnoughArguments(cmdStr string) error {
	return errors.New("not enough arguments (see: --help " + cmdStr + ")")
}

func errInputEmpty() error {
	return errors.New("input path is empty")
}

func errInputDirEmpty() error {
	return errors.New("input directory is empty")
}

func errInputFileNotExist(inputPath string) error {
	return errors.New("input file does not exist \"" + inputPath + "\"")
}

func errInputDirNotExist(inputDir string) error {
	return errors.New("input directory does not exist \"" + inputDir + "\"")
}

func errInputFileNotAFile(inputPath string) error {
	return errors.New("input file is not a regular file \"" + inputPath + "\"")
}

func errInputFileNotADir(inputDir string) error {
	return errors.New("input directory is not a directory \"" + inputDir + "\"")
}

func errInputFileWrongPathSyntax(inputPath string) error {
	return errors.New("input file has wrong path syntax \"" + inputPath + "\"")
}

func errInputDirWrongPathSyntax(inputDir string) error {
	return errors.New("input directory has wrong path syntax \"" + inputDir + "\"")
}

func errInputFileCantRead(inputPath string) error {
	return errors.New("can't read input file \"" + inputPath + "\"")
}

func errInputDirCantRead(inputDir string) error {
	return errors.New("can't read input directory \"" + inputDir + "\"")
}

func errOutputFileExists(outputPath string) error {
	return errors.New("output input file already exists \"" + outputPath + "\"")
}

func errOutputFileWrongPathSyntax(outputPath string) error {
	return errors.New("output file has wrong path syntax \"" + outputPath + "\"")
}

func errOutputEmpty() error {
	return errors.New("output path is empty")
}

func errUnknownState() error {
	return errors.New("unknown state")
}

func errMissingArgValue(arg string) error {
	return errors.New("missing argument value \"" + arg + "\"")
}

func errArgNotInteger(arg, val string) error {
	return errors.New("argument \"" + arg + "\" must be integer, is \"" + val + "\"")
}

func inList(list []string, value string) bool {
	for _, val := range list {
		if val == value {
			return true
		}
	}
	return false
}
