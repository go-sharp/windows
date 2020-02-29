package ps

import (
	"fmt"
	"os"
	"os/exec"
	"testing"
	"time"
)

// This builds the spawner binary for the tests.
func TestMain(m *testing.M) {
	fmt.Println("Building spawner binary...")
	if err := exec.Command("go", "build", "-tags", "forceposix", "-o", "spawner.exe", "./spawner").Run(); err != nil {
		panic(err)
	}

	killAllSpawner()
	code := m.Run()

	// Need some time to release binary
	killAllSpawner()
	if err := os.Remove("spawner.exe"); err != nil {
		fmt.Println("Failed to remove spawner:", err)
	}
	os.Exit(code)
}

func TestKillChildProcesses(t *testing.T) {
	type args struct {
		spawn     string
		recursive bool
	}
	tests := []struct {
		name          string
		args          args
		wantErr       bool
		wantProcesses int
	}{
		{
			name:          "Kill only direct children",
			args:          args{spawn: "3", recursive: false},
			wantProcesses: 12,
			wantErr:       false,
		},
		{
			name:          "Kill all children",
			args:          args{spawn: "3", recursive: true},
			wantProcesses: 0,
			wantErr:       false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			killAllSpawner()
			cmd := exec.Command("spawner.exe", "--spawn", tt.args.spawn)
			if err := cmd.Start(); err != nil {
				t.Fatalf("FindChildProcesses() can't spawn processes: %v", err)
			}
			// Give time to start processes
			time.Sleep(100 * time.Millisecond)

			if err := KillChildProcesses(uint32(cmd.Process.Pid), tt.args.recursive); (err != nil) != tt.wantErr {
				t.Errorf("KillChildProcesses() error = %v, wantErr %v", err, tt.wantErr)
			}

			pids, err := FindProcessesByName("spawner.exe", false)
			if err != nil {
				t.Fatalf("KillChildProcesses() error = %v", err)
			}

			if len(pids) != tt.wantProcesses {
				t.Fatalf("KillChildProcesses() processes = %v, wantProcesses %v", len(pids), tt.wantProcesses)
			}
		})
	}
}

func TestFindChildProcesses(t *testing.T) {
	type args struct {
		spawn     string
		recursive bool
	}
	tests := []struct {
		name          string
		args          args
		wantProcesses int
		wantErr       bool
	}{
		{
			name:          "Only direct children",
			args:          args{spawn: "3", recursive: false},
			wantProcesses: 3,
			wantErr:       false,
		},
		{
			name:          "All children",
			args:          args{spawn: "3", recursive: true},
			wantProcesses: 15,
			wantErr:       false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			killAllSpawner()
			cmd := exec.Command("spawner.exe", "--spawn", tt.args.spawn)
			if err := cmd.Start(); err != nil {
				t.Fatalf("FindChildProcesses() can't spawn processes: %v", err)
			}
			// Give time to start processes
			time.Sleep(100 * time.Millisecond)

			pids, err := FindChildProcesses(uint32(cmd.Process.Pid), tt.args.recursive)
			if (err != nil) != tt.wantErr {
				t.Fatalf("FindChildProcesses() error = %v, wantErr %v", err, tt.wantErr)
			}

			if len(pids) != tt.wantProcesses {
				t.Fatalf("FindChildProcesses() processes = %v, wantProcesses %v", len(pids), tt.wantProcesses)
			}
		})
	}
}

func TestFindProcessesByName(t *testing.T) {
	type args struct {
		spawn      string
		ignoreCase bool
		name       string
	}
	tests := []struct {
		name     string
		args     args
		wantPids int
		wantErr  bool
	}{
		{
			name:     "Find all 16 processes",
			args:     args{spawn: "3", name: "spawner.exe", ignoreCase: false},
			wantPids: 16,
			wantErr:  false,
		},
		{
			name:     "Find zero processes with wrong name",
			args:     args{spawn: "3", name: "spawner0.exe", ignoreCase: false},
			wantPids: 0,
			wantErr:  false,
		},
		{
			name:     "Find all processes with upper and lower letters in name",
			args:     args{spawn: "3", name: "SpaWnEr.exe", ignoreCase: true},
			wantPids: 16,
			wantErr:  false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			killAllSpawner()
			cmd := exec.Command("spawner.exe", "--spawn", tt.args.spawn)
			if err := cmd.Start(); err != nil {
				t.Fatalf("FindChildProcesses() can't spawn processes: %v", err)
			}
			// Give time to start processes
			time.Sleep(100 * time.Millisecond)

			gotPids, err := FindProcessesByName(tt.args.name, tt.args.ignoreCase)
			if (err != nil) != tt.wantErr {
				t.Errorf("FindProcessesByName() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if len(gotPids) != tt.wantPids {
				t.Errorf("FindProcessesByName() = %v, want %v", len(gotPids), tt.wantPids)
			}
		})
	}
}

func killAllSpawner() {
	exec.Command("taskkill", "/IM", "spawner.exe", "/F").Run()
}
