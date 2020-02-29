package main

import (
	"fmt"
	"os"

	"github.com/go-sharp/windows/pkg/ps"
	"github.com/jessevdk/go-flags"
)

type Config struct {
	Name *string `short:"n" long:"name" description:"List all process ids with the given name."`
	Kill bool    `short:"k" long:"kill" description:"Kill the processes specified with name or pid."`
	Pid  *uint32 `short:"p" long:"pid" description:"List all child processes of the given pid."`
}

func main() {
	var conf Config
	if _, err := flags.Parse(&conf); err != nil {
		os.Exit(1)
	}

	if conf.Name != nil {
		pids, err := ps.FindProcessesByName(*conf.Name, true)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		for _, i := range pids {
			if conf.Kill {
				killProcess(i)
			} else {
				fmt.Println(*conf.Name, i)
			}
		}
		os.Exit(0)
	}

	if conf.Pid != nil {
		pids, err := ps.FindChildProcesses(*conf.Pid, true)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		pids = append([]uint32{*conf.Pid}, pids...)
		for _, i := range pids {
			if conf.Kill {
				killProcess(i)
			} else {
				fmt.Print(i, " ")
			}
		}
		os.Exit(0)
	}
}

func killProcess(pid uint32) {
	// if p, err := os.FindProcess(int(pid)); err == nil {
	// 	if err := p.Kill(); err != nil {
	// 		fmt.Println(err)
	// 	}
	// }
	// h, err := syscall.OpenProcess(syscall.PROCESS_TERMINATE, false, pid)
	// if err != nil {
	// 	fmt.Println(err)
	// 	os.Exit(1)
	// }
	// defer syscall.CloseHandle(h)
	// if err := syscall.TerminateProcess(h, 1); err != nil {
	// 	fmt.Println(err)
	// }
	if err := ps.KillChildProcesses(pid, true); err != nil {
		fmt.Println(err)
	}
}
