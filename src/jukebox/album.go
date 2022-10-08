package jukebox

type Album struct {
   AlbumUid string
   AlbumName string
   ArtistUid string
   GenreUid string
   AlbumDescription string
}

func NewAlbum(albumUid string,
              albumName string,
	      artistUid string,
	      genreUid string,
	      albumDescription string) (*Album) {
   var album Album
   album.AlbumUid = albumUid
   album.AlbumName = albumName
   album.ArtistUid = artistUid
   album.GenreUid = genreUid
   album.AlbumDescription = albumDescription
   return &album
}

