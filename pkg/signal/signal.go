package signal

//go:generate go run golang.org/x/sys/windows/mkwinsyscall -output zsignal_windows.go signal.go
import (
	"errors"
	"fmt"
	"unsafe"

	"golang.org/x/sys/windows"
)

//sys postMessage(hWnd HWND, msg UINT, wParam WPARAM, lParam LPARAM) (err error) = user32.PostMessageW
//sys enumWindows(cb uintptr, lParam LPARAM) (err error) = user32.EnumWindows
//sys enumChildWindows(hWndParent HWND, cb uintptr, lParam LPARAM) (err error) = user32.EnumChildWindows
//sys getWindowThreadProcessId(hwnd HWND, lpdwProcessId LPDWORD) (id DWORD) = user32.GetWindowThreadProcessId
//sys attachConsole(dwProcessId DWORD) (err error) = kernel32.AttachConsole
//sys freeConsole() (err error) = kernel32.FreeConsole
//sys setConsoleCtrlHandler(handler uintptr, add bool) (err error) = kernel32.SetConsoleCtrlHandler

// Windows types
type (
	// HWND is a windows handle.
	HWND windows.Handle
	// UINT is a 32-bit unsigned int.
	UINT uint32
	// WPARAM is a uint ptr see: https://docs.microsoft.com/en-us/windows/win32/winprog/windows-data-types
	WPARAM uintptr
	// LPARAM is a long ptr see: https://docs.microsoft.com/en-us/windows/win32/winprog/windows-data-types
	LPARAM uintptr
	// LPDWORD is a pointer to a uint32.
	LPDWORD unsafe.Pointer
	// DWORD is a 32-bit unsigend int.
	DWORD uint32
)

var (
	// ErrNoWndHandle return if no valid window handle found.
	ErrNoWndHandle = errors.New("no window handle found for the given process id")
)

// WindowsMessage defines Windows system messages. See for example https://docs.microsoft.com/en-us/windows/win32/winmsg/wm-close
type WindowsMessage uint64

const (
	// WmClose sent as a signal that a window or an application should terminate.
	WmClose WindowsMessage = 0x0010
	// WmQuit indicates a request to terminate an application, and is generated when the application calls the PostQuitMessage function.
	WmQuit WindowsMessage = 0x0012
)

func (w WindowsMessage) String() string {
	switch w {
	case WmClose:
		return "WM_CLOSE"
	case WmQuit:
		return "WM_QUIT"
	default:
		return "Unknown Message"
	}
}

// CtrlEvent is a Window control event.
type CtrlEvent uint32

// Defines the different control events.
const (
	CtrlCEvent        = CtrlEvent(0)
	CtrlBreakEvent    = CtrlEvent(1)
	CtrlCloseEvent    = CtrlEvent(2)
	CtrlLogoffEvent   = CtrlEvent(5)
	CtrlShutdownEvent = CtrlEvent(6)
)

// SendCtrlEvent sends a windows control event. Caveat: If the
// the process that calls this functions has already a console attached,
// this call will fail. This call will also call SetConsoleCtrlHandler(0, true)
// to prevent ourself from receiving the event.
func SendCtrlEvent(pid uint32, ctrlEvent CtrlEvent) error {
	if err := attachConsole(DWORD(pid)); err != nil {
		return fmt.Errorf("failed to attach console: %w", err)
	}
	defer freeConsole()

	setConsoleCtrlHandler(0, true)
	if err := windows.GenerateConsoleCtrlEvent(uint32(ctrlEvent), 0); err != nil {
		return fmt.Errorf("failed to send event: %w", err)
	}

	return nil
}

// SendSignal sends a Windows message to all windows handles that belongs to a given process.
func SendSignal(pid uint32, msg WindowsMessage) error {

	var handlePid uint32
	var handles []HWND
	var callback uintptr
	callback = windows.NewCallback(func(hwnd HWND, lParam LPARAM) uintptr {
		getWindowThreadProcessId(hwnd, LPDWORD(&handlePid))
		if pid == handlePid {
			handles = append(handles, hwnd)
		}
		enumChildWindows(hwnd, callback, 0)
		return 1
	})

	enumWindows(callback, 0)

	if len(handles) == 0 {
		return ErrNoWndHandle
	}

	for _, h := range handles {
		if err := postMessage(HWND(h), UINT(msg), 0, 0); err != nil {
			return fmt.Errorf("failed to send signal %v: %w", msg, err)
		}
	}

	return nil
}
