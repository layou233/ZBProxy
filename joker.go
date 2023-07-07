//go:build joker

package main

import "C"

// This is a helper code to generate executable library
// for Joker project. CGO is required.
// This is for convenience only and has nothing to do
// with ZBProxy itself. You can ignore this if you
// are not interested in Joker anyway.

// Compile command: go build -buildmode=c-shared -tags=joker -v -ldflags="-s -w -buildid="
// And then add library filename extension of your platform
// to the output file (like .so, .dll, etc).

//export JokerMain
func JokerMain() {
	main()
}
