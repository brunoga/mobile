package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

func processDirectory(srcDir string, processFunc func(srcFile, dstFile string) error) error {
	// Resolve symlinks.
	srcDir, err := filepath.EvalSymlinks(srcDir)
	if err != nil {
		if os.IsNotExist(err) {
			// Do nothing if the directory does not exist.
			return nil
		} else {
			return err
		}
	}

	fi, err := os.Stat(srcDir)
	if err != nil {
		return err
	}

	if !fi.IsDir() {
		return fmt.Errorf("%s is not a directory", srcDir)
	}

	return filepath.WalkDir(srcDir, func(srcPath string, d os.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if name := filepath.Base(srcPath); strings.HasPrefix(name, ".") {
			// Do not include hidden files.
			return nil
		}

		if d.IsDir() {
			return nil
		}

		return processFunc(srcPath, srcPath[len(srcDir)+1:])
	})
}
