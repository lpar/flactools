# flactools

_Go-based tools for FLAC music collection maintenance_

This repository contains tools I use to manage FLAC music files on my
Synology NAS.

At the moment, the two utilities are `flaccat` and `flacdiff`. Used together,
they offer a way to check your music library for additions, deletions, and file
moves.

    % flaccat /volume1/FLAC > ~/flaclist.old
    [ some weeks later ]
    % flaccat /volume1/FLAC > ~/flaclist.new
    % flaccat ~/flaclist.old ~/flaclist.new

## Installation

Because the utilities are all written in pure Go, and use only pure Go
libraries, it should be easy to cross-compile them for any Synology box,
whether it's Intel-based or ARM-based.

For example, to build a 64-bit ARM binary of flaccat:

    cd flaccat
    GOOS=linux GOARCH=arm64 go build

## License

GPL v3.

## Contact

mathew  
<meta@pobox.com>

