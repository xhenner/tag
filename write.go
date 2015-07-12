// Copyright 2015, Xavier Henner
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package tag

import (
	"errors"
	"io"
)

func WriteTo(rw io.ReadWriteSeeker) (TagUpdater, error) {
	m, err := ReadFrom(rw)
	if err != nil {
		return nil, err
	}
	switch m.Format() {
	case ID3v2_3:
		return WriteID3v2Tags(rw, m)
	}
	return nil, errors.New("Write Support not available yet")
}

type TagUpdater interface {
	// Add or Replace the title of the track.
	Title(string) error

	// Add or Replace the album name of the track.
	Album(string) error

	// Add or Replace the artist name of the track.
	Artist(string) error

	// Add or Replace the album artist name of the track.
	AlbumArtist(string) error

	// Add or Replace the composer of the track.
	Composer(string) error

	// Add or Replace the year of the track.
	Year(int) error

	// Add or Replace the genre of the track.
	Genre(string) error

	// Add or Replace the track number and total tracks
	Track(int, int) error

	// Add or Replace the disc number and total discs
	Disc(int, int) error

	// Add or Replace the Picture. Remove all old Picture informations
	Picture(*Picture) error

	// Add or Replace Lyrics. Reomve all old Lyrics informations
	Lyrics(string) error

	// Write the new tags
	Commit() error
}
