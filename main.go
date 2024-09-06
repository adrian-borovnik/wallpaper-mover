package main

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"regexp"
)

type MoveConfig struct {
	Source      string
	Destination string
	Pattern     string
}

func MoveFile(source string, destination string) error {
	sourceFile, err := os.Open(source)
	if err != nil {
		return fmt.Errorf("could not open source file: %v", err)
	}
	defer sourceFile.Close()

	fi, err := sourceFile.Stat()
	if err != nil {
		return err
	}

	flag := os.O_WRONLY | os.O_CREATE | os.O_TRUNC
	perm := fi.Mode() & os.ModePerm

	destinationFile, err := os.OpenFile(destination, flag, perm)
	if err != nil {
		return fmt.Errorf("could not create destination file: %v", err)
	}
	defer destinationFile.Close()

	_, err = io.Copy(destinationFile, sourceFile)
	if err != nil {
		destinationFile.Close()
		os.Remove(destination)
		return fmt.Errorf("could not copy file contents from source to destination: %v", err)
	}

	sourceFile.Close()

	err = os.Remove(source)
	if err != nil {
		return fmt.Errorf("could not remove source file: %v", err)
	}

	return nil
}

func MoveFiles(config MoveConfig) error {
	entries, err := os.ReadDir(config.Source)
	if err != nil {
		return err
	}

	for _, e := range entries {
		if e.IsDir() {
			continue
		}

		fileName := e.Name()
		r, _ := regexp.Compile(config.Pattern)

		if !r.MatchString(fileName) {
			continue
		}

		sourcePath := filepath.Join(config.Source, fileName)
		destinationPath := filepath.Join(config.Destination, fileName)
		err := MoveFile(sourcePath, destinationPath)
		if err != nil {
			return err
		}
	}

	return nil
}

func main() {
	config := MoveConfig{
		Source:      "/Users/adrianborovnik/Downloads",
		Destination: "/Users/adrianborovnik/Wallpapers",
		Pattern:     ".*unsplash.*\\.jpg",
	}

	if err := MoveFiles(config); err != nil {
		fmt.Println(err)
	}
}
