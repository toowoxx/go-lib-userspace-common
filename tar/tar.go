package tar

import (
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/alessio/shellescape"
	"github.com/toowoxx/go-lib-fs/fs"

	"github.com/toowoxx/go-lib-userspace-common/cmds"
)

type Format string

const (
	FormatGNU    Format = "gnu"
	FormatOldGNU Format = "oldgnu"
	FormatPOSIX  Format = "posix"
)

const defaultFormat = FormatPOSIX

type CreateOptions struct {
	Archive            bool
	Seekable           bool
	Directory          string
	ExcludePatterns    []string
	ExcludeVCS         bool
	ExcludeVCSIgnores  bool
	NumericOwner       bool
	ACLs               bool
	SELinux            bool
	ExtendedAttributes bool
	Format             Format
	FullTime           bool

	Verbose bool
}

func (f Format) IsValid() bool {
	switch f {
	case FormatGNU, FormatOldGNU, FormatPOSIX:
		return true
	default:
		return false
	}
}

func (f Format) String() string {
	return string(f)
}

func Create(options CreateOptions) (io.Reader, chan error, error) {
	command := "tar"
	var arguments []string

	for _, pattern := range options.ExcludePatterns {
		arguments = append(arguments, "--exclude", pattern)
	}

	arguments = append(arguments, "-c")
	if options.Verbose {
		arguments = append(arguments, "-v")
	}
	if options.Archive {
		arguments = append(arguments, "-a")
	}

	if len(options.Format) == 0 {
		options.Format = defaultFormat
	}
	if !options.Format.IsValid() {
		return nil, nil, fmt.Errorf("format \"%s\" is invalid", options.Format)
	}
	arguments = append(arguments, fmt.Sprintf("--format=%s", options.Format))

	if options.ExcludeVCS {
		arguments = append(arguments, "--exclude-vcs")
	}
	if options.ExcludeVCSIgnores {
		arguments = append(arguments, "--exclude-vcs-ignores")
	}
	if options.Seekable {
		arguments = append(arguments, "--seek")
	}
	if options.NumericOwner {
		arguments = append(arguments, "--numeric-owner")
	}
	if options.ACLs {
		arguments = append(arguments, "--acls")
	} else {
		arguments = append(arguments, "--no-acls")
	}
	if options.SELinux {
		arguments = append(arguments, "--selinux")
	} else {
		arguments = append(arguments, "--no-selinux")
	}
	if options.ExtendedAttributes {
		arguments = append(arguments, "--xattrs")
	} else {
		arguments = append(arguments, "--no-xattrs")
	}

	if options.FullTime {
		arguments = append(arguments, "--full-time")
	}

	if len(options.Directory) > 0 {
		if !fs.DirectoryExists(options.Directory) {
			return nil, nil, fmt.Errorf("directory \"%s\" to change to does not exist", options.Directory)
		}
		arguments = append(arguments, "-C", options.Directory)
	}

	arguments = append(arguments, ".")

	cmd := exec.Command(command, arguments...)
	cmd.Dir, _ = filepath.Abs(".")

	if options.Verbose {
		_, _ = fmt.Fprintln(os.Stderr, "Running command", command, shellescape.QuoteCommand(arguments))
		_, _ = fmt.Fprintln(os.Stderr, "  in directory", cmd.Dir, "- command path:", cmd.Path)
	}

	return cmds.StartCommandWithPipes(cmd, nil)
}
