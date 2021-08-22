package main

import (
	"bytes"
	"encoding/binary"
	"io"
	"unicode/utf16"
)

type DAT struct {
	*bytes.Reader
}

func NewDAT(input []byte) DAT {
	return DAT{bytes.NewReader(input)}
}

func (d *DAT) ReadOffset(offset uint32) []rune {
	d.Seek(int64(offset), io.SeekStart)

	var full []uint16

	for {
		one, _ := d.ReadByte()
		two, _ := d.ReadByte()

		current := binary.BigEndian.Uint16([]byte{one, two})
		if current == 0 {
			// It seems our string has been terminated.
			break
		}

		full = append(full, current)
	}

	return utf16.Decode(full)
}
