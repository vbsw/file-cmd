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
	idxOptSplitOvierwrite
	idxOptSplitTotal
)

const (
	idxOptConcatInput = iota
	idxOptConcatOutput
	idxOptConcatOvierwrite
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
