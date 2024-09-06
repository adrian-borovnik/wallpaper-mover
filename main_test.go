package main

import (
	"os"
	"path/filepath"
	"testing"
)

func TestMoveFile(t *testing.T) {
	const source = "./foo.txt"
	const destination = "./bar.txt"
	content := []byte("test content")

	// Create the source file and write some content to it
	file, err := os.Create(source)
	if err != nil {
		t.Fatalf("Failed to create source file: %v", err)
	}
	defer os.Remove(source) // Cleanup if test fails early
	_, err = file.Write(content)
	if err != nil {
		t.Fatalf("Failed to write to source file: %v", err)
	}
	file.Close()

	// Perform the file move operation
	err = MoveFile(source, destination)
	if err != nil {
		t.Fatalf("Error while moving the file: %v", err)
	}
	defer os.Remove(destination) // Cleanup destination file at the end

	// Check if the destination file exists and contains the correct content
	destContent, err := os.ReadFile(destination)
	if err != nil {
		t.Fatalf("Destination file has not been created: %v", err)
	}
	if string(destContent) != string(content) {
		t.Errorf("Destination file content mismatch: got %q, want %q", string(destContent), string(content))
	}

	// Check that the source file no longer exists
	if _, err := os.Stat(source); err == nil {
		t.Errorf("Source file has not been deleted: %s", source)
	}
}

func TestMoveFiles(t *testing.T) {
	// Create temporary source and destination directories
	sourceDir, err := os.MkdirTemp("", "source")
	if err != nil {
		t.Fatalf("Failed to create source temp directory: %v", err)
	}
	defer os.RemoveAll(sourceDir)

	destinationDir, err := os.MkdirTemp("", "destination")
	if err != nil {
		t.Fatalf("Failed to create destination temp directory: %v", err)
	}
	defer os.RemoveAll(destinationDir)

	// Create test files
	testFiles := []struct {
		name    string
		content string
	}{
		{"file1_unsplash.jpg", "content1"},
		{"file2_unsplash.jpg", "content2"},
		{"file3.txt", "content3"},
	}

	for _, f := range testFiles {
		filePath := filepath.Join(sourceDir, f.name)
		err = os.WriteFile(filePath, []byte(f.content), 0644)
		if err != nil {
			t.Fatalf("Failed to create test file: %v", err)
		}
	}

	// Create MoveConfig
	config := MoveConfig{
		Source:      sourceDir,
		Destination: destinationDir,
		Pattern:     ".*unsplash.*\\.jpg",
	}

	// Call MoveFiles
	err = MoveFiles(config)
	if err != nil {
		t.Errorf("MoveFiles failed: %v", err)
	}

	// Check that the correct files were moved
	movedFiles := []string{"file1_unsplash.jpg", "file2_unsplash.jpg"}
	for _, fileName := range movedFiles {
		sourceFilePath := filepath.Join(sourceDir, fileName)
		if _, err := os.Stat(sourceFilePath); !os.IsNotExist(err) {
			t.Errorf("Expected file %s to be moved, but it still exists in source", fileName)
		}

		destinationFilePath := filepath.Join(destinationDir, fileName)
		if _, err := os.Stat(destinationFilePath); os.IsNotExist(err) {
			t.Errorf("Expected file %s to be moved, but it does not exist in destination", fileName)
		}
	}

	// Check that non-matching files were not moved
	nonMovedFiles := []string{"file3.txt"}
	for _, fileName := range nonMovedFiles {
		sourceFilePath := filepath.Join(sourceDir, fileName)
		if _, err := os.Stat(sourceFilePath); os.IsNotExist(err) {
			t.Errorf("Expected file %s to remain in source, but it was moved", fileName)
		}
	}
}
