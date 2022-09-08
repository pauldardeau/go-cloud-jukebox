package jukebox

type Album struct {
   Album_uid string
   Album_name string
   Artist_uid string
   Genre_uid string
   Album_description string
}

func NewAlbum(album_uid string,
              album_name string,
	      artist_uid string,
	      genre_uid string,
	      album_description string) (*Album) {
   var album Album
   album.Album_uid = album_uid
   album.Album_name = album_name
   album.Artist_uid = artist_uid
   album.Genre_uid = genre_uid
   album.Album_description = album_description
   return &album
}

