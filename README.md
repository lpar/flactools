# flactools

_Go-based tools for FLAC music collection maintenance_

This repository contains tools I use to manage FLAC music files on my
Synology NAS.

At the moment, the two main utilities are `flaccat` and `flacdiff`. Used
together, they offer a way to check your music library for additions,
deletions, and file moves.

    % flaccat /volume1/FLAC > ~/flaclist.old
    [ some weeks later ]
    % flaccat /volume1/FLAC > ~/flaclist.new
    % flacdiff ~/flaclist.old ~/flaclist.new > report.txt

The third utility is `flacsum`, which checks the actual data of the FLAC files
against their recorded MD5 checksums. It's very slow, so you probably won't
want to run it very often.

The fourth utility is `flacmv`, which uses FLAC metadata to move folders around
based on artist and album name, putting them under `Artist\Album`. It can also
write out a shell script to perform the moves, so you can examine what it
suggests before doing it.

Check the source of each program for additional information that would be in
the man page if Synology boxes had man pages.

## Installation

Because the utilities are all written in pure Go, and use only pure Go
libraries, it should be easy to cross-compile them for any Synology box,
whether it's Intel-based or ARM-based.

For example, to build a 64-bit ARM binary of flaccat:

    cd flaccat
    GOOS=linux GOARCH=arm64 go build

Or to build a binary for an Intel-based Synology box on my Mac:

    GOOS=linux GOARCH=amd64 go build

## Bugs and limitations

For flacdiff, the entire list of FLAC files and their MD5 checksums must fit
into RAM. That hasn't been a problem for me yet, and I have 37,000 tracks.
If it becomes a problem I may look at using some sort of database, but for now
the performance of doing it all in RAM is too good to pass up.

## Credits

The heavy lifting is performed by https://github.com/mewkiz/flac

I took some Ruby scripts I had been using on my old Linux box, and rewrote them
in Go using the above library, so that I could easily build binaries for my
Syno box and get native code performance.

## License

GPL v3.

## Contact

mathew  
<meta@pobox.com>

