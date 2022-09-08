package jukebox


type Playlist struct {
   Uid string
   Name string
   Description string
   Songs []*SongMetadata
}

func NewPlaylist(playlist_uid string,
                 playlist_name string,
                 playlist_description string) (*Playlist) {
   var pl Playlist
   pl.Uid = playlist_uid
   pl.Name = playlist_name
   pl.Description = playlist_description
   pl.Songs = []*SongMetadata{}
   return &pl
}
