package jukebox

import "strings"

const DoubleDashes string = "--"

func DecodeValue(encodedValue string) string {
	return strings.Replace(encodedValue, "-", " ", -1)
}

func EncodeValue(value string) string {
	cleanValue := RemovePunctuation(value)
	return strings.Replace(cleanValue, " ", "-", -1)
}

func EncodeArtistAlbum(artist string, album string) string {
	return EncodeValue(artist) + DoubleDashes + EncodeValue(album)
}

func EncodeArtistAlbumSong(artist string, album string, song string) string {
	return EncodeArtistAlbum(artist, album) + DoubleDashes + EncodeValue(song)
}

func RemovePunctuation(s string) string {
	if strings.Contains(s, "'") {
		s = strings.Replace(s, "'", "", -1)
	}

	if strings.Contains(s, "!") {
		s = strings.Replace(s, "!", "", -1)
	}

	if strings.Contains(s, "?") {
		s = strings.Replace(s, "?", "", -1)
	}

	return s
}
