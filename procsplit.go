/*
 *          Copyright 2026 Vitali Baumtrok.
 * Distributed under the Boost Software License, Version 1.0.
 *     (See accompanying file LICENSE or copy at
 *        http://www.boost.org/LICENSE_1_0.txt)
 */

package main

import (
	"errors"
	"fmt"
	"github.com/vbsw/go-lib/fs"
	"os"
)

type tLines struct {
	lineEnd int
	line    int
	rPrev   bool
}

func processSplit(command *tCommand) {
	var file fs.FileReader
	defer func() {
		file.Close()
		command.err = file.Err
	}()
	if file.Open(command.inputPath) {
		fileSize := file.Info.Size()
		if fileSize >= splitBufferSize {
			file.Buffer = make([]byte, splitBufferSize)
		} else {
			file.Buffer = make([]byte, fileSize)
		}
		chunks := getChunks(command, &file, fileSize)
		if file.Err == nil {
			outputPaths := getOutputPaths(command.outputPath, chunks)
			if ensureOutput(command, outputPaths) {
				if file.Seek(0) {
					for i, outputPath := range outputPaths {
						if !file.CopyNTo(outputPath, chunks[i]) {
							break
						}
					}
				}
			} else {
				file.Err = command.err
			}
		}
	}
}

func getChunks(command *tCommand, file *fs.FileReader, fileSize int64) []int64 {
	var chunks []int64
	if command.parts > 0 {
		chunkSize := getChunkSize(int64(command.parts), fileSize)
		chunks = getChunksFromChunkSize(int64(command.parts), chunkSize, fileSize)
	} else if command.bytes > 0 {
		parts := getParts(int64(command.bytes), fileSize)
		chunks = getChunksFromChunkSize(parts, int64(command.bytes), fileSize)
	} else {
		chunks = getChunksFromLines(file, command.lines, fileSize)
	}
	return chunks
}

func getChunkSize(parts, fileSize int64) int64 {
	chunkSize := fileSize / parts
	if fileSize <= maxInt64-(parts-1) {
		chunkSize = (fileSize + parts - 1) / parts
	} else {
		chunkSize = fileSize / parts
		finalFileSize := chunkSize * parts
		if finalFileSize < fileSize {
			chunkSize++
		}
	}
	return chunkSize
}

func getChunksFromChunkSize(parts, chunkSize, fileSize int64) []int64 {
	chunks := make([]int64, parts)
	for i := range chunks {
		if i+1 < len(chunks) {
			chunks[i] = chunkSize
		} else {
			chunks[i] = fileSize - chunkSize*int64(i)
		}
	}
	return chunks
}

func getParts(chunkSize, fileSize int64) int64 {
	parts := fileSize / chunkSize
	if parts*chunkSize < fileSize {
		return parts + 1
	}
	return parts
}

func getChunksFromLines(file *fs.FileReader, maxLines int, fileSize int64) []int64 {
	var lines tLines
	var lastChunkOffset int64
	chunks := make([]int64, 0, 16)
	checking := file.Read(0)
	for checking {
		bytes := file.Buffer[:file.NRead]
		for lines.nextLine(bytes) {
			if lines.line >= maxLines {
				chunkSize := file.Offset - int64(file.NRead) + int64(lines.lineEnd) - lastChunkOffset
				lastChunkOffset += chunkSize
				chunks = append(chunks, chunkSize)
				lines.line = 0
			}
		}
		checking = file.Read(0)
	}
	if file.Offset-lastChunkOffset > 0 {
		chunks = append(chunks, file.Offset-lastChunkOffset)
	}
	return chunks
}

func (lines *tLines) nextLine(bytes []byte) bool {
	if !lines.rPrev {
		for i := lines.lineEnd; i < len(bytes); i++ {
			b := bytes[i]
			if b == '\r' {
				if i+1 < len(bytes) {
					lines.line++
					if bytes[i+1] == '\n' {
						lines.lineEnd = i + 2
						return true
					}
					lines.lineEnd = i + 1
					return true
				}
				lines.rPrev = true
				return false
			} else if b == '\n' {
				lines.lineEnd = i + 1
				lines.line++
				return true
			}
		}
	} else if lines.lineEnd < len(bytes) && bytes[lines.lineEnd] == '\n' {
		lines.lineEnd++
		lines.line++
		lines.rPrev = false
		return true
	} else {
		lines.line++
		lines.rPrev = false
		return true
	}
	lines.lineEnd = 0
	return false
}

func getOutputPaths(outputPath string, chunks []int64) []string {
	var decimals int
	outputPaths := make([]string, len(chunks))
	for i := len(chunks); i > 0; i /= 10 {
		decimals++
	}
	if decimals < 3 {
		decimals = 3
	}
	for i := range outputPaths {
		outputPaths[i] = outputPath + fmt.Sprintf(".%0*d", decimals, i)
	}
	return outputPaths
}

func ensureOutput(command *tCommand, outputPaths []string) bool {
	if !command.overwrite {
		for _, outputPath := range outputPaths {
			_, err := os.Stat(outputPath)
			if err == nil || !errors.Is(err, os.ErrNotExist) {
				command.err = errFileExists("output", outputPath)
				return false
			}
		}
	}
	return true
}
