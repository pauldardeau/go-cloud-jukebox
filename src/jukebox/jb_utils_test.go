package jukebox

import "testing"

func TestDecodeValue(t *testing.T) {
	artist := DecodeValue("The-Who")
	if artist != "The Who" {
		t.Fail()
	}

	album := DecodeValue("Whos-Next")
	if album != "Whos Next" {
		t.Fail()
	}

	song := DecodeValue("My-Wife")
	if song != "My Wife" {
		t.Fail()
	}
}

func TestEncodeValue(t *testing.T) {
	artist := EncodeValue("The Who")
	if artist != "The-Who" {
		t.Fail()
	}

	album := EncodeValue("Whos Next")
	if album != "Whos-Next" {
		t.Fail()
	}

	song := EncodeValue("My Wife")
	if song != "My-Wife" {
		t.Fail()
	}
}
