package xdelta

import (
	"fmt"
	"os/exec"

	"github.com/pkg/errors"
	"github.com/toowoxx/go-lib-fs/fs"

	"github.com/toowoxx/go-lib-userspace-common/cmds"
)

func Delta(sourceFile, targetFile, deltaFile string) (foundDiff bool, err error) {
	if !fs.FileExists(sourceFile) {
		return false, errors.New(fmt.Sprintf("source file %s does not exist", sourceFile))
	}
	if !fs.FileExists(targetFile) {
		return false, errors.New(fmt.Sprintf("target file %s does not exist", sourceFile))
	}
	exitCode, err := cmds.RunCommandReturnExitStatus(
		"xdelta", "delta", sourceFile, targetFile, deltaFile,
	)
	if exitCode == 1 {
		foundDiff = true
		err = nil
	} else if err != nil {
		return false, errors.Wrap(err, "could not run xdelta delta command")
	}
	if !fs.FileExists(deltaFile) {
		return foundDiff, errors.New("command ran successfully but did not produce delta file")
	}
	return
}

func Patch(deltaFile, sourceFile, targetFile string) error {
	if !fs.FileExists(sourceFile) {
		return errors.New(fmt.Sprintf("source file %s does not exist", sourceFile))
	}
	if !fs.FileExists(deltaFile) {
		return errors.New(fmt.Sprintf("delta file %s does not exist", sourceFile))
	}
	err := cmds.RunCommand(
		"xdelta", "patch", sourceFile, deltaFile, targetFile,
	)
	if err != nil {
		return errors.Wrap(err, "could not run xdelta patch command")
	}
	if !fs.FileExists(deltaFile) {
		return errors.New("command ran successfully but did not produce delta file")
	}
	return nil
}

func Info(deltaFile string) (string, error) {
	if !fs.FileExists(deltaFile) {
		return "", errors.New(fmt.Sprintf("delta file %s does not exist", deltaFile))
	}
	output, err := exec.Command("xdelta", "info", deltaFile).CombinedOutput()
	if err != nil {
		return "", errors.Wrap(err, "could not run xdelta info command")
	}
	return string(output), nil
}
