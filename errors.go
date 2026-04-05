/*
 *          Copyright 2026 Vitali Baumtrok.
 * Distributed under the Boost Software License, Version 1.0.
 *     (See accompanying file LICENSE or copy at
 *        http://www.boost.org/LICENSE_1_0.txt)
 */

package main

import "errors"

func errWrongArgumentUsage() error {
	return errors.New("wrong argument usage")
}

func errUnknownCommand(arg string) error {
	return errors.New("unknown command \"" + arg + "\"")
}

func errTooManyArguments() error {
	return errors.New("too many arguments")
}
