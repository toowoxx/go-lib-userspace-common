package cmds

import (
	"fmt"
	"io"
	"os"
	"os/exec"
)

func RunCommand(path string, args ...string) error {
	cmd := exec.Command(path, args...)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

func RunCommandReturnOutput(path string, args ...string) ([]byte, error) {
	return exec.Command(path, args...).Output()
}

func RunCommandInDirWithEnvReturnOutput(name string, dir string, env map[string]string, params ...string) ([]byte, error) {
	cmd := exec.Command(name, params...)
	if dir != "." {
		cmd.Dir = dir
	}
	for key, value := range env {
		cmd.Env = append(cmd.Env, fmt.Sprintf("%s=%s", key, value))
	}
	return cmd.Output()
}

func RunCommandInDirReturnOutput(name string, dir string, params ...string) ([]byte, error) {
	cmd := exec.Command(name, params...)
	if dir != "." {
		cmd.Dir = dir
	}
	return cmd.Output()
}

func RunCommandInDirReturnCombinedOutput(name string, dir string, params ...string) ([]byte, error) {
	cmd := exec.Command(name, params...)
	if dir != "." {
		cmd.Dir = dir
	}
	return cmd.CombinedOutput()
}

func RunCommandWithEnvReturnOutput(path string, env map[string]string, args ...string) ([]byte, error) {
	cmd := exec.Command(path, args...)
	for key, value := range env {
		cmd.Env = append(cmd.Env, fmt.Sprintf("%s=%s", key, value))
	}
	return cmd.Output()
}

func RunCommandWithEnv(path string, env map[string]string, args ...string) error {
	cmd := exec.Command(path, args...)
	for key, value := range env {
		cmd.Env = append(cmd.Env, fmt.Sprintf("%s=%s", key, value))
	}
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

func RunCommandReturnExitStatus(path string, args ...string) (int, error) {
	cmd := exec.Command(path, args...)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		if exitError, ok := err.(*exec.ExitError); ok {
			return exitError.ExitCode(), nil
		} else {
			return -1, err
		}
	}
	return 0, nil
}

func StartCommandWithPipes(cmd *exec.Cmd, input io.Reader) (io.ReadCloser, chan error, error) {
	pipeR, pipeW := io.Pipe()
	cmd.Stdin = input
	cmd.Stdout = pipeW
	cmd.Stderr = os.Stderr

	err := cmd.Start()
	if err != nil {
		return nil, nil, err
	}

	errChan := make(chan error)

	go func() {
		err := cmd.Wait()
		_ = pipeW.Close()
		if err != nil {
			errChan <- err
		}
		close(errChan)
	}()

	return pipeR, errChan, nil
}
