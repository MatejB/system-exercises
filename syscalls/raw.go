// +build darwin,amd64

package main

import (
	"fmt"
	"log"
	"syscall"
	"unsafe"
)

func main() {
	buf := syscall.Stat_t{}
	err := statfs(&buf)
	if err != nil {
		log.Println(err)
	}
	fmt.Printf("%+v", buf)
}

func statfs(buf *syscall.Stat_t) (err error) {
	var _p0 *byte
	_p0, err = syscall.BytePtrFromString("raw.go")
	if err != nil {
		return
	}
	_, _, e := syscall.Syscall(338, uintptr(unsafe.Pointer(_p0)), uintptr(unsafe.Pointer(buf)), 0)
	if e != 0 {
		return fmt.Errorf("Syscall: %s", e)
	}

	return nil
}
