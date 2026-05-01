# file-cmd

[![Go Reference](https://pkg.go.dev/badge/github.com/vbsw/file-cmd.svg)](https://pkg.go.dev/github.com/vbsw/file-cmd) [![Go Report Card](https://goreportcard.com/badge/github.com/vbsw/file-cmd)](https://goreportcard.com/report/github.com/vbsw/file-cmd)

## About
file-cmd executes various commands on files. file-cmd is published on <https://github.com/vbsw/file-cmd>.

## Copyright
Copyright 2026, Vitali Baumtrok (vbsw@mailbox.org).

file-cmd is distributed under the Boost Software License, version 1.0. (See accompanying file LICENSE or copy at http://www.boost.org/LICENSE_1_0.txt)

file-cmd is distributed in the hope that it will be useful, but WITHOUT ANY WARRANTY; without even the implied warranty of MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the Boost Software License for more details.

## Usage

	USAGE
	    file-cmd ( INFO | HELP | EXAMPLE | COMMAND )
	INFO
	    -v, --version     print version
	    -c, --copyright   print copyright
	HELP
	    ( -h | --help ) [COMMAND]
	               print this help, or print help for COMMAND
	EXAMPLE
	    ( -e | --example ) [COMMAND]
	               print available examples, or print expample for COMMAND
	COMMAND
	    split      split file into multiple files
	    concat     concatenate files
	    ls         list files
	    count      count files
	    cp         copy files
	    mv         move files
	    rm         delete files
	    clean      remove empty folders and files
	    text       generate random text

## References
- https://golang.org/doc/install
- https://git-scm.com/book/en/v2/Getting-Started-Installing-Git
