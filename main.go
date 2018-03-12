package main

import (
	"bufio"
	"crypto/md5"
	"crypto/sha1"
	"crypto/sha256"
	"encoding/json"
	"flag"
	"fmt"
	"hash"
	"io"
	"os"

	"github.com/northbright/pathhelper"
)

const (
	// incorrectArgsMsg is the message while arguments error occurs.
	incorrectArgsMsg string = "Incorrect arguments, please see usage:\n"
	// usage is the message of checksum-calc usage.
	usage string = "usage:\nchecksum-calc -f=<file>\nEx: checksum-calc -f='my-cd.iso'"
	// bufSize is buffer size of reading file.
	bufSize = 1 * 1024 * 1024
)

// ComputeChecksums reads bytes and compute checksums.
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
	hashes["SHA-256"] = sha256.New()

	buf := make([]byte, bufSize)
	for {
		n, err := reader.Read(buf)
		if err != nil && err != io.EOF {
			return checksums, err
		}

		if n == 0 {
			break
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

	var jsonOut bool

	flag.StringVar(&inputFile, "f", "", "File to calculate MD5 / SHA-1 hash. Ex: -f='my-cd.iso'")
	flag.BoolVar(&jsonOut, "jsonout", false, "dump json format output")
	flag.Parse()

	if inputFile == "" {
		fmt.Printf("%s\n", incorrectArgsMsg)
		flag.PrintDefaults()
		fmt.Printf("%s\n", usage)
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

	if !jsonOut {
		fmt.Printf("Computing checksums...\n")
	}

	checksums, err := ComputeChecksums(f)
	if err != nil {
		fmt.Printf("ComputeChecksums() error: %v\n", err)
		return
	}

	if !jsonOut {
		fmt.Printf("Done.\n")
		s := "--------------------------------------\n"
		for k, v := range checksums {
			s += fmt.Sprintf("%s: %s\n", k, v)
		}
		fmt.Printf("%s\n", s)
	} else {
		hashes := struct {
			Md5    string
			Sha256 string
			Sha1   string
		}{
			string(checksums["MD5"]),
			string(checksums["SHA-256"]),
			string(checksums["SHA-1"]),
		}

		jsonStr, _ := json.Marshal(hashes)

		fmt.Print(string(jsonStr))

	}
}
