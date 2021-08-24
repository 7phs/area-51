package config

import (
	"errors"
	"os"
)

func isFileValid(filePath string) error {
	stat, err := os.Stat(filePath)
	switch {
	case errors.Is(err, os.ErrNotExist):
		// it is expected if a file is not exists. Will wait for it
		return nil

	case errors.Is(err, os.ErrPermission):
		return ErrPermissionDenied(filePath)

	case err != nil:
		return ErrUnexpected(filePath, err)

	case stat.IsDir():
		return ErrIsDir(filePath)
	}

	return nil
}

func isDirValid(filePath string) error {
	stat, err := os.Stat(filePath)
	switch {
	case errors.Is(err, os.ErrNotExist):
		return ErrNotExists(filePath)

	case errors.Is(err, os.ErrPermission):
		return ErrPermissionDenied(filePath)

	case err != nil:
		return ErrUnexpected(filePath, err)

	case !stat.IsDir():
		return ErrIsDir(filePath)
	}

	return nil
}
