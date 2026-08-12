package main

import (
	"encoding/json"
	"fmt"
	"strings"
)

type JsonSourceMapContent struct {
	Version        int      `json:"version"`
	Sources        []string `json:"sources"`
	SourcesContent []string `json:"sourcesContent"`
	Mappings       string   `json:"mappings"`
	Names          []string `json:"names"`
}

type VLQMapping struct {
	Column       int
	FileIndex    int
	NameIndex    int
	SourceRow    int
	SourceColumn int
}

type SourceMap struct {
	JsonSourceMapContent
	Mappings []VLQMapping
}

var vlqAlphabet = map[byte]byte{
	'A': 0,
	'B': 1,
	'C': 2,
	'D': 3,
	'E': 4,
	'F': 5,
	'G': 6,
	'H': 7,
	'I': 8,
	'J': 9,
	'K': 10,
	'L': 11,
	'M': 12,
	'N': 13,
	'O': 14,
	'P': 15,
	'Q': 16,
	'R': 17,
	'S': 18,
	'T': 19,
	'U': 20,
	'V': 21,
	'W': 22,
	'X': 23,
	'Y': 24,
	'Z': 25,
	'a': 26,
	'b': 27,
	'c': 28,
	'd': 29,
	'e': 30,
	'f': 31,
	'g': 32,
	'h': 33,
	'i': 34,
	'j': 35,
	'k': 36,
	'l': 37,
	'm': 38,
	'n': 39,
	'o': 40,
	'p': 41,
	'q': 42,
	'r': 43,
	's': 44,
	't': 45,
	'u': 46,
	'v': 47,
	'w': 48,
	'x': 49,
	'y': 50,
	'z': 51,
	'0': 52,
	'1': 53,
	'2': 54,
	'3': 55,
	'4': 56,
	'5': 57,
	'6': 58,
	'7': 59,
	'8': 60,
	'9': 61,
	'+': 62,
	'/': 63,
}

func DecodeVLQ(chunk string) ([5]int, error) {
	vlq := [5]int{0, 0, 0, 0, -1}

	var i, j, k int
	for j < len(chunk) {
		b, ok := vlqAlphabet[chunk[j]]
		if !ok {
			return vlq, fmt.Errorf("Invalid character in VLQ chunk.")
		}

		value := int(b & 0b011110 >> 1)
		sign := 1
		if b&0b000001 != 0 {
			sign = -1
		}

		for k = 1; b&0b100000 != 0; k++ {
			b, ok = vlqAlphabet[chunk[j+k]]
			if !ok {
				return vlq, fmt.Errorf("Invalid continuation character in VLQ chunk.")
			}
			value += int(b&0b011111) << (k * 4)
		}
		j += k

		vlq[i] = value * sign
		i++
	}

	return vlq, nil
}

func ParseSourceMap(data []byte) (sourceMap *SourceMap, err error) {
	jsonContent := &JsonSourceMapContent{}
	err = json.Unmarshal(data, jsonContent)
	if err != nil {
		return nil, err
	}

	sourceMap = &SourceMap{
		JsonSourceMapContent: *jsonContent,
		Mappings:             []VLQMapping{},
	}

	for a := range strings.SplitSeq(jsonContent.Mappings, ";") {
		for b := range strings.SplitSeq(a, ",") {
			mapping := VLQMapping{}
			vlq, _ := DecodeVLQ(b)
			mapping.Column = vlq[0]
			mapping.FileIndex = vlq[1]
			mapping.SourceRow = vlq[2]
			mapping.SourceColumn = vlq[3]
			mapping.NameIndex = vlq[4]
			sourceMap.Mappings = append(sourceMap.Mappings, mapping)
		}
	}

	return sourceMap, nil
}
