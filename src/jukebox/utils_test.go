package jukebox

import (
	"testing"
)

func TestUnencodeValue(t *testing.T) {
	artist := UnencodeValue("The-Who")
	if artist != "The Who" {
		t.Fail()
	}

	album := UnencodeValue("Whos-Next")
	if album != "Whos Next" {
		t.Fail()
	}

	song := UnencodeValue("My-Wife")
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

func TestFileExists(t *testing.T) {
}

func TestDeleteFile(t *testing.T) {
}

func TestDeleteFilesInDirectory(t *testing.T) {
}

func TestDirectoryExists(t *testing.T) {
}

func TestCreateDirectory(t *testing.T) {
}

func TestListDirsInDirectory(t *testing.T) {
}

func TestDirectoryDeleteDirectory(t *testing.T) {
}

func TestListFilesInDirectory(t *testing.T) {
}

func TestPathJoin(t *testing.T) {
}

func TestPathSplitExt(t *testing.T) {
	var root string
	var ext string

	// splitext("bar") -> ("bar", "")
	root, ext = PathSplitExt("bar")
	if root != "bar" || ext != "" {
		t.Fail()
	}

	// splitext("foo.bar.exe") -> ("foo.bar", ".exe")
	root, ext = PathSplitExt("foo.bar.exe")
	if root != "foo.bar" || ext != ".exe" {
		t.Fail()
	}

	// splitext("/foo/bar.exe") -> ("/foo/bar", ".exe")
	root, ext = PathSplitExt("/foo/bar.exe")
	if root != "/foo/bar" || ext != ".exe" {
		t.Fail()
	}

	// splitext(".cshrc") -> (".cshrc", "")
	root, ext = PathSplitExt(".cshrc")
	if root != ".cshrc" || ext != "" {
		t.Fail()
	}

	// splitext("/foo/....jpg") -> ("/foo/....jpg", "")
	root, ext = PathSplitExt("/foo/....jpg")
	if root != "/foo/....jpg" || ext != "" {
		t.Fail()
	}
}

func TestPathGetMtime(t *testing.T) {
}

func TestGetFileSize(t *testing.T) {
}

func TestFileReadAllText(t *testing.T) {
}

func TestFileWriteAllText(t *testing.T) {
}

func TestFileWriteAllBytes(t *testing.T) {
}

func TestFileReadAllBytes(t *testing.T) {
}

func TestMd5ForFile(t *testing.T) {
}

func TestTimeSleepSeconds(t *testing.T) {
}
