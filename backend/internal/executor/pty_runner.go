package executor

import (
	"bufio"
	"context"
	"os/exec"

	"github.com/creack/pty"
)

type PtyRunner struct{}

// NewPtyRunner 创建并返回对应的组件。
func NewPtyRunner() *PtyRunner {
	return &PtyRunner{}
}

// Run 执行当前组件的主循环或工作流。
func (r *PtyRunner) Run(ctx context.Context, cmd *exec.Cmd, onOutput func([]byte)) error {
	ptmx, err := pty.Start(cmd)
	if err != nil {
		return err
	}
	defer func() { _ = ptmx.Close() }()

	done := make(chan error, 1)
	go func() {
		scanner := bufio.NewScanner(ptmx)
		for scanner.Scan() {
			line := append([]byte(nil), scanner.Bytes()...)
			onOutput(line)
		}
		if err := scanner.Err(); err != nil {
			done <- err
			return
		}
		done <- cmd.Wait()
	}()

	select {
	case <-ctx.Done():
		_ = cmd.Process.Kill()
		return ctx.Err()
	case err := <-done:
		return err
	}
}
