package main

import (
	"os"
	"path/filepath"
	"testing"
)

func TestBuildFrameworks_Copy(t *testing.T) {
	// Create a temporary directory for the Xcode project
	xcodeProjDir, err := os.MkdirTemp("", "xcodeProj")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(xcodeProjDir)

	// Create a temporary directory for the frameworks
	srcFrameworks, err := os.MkdirTemp("", "Frameworks")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(srcFrameworks)

	// Create some mock frameworks
	frameworks := []string{"FrameworkA.framework", "FrameworkB.framework"}
	for _, name := range frameworks {
		frameworkDir := filepath.Join(srcFrameworks, name)
		if err := os.Mkdir(frameworkDir, 0755); err != nil {
			t.Fatalf("Failed to create temp dir: %v", err)
		}

		// Add a dummy file so it will be copied
		dummyFile := filepath.Join(frameworkDir, "dummy.txt")
		if _, err := os.Create(dummyFile); err != nil {
			t.Fatalf("Failed to create dummy file: %v", err)
		}
	}

	// Create a file in the Frameworks dir that should be skipped
	skipFile := filepath.Join(srcFrameworks, "skip.txt")
	if _, err := os.Create(skipFile); err != nil {
		t.Fatalf("Failed to create skip file: %v", err)
	}

	// And now a directory that should be skipped
	skipDir := filepath.Join(srcFrameworks, "skipDir")
	if err := os.Mkdir(skipDir, 0755); err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}

	// Collect frameworks
	bf, err := CollectBuildFrameworks(filepath.Dir(srcFrameworks))
	if err != nil {
		t.Fatalf("CollectBuildFrameworks failed: %v", err)
	}

	// Call the Copy method
	err = bf.Copy(xcodeProjDir)
	if err != nil {
		t.Fatalf("Copy method failed: %v", err)
	}

	// Verify that the frameworks are copied to the correct location
	for _, framework := range bf.frameworks {
		expectedPath := filepath.Join(xcodeProjDir, "Frameworks", framework.Name)
		if _, err := os.Stat(expectedPath); os.IsNotExist(err) {
			t.Errorf("Expected framework %s to be copied to %s, but it was not", framework.Name, expectedPath)
		}
	}

	// Verify that the file and directory that needed to be skipped were skipped.
	if _, err := os.Stat(filepath.Join(xcodeProjDir, "Frameworks", "skip.txt")); !os.IsNotExist(err) {
		t.Errorf("Expected skip.txt to be skipped, but it was not")
	}
	if _, err := os.Stat(filepath.Join(xcodeProjDir, "Frameworks", "skipDir")); !os.IsNotExist(err) {
		t.Errorf("Expected skipDir to be skipped, but it was not")
	}
}
