// Copyright 2015, Xavier Henner
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package tag

import (
	"bytes"
	"encoding/binary"
	"errors"
	"fmt"
	"io"
	"os"
	"strconv"
	"unicode/utf16"
)

func WriteID3v2Tags(rw io.ReadWriteSeeker, m Metadata) (*ID3v2TagUpdater, error) {
	var w ID3v2TagUpdater
	pos, err := rw.Seek(0, os.SEEK_CUR)
	if err != nil {
		return nil, err
	}
	w.usedSize = pos
	if id3, ok := (m).(metadataID3v2); !ok {
		return nil, errors.New("Not an ID3 ???")
	} else {
		w.m = id3
	}
	w.rw = rw
	w.updated = make(map[string]bool)
	return &w, nil
}

type ID3v2TagUpdater struct {
	rw       io.ReadWriteSeeker
	m        metadataID3v2
	updated  map[string]bool
	usedSize int64
}

func writeID3v2FrameHeader(format Format, tag string, length int) ([]byte, error) {
	r := []byte(frames.Name(tag, format))
	switch format {
	case ID3v2_2:
		size, err := WriteInt(length, 3)
		if err != nil {
			return r, err
		}
		r = append(r, size...)
	default:
		size, err := WriteInt(length, 4)
		if err != nil {
			return r, err
		}
		r = append(r, size...)
		r = append(r, 0x00, 0x00)
	}
	return r, nil
}

func (w *ID3v2TagUpdater) updateTFrame(s string, f string) error {
	if w.m.frames[frames.Name(f, w.m.Format())] == s {
		// not an update
		return nil
	}
	var b bytes.Buffer
	b.Write([]byte{0x01, 0xFF, 0xFE})
	binary.Write(&b, binary.LittleEndian, utf16.Encode([]rune(s)))
	b.Write([]byte{0x00, 0x00})
	frame, err := writeID3v2FrameHeader(w.m.Format(), f, b.Len())
	if err != nil {
		return nil
	}

	frame = append(frame, b.Bytes()...)

	w.m.frames[frames.Name(f, w.m.Format())] = frame
	w.updated[frames.Name(f, w.m.Format())] = true
	return nil
}

func (w *ID3v2TagUpdater) Title(s string) error {
	return w.updateTFrame(s, "title")
}
func (w *ID3v2TagUpdater) Album(s string) error {
	return w.updateTFrame(s, "album")
}
func (w *ID3v2TagUpdater) Artist(s string) error {
	return w.updateTFrame(s, "artist")
}
func (w *ID3v2TagUpdater) AlbumArtist(s string) error {
	return w.updateTFrame(s, "album_artist")
}
func (w *ID3v2TagUpdater) Composer(s string) error {
	return w.updateTFrame(s, "composer")
}
func (w *ID3v2TagUpdater) Year(s int) error {
	return w.updateTFrame(strconv.Itoa(s), "year")
}
func (w *ID3v2TagUpdater) Genre(s string) error {
	return w.updateTFrame(s, "genre")
}
func (w *ID3v2TagUpdater) Track(t int, total int) error {
	s := strconv.Itoa(t)
	if total > 1 {
		s += "/" + strconv.Itoa(total)
	}
	return w.updateTFrame(s, "disc")
}
func (w *ID3v2TagUpdater) Disc(d int, total int) error {
	s := strconv.Itoa(d)
	if total > 1 {
		s += "/" + strconv.Itoa(total)
	}
	return w.updateTFrame(s, "track")
}
func (w *ID3v2TagUpdater) Picture(*Picture) error {
	return errors.New("need some works")
	return nil
}
func (w *ID3v2TagUpdater) Lyrics(s string) error {
	return errors.New("need some works")
	return w.updateTFrame(s, "lyrics")
}
func (u *ID3v2TagUpdater) Commit() error {
	var b bytes.Buffer
	var c io.Reader
	// we have to rewind to the start
	_, err := u.rw.Seek(0, os.SEEK_SET)
	if err != nil {
		return err
	}
	h, err := readID3v2Header(u.rw)
	if err != nil {
		return err
	}

	if h.Unsynchronisation {
		c = &unsynchroniser{Reader: u.rw}
	} else {
		c = u.rw
	}
	// store the current position
	offset := 10

	for _, k := range sortMapByValue(u.m.positions) {
		if d, err := readBytes(c, k.Value-offset); err == nil && !u.updated[k.Key] {
			b.Write(d)
		}
		offset = k.Value
	}
	for k := range u.updated {
		b.Write((u.m.frames[k]).([]byte))
	}

	fmt.Println(b.Bytes())

	return nil
}
