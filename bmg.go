package main

import (
	"bytes"
	"encoding/binary"
	"encoding/xml"
	"errors"
	"io"
	"log"
	"strings"
)

var (
	// FileMagic is the byte representation of "MESGbmg1".
	FileMagic = []byte{'M', 'E', 'S', 'G', 'b', 'm', 'g', '1'}

	ErrInvalidMagic        = errors.New("provided BMG has an invalid magic")
	ErrUnsupportedEncoding = errors.New("provided BMG uses a text encoding not supported")
	ErrInvalidSection      = errors.New("provided BMG has an invalid section size")
)

// CharsetTypes represents the enum for possible charsets within a BMG.
type CharsetTypes byte

const (
	CharsetUndefined CharsetTypes = iota
	CharsetCP1252
	CharsetUTF16
	CharsetShiftJIS
	CharsetUTF8
)

const (
	NullStringPlaceholder  = "==== THIS STRING INTENTIONALLY LEFT NULL ===="
	LessThanPlaceholder    = "##LESS_THAN_SYMBOL##"
	GreaterThanPlaceholder = "##GREATER_THAN_SYMBOL##"
)

// BMG represents the internal structure of our BMG.
type BMG struct {
	INF *INF
	DAT DAT
	MID []MID
}

// SectionTypes are known parts of a BMG.
type SectionTypes [4]byte

var (
	SectionTypeINF1 SectionTypes = [4]byte{'I', 'N', 'F', '1'}
	SectionTypeDAT1 SectionTypes = [4]byte{'D', 'A', 'T', '1'}
	SectionTypeMID1 SectionTypes = [4]byte{'M', 'I', 'D', '1'}
)

// BMGHeader is taken from http://wiki.tockdom.com/w/index.php?title=BMG_%28File_Format%29.
type BMGHeader struct {
	// Magic is only "MESG" followed by "bmg1",
	// but bmg1 following sequentially is the only value.
	Magic        [8]byte
	FileSize     uint32
	SectionCount uint32
	Charset      CharsetTypes
	// Appears to be padding.
	_ [15]byte
}

// SectionHeader allows us to read section header info.
type SectionHeader struct {
	Type SectionTypes
	Size uint32
}

// XMLFormat specifies XML necessities for marshalling and unmarshalling.
type XMLFormat struct {
	MessageID  MID    `xml:"key,attr"`
	Attributes uint32 `xml:"attributes,attr"`
	String     string `xml:",innerxml"`
}

func (b BMG) ReadString(entry INFEntry) []rune {
	return b.DAT.ReadOffset(entry.Offset)
}

func parseBMG(data []byte) ([]byte, error) {
	// Create a new reader for serialization
	readable := bytes.NewReader(data)

	var header BMGHeader
	err := binary.Read(readable, binary.BigEndian, &header)
	if err != nil {
		return nil, err
	}

	// Validate header
	if !bytes.Equal(FileMagic[:], header.Magic[:]) {
		return nil, ErrInvalidMagic
	}
	if readable.Size() != int64(header.FileSize) {
		return nil, io.ErrUnexpectedEOF
	}
	if header.Charset != CharsetUTF16 {
		return nil, ErrUnsupportedEncoding
	}

	var currentBMG BMG

	// Read sections
	for count := header.SectionCount; count != 0; count-- {
		var sectionHeader SectionHeader
		err = binary.Read(readable, binary.BigEndian, &sectionHeader)
		if err != nil {
			return nil, err
		}

		// Subtract the header size
		sectionSize := int(sectionHeader.Size) - 8
		temp := make([]byte, sectionSize)
		_, err = readable.Read(temp)
		if err != nil {
			return nil, err
		}

		// Add to sections
		switch sectionHeader.Type {
		case SectionTypeINF1:
			currentBMG.INF, err = NewINF(temp)
		case SectionTypeDAT1:
			currentBMG.DAT = NewDAT(temp)
		case SectionTypeMID1:
			currentBMG.MID, err = NewMID(temp)
		default:
			log.Println("unhandled type", string(sectionHeader.Type[:]))
		}

		if err != nil {
			return nil, err
		}
	}

	if len(currentBMG.INF.Entries) != len(currentBMG.MID) {
		return nil, ErrInvalidSection
	}

	var output []XMLFormat
	for index, entry := range currentBMG.INF.Entries {
		currentString := string(currentBMG.ReadString(entry))
		currentString = strings.ReplaceAll(currentString, "<", LessThanPlaceholder)
		currentString = strings.ReplaceAll(currentString, ">", GreaterThanPlaceholder)
		if currentString == "" {
			currentString = NullStringPlaceholder
		}

		xmlNode := XMLFormat{
			MessageID:  currentBMG.MID[index],
			Attributes: binary.BigEndian.Uint32(entry.Attributes[:]),
			String:     currentString,
		}
		output = append(output, xmlNode)
	}

	type Translations struct {
		XMLName     xml.Name    `xml:"root"`
		Translation []XMLFormat `xml:"str"`
	}

	return xml.MarshalIndent(Translations{Translation: output}, "", "\t")
}

func createBMG(input []byte) ([]byte, error) {
	return nil, nil
}
