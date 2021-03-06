// Copyright 2015, David Howden
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package tag provides MP3 (ID3: v1, 2.2, 2.3 and 2.4), MP4, FLAC and OGG metadata detection,
// parsing and artwork extraction.
//
// Detect and parse tag metadata from an io.ReadSeeker (i.e. an *os.File):
// 	m, err := tag.ReadFrom(f)
// 	if err != nil {
// 		log.Fatal(err)
// 	}
// 	log.Print(m.Format()) // The detected format.
// 	log.Print(m.Title())  // The title of the track (see Metadata interface for more details).
package tag

import (
	"errors"
	"fmt"
	"io"
	"os"
)

// ErrNoTagsFound is the error returned by ReadFrom when the metadata format
// cannot be identified.
var ErrNoTagsFound = errors.New("no tags found")

// ReadFrom detects and parses audio file metadata tags (currently supports ID3v1,2.{2,3,4}, MP4, FLAC/OGG).
// Returns non-nil error if the format of the given data could not be determined, or if there was a problem
// parsing the data.
func ReadFrom(r io.ReadSeeker) (Metadata, error) {
	b, err := readBytes(r, 11)
	if err != nil {
		return nil, err
	}

	_, err = r.Seek(-11, os.SEEK_CUR)
	if err != nil {
		return nil, fmt.Errorf("could not seek back to original position: %v", err)
	}

	switch {
	case string(b[0:4]) == "fLaC":
		return ReadFLACTags(r)

	case string(b[0:4]) == "OggS":
		return ReadOGGTags(r)

	case string(b[4:11]) == "ftypM4A":
		return ReadAtoms(r)

	case string(b[0:3]) == "ID3":
		return ReadID3v2Tags(r)
	}

	m, err := ReadID3v1Tags(r)
	if err != nil {
		if err == ErrNotID3v1 {
			err = ErrNoTagsFound
		}
		return nil, err
	}
	return m, nil
}

// Format is an enumeration of metadata types supported by this package.
type Format string

// Supported tag formats.
const (
	UnknownFormat Format = ""        // Unknown Format.
	ID3v1                = "ID3v1"   // ID3v1 tag format.
	ID3v2_2              = "ID3v2.2" // ID3v2.2 tag format.
	ID3v2_3              = "ID3v2.3" // ID3v2.3 tag format (most common).
	ID3v2_4              = "ID3v2.4" // ID3v2.4 tag format.
	MP4                  = "MP4"     // MP4 tag (atom) format.
	VORBIS               = "VORBIS"  // Vorbis Comment tag format.
)

// FileType is an enumeration of the audio file types supported by this package, in particular
// there are audio file types which share metadata formats, and this type is used to distinguish
// between them.
type FileType string

// Supported file types.
const (
	UnknownFileType FileType = ""     // Unknown FileType.
	MP3                      = "MP3"  // MP3 file
	AAC                      = "AAC"  // M4A file (MP4)
	ALAC                     = "ALAC" // Apple Lossless file FIXME: actually detect this
	FLAC                     = "FLAC" // FLAC file
	OGG                      = "OGG"  // OGG file
)

// Metadata is an interface which is used to describe metadata retrieved by this package.
type Metadata interface {
	// Format returns the metadata Format used to encode the data.
	Format() Format

	// FileType returns the file type of the audio file.
	FileType() FileType

	// Title returns the title of the track.
	Title() string

	// Album returns the album name of the track.
	Album() string

	// Artist returns the artist name of the track.
	Artist() string

	// AlbumArtist returns the album artist name of the track.
	AlbumArtist() string

	// Composer returns the composer of the track.
	Composer() string

	// Year returns the year of the track.
	Year() int

	// Genre returns the genre of the track.
	Genre() string

	// Track returns the track number and total tracks, or zero values if unavailable.
	Track() (int, int)

	// Disc returns the disc number and total discs, or zero values if unavailable.
	Disc() (int, int)

	// Picture returns a picture, or nil if not available.
	Picture() *Picture

	// Lyrics returns the lyrics, or an empty string if unavailable.
	Lyrics() string

	// Raw returns the raw mapping of retrieved tag names and associated values.
	// NB: tag/atom names are not standardised between formats.
	Raw() map[string]interface{}
}
