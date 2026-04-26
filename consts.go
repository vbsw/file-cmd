/*
 *          Copyright 2026 Vitali Baumtrok.
 * Distributed under the Boost Software License, Version 1.0.
 *     (See accompanying file LICENSE or copy at
 *        http://www.boost.org/LICENSE_1_0.txt)
 */

package main

const (
	defaultBufferSize     = 8 * 1024 * 1024
	procInitialBufferSize = defaultBufferSize
)

const (
	textGenBufferSize  = defaultBufferSize
	maxWordsPerLine    = 40
	newLineProbability = 0.05
	minWordLength      = 2
	maxWordLength      = 20
)

const (
	maxInt64        = int64((^uint64(0)) >> 1)
	maxInt32        = int32((^uint32(0)) >> 1)
	splitBufferSize = defaultBufferSize
)

const (
	idxInfoHelp = iota
	idxInfoVersion
	idxInfoExample
	idxInfoCopyright
	idxInfoTotal
)

const (
	idxCmdSplit = iota
	idxCmdConcat
	idxCmdList
	idxCmdCount
	idxCmdCopy
	idxCmdMove
	idxCmdRemove
	idxCmdClean
	idxCmdText
	idxCmdTotal
)

const (
	idxOptSplitInput = iota
	idxOptSplitOutput
	idxOptSplitParts
	idxOptSplitSize
	idxOptSplitLines
	idxOptSplitOverwrite
	idxOptSplitTotal
)

const (
	idxOptConcatInput = iota
	idxOptConcatOutput
	idxOptConcatOverwrite
	idxOptConcatTotal
)

const (
	idxOptListInput = iota
	idxOptListOr
	idxOptListRecursive
	idxOptListVerbose
	idxOptListFilter
	idxOptListDelimiter
	idxOptListThreads
	idxOptListTotal
)

const (
	idxOptCountInput = iota
	idxOptCountOr
	idxOptCountRecursive
	idxOptCountVerbose
	idxOptCountFilter
	idxOptCountDelimiter
	idxOptCountThreads
	idxOptCountTotal
)

const (
	idxOptCopyInput = iota
	idxOptCopyOutput
	idxOptCopyOr
	idxOptCopyOverwrite
	idxOptCopyRecursive
	idxOptCopyVerbose
	idxOptCopyFilter
	idxOptCopyDelimiter
	idxOptCopyThreads
	idxOptCopyTotal
)

const (
	idxOptMoveInput = iota
	idxOptMoveOutput
	idxOptMoveOr
	idxOptMoveOverwrite
	idxOptMoveRecursive
	idxOptMoveVerbose
	idxOptMoveFilter
	idxOptMoveDelimiter
	idxOptMoveThreads
	idxOptMoveTotal
)

const (
	idxOptRemoveInput = iota
	idxOptRemoveOr
	idxOptRemoveRecursive
	idxOptRemoveVerbose
	idxOptRemoveFilter
	idxOptRemoveDelimiter
	idxOptRemoveThreads
	idxOptRemoveTotal
)

const (
	idxOptCleanInput = iota
	idxOptCleanAll
	idxOptCleanRecursive
	idxOptCleanVerbose
	idxOptCleanList
	idxOptCleanTotal
)

const (
	idxOptTextOutput = iota
	idxOptTextOverwrite
	idxOptTextSize
	idxOptTextTerminator
	idxOptTextDelimiter
	idxOptTextFormat
	idxOptTextThreads
	idxOptTextVerbose
	idxOptTextTotal
)

const (
	cmdNone = iota
	cmdInfo
	cmdInfoExample
	cmdInfoSplit
	cmdInfoConcat
	cmdInfoList
	cmdInfoCount
	cmdInfoCopy
	cmdInfoMove
	cmdInfoRemove
	cmdInfoClean
	cmdInfoText
	cmdVersion
	cmdExampleSplit
	cmdExampleConcat
	cmdExampleList
	cmdExampleCount
	cmdExampleCopy
	cmdExampleMove
	cmdExampleRemove
	cmdExampleClean
	cmdExampleText
	cmdCopyright
	cmdSplit
	cmdConcat
	cmdList
	cmdCount
	cmdCopy
	cmdMove
	cmdRemove
	cmdClean
	cmdText
	cmdTotal
)
