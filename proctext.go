/*
 *          Copyright 2026 Vitali Baumtrok.
 * Distributed under the Boost Software License, Version 1.0.
 *     (See accompanying file LICENSE or copy at
 *        http://www.boost.org/LICENSE_1_0.txt)
 */

package main

import (
	"fmt"
	"github.com/vbsw/go-lib/fs"
	"math/rand"
	"time"
	"unsafe"
)

const (
	textGenBufferSize  = 8 * 1024 * 1024
	maxWordsPerLine    = 40
	newLineProbability = 0.05
	minWordLength      = 2
	maxWordLength      = 20
)

type tSender struct {
	bytesSent int64
}

func processText(command *tCommand) {
	if command.bytes > 0 {
		if command.outputPath > "" {
			processTextToFile(command)
		} else {
			processTextToStd(command)
		}
	}
}

func processTextToFile(command *tCommand) {
	var file fs.FileWriter
	defer func() {
		file.Close()
		command.err = file.Err
	}()
	if file.Open(command.outputPath) {
		var sender tSender
		var bytesWritten int64
		adjustThreadCount(command)
		toGenerator := make(chan []byte, command.threads*2)
		toWriter := make(chan []byte, command.threads*2)
		generateFunc := getGenerateFunc(command)
		trmt := unsafe.Slice(unsafe.StringData(command.lineTerminator), len(command.lineTerminator))
		dlmt := unsafe.Slice(unsafe.StringData(command.delimiter), len(command.delimiter))
		for i := 0; i < command.threads; i++ {
			if sender.bytesSent < command.bytes {
				go generatorGo(toGenerator, toWriter, int64(i), generateFunc, trmt, dlmt)
				sender.sendNewBytes(command.bytes-sender.bytesSent, toGenerator)
			}
			if sender.bytesSent < command.bytes {
				sender.sendNewBytes(command.bytes-sender.bytesSent, toGenerator)
			}
		}
		for sender.bytesSent < command.bytes && file.Err == nil {
			bytes := <-toWriter
			file.Write(bytes)
			bytesWritten += int64(file.NWritten)
			sender.sendOldBytes(bytes, command.bytes-sender.bytesSent, toGenerator)
		}
		for bytesWritten < command.bytes && file.Err == nil {
			bytes := <-toWriter
			file.Write(bytes)
			bytesWritten += int64(file.NWritten)
		}
	}
}

func processTextToStd(command *tCommand) {
	var file fs.FileWriter
	var sender tSender
	var bytesWritten int64
	adjustThreadCount(command)
	toGenerator := make(chan []byte, command.threads*2)
	toWriter := make(chan []byte, command.threads*2)
	generateFunc := getGenerateFunc(command)
	trmt := unsafe.Slice(unsafe.StringData(command.lineTerminator), len(command.lineTerminator))
	dlmt := unsafe.Slice(unsafe.StringData(command.delimiter), len(command.delimiter))
	for i := 0; i < command.threads; i++ {
		if sender.bytesSent < command.bytes {
			go generatorGo(toGenerator, toWriter, int64(i), generateFunc, trmt, dlmt)
			sender.sendNewBytes(command.bytes-sender.bytesSent, toGenerator)
		}
		if sender.bytesSent < command.bytes {
			sender.sendNewBytes(command.bytes-sender.bytesSent, toGenerator)
		}
	}
	for sender.bytesSent < command.bytes && file.Err == nil {
		bytes := <-toWriter
		file.StdoutWrite(bytes)
		bytesWritten += int64(file.NWritten)
		sender.sendOldBytes(bytes, command.bytes-sender.bytesSent, toGenerator)
	}
	for bytesWritten < command.bytes && file.Err == nil {
		bytes := <-toWriter
		file.StdoutWrite(bytes)
		bytesWritten += int64(file.NWritten)
	}
	if file.Err == nil {
		file.StdoutWrite([]byte{'\n'})
	} else {
		fmt.Println()
	}
}

func (sender *tSender) sendNewBytes(rest int64, toGenerator chan []byte) {
	var bytes []byte
	if rest >= textGenBufferSize {
		bytes = make([]byte, textGenBufferSize)
	} else {
		bytes = make([]byte, rest)
	}
	toGenerator <- bytes
	sender.bytesSent += int64(len(bytes))
}

func (sender *tSender) sendOldBytes(bytes []byte, rest int64, toGenerator chan []byte) {
	if rest < textGenBufferSize {
		bytes = bytes[:rest]
	}
	toGenerator <- bytes
	sender.bytesSent += int64(len(bytes))
}

