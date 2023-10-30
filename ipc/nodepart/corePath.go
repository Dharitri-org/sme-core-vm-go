package nodepart

import (
	"os"
	"path"

	"github.com/Dharitri-org/sme-core-vm-go/ipc/common"
)

func (driver *CoreDriver) getCorePath() (string, error) {
	corePath, err := driver.getCorePathInCurrentDirectory()
	if err == nil {
		return corePath, nil
	}

	corePath, err = driver.getCorePathFromEnvironment()
	if err == nil {
		return corePath, nil
	}

	return "", common.ErrCoreNotFound
}

func (driver *CoreDriver) getCorePathInCurrentDirectory() (string, error) {
	cwd, err := os.Getwd()
	if err != nil {
		return "", err
	}

	corePath := path.Join(cwd, "core")
	if fileExists(corePath) {
		return corePath, nil
	}

	return "", common.ErrCoreNotFound
}

func (driver *CoreDriver) getCorePathFromEnvironment() (string, error) {
	corePath := os.Getenv(common.EnvVarCorePath)
	if fileExists(corePath) {
		return corePath, nil
	}

	return "", common.ErrCoreNotFound
}

func fileExists(filename string) bool {
	info, err := os.Stat(filename)
	if os.IsNotExist(err) {
		return false
	}

	return !info.IsDir()
}
