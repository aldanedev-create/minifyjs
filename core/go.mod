// The core MinifyJS Go module. Everything that ships in the native
// binary lives under this module. The tools/ module is separate so
// developer utilities never leak into the released binary's
// dependency graph.
module github.com/minifyjs/minifyjs/core

go 1.22

require (
	github.com/evanw/esbuild v0.28.2
)

require (
	golang.org/x/sync v0.10.0 // indirect
	golang.org/x/sys v0.28.0 // indirect
)