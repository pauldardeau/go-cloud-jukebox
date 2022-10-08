package jukebox

type Artist struct {
   ArtistUid string
   ArtistName string
   ArtistDescription string
}

func NewArtist(artistUid string,
               artistName string,
               artistDescription string) (*Artist) {
   var artist Artist
   artist.ArtistUid = artistUid
   artist.ArtistName = artistName
   artist.ArtistDescription = artistDescription
   return &artist
}