func generatorGo(toGenerator, toWriter chan []byte, seedAdd int64, generateFunc func(*rand.Rand, []byte), trmt, dlmt []byte) {
	var words int
	random := rand.New(rand.NewSource(time.Now().UnixNano() + seedAdd))
	for {
		var bytesWritten int
		bytes := <-toGenerator
		for bytesWritten < len(bytes) {
			var lineBreak bool
			wordLength := getWordLength(random, len(bytes)-bytesWritten)
			generateFunc(random, bytes[bytesWritten:bytesWritten+wordLength])
			bytesWritten += wordLength
			if len(bytes)-bytesWritten >= len(trmt) {
				if getLineBreak(random, words) {
					lineBreak = true
					for i, b := range trmt {
						bytes[bytesWritten+i] = b
					}
					bytesWritten += len(trmt)
					words = 0
				}
			}
			if !lineBreak && len(bytes)-bytesWritten >= len(dlmt) {
				for i, b := range dlmt {
					bytes[bytesWritten+i] = b
				}
				bytesWritten += len(dlmt)
				words++
			}
		}
		toWriter <- bytes
	}
}

func adjustThreadCount(command *tCommand) {
	if command.threads > 1 {
		threads := ((command.bytes+textGenBufferSize-1)/textGenBufferSize + 1) / 2
		if int64(command.threads) > threads {
			command.threads = int(threads)
			if command.verbose {
				printWarningInfo(command, errThreadsReduced(command.threads))
			}
		}
	}
}

func getGenerateFunc(command *tCommand) func(*rand.Rand, []byte) {
	if command.lettersOnly {
		if command.upperCase {
			return generateU
		} else if command.lowerCase {
			return generateL
		}
		return generateA
	} else if command.upperCase {
		return generateU
	} else if command.lowerCase {
		return generateL
	}
	return generateZ
}

func getWordLength(random *rand.Rand, restLength int) int {
	lengthFloat := random.Float32() * float32(maxWordLength-minWordLength+1)
	wordLength := int(lengthFloat) + minWordLength
	if wordLength < restLength {
		return wordLength
	}
	return restLength
}

func getLineBreak(random *rand.Rand, words int) bool {
	if words < maxWordsPerLine {
		if random.Float32() > newLineProbability {
			return false
		}
	}
	return true
}

func generateA(random *rand.Rand, bytes []byte) {
	for i := range bytes {
		randomFloat := random.Float32()
		numberFloat := randomFloat * float32((90-65)*2+2) // (Z-A)*2+2
		letter := byte(numberFloat)
		if letter > 90-65 { // Z-A
			bytes[i] = letter + 65 + 6 // letter+A+6 (lower case)
		} else {
			bytes[i] = letter + 65 // letter+A (upper case)
		}
	}
}

func generateL(random *rand.Rand, bytes []byte) {
	for i := range bytes {
		randomFloat := random.Float32()
		numberFloat := randomFloat * float32(122-97+1) // z-a+1
		letter := byte(numberFloat)
		bytes[i] = letter + 97 // letter+a (lower case)
	}
}

func generateU(random *rand.Rand, bytes []byte) {
	for i := range bytes {
		randomFloat := random.Float32()
		numberFloat := randomFloat * float32(90-65+1) // Z-A+1
		letter := byte(numberFloat)
		bytes[i] = letter + 65 // letter+A (upper case)
	}
}

func generateZ(random *rand.Rand, bytes []byte) {
	for i := range bytes {
		randomFloat := random.Float32()
		numberFloat := randomFloat * float32(126-33+1-5) // ~-!+1-5
		letter := byte(numberFloat)
		if letter > 95-33+1-5 { // _-!+1-5
			bytes[i] = letter + 33 + 5 // letter+!+5
		} else if letter > 93-33+1-4 { // ]-!+1-4
			bytes[i] = letter + 33 + 4 // letter+!+4
		} else if letter > 91-33+1-3 { // [-!+1-3
			bytes[i] = letter + 33 + 3 // letter+!+3
		} else if letter > 38-33+1-2 { // &-!+1-2
			bytes[i] = letter + 33 + 2 // letter+!+2
		} else if letter > 33-33+1-1 { // !-!+1-1
			bytes[i] = letter + 33 + 1 // letter+!+1
		} else {
			bytes[i] = letter + 33 // letter+!
		}
	}
}
