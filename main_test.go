package main

import (
	"os"
	"testing"
)

func TestMoveFile(t *testing.T) {
	const source = "./foo.txt"
	const destination = "./bar.txt"

	file, err := os.Create(source)
	if err != nil {
		t.Error(err)
	}
	file.Close()

	err = MoveFile(source, destination)
	if err != nil {
		t.Errorf("Error while moving a file: %v", err)
	}

	if _, err := os.Stat(destination); err != nil {
		t.Errorf("Destination file has not been created: %s", destination)
	}

	if _, err := os.Stat(source); err == nil {
		t.Errorf("Source file has not been deleted: %s", source)
	}

	err = os.Remove(destination)
	if err != nil {
		t.Error(err)
	}
}
