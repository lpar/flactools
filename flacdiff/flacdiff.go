// flacdiff
//
// Usage:
//   flacdiff <old report> <new report>
//
// Compares two FLAC file catalogs generated by the flaccat utility. Produces a
// plain text report listing changes to FLAC files between the time the two
// catalog files were generated. The report lists:
//
//  1. FLAC files which were deleted.
//  2. FLAC files which were added.
//  3. FLAC files which were moved.
//
// For item 3, the files are listed in the format:
//
// oldfile
//  ↳ newfile
//
// File paths listed are shortened by computing them relative to a top-level
// directory common to all FLAC files.

package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"path"
	"strings"
)

type Move struct {
	From string
	To   string
}

type Report struct {
	FileMap map[string]string
	Added   []string
	Moved   []Move
	Prefix  string
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func common(s1 string, s2 string) string {
	if s1 == "\x09" {
		return s2
	}
	b1 := []byte(s1)
	b2 := []byte(s2)
	n := min(len(b1), len(b2))
	for i := 0; i < n; i++ {
		if b1[i] != b2[i] {
			return string(s1[:i])
		}
	}
	return s1
}

func NewReport() *Report {
	return &Report{
		FileMap: make(map[string]string),
		Prefix:  "\x09",
	}
}

func (r *Report) addOldFile(fspc string, md5 string) {
	r.Prefix = common(r.Prefix, fspc)
	r.FileMap[md5] = fspc
}

func (r *Report) addNewFile(fspc string, md5 string) {
	r.Prefix = common(r.Prefix, fspc)
	old, ok := r.FileMap[md5]
	if !ok {
		r.Added = append(r.Added, fspc)
		return
	}
	delete(r.FileMap, md5)
	if old != fspc {
		r.Moved = append(r.Moved, Move{From: old, To: fspc})
	}
}

func (r *Report) Deleted() []string {
	dels := make([]string, len(r.FileMap))
	for _, v := range r.FileMap {
		dels = append(dels, v)
	}
	return dels
}

func (r *Report) process(filename string, procfunc func(string, string)) error {
	list, err := os.Open(filename)
	if err != nil {
		return fmt.Errorf("unable to open input file %s: %s", filename, err)
	}
	defer list.Close()
	scanner := bufio.NewScanner(list)
	scanner.Split(bufio.ScanLines)
	for scanner.Scan() {
		md5, file, err := split(scanner.Text())
		if err != nil {
			return fmt.Errorf("bad input in %s: %s", filename, err)
		}
		procfunc(file, md5)
	}
	return nil
}

func split(line string) (string, string, error) {
	s := strings.IndexRune(line, ' ')
	if s == -1 {
		return "", "", fmt.Errorf("bad input line '%s'", line)
	}
	if s < 32 {
		return "", "", fmt.Errorf("bad md5 value in line '%s'", line)
	}
	return line[:s-1], line[s+1:], nil
}

func main() {
	flag.Parse()
	if len(flag.Args()) != 2 {
		fmt.Fprintln(os.Stderr, "Usage:\n  flacdiff previous-list current-list")
		return
	}
	a := flag.Args()
	r := NewReport()
	r.process(a[0], r.addOldFile)
	r.process(a[1], r.addNewFile)

	base := path.Dir(r.Prefix)

	fmt.Println("BASE DIRECTORY\n")
	fmt.Println(base)

	base = base + "/"

	fmt.Println("\nDELETED FILES:\n")
	for _, f := range r.Deleted() {
		if f != "" {
			fmt.Println(strings.TrimPrefix(f, base))
		}
	}
	fmt.Println("\nADDED FILES:\n")
	for _, f := range r.Added {
		fmt.Println(strings.TrimPrefix(f, base))
	}
	fmt.Println("\nMOVED FILES:\n")
	for _, m := range r.Moved {
		fmt.Println(strings.TrimPrefix(m.From, base))
		fmt.Println(" ↳ " + strings.TrimPrefix(m.To, base))
	}

}
