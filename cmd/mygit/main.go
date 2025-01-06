package main

import (
	"bytes"
	"compress/zlib"
	"fmt"
	"io"
	"os"
)

func main() {
	fmt.Fprintf(os.Stderr, "Logs from your program will appear here!\n")

	if len(os.Args) < 2 {
		fmt.Fprintf(os.Stderr, "usage: mygit <command> [<args>...]\n")
		os.Exit(1)
	}
	switch command := os.Args[1]; command {
	case "init":
		for _, dir := range []string{".git", ".git/objects", ".git/refs"} {
			if err := os.MkdirAll(dir, 0755); err != nil {
				fmt.Fprintf(os.Stderr, "Error creating directory: %s\n", err)
			}
		}

		headFileContents := []byte("ref: refs/heads/main\n")
		if err := os.WriteFile(".git/HEAD", headFileContents, 0644); err != nil {
			fmt.Fprintf(os.Stderr, "Error writing file: %s\n", err)
		}

		fmt.Println("Initialized git directory")
	case "cat-file":
		objectHash := os.Args[3]
		dirName := ".git/objects/" + objectHash[0:2]
		entries, err := os.ReadDir(dirName)
		if err != nil {
			fmt.Println(err)
			return
		}

		for _, entry := range entries {
			_, err := entry.Info()
			if err != nil {
				fmt.Println(err)
				return
			}
			file, err := os.Open(dirName + "/" + entry.Name())
			if err != nil {
				fmt.Println("error opening file", err)
				return
			}
			defer file.Close()
			var b bytes.Buffer
			_, err = b.ReadFrom(file)
			if err != nil {
				fmt.Println(err)
			}
			r, err := zlib.NewReader(&b)
			if err != nil {
				fmt.Println("error decompressing ", err)
				return
			}
			content := new(bytes.Buffer)
			io.Copy(content, r)
			data := content.Bytes()
			nullIndex := bytes.IndexByte(data, 0)
			if nullIndex != -1 {
				fmt.Print(string(data[nullIndex+1:]))
			}

			defer r.Close()

		}

	default:
		fmt.Fprintf(os.Stderr, "Unknown command %s\n", command)
		os.Exit(1)
	}
}
