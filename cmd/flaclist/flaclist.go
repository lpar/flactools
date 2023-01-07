// flaccat
//
// Usage:
//
//	flaclist <path>
//
// Scans all directories under the specified path for FLAC files, and outputs
// a listing of artist, album
package main

import (
	"encoding/hex"
	"flag"
	"fmt"
	"io/fs"
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

func examine(path string, info fs.DirEntry, err error) error {
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
	flag.Usage = func() {
		cmdname := filepath.Base(os.Args[0])
		fmt.Fprintf(os.Stderr, "%s - scan for FLAC files and list artists and album titles\nusage: %s path ...\n", cmdname, cmdname)
		flag.PrintDefaults()
	}
	flag.Parse()
	for _, fname := range flag.Args() {
		if err := filepath.WalkDir(fname, examine); err != nil {
			fmt.Fprintf(os.Stderr, "%s: %v\n", fname, err)
		}
	}
}
