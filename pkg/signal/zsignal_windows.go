// Code generated by 'go generate'; DO NOT EDIT.

package signal

import (
	"syscall"
	"unsafe"

	"golang.org/x/sys/windows"
)

var _ unsafe.Pointer

// Do the interface allocations only once for common
// Errno values.
const (
	errnoERROR_IO_PENDING = 997
)

var (
	errERROR_IO_PENDING error = syscall.Errno(errnoERROR_IO_PENDING)
	errERROR_EINVAL     error = syscall.EINVAL
)

// errnoErr returns common boxed Errno values, to prevent
// allocations at runtime.
func errnoErr(e syscall.Errno) error {
	switch e {
	case 0:
		return errERROR_EINVAL
	case errnoERROR_IO_PENDING:
		return errERROR_IO_PENDING
	}
	// TODO: add more here, after collecting data on the common
	// error values see on Windows. (perhaps when running
	// all.bat?)
	return e
}

var (
	modkernel32 = windows.NewLazySystemDLL("kernel32.dll")
	moduser32   = windows.NewLazySystemDLL("user32.dll")

	procAttachConsole            = modkernel32.NewProc("AttachConsole")
	procFreeConsole              = modkernel32.NewProc("FreeConsole")
	procSetConsoleCtrlHandler    = modkernel32.NewProc("SetConsoleCtrlHandler")
	procEnumChildWindows         = moduser32.NewProc("EnumChildWindows")
	procEnumWindows              = moduser32.NewProc("EnumWindows")
	procGetWindowThreadProcessId = moduser32.NewProc("GetWindowThreadProcessId")
	procPostMessageW             = moduser32.NewProc("PostMessageW")
)

func attachConsole(dwProcessId DWORD) (err error) {
	r1, _, e1 := syscall.Syscall(procAttachConsole.Addr(), 1, uintptr(dwProcessId), 0, 0)
	if r1 == 0 {
		err = errnoErr(e1)
	}
	return
}

func freeConsole() (err error) {
	r1, _, e1 := syscall.Syscall(procFreeConsole.Addr(), 0, 0, 0, 0)
	if r1 == 0 {
		err = errnoErr(e1)
	}
	return
}

func setConsoleCtrlHandler(handler uintptr, add bool) (err error) {
	var _p0 uint32
	if add {
		_p0 = 1
	}
	r1, _, e1 := syscall.Syscall(procSetConsoleCtrlHandler.Addr(), 2, uintptr(handler), uintptr(_p0), 0)
	if r1 == 0 {
		err = errnoErr(e1)
	}
	return
}

func enumChildWindows(hWndParent HWND, cb uintptr, lParam LPARAM) (err error) {
	r1, _, e1 := syscall.Syscall(procEnumChildWindows.Addr(), 3, uintptr(hWndParent), uintptr(cb), uintptr(lParam))
	if r1 == 0 {
		err = errnoErr(e1)
	}
	return
}

func enumWindows(cb uintptr, lParam LPARAM) (err error) {
	r1, _, e1 := syscall.Syscall(procEnumWindows.Addr(), 2, uintptr(cb), uintptr(lParam), 0)
	if r1 == 0 {
		err = errnoErr(e1)
	}
	return
}

func getWindowThreadProcessId(hwnd HWND, lpdwProcessId LPDWORD) (id DWORD) {
	r0, _, _ := syscall.Syscall(procGetWindowThreadProcessId.Addr(), 2, uintptr(hwnd), uintptr(lpdwProcessId), 0)
	id = DWORD(r0)
	return
}

func postMessage(hWnd HWND, msg UINT, wParam WPARAM, lParam LPARAM) (err error) {
	r1, _, e1 := syscall.Syscall6(procPostMessageW.Addr(), 4, uintptr(hWnd), uintptr(msg), uintptr(wParam), uintptr(lParam), 0, 0)
	if r1 == 0 {
		err = errnoErr(e1)
	}
	return
}
