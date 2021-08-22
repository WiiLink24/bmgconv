package main

import (
	"bytes"
	"encoding/binary"
)

type MID uint32

type MIDHeader struct {
	SectionCount uint16
	Format       uint8
	Info         uint8
	_            [4]byte
}

func NewMID(input []byte) ([]MID, error) {
	if len(input)%4 != 0 {
		return nil, ErrInvalidSection
	}

	readable := bytes.NewReader(input)

	var header MIDHeader
	err := binary.Read(readable, binary.BigEndian, &header)
	if err != nil {
		return nil, err
	}

	var contents []MID
	currentBytes := []byte{0, 0, 0, 0}
	for count := header.SectionCount; count != 0; count-- {
		_, err = readable.Read(currentBytes)
		if err != nil {
			return nil, err
		}

		current := binary.BigEndian.Uint32(currentBytes)
		contents = append(contents, MID(current))
	}

	return contents, nil
}
