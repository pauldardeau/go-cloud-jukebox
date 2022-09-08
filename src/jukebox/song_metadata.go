package jukebox

import (
)

type SongMetadata struct {
   Fm *FileMetadata
   Artist_uid string
   Artist_name string
   Album_uid string
   Song_name string
}

func (sm *SongMetadata) Equals(other *SongMetadata) bool {
   if sm.Fm != nil {
      if ! sm.Fm.Equals(other.Fm) {
         return false
      }
   } else {
      if other.Fm != nil {
         return false
      }
   }
   return sm.Artist_uid == other.Artist_uid &&
          sm.Artist_name == other.Artist_name &&
	  sm.Album_uid == other.Album_uid &&
	  sm.Song_name == other.Song_name
}

func NewSongMetadata() *SongMetadata {
   var sm SongMetadata
   sm.Fm = nil
   sm.Artist_uid = ""
   sm.Artist_name = ""  // keep temporarily until artist_uid is hooked up to artist table
   sm.Album_uid = ""
   sm.Song_name = ""
   return &sm
}

func (sm *SongMetadata) From_Dictionary(dictionary map[string]string) {
   sm.From_Dictionary_With_Prefix(dictionary, "")
}

func (sm *SongMetadata) From_Dictionary_With_Prefix(dictionary map[string]string,
                                                    prefix string) {
   sm.Fm = NewFileMetadata()
   sm.Fm.From_Dictionary_With_Prefix(dictionary, prefix)

   if value, isPresent := dictionary[prefix + "artist_uid"]; isPresent {
      sm.Artist_uid = value
   }

   if value, isPresent := dictionary[prefix + "artist_name"]; isPresent {
      sm.Artist_name = value
   }

   if value, isPresent := dictionary[prefix + "album_uid"]; isPresent {
      sm.Album_uid = value
   }

   if value, isPresent := dictionary[prefix + "song_name"]; isPresent {
      sm.Song_name = value
   }
}

func (sm *SongMetadata) To_Dictionary() map[string]string {
   return sm.To_Dictionary_With_Prefix("")
}

func (sm *SongMetadata) To_Dictionary_With_Prefix(prefix string) map[string]string {
   fm_dict := make(map[string]string)
   sm.Fm.To_Dictionary_With_Prefix(prefix)

   sm_dict := map[string]string {
           prefix + "artist_uid": sm.Artist_uid,
           prefix + "artist_name": sm.Artist_name,
           prefix + "album_uid": sm.Album_uid,
           prefix + "song_name": sm.Song_name}

   for key, value := range fm_dict {
      sm_dict[prefix + key] = value
   }

   return sm_dict
}
