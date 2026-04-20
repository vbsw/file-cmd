/*
 *          Copyright 2026 Vitali Baumtrok.
 * Distributed under the Boost Software License, Version 1.0.
 *     (See accompanying file LICENSE or copy at
 *        http://www.boost.org/LICENSE_1_0.txt)
 */

package main

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
	idxOptSplitBytes
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
	idxOptListSilent
	idxOptListFilter
	idxOptListDelimiter
	idxOptListThreads
	idxOptListTotal
)

const (
	idxOptCountInput = iota
	idxOptCountOr
	idxOptCountRecursive
	idxOptCountSilent
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
	idxOptCopySilent
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
	idxOptMoveSilent
	idxOptMoveFilter
	idxOptMoveDelimiter
	idxOptMoveThreads
	idxOptMoveTotal
)

const (
	idxOptRemoveInput = iota
	idxOptRemoveOr
	idxOptRemoveRecursive
	idxOptRemoveSilent
	idxOptRemoveFilter
	idxOptRemoveDelimiter
	idxOptRemoveThreads
	idxOptRemoveTotal
)

const (
	idxOptCleanInput = iota
	idxOptCleanRecursive
	idxOptCleanSilent
	idxOptCleanTotal
)

const (
	idxOptTextOutput = iota
	idxOptTextOverwrite
	idxOptTextSize
	idxOptTextTerminator
	idxOptTextDelimiter
	idxOptTextFormat
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
