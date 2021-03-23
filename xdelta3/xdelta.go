package xdelta3

import (
	"fmt"
	"io"
	"os/exec"

	"github.com/pkg/errors"
	"github.com/toowoxx/go-lib-fs/fs"

	"github.com/toowoxx/go-lib-userspace-common/cmds"
)

func Delta(sourceFile, targetFile, deltaFile string) (err error) {
	if !fs.FileExists(sourceFile) {
		return errors.New(fmt.Sprintf("source file %s does not exist", sourceFile))
	}
	if !fs.FileExists(targetFile) {
		return errors.New(fmt.Sprintf("target file %s does not exist", sourceFile))
	}
	err = cmds.RunCommand(
		"xdelta3", "-s", sourceFile, targetFile, deltaFile,
	)
	if err != nil {
		return errors.Wrap(err, "could not run xdelta delta command")
	}
	if !fs.FileExists(deltaFile) {
		return errors.New("command ran successfully but did not produce delta file")
	}
	return
}

func DeltaStream(sourceFile string, r io.Reader) (io.ReadCloser, chan error, error) {
	if !fs.FileExists(sourceFile) {
		return nil, nil, errors.New(fmt.Sprintf("source file %s does not exist", sourceFile))
	}

	cmd := exec.Command("xdelta3", "-s", sourceFile)

	return cmds.StartCommandWithPipes(cmd, r)
}

func Patch(deltaFile, sourceFile, targetFile string) error {
	if !fs.FileExists(sourceFile) {
		return errors.New(fmt.Sprintf("source file %s does not exist", sourceFile))
	}
	if !fs.FileExists(deltaFile) {
		return errors.New(fmt.Sprintf("delta file %s does not exist", sourceFile))
	}
	err := cmds.RunCommand(
		"xdelta3", "-d", "-s", sourceFile, deltaFile, targetFile,
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
	output, err := exec.Command("xdelta3", "printdelta", deltaFile).CombinedOutput()
	if err != nil {
		return "", errors.Wrap(err, "could not run xdelta info command")
	}
	return string(output), nil
}
