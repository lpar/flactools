package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"unicode"

	"github.com/mewkiz/flac"
	"github.com/mewkiz/flac/meta"
)

func examine(path string, info os.FileInfo, err error) error {
	if info.IsDir() {
		err = processDir(path, info)
		if err != nil {
			fmt.Printf("%v - skipped\n", err)
		}
	}
	return nil
}

func processDir(path string, info os.FileInfo) error {
	files, err := ioutil.ReadDir(path)
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
		return fmt.Errorf("can't move %s: %v", path, err)
	}
	return nil
}

func scanFile(fname string) (string, string, error) {
	stream, err := flac.ParseFile(fname)
	if err != nil {
		return "","", fmt.Errorf("can't parse %s: %w", fname, err)
	}
	defer stream.Close()
	mdblocks := stream.Blocks
	for _, blk := range mdblocks {
		if blk.Type == meta.TypeVorbisComment {
			vc := blk.Body.(*meta.VorbisComment)
			artist, album := pickTags(vc.Tags)
			if artist == "" || album == "" {
				return "", "", fmt.Errorf("missing artist or album for %s", fname)
			}
			return artist, album, nil
		}
	}
	return "","", fmt.Errorf("no metadata found in %s", fname)
}

func coalesce(args... string) string {
	for _, arg := range args {
		if arg != "" {
			return arg
		}
	}
	return ""
}

func pickTags(tuplelist [][2]string) (string, string) {
	tags := make(map[string]string)
	for _, tuple := range tuplelist {
		k := tuple[0]
		v := tuple[1]
		tags[strings.ToUpper(k)] = v
	}
	artist := coalesce(tags["ALBUMARTIST"], tags["ARTIST"])
	album := tags["ALBUM"]
	return artist, album
}

func cleanName(x string) string {
	var dn strings.Builder
	spc := false
	for _, r := range x {
		if unicode.IsDigit(r) || unicode.IsLetter(r) {
			dn.WriteRune(r)
			spc = false
		}
		if r == ' ' && !spc {
			dn.WriteRune('_')
			spc = true
		}
		if r == '&' {
			dn.WriteString("and")
		}
	}
	return dn.String()
}

func renameDirectory(path string, artist string, album string) error {
	artdir := filepath.Join(*destDir, cleanName(artist))
	albdir := filepath.Join(artdir, cleanName(album))
	if *shell {
		fmt.Printf("mkdir -p \"%s\"\nmv \"%s\" \"%s\"\n", artdir, path, albdir)
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

var destDir = flag.String("dest", "/volume1/music/FLAC", "destination directory")
var dryRun = flag.Bool("test", false, "test mode, don't actually move files")
var shell = flag.Bool("shell", false, "output shell script instead of moving files")

func main() {
	flag.Parse()
	for _, fname := range flag.Args() {
		if err := filepath.Walk(fname, examine); err != nil {
			fmt.Printf("error scanning %s: %v\n", fname, err)
		}
	}
}
