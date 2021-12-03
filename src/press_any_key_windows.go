//go:build windows
// +build windows

package main

import (
	"fmt"
	"os"
	"syscall"
)

func PressAnyKey() {
	var (
		h = os.Stdin.Fd()
		m uint32
	)

	dll := syscall.MustLoadDLL("kernel32")
	proc := dll.MustFindProc("SetConsoleMode")
	if proc != nil {
		proc.Call(uintptr(h), uintptr(m))
	}

	fmt.Print("press any key to exit... ")
	b := make([]byte, 10)
	os.Stdin.Read(b)
}
