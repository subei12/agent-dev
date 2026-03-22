package executor

import (
	"bufio"
	"context"
	"io"
	"os/exec"
	"sync"
)

type StreamRunner struct{}

// NewStreamRunner 创建并返回对应的组件。
func NewStreamRunner() *StreamRunner {
	return &StreamRunner{}
}

// Run 执行命令并以流式方式转发 stdout/stderr 内容。
func (r *StreamRunner) Run(ctx context.Context, cmd *exec.Cmd, onOutput func([]byte)) error {
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return err
	}
	stderr, err := cmd.StderrPipe()
	if err != nil {
		return err
	}

	if err := cmd.Start(); err != nil {
		return err
	}

	var wg sync.WaitGroup
	readStream := func(reader io.Reader) {
		defer wg.Done()

		scanner := bufio.NewScanner(reader)
		for scanner.Scan() {
			line := append([]byte(nil), scanner.Bytes()...)
			onOutput(line)
		}
	}

	wg.Add(2)
	go readStream(stdout)
	go readStream(stderr)

	done := make(chan error, 1)
	go func() {
		wg.Wait()
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
