// flacsum
//
// Usage:
//
//	flacsum <path>
//
// For each FLAC file somewhere under the specified path, checks the MD5
// checksum stored in the file against the actual audio data in the file,
// to detect any file corruption.
//
// Output consists of a series of lines consisting of either PASS or FAIL
// followed by the filename in question. FAIL lines are sent to stderr, PASS
// lines to stdout.
//
// Note that the process of computing the MD5 checksums is slow.
package main

import (
	"bufio"
	"bytes"
	"crypto/md5"
	"encoding/hex"
	"flag"
	"fmt"
	"io"
	"io/fs"
	"os"
	"path/filepath"
	"strings"

	"github.com/mewkiz/flac"
)

// Based on https://godoc.org/github.com/mewkiz/flac
func checksum(fspc string) error {
	fp, err := os.Open(fspc)
	if err != nil {
		return err
	}
	defer fp.Close()
	bfp := bufio.NewReader(fp)
	stream, err := flac.New(bfp)
	if err != nil {
		return err
	}
	defer stream.Close()
	md5sum := md5.New()
	for {
		frame, err := stream.ParseNext()
		if err != nil {
			if err == io.EOF {
				break
			}
			return fmt.Errorf("file parse error: %s", err)
		}
		frame.Hash(md5sum)
	}
	got, want := md5sum.Sum(nil), stream.Info.MD5sum[:]
	if !bytes.Equal(got, want) {
		return fmt.Errorf("bad checksum, wanted %s got %s", hex.EncodeToString(want), hex.EncodeToString(got))
	}
	return nil
}

func examine(path string, info fs.DirEntry, err error) error {
	if info.IsDir() {
		if strings.HasSuffix(path, "/#recycle") || strings.HasSuffix(path, "/@eaDir") {
			return filepath.SkipDir
		}
	} else {
		lcpath := strings.ToLower(path)
		if strings.HasSuffix(lcpath, ".flac") {
			ferr := checksum(path)
			if ferr != nil {
				fmt.Fprintf(os.Stderr, "FAIL %s: %s\n", path, ferr)
			} else {
				fmt.Printf("PASS %s\n", path)
			}
		}
	}
	return err
}

func main() {
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "%s - check FLAC file integrity\nusage: flacsum path ...\n", filepath.Base(os.Args[0]))
	}
	flag.Parse()
	for _, fname := range flag.Args() {
		if err := filepath.WalkDir(fname, examine); err != nil {
			fmt.Fprintf(os.Stderr, "%s: %v\n", fname, err)
		}
	}
}
