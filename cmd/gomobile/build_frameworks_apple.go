package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/google/uuid"
)

type Framework struct {
	Name              string
	FileRefID         string
	FrameworksID      string
	EmbedFrameworksID string
	Path              string
}

type BuildFrameworks struct {
	frameworks []Framework
}

func CollectBuildFrameworks(projectDir string) (*BuildFrameworks, error) {
	srcFrameworks := filepath.Join(projectDir, "Frameworks")
	fi, err := os.Stat(srcFrameworks)
	if err != nil {
		if os.IsNotExist(err) {
			return &BuildFrameworks{}, nil
		}

		return nil, err
	}

	if !fi.IsDir() {
		return nil, fmt.Errorf("stat error: Frameworks must be a directory")
	}

	frameworks := make([]Framework, 0)
	err = filepath.WalkDir(srcFrameworks, func(path string, d os.DirEntry, err error) error {
		if !d.IsDir() {
			// Skip non-dir entries as we are looking for .framework directories.
			return nil
		}

		if path == srcFrameworks {
			// Skip the root directory.
			return nil
		}

		relPath, err := filepath.Rel(srcFrameworks, path)
		if err != nil {
			// Should never happen.
			return err
		}

		pathComponents := filepath.SplitList(relPath)

		matched, err := filepath.Match("*.framework", pathComponents[0])
		if err != nil {
			return err
		}
		if !matched {
			// Skip non-framework directories.
			return filepath.SkipDir
		}

		// Now we have a .framework directory.
		frameworks = append(frameworks, Framework{
			Name:              pathComponents[0],
			FileRefID:         strings.ToUpper(strings.ReplaceAll(uuid.New().String(), "-", "")),
			FrameworksID:      strings.ToUpper(strings.ReplaceAll(uuid.New().String(), "-", "")),
			EmbedFrameworksID: strings.ToUpper(strings.ReplaceAll(uuid.New().String(), "-", "")),
			Path:              path,
		})

		// And we can skip any other files here.
		return filepath.SkipDir
	})
	if err != nil {
		return nil, err
	}

	return &BuildFrameworks{frameworks}, nil
}

func (bf *BuildFrameworks) Copy(xcodeProjDir string) error {
	dstFrameworks := xcodeProjDir + "/Frameworks"
	if err := mkdir(dstFrameworks); err != nil {
		return err
	}

	for _, framework := range bf.frameworks {
		err := processDirectory(framework.Path, func(srcPath, dstPath string) error {
			return copyFile(filepath.Join(dstFrameworks, framework.Name, dstPath), srcPath)
		})
		if err != nil {
			return err
		}
	}

	return nil
}

func (bf *BuildFrameworks) AsSlice() []Framework {
	return bf.frameworks
}
