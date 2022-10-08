package jukebox


type Playlist struct {
   Uid string
   Name string
   Description string
   Songs []*SongMetadata
}

func NewPlaylist(playlistUid string,
                 playlistName string,
                 playlistDescription string) (*Playlist) {
   var pl Playlist
   pl.Uid = playlistUid
   pl.Name = playlistName
   pl.Description = playlistDescription
   pl.Songs = []*SongMetadata{}
   return &pl
}
