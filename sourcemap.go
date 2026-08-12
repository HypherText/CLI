package main

import (
	"encoding/json"
	"fmt"

	"github.com/kstenerud/go-vlq"
)

type SourceMap struct {
	Version        int      `json:"version"`
	Sources        []string `json:"sources"`
	SourcesContent []string `json:"sourcesContent"`
	Mappings       string   `json:"mappings"`
	Names          []string `json:"names"`
}

func ParseSourceMap(data []byte) (*SourceMap, error) {
	var sm *SourceMap
	err := json.Unmarshal(data, sm)
	if err != nil {
		return nil, err
	}

	return sm, nil
}

func (sm *SourceMap) ParseMappings() {
	fmt.Println(sm.Mappings)
	decoded, _, _ := vlq.DecodeRvlqFrom([]byte(sm.Mappings))
	fmt.Println(decoded)
}
