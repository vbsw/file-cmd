module github.com/vbsw/file-cmd

go 1.22.2

require (
	github.com/vbsw/go-lib/cl v0.1.0
	github.com/vbsw/go-lib/match v0.1.0
	github.com/vbsw/go-lib/fs v0.1.0
)

replace (
	github.com/vbsw/go-lib/cl => ../go-lib/cl
	github.com/vbsw/go-lib/match => ../go-lib/match
	github.com/vbsw/go-lib/fs => ../go-lib/fs
)
