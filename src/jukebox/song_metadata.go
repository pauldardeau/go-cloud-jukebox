package jukebox

type SongMetadata struct {
	Fm         *FileMetadata
	ArtistUid  string
	ArtistName string
	AlbumUid   string
	SongName   string
}

func (sm *SongMetadata) Equals(other *SongMetadata) bool {
	if sm.Fm != nil {
		if !sm.Fm.Equals(other.Fm) {
			return false
		}
	} else {
		if other.Fm != nil {
			return false
		}
	}
	return sm.ArtistUid == other.ArtistUid &&
		sm.ArtistName == other.ArtistName &&
		sm.AlbumUid == other.AlbumUid &&
		sm.SongName == other.SongName
}

func NewSongMetadata() *SongMetadata {
	var sm SongMetadata
	sm.Fm = nil
	sm.ArtistUid = ""
	sm.ArtistName = "" // keep temporarily until ArtistUid is hooked up to artist table
	sm.AlbumUid = ""
	sm.SongName = ""
	return &sm
}

func (sm *SongMetadata) FromDictionary(dictionary map[string]string) {
	sm.FromDictionaryWithPrefix(dictionary, "")
}

func (sm *SongMetadata) FromDictionaryWithPrefix(dictionary map[string]string,
	prefix string) {
	sm.Fm = NewFileMetadata()
	sm.Fm.FromDictionaryWithPrefix(dictionary, prefix)

	if value, isPresent := dictionary[prefix+"ArtistUid"]; isPresent {
		sm.ArtistUid = value
	}

	if value, isPresent := dictionary[prefix+"ArtistName"]; isPresent {
		sm.ArtistName = value
	}

	if value, isPresent := dictionary[prefix+"AlbumUid"]; isPresent {
		sm.AlbumUid = value
	}

	if value, isPresent := dictionary[prefix+"SongName"]; isPresent {
		sm.SongName = value
	}
}

func (sm *SongMetadata) ToDictionary() map[string]string {
	return sm.ToDictionaryWithPrefix("")
}

func (sm *SongMetadata) ToDictionaryWithPrefix(prefix string) map[string]string {
	fmDict := make(map[string]string)
	sm.Fm.ToDictionaryWithPrefix(prefix)

	smDict := map[string]string{
		prefix + "ArtistUid":  sm.ArtistUid,
		prefix + "ArtistName": sm.ArtistName,
		prefix + "AlbumUid":   sm.AlbumUid,
		prefix + "SongName":   sm.SongName}

	for key, value := range fmDict {
		smDict[prefix+key] = value
	}

	return smDict
}
