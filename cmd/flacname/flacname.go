package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"strconv"

	"github.com/lpar/flactools"
)

func examine(path string, info os.FileInfo, err error) error {
	if err != nil {
		return err
	}
	if info.IsDir() {
		return nil
	}
	return processFile(path)
}

func toInt(x string) int {
	n, err := strconv.Atoi(x)
	if err != nil {
		return 1
	}
	return n
}

func processFile(fpath string) error {
	dir := filepath.Dir(fpath)
	ext := filepath.Ext(fpath)
	if ext != ".flac" {
		return nil
	}
	tags, err := flactools.GetTags(fpath)
	if err != nil {
		return err
	}
	title := tags["TITLE"]
	tracknum := toInt(tags["TRACKNUMBER"])
	discnum := toInt(tags["DISCNUMBER"])
	discs := toInt(tags["DISCTOTAL"])
	var nfn string
	if discs > 1 {
		nfn = fmt.Sprintf("%02d-%02d-%s%s", discnum, tracknum, flactools.CleanName(title), ext)
	} else {
		nfn = fmt.Sprintf("%02d-%s%s", tracknum, flactools.CleanName(title), ext)
	}
	nfpath := filepath.Join(dir, nfn)
	if *shell {
		fmt.Printf("mv -d \"%s\" \"%s\"\n", fpath, nfpath)
		return nil
	}
	fmt.Printf("%s -> %s\n", fpath, nfpath)
	_, err = os.Stat(nfpath)
	if !os.IsNotExist(err) {
		return fmt.Errorf("can't rename %s, %s already exists", fpath, nfpath)
	}
	if !*dryRun {
		return os.Rename(fpath, nfpath)
	}
	return nil
}

var dryRun = flag.Bool("test", false, "test mode, don't actually rename files")
var shell = flag.Bool("shell", false, "output shell script instead of moving files")

func main() {
	flag.Usage = func() {
		cmdname := filepath.Base(os.Args[0])
		fmt.Fprintf(os.Stderr, "%s - rename FLAC files consistently but don't move them\nusage: %s path ...\n", cmdname, cmdname)
		flag.PrintDefaults()
	}
	flag.Parse()
	for _, fname := range flag.Args() {
		if err := filepath.Walk(fname, examine); err != nil {
			fmt.Printf("error scanning %s: %v\n", fname, err)
		}
	}
}
