package jukebox

type Genre struct {
   Genre_uid string
   Genre_name string
   Genre_description string
}

func NewGenre(genre_uid string,
              genre_name string,
	      genre_description string) (*Genre) {
   var genre Genre
   genre.Genre_uid = genre_uid
   genre.Genre_name = genre_name
   genre.Genre_description = genre_description
   return &genre
}

