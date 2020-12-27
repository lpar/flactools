package flactools

import (
	"fmt"
	"strings"
	"unicode"

	"github.com/mewkiz/flac"
	"github.com/mewkiz/flac/meta"
)

func Coalesce(args ...string) string {
	for _, arg := range args {
		if arg != "" {
			return arg
		}
	}
	return ""
}

func GetTags(fname string) (map[string]string, error) {
	vc, err := GetVorbisComment(fname)
	if err != nil {
		return nil, err
	}
	return TagMap(vc), nil
}

func GetVorbisComment(fname string) (*meta.VorbisComment, error) {
	stream, err := flac.ParseFile(fname)
	if err != nil {
		return nil, fmt.Errorf("can't parse %s: %w", fname, err)
	}
	defer stream.Close()
	mdblocks := stream.Blocks
	for _, blk := range mdblocks {
		if blk.Type == meta.TypeVorbisComment {
			vc := blk.Body.(*meta.VorbisComment)
			return vc, nil
		}
	}
	return nil, fmt.Errorf("no tags found in %s", fname)
}

func TagMap(vc *meta.VorbisComment) map[string]string {
	tuplelist := vc.Tags
	tags := make(map[string]string)
	for _, tuple := range tuplelist {
		k := tuple[0]
		v := tuple[1]
		tags[strings.ToUpper(k)] = v
	}
	return tags
}

func CleanName(x string) string {
	var dn strings.Builder
	spc := false // was the last character a space?
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
			if !spc {
				dn.WriteRune('_')
			}
			dn.WriteString("and_")
			spc = true
		}
	}
	return dn.String()
}
