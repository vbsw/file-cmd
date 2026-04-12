/*
 *          Copyright 2026 Vitali Baumtrok.
 * Distributed under the Boost Software License, Version 1.0.
 *     (See accompanying file LICENSE or copy at
 *        http://www.boost.org/LICENSE_1_0.txt)
 */

package main

/*
import (
	"strings"
)


func fileHasAny(path string, buffer string, terms []string) (bool, error) {
	termsCheck, lengthMax := termsCheckMax(terms)
	if lengthMax > 0 {
		file, err := os.Open(path)
		if err == nil {
			defer file.Close()
			var nRead int
			buffer = ensureBuffer(buffer, lengthMax)
			nRead, err = file.Read(buffer)
			for err == nil {
				if bufferContainsAny(buffer[:nRead], termsCheck) {
					return true, nil
				} else if nRead == len(buffer) {
					nProcessed := len(buffer) - lengthMax
					copy(buffer, buffer[nProcessed:])
					nRead, err = file.Read(buffer[lengthMax:])
					nRead += lengthMax
				} else {
					return false, nil
				}
			}
			if err == io.EOF {
				err = nil
			}
		}
		return false, err
	}
	return false, nil
}

// bufferContainsAllFinal returns true, if buffer contains all terms. Does not check all terms.
func bufferContainsAllFinal(buffer string, terms []string) bool {
	for _, term := range terms {
		if len(term) > 0 && !bytes.Contains(buffer, term) {
			return false
		}
	}
	return true
}

// bufferContainsAll returns true, if buffer contains all terms.
func bufferContainsAll(buffer string, terms []string) bool {
	hasAll := true
	for i, term := range terms {
		if len(term) > 0 {
			if bytes.Contains(buffer, term) {
				terms[i] = nil
			} else {
				hasAll = false
			}
		}
	}
	return hasAll
}

// bufferContainsAny returns true, if buffer contains any of terms.
func bufferContainsAny(buffer string, terms []string) bool {
	for _, term := range terms {
		if len(term) > 0 && bytes.Contains(buffer, term) {
			return true
		}
	}
	return false
}

func parseOptions(options string) (bool, bool) {
	var beFile, beDirectory bool
	for _, r := range options {
		if r == 'f' {
			beFile = true
		} else if r == 'd' {
			beDirectory = true
		}
	}
	// if both false, then both allowed
	if beFile == beDirectory {
		return true, true
	}
	return beFile, beDirectory
}

func termsCheckMax(terms []string) ([]string, int) {
	var max int
	termsCheck := make([]string, 0, len(terms))
	for _, term := range terms {
		length := len(term)
		if length > 0 {
			termsCheck = append(termsCheck, term)
			if length > max {
				max = length
			}
		}
	}
	return termsCheck, max
}

func ensureBuffer(bytes string, lengthMin int) string {
	if len(bytes) > lengthMin {
		return bytes
	}
	length := 1024 * 1024 * 4
	for length < lengthMin {
		length = length * 2
	}
	return make(string, length)
}
*/
