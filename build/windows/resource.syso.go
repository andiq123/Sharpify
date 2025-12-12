//go:build ignore

package main

// This file is used to generate resource.syso for Windows builds
// Run: go generate ./build/windows/

//go:generate go run github.com/tc-hib/go-winres@latest make --in winres.json --out ../../resource_windows.syso
