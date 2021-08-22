package main

import (
	"bytes"
	"encoding/binary"
)

// INF represents the INF1/INFO format within a BMG.
type INF struct {
	INFHeader
	Entries []INFEntry
}

// INFHeader represents the initial header of an INF1 block.
type INFHeader struct {
	EntryCount uint16
	// This is the length of an INFEntry.
	// We assume a size of 8 and hardcode such.
	// It's unknown where values may differ.
	EntryLength  uint16
	GroupID      uint16
	DefaultColor uint8
	_            uint8
}

// INFEntry represents an info entry within the INF1 block.
// It's assumed every entry's size is 8. Tweak accordingly if such is fale.
type INFEntry struct {
	Offset     uint32
	Attributes [4]byte
}

func NewINF(data []byte) (*INF, error) {
	// Parse initial header
	readable := bytes.NewReader(data)

	var header INFHeader
	err := binary.Read(readable, binary.BigEndian, &header)
	if err != nil {
		return nil, err
	}

	entries := make([]INFEntry, header.EntryCount)
	err = binary.Read(readable, binary.BigEndian, &entries)
	if err != nil {
		return nil, err
	}

	return &INF{
		header,
		entries,
	}, nil
}
