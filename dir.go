// Credit for implementation:
// https://stackoverflow.com/questions/17946269/call-windows-function-getting-the-font-directory

package main

import (
	"syscall"
	"unsafe"
)

type guid struct {
	Data1 uint32
	Data2 uint16
	Data3 uint16
	Data4 [8]byte
}

var (
	documentsFolderID = guid{0xFDD39AD0, 0x238F, 0x46AF, [8]byte{0xAD, 0xB4, 0x6C, 0x85, 0x48, 0x03, 0x69, 0xC7}}
)

var (
	modShell32               = syscall.NewLazyDLL("Shell32.dll")
	modOle32                 = syscall.NewLazyDLL("Ole32.dll")
	procSHGetKnownFolderPath = modShell32.NewProc("SHGetKnownFolderPath")
	procCoTaskMemFree        = modOle32.NewProc("CoTaskMemFree")
)

func shGetKnownFolderPath(rfid *guid, dwFlags uint32, hToken syscall.Handle, pszPath *uintptr) (retval error) {
	r0, _, _ := syscall.Syscall6(
		procSHGetKnownFolderPath.Addr(),
		4,
		uintptr(unsafe.Pointer(rfid)),
		uintptr(dwFlags),
		uintptr(hToken),
		uintptr(unsafe.Pointer(pszPath)),
		0,
		0,
	)

	if r0 != 0 {
		return syscall.Errno(r0)
	}

	return nil
}

func coTaskMemFree(pv uintptr) {
	syscall.Syscall(procCoTaskMemFree.Addr(), 1, uintptr(pv), 0, 0)
	return
}

func DocumentsFolder() (string, error) {
	var path uintptr

	err := shGetKnownFolderPath(&documentsFolderID, 0, 0, &path)

	if err != nil {
		return "", err
	}

	defer coTaskMemFree(path)

	folder := syscall.UTF16ToString((*[1 << 16]uint16)(unsafe.Pointer(path))[:])

	return folder, nil
}
