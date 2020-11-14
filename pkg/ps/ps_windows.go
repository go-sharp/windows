package ps

import (
	"log"
	"strings"
	"syscall"
	"unsafe"

	"golang.org/x/sys/windows"
)

// FindParentProcesses takes a process id and searches for parent processes.
// If recursive is true, FindParentProcesses will walk up the whole process tree
// and return all parent process ids.
func FindParentProcesses(pid uint32, recursive bool) (pids []uint32, err error) {
	procs, err := getAllProcesses()
	if err != nil {
		return pids, err
	}

	var parentPid = pid
	for {
		hasMore := false
	loop:
		for i := range procs {
			if procs[i].ProcessID == parentPid {
				parentPid = procs[i].ParentProcessID

				if parentPid > 0 {
					pids = append(pids, parentPid)
					if recursive {
						hasMore = true
					}
				}
				break loop
			}
		}

		if !hasMore {
			break
		}
	}
	return pids, nil
}

// FindChildProcesses takes a process id and searches for child processes.
// If recursive is true, FindChildProcess will walk down the whole process tree
// and return pids of all children whether they are direct children or grandchildren.
func FindChildProcesses(pid uint32, recursive bool) (pids []uint32, err error) {
	procs, err := getAllProcesses()
	if err != nil {
		return pids, err
	}

	tmpProcs := make([]windows.ProcessEntry32, len(procs))
	toProcess := []uint32{pid}
	for len(toProcess) > 0 {
		p := toProcess[0]
		toProcess = toProcess[1:]

		for i := range procs {
			if procs[i].ParentProcessID == p {
				if recursive {
					toProcess = append(toProcess, procs[i].ProcessID)
				}
				pids = append(pids, procs[i].ProcessID)
				continue
			}
			tmpProcs = append(tmpProcs, procs[i])
		}
		procs = tmpProcs[0:]
		tmpProcs = tmpProcs[:0]
	}
	return pids, nil
}

// FindProcessesByName takes a name and searches for any process with that name.
func FindProcessesByName(name string, ignoreCase bool) (pids []uint32, err error) {
	procs, err := getAllProcesses()
	if err != nil {
		return pids, err
	}

	if ignoreCase {
		name = strings.ToLower(name)
	}

	for i := range procs {
		n := windows.UTF16ToString(procs[i].ExeFile[:])
		if ignoreCase {
			n = strings.ToLower(n)
		}

		if n == name {
			pids = append(pids, procs[i].ProcessID)
		}
	}

	return pids, nil
}

// KillChildProcesses kills the process and his child processes.
// If recursive is true, KillChildProcesses will kill all child processes
// as well as its grandchildren.
func KillChildProcesses(pid uint32, recursive bool) error {
	// Kill parent process first.
	// Unfortunately killProcess returns error despite the process
	// was killed.
	killProcess(pid)

	pids, err := FindChildProcesses(pid, recursive)
	if err != nil {
		return err
	}

	for _, i := range pids {
		// We ignoring any error, because if we kill
		// a process without a living parent process,
		// we get an error.
		killProcess(i)
	}
	return nil
}

func killProcess(pid uint32) error {
	h, err := syscall.OpenProcess(syscall.PROCESS_TERMINATE, false, pid)
	if err != nil {
		return err
	}
	defer syscall.CloseHandle(h)

	if err := syscall.TerminateProcess(h, 1); err != nil {
		return err
	}
	return nil
}

func getAllProcesses() (procs []windows.ProcessEntry32, err error) {
	handle, err := windows.CreateToolhelp32Snapshot(windows.TH32CS_SNAPALL, 0)
	if err != nil {
		return procs, err
	}
	defer windows.CloseHandle(handle)

	var proc windows.ProcessEntry32
	proc.Size = uint32(unsafe.Sizeof(proc))
	if err := windows.Process32First(handle, &proc); err != nil {
		log.Fatalln(err, windows.GetLastError())
	}

	procs = append(procs, proc)

	for {
		var p windows.ProcessEntry32
		p.Size = uint32(unsafe.Sizeof(p))
		if err := windows.Process32Next(handle, &p); err != nil {
			if err != windows.ERROR_NO_MORE_FILES {
				return nil, err
			}
			break
		}
		procs = append(procs, p)
	}
	return procs, nil
}
