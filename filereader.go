/*
 *          Copyright 2026 Vitali Baumtrok.
 * Distributed under the Boost Software License, Version 1.0.
 *     (See accompanying file LICENSE or copy at
 *        http://www.boost.org/LICENSE_1_0.txt)
 */

package main

import (
	"os"
)

type tFileReader struct {
	err       error
	fileInfo  os.FileInfo
	file      *os.File
	buffer    []byte
	fileSize  int64
	totalRead int64
	nRead     int
}

func (reader *tFileReader) openFile(path string) {
	reader.totalRead = 0
	reader.file, reader.err = os.Open(path)
	if reader.err == nil {
		reader.fileInfo, reader.err = reader.file.Stat()
		if reader.err == nil {
			reader.fileSize = reader.fileInfo.Size()
		} else {
			reader.file.Close()
		}
	}
}

func (reader *tFileReader) readToBuffer(offset int) {
	if offset == 0 {
		reader.nRead, reader.err = reader.file.Read(reader.buffer)
	} else {
		buffer := reader.buffer
		copy(buffer, buffer[len(buffer)-offset:])
		reader.nRead, reader.err = reader.file.Read(buffer[offset:])
		reader.nRead += offset
	}
	reader.totalRead += int64(reader.nRead)
}

func (reader *tFileReader) closeFile() {
	if reader.err == nil {
		reader.err = reader.file.Close()
	} else {
		reader.file.Close()
	}
}
