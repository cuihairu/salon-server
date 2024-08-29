package utils

import (
	"os"
	"path/filepath"
)

func CreateDirIfNotExist(paths []string) error {
	for _, path := range paths {
		dir := filepath.Dir(path)
		if _, err := os.Stat(dir); os.IsNotExist(err) {
			if err = os.MkdirAll(dir, os.ModePerm); err != nil {
				return err
			}
		}
	}
	return nil
}
