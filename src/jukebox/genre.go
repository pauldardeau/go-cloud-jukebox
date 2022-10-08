package jukebox

type Genre struct {
   GenreUid string
   GenreName string
   GenreDescription string
}

func NewGenre(genreUid string,
              genreName string,
              genreDescription string) *Genre {
   var genre Genre
   genre.GenreUid = genreUid
   genre.GenreName = genreName
   genre.GenreDescription = genreDescription
   return &genre
}

