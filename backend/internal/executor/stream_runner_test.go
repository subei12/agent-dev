package executor

import (
	"context"
	"os/exec"
	"testing"
)

// TestStreamRunnerReadsStdoutAndStderr 验证流式执行器能同时读取标准输出和标准错误。
func TestStreamRunnerReadsStdoutAndStderr(t *testing.T) {
	runner := NewStreamRunner()
	cmd := exec.Command("bash", "-lc", "printf 'stdout line\n'; printf 'stderr line\n' 1>&2")

	var lines []string
	err := runner.Run(context.Background(), cmd, func(line []byte) {
		lines = append(lines, string(line))
	})
	if err != nil {
		t.Fatalf("run command: %v", err)
	}

	if len(lines) != 2 {
		t.Fatalf("expected 2 output lines, got %d", len(lines))
	}
	lineSet := map[string]bool{}
	for _, line := range lines {
		lineSet[line] = true
	}
	if !lineSet["stdout line"] {
		t.Fatalf("expected stdout line in %v", lines)
	}
	if !lineSet["stderr line"] {
		t.Fatalf("expected stderr line in %v", lines)
	}
}
