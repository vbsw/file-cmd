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
	idxOptListTotal
)

const (
	idxOptCountInput = iota
	idxOptCountOr
	idxOptCountRecursive
	idxOptCountSilent
	idxOptCountFilter
	idxOptCountDelimiter
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
	idxOptMoveTotal
)

const (
	idxOptRemoveInput = iota
	idxOptRemoveOr
	idxOptRemoveRecursive
	idxOptRemoveSilent
	idxOptRemoveFilter
	idxOptRemoveDelimiter
	idxOptRemoveTotal
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
	cmdInfoText
	cmdVersion
	cmdExampleSplit
	cmdExampleConcat
	cmdExampleList
	cmdExampleCount
	cmdExampleCopy
	cmdExampleMove
	cmdExampleRemove
	cmdExampleText
	cmdCopyright
	cmdSplit
	cmdConcat
	cmdList
	cmdCount
	cmdCopy
	cmdMove
	cmdRemove
	cmdText
	cmdTotal
)
