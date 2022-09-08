package jukebox

type Artist struct {
   Artist_uid string
   Artist_name string
   Artist_description string
}

func NewArtist(artist_uid string,
               artist_name string,
	       artist_description string) (*Artist) {
   var artist Artist
   artist.Artist_uid = artist_uid
   artist.Artist_name = artist_name
   artist.Artist_description = artist_description
   return &artist
}
