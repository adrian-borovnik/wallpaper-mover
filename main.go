package main

import (
	"fmt"
	"io"
	"log"
	"os"
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
	return nil
}

func main() {
	fmt.Println("Zdravo!")

	config := MoveConfig{
		Source:      "~/Downloads/",
		Destination: "~/Wallpapers",
		Pattern:     "*unsplash*.jpg",
	}

	if err := MoveFiles(config); err != nil {
		log.Fatalln(err)
	}

}
