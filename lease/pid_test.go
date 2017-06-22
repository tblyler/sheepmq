package lease

import (
	"os/exec"
	"testing"
)

func TestPIDValid(t *testing.T) {
	// create a sleep command for 1 second
	cmd := exec.Command("sleep", "1")
	err := cmd.Start()
	if err != nil {
		t.Fatal("Failed to create a sleep process:", err)
	}

	pid := NewPID(cmd.Process.Pid)
	if !pid.Valid() {
		t.Error("PID died too soon")
	}

	cmd.Wait()

	if pid.Valid() {
		t.Error("PID didn't die somehow")
	}

	pid = NewPID(-1)
	if pid.Valid() {
		t.Error("Negative PIDS are not a thing")
	}
}
