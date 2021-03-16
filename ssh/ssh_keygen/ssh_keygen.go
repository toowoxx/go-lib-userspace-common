package ssh_keygen

import (
	"os"
	"path"

	"github.com/pkg/errors"
	"github.com/toowoxx/go-lib-fs/fs"

	"github.com/toowoxx/go-lib-userspace-common/cmds"
)

func GenerateKeyPair(sshKeyPath string) error {
	var err error
	dir := path.Dir(sshKeyPath)
	err = os.MkdirAll(dir, 0700)
	if err != nil {
		return errors.Wrapf(err, "could not create directory \"%s\"", dir))
	}
	err = cmds.RunCommand("ssh-keygen", "-t", "ed25519", "-f", sshKeyPath, "-N", "")
	if err != nil {
		return errors.Wrap(err, "could not generate ssh keypair")
	}
	if !fs.FileExists(sshKeyPath) {
		return errors.New("ssh-keygen command was successful but the key does not exist")
	}
	return nil
}
