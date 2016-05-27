package main

import (
	"bufio"
	"crypto/md5"
	"crypto/sha1"
	"flag"
	"fmt"
	"hash"
	"io"
	"os"

	"github.com/northbright/pathhelper"
)

const (
	// Usage of checksum-calc
	Usage string = "checksum-calc -f='my-cd.iso'"
	// Buffer size of reading file
	BufSize = 1 * 1024 * 1024
)

// ComputeChecksums() reads bytes and compute checksums.
//
//   Params:
//     r: io.Reader to read bytes from.
//   Returns:
//     checksums: map contains calculated checksums. Key: hash type(e.g. "MD5", "SHA-1"), Value: Checksum string.
func ComputeChecksums(r io.Reader) (checksums map[string]string, err error) {
	checksums = make(map[string]string)

	reader := bufio.NewReader(r)
	if err != nil {
		return checksums, err
	}

	hashes := make(map[string]hash.Hash)

	hashes["MD5"] = md5.New()
	hashes["SHA-1"] = sha1.New()

	buf := make([]byte, BufSize)
	for {
		n, err := reader.Read(buf)
		if err != nil && err != io.EOF {
			return checksums, err
		} else {
			if n == 0 {
				break
			}
		}

		// Adds more data to the running hash.
		for _, w := range hashes {
			if n, err = w.Write(buf[:n]); err != nil {
				return checksums, err
			}
		}
	}

	for k, v := range hashes {
		checksums[k] = fmt.Sprintf("%X", v.Sum(nil))
	}

	return checksums, nil
}

func main() {
	inputFile := ""

	flag.StringVar(&inputFile, "f", "", "File to calculate MD5 / SHA-1 hash. Ex: -f='my-cd.iso'")
	flag.Parse()

	if inputFile == "" {
		flag.PrintDefaults()
		fmt.Printf("%s\n", Usage)
		return
	}

	absFilePath, err := pathhelper.GetAbsPath(inputFile)
	if err != nil {
		fmt.Printf("pathhelper.GetAbsPath(%v) error: %v\n", inputFile, err)
		return
	}

	f, err := os.Open(absFilePath)
	if err != nil {
		fmt.Printf("os.Open(%v) error: %v\n", absFilePath, err)
		return
	}
	defer f.Close()

	fmt.Printf("Starting computing checksums...\n")

	checksums, err := ComputeChecksums(f)
	if err != nil {
		fmt.Printf("ComputeChecksums() error: %v\n", err)
		return
	}

	fmt.Printf("Checksums has been computed successfully.\n")
	s := "--------------------------------------\n"
	for k, v := range checksums {
		s += fmt.Sprintf("%s: %s\n", k, v)
	}
	fmt.Printf("%s\n", s)
}
