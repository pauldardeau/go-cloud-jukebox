package jukebox

import (
	"fmt"
)

type JukeboxOptions struct {
	DebugMode                bool
	CheckDataIntegrity       bool
	FileCacheCount           int
	NumberSongs              int
	SuppressMetadataDownload bool
}

func NewJukeboxOptions() *JukeboxOptions {
	var o JukeboxOptions
	o.DebugMode = false
	o.CheckDataIntegrity = false
	o.FileCacheCount = 3
	o.NumberSongs = 0
	o.SuppressMetadataDownload = false
	return &o
}

func printBoolValue(varName string, boolValue bool) {
	if boolValue {
		fmt.Printf("%s = true\n", varName)
	} else {
		fmt.Printf("%s = false\n", varName)
	}
}

func (o *JukeboxOptions) Show() {
	fmt.Println("========= Start JukeboxOptions ========")
	printBoolValue("DebugMode", o.DebugMode)
	printBoolValue("CheckDataIntegrity", o.CheckDataIntegrity)
	fmt.Printf("FileCacheCount = %d\n", o.FileCacheCount)
	fmt.Printf("NumberSongs = %d\n", o.NumberSongs)
	printBoolValue("SuppressMetadataDownload", o.SuppressMetadataDownload)
	fmt.Println("========= End JukeboxOptions =========")
}

func (o *JukeboxOptions) ValidateOptions() bool {
	if o.FileCacheCount < 0 {
		fmt.Println("error: file cache count must be non-negative integer value")
		return false
	}

	return true
}
