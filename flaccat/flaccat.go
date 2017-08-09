// flaccat
//
// Usage:
//   flaccat <path>
//
// Scans all directories under the specified path for FLAC files, and outputs
// a catalog file to standard output consisting of the MD5 checksum for each
// FLAC file, followed by the path to the file.
//
// By comparing successive versions of the output, you can determine which
// flac files have been added, removed, or moved. The flacdiff utility will
// perform that function.
//
package main

import (
	"encoding/hex"
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/mewkiz/flac"
)

func checksum(fspc string) (string, error) {
	stream, err := flac.ParseFile(fspc)
	if err != nil {
		return "", err
	}
	defer stream.Close()
	info := stream.Info
	md5 := hex.EncodeToString(info.MD5sum[:])
	return md5, nil
}

func examine(path string, info os.FileInfo, err error) error {
	if info.IsDir() {
		if strings.HasSuffix(path, "/#recycle") || strings.HasSuffix(path, "/@eaDir") {
			return filepath.SkipDir
		}
	} else {
		lcpath := strings.ToLower(path)
		if strings.HasSuffix(lcpath, ".flac") {
			md5, ferr := checksum(path)
			if ferr != nil {
				fmt.Fprintf(os.Stderr, "error reading %s: %s\n", path, ferr)
			} else {
				fmt.Printf("%s %s\n", md5, path)
			}
		}
	}
	return err
}

func main() {
	flag.Parse()
	for _, fname := range flag.Args() {
		filepath.Walk(fname, examine)
	}
}
