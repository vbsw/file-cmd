# file-cmd

[![Go Reference](https://pkg.go.dev/badge/github.com/vbsw/file-cmd.svg)](https://pkg.go.dev/github.com/vbsw/file-cmd) [![Go Report Card](https://goreportcard.com/badge/github.com/vbsw/file-cmd)](https://goreportcard.com/report/github.com/vbsw/file-cmd) [![Stability: Experimental](https://masterminds.github.io/stability/experimental.svg)](https://masterminds.github.io/stability/experimental.html)

## About
file-cmd executes various commands on files. file-cmd is published on <https://github.com/vbsw/file-cmd>.

## Copyright
Copyright 2026, Vitali Baumtrok (vbsw@mailbox.org).

file-cmd is distributed under the Boost Software License, version 1.0. (See accompanying file LICENSE or copy at http://www.boost.org/LICENSE_1_0.txt)

file-cmd is distributed in the hope that it will be useful, but WITHOUT ANY WARRANTY; without even the implied warranty of MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the Boost Software License for more details.

## Usage

	file-cmd ( INFO | SPLIT/CONCATENATE )

	INFO
		-h, --help    print this help
		-v, --version print version
		--copyright   print copyright

	SPLIT/CONCATENATE
		[COMMAND] INPUT-FILE [OUTPUT-FILE]

	COMMAND
		-p=N          split file into N parts (chunks)
		-b=N[U]       split file into N bytes per chunk, U = unit (k/K, m/M or g/G)
		-l=N          split file into N lines per chunk
		-c            concatenate files (INPUT-FILE is only one file, the first one)

## References
- https://golang.org/doc/install
- https://git-scm.com/book/en/v2/Getting-Started-Installing-Git
