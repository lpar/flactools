package main

import (
	"flag"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"

	"github.com/lpar/flactools"
)

func examine(path string, info fs.DirEntry, err error) error {
	if err != nil {
		return err
	}
	if info.IsDir() {
		err = processDir(path)
		if err != nil {
			fmt.Printf("%v - skipped\n", err)
		}
	}
	return nil
}

func processDir(path string) error {
	files, err := os.ReadDir(path)
	if err != nil {
		return fmt.Errorf("can't read %s: %w", path, err)
	}
	artists := make(map[string]struct{})
	albums := make(map[string]struct{})
	var artist string
	var album string
	for _, file := range files {
		fname := file.Name()
		ext := filepath.Ext(fname)
		if strings.EqualFold(ext, ".flac") && !file.IsDir() {
			artist, album, err = scanFile(filepath.Join(path, fname))
			if err != nil {
				return fmt.Errorf("can't scan %s: %w", fname, err)
			}
			artists[artist] = struct{}{}
			albums[album] = struct{}{}
		}
	}
	if len(artists) == 0 && len(albums) == 0 {
		return nil
	}
	if len(artists) == 0 {
		return fmt.Errorf("%s contains albums but no artists", path)
	}
	if len(albums) == 0 {
		return fmt.Errorf("%s contains artists but no albums", path)
	}
	if len(albums) != 1 {
		return fmt.Errorf("%s contains %d albums", path, len(albums))
	}
	if len(artists) > 1 {
		artist = "Various Artists"
	}
	if err := renameDirectory(path, artist, album); err != nil {
		return fmt.Errorf("can't move %s: %w", path, err)
	}
	return nil
}

func scanFile(fname string) (string, string, error) {
	tags, err := flactools.GetTags(fname)
	if err != nil {
		return "", "", err
	}
	artist := flactools.Coalesce(tags["ALBUMARTIST"], tags["ARTIST"])
	album := tags["ALBUM"]
	if artist == "" || album == "" {
		return "", "", fmt.Errorf("missing artist or album for %s", fname)
	}
	return artist, album, nil
}

func renameDirectory(path string, artist string, album string) error {
	artdir := filepath.Join(*destDir, flactools.CleanName(artist))
	albdir := filepath.Join(artdir, flactools.CleanName(album))
	if path == albdir {
		return nil // already in the right place
	}
	if *shell {
		fmt.Printf("mkdir -p \"%s\"\nmv -n \"%s\" \"%s\"\n", artdir, path, albdir)
		return nil
	}
	fmt.Printf("%s -> %s\n", path, albdir)
	if !*dryRun {
		os.MkdirAll(artdir, 0755)
		err := os.Rename(path, albdir)
		return err
	}
	return nil
}

var destDir = flag.String("dest", "/volume1/FLAC", "destination directory")
var dryRun = flag.Bool("test", false, "test mode, display what would be done but don't actually move files")
var shell = flag.Bool("shell", false, "output shell script instead of moving files")

func main() {
	flag.Usage = func() {
		cmdname := filepath.Base(os.Args[0])
		fmt.Fprintf(os.Stderr, "%s - file FLAC files by artist and album, keeping them in the same directory\nusage: %s file1 ...\n", cmdname, cmdname)
		flag.PrintDefaults()
	}
	flag.Parse()
	for _, fname := range flag.Args() {
		if err := filepath.WalkDir(fname, examine); err != nil {
			fmt.Printf("error scanning %s: %v\n", fname, err)
		}
	}
}
