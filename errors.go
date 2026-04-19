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
		} else if !isInList(args[:i], arg) {
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

func errUnknownFormat(format string) error {
	return errors.New("unknown format \"" + format + "\"")
}

func errIncompatibleFormats(formatA, formatB string) error {
	return errors.New("incompatible formats \"" + formatA + "\" and \"" + formatB + "\"")
}

func errNotEnoughArguments(cmdStr string) error {
	return errors.New("not enough arguments (see: --help " + cmdStr + ")")
}

func errPathEmpty(io string) error {
	return errors.New(io + " path is empty")
}

func errDirEmpty(io string) error {
	return errors.New(io + " directory is empty")
}

func errFileNotExist(io, path string) error {
	return errors.New(io + " file does not exist \"" + path + "\"")
}

func errDirNotExist(io, dir string) error {
	return errors.New(io + " directory does not exist \"" + dir + "\"")
}

func errFileNotAFile(io, path string) error {
	return errors.New(io + " file is not a regular file \"" + path + "\"")
}

func errFileNotADir(io, dir string) error {
	return errors.New(io + " directory is not a directory \"" + dir + "\"")
}

func errFileWrongPathSyntax(io, path string) error {
	return errors.New(io + " file has wrong path syntax \"" + path + "\"")
}

func errDirWrongPathSyntax(io, dir string) error {
	return errors.New(io + " directory has wrong path syntax \"" + dir + "\"")
}

func errFileCantRead(io, path string) error {
	return errors.New("can't read " + io + " file \"" + path + "\"")
}

func errDirCantRead(io, dir string) error {
	return errors.New("can't read " + io + " directory \"" + dir + "\"")
}

func errFileExists(io, path string) error {
	return errors.New(io + " input file already exists \"" + path + "\"")
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

func isInList(list []string, value string) bool {
	for _, val := range list {
		if val == value {
			return true
		}
	}
	return false
}
