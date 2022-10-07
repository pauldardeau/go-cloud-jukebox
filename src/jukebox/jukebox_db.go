package jukebox

import (
    "database/sql"
    "fmt"

    _ "github.com/mattn/go-sqlite3"
)

// https://pkg.go.dev/database/sql

type JukeboxDB struct {
   debug_print bool
   use_encryption bool
   use_compression bool
   db_connection *sql.DB
   metadata_db_file_path string
}


func NewJukeboxDB(metadata_db_file_path string,
                  use_encryption bool,
		  use_compression bool,
                  debug_print bool) *JukeboxDB {
   var jukeboxDB JukeboxDB
   jukeboxDB.debug_print = true //debug_print
   jukeboxDB.use_encryption = use_encryption
   jukeboxDB.use_compression = use_compression
   jukeboxDB.db_connection = nil
   if len(metadata_db_file_path) > 0 {
      jukeboxDB.metadata_db_file_path = metadata_db_file_path
   } else {
      jukeboxDB.metadata_db_file_path = "jukebox_db.sqlite3"
   }
   return &jukeboxDB
}

func (jukeboxDB *JukeboxDB) is_open() bool {
   return jukeboxDB.db_connection != nil
}

func (jukeboxDB *JukeboxDB) open() bool {
   jukeboxDB.close()
   open_success := false
   db, err := sql.Open("sqlite3", jukeboxDB.metadata_db_file_path)
   if err != nil {
      fmt.Printf("error: unable to open SQLite db: %v\n", err)
   } else {
      jukeboxDB.db_connection = db
      if !jukeboxDB.have_tables() {
         open_success = jukeboxDB.create_tables()
         if !open_success {
            fmt.Println("error: unable to create all tables")
         }
      } else {
         open_success = true
      }
   }
   return open_success
}

func (jukeboxDB *JukeboxDB) close() bool {
   did_close := false
   if jukeboxDB.db_connection != nil {
      jukeboxDB.db_connection.Close()
      jukeboxDB.db_connection = nil
      did_close = true
   }
   return did_close
}

func (jukeboxDB *JukeboxDB) enter() bool {
    // look for stored metadata in the storage system
    if jukeboxDB.open() {
        if jukeboxDB.db_connection != nil {
            if jukeboxDB.debug_print {
                fmt.Println("have db connection")
            }
        }
    } else {
        fmt.Println("unable to connect to database")
	jukeboxDB.db_connection = nil
    }

    return jukeboxDB.db_connection != nil
}

func (jukeboxDB *JukeboxDB) exit() {
    if jukeboxDB.db_connection != nil {
        jukeboxDB.db_connection.Close()
        jukeboxDB.db_connection = nil
    }
}

func (jukeboxDB *JukeboxDB) create_table(sqlStatement string) bool {
    if jukeboxDB.db_connection != nil {
        stmt, err := jukeboxDB.db_connection.Prepare(sqlStatement)
	if err != nil {
            fmt.Printf("prepare of statement failed: %s\n", sqlStatement)
	    fmt.Printf("error: %v\n", err)
	    return false
	}
	defer stmt.Close()

        _, err_stmt_exec := stmt.Exec()
        if err_stmt_exec != nil {
            fmt.Println("creation of table failed")
            fmt.Print(sqlStatement)
            return false
        } else {
            return true
        }
    } else {
        return false
    }
}

func (jukeboxDB *JukeboxDB) create_tables() bool {
    if jukeboxDB.db_connection != nil {
        if jukeboxDB.debug_print {
            fmt.Println("creating tables")
        }

	create_genre_table := "CREATE TABLE genre (" +
                              "genre_uid TEXT UNIQUE NOT NULL, " +
                              "genre_name TEXT UNIQUE NOT NULL, " +
                              "genre_description TEXT);"

        create_artist_table := "CREATE TABLE artist (" +
                              "artist_uid TEXT UNIQUE NOT NULL," +
                              "artist_name TEXT UNIQUE NOT NULL," +
                              "artist_description TEXT)"

        create_album_table := "CREATE TABLE album (" +
                             "album_uid TEXT UNIQUE NOT NULL," +
                             "album_name TEXT UNIQUE NOT NULL," +
                             "album_description TEXT," +
                             "artist_uid TEXT NOT NULL REFERENCES artist(artist_uid)," +
                             "genre_uid TEXT REFERENCES genre(genre_uid))"

        create_song_table := "CREATE TABLE song (" +
                            "song_uid TEXT UNIQUE NOT NULL," +
                            "file_time TEXT," +
                            "origin_file_size INTEGER," +
                            "stored_file_size INTEGER," +
                            "pad_char_count INTEGER," +
                            "artist_name TEXT," +
                            "artist_uid TEXT REFERENCES artist(artist_uid)," +
                            "song_name TEXT NOT NULL," +
                            "md5_hash TEXT NOT NULL," +
                            "compressed INTEGER," +
                            "encrypted INTEGER," +
                            "container_name TEXT NOT NULL," +
                            "object_name TEXT NOT NULL," +
                            "album_uid TEXT REFERENCES album(album_uid))"

        create_playlist_table := "CREATE TABLE playlist (" +
                                "playlist_uid TEXT UNIQUE NOT NULL," +
                                "playlist_name TEXT UNIQUE NOT NULL," +
                                "playlist_description TEXT)"

        create_playlist_song_table := "CREATE TABLE playlist_song (" +
                                     "playlist_song_uid TEXT UNIQUE NOT NULL," +
                                     "playlist_uid TEXT NOT NULL REFERENCES playlist(playlist_uid)," +
                                     "song_uid TEXT NOT NULL REFERENCES song(song_uid))"

        return jukeboxDB.create_table(create_genre_table) &&
               jukeboxDB.create_table(create_artist_table) &&
               jukeboxDB.create_table(create_album_table) &&
               jukeboxDB.create_table(create_song_table) &&
               jukeboxDB.create_table(create_playlist_table) &&
               jukeboxDB.create_table(create_playlist_song_table)
    }

    return false
}

func (jukeboxDB *JukeboxDB) have_tables() bool {
   have_tables_in_db := false
   if jukeboxDB.db_connection != nil {
      sqlQuery := "SELECT name " +
                  "FROM sqlite_master " +
                  "WHERE type='table' AND name='song'"
      stmt, err := jukeboxDB.db_connection.Prepare(sqlQuery)
      if err != nil {
         fmt.Printf("error: unable to prepare sql: %s\n", sqlQuery)
	 fmt.Printf("error: %v\n", err)
	 return false
      }
      defer stmt.Close()

      var name string
      err = stmt.QueryRow().Scan(&name)
      if err == nil {
         have_tables_in_db = true
      }
   }

   return have_tables_in_db
}

func (jukeboxDB *JukeboxDB) get_playlist(playlist_name string) *string {
    var pl_object string
    if len(playlist_name) > 0 {
        sqlQuery := "SELECT playlist_uid FROM playlist WHERE playlist_name = ?"
        stmt, err := jukeboxDB.db_connection.Prepare(sqlQuery)
        if err != nil {
        }
        defer stmt.Close()

        err = stmt.QueryRow(playlist_name).Scan(&pl_object)
        if err != nil {
        }
    }
    return &pl_object
}

func (jukeboxDB *JukeboxDB) songs_for_query_results(rows *sql.Rows) []*SongMetadata {
    result_songs := []*SongMetadata{}

    for rows.Next() {
        var file_uid string
        var file_time string
        var o_file_size int64
        var s_file_size int64
        var pad_count int
        var artist_name string
        var artist_uid string
        var song_name string
        var md5_hash string
        var compressed int
        var encrypted int
        var container_name string
        var object_name string
        var album_uid string

        err := rows.Scan(&file_uid, &file_time, &o_file_size, &s_file_size,
                         &pad_count, &artist_name, &artist_uid, &song_name,
                         &md5_hash, &compressed, &encrypted, &container_name,
                         &object_name, &album_uid)

        if err != nil {
           fmt.Printf("error: scan of row values failed\n")
	   fmt.Printf("error: %v\n", err)
	   return result_songs
        }

        song := NewSongMetadata()
        song.Fm = NewFileMetadata()
        song.Fm.File_uid = file_uid
        song.Fm.File_time = file_time
        song.Fm.Origin_file_size = o_file_size
        song.Fm.Stored_file_size = s_file_size
        song.Fm.Pad_char_count = pad_count
        song.Artist_name = artist_name
        song.Artist_uid = artist_uid
        song.Song_name = song_name
        song.Fm.Md5_hash = md5_hash
        song.Fm.Compressed = compressed == 1
        song.Fm.Encrypted = encrypted == 1
        song.Fm.Container_name = container_name
        song.Fm.Object_name = object_name
        song.Album_uid = album_uid

        result_songs = append(result_songs, song)
    }

    return result_songs
}

func (jukeboxDB *JukeboxDB) retrieve_song(file_name string) *SongMetadata {
    if jukeboxDB.db_connection != nil {
        sqlQuery := `
            SELECT song_uid,
            file_time,
            origin_file_size,
            stored_file_size,
            pad_char_count,
            artist_name,
            artist_uid,
            song_name,
            md5_hash,
            compressed,
            encrypted,
            container_name,
            object_name,
            album_uid
            FROM song WHERE song_uid = ?
        `
        stmt, err := jukeboxDB.db_connection.Prepare(sqlQuery)
        if err != nil {
            fmt.Printf("error: unable to prepare statement '%s'\n", sqlQuery)
	    fmt.Printf("error: %v\n", err)
            return nil
        }
        defer stmt.Close()

	rows, err_rows := stmt.Query(file_name)
	if err_rows != nil {
            return nil
        }
	song_results := jukeboxDB.songs_for_query_results(rows)
        if song_results != nil && len(song_results) > 0 {
            return song_results[0]
        }
    }
    return nil
}

func (jukeboxDB *JukeboxDB) insert_playlist(pl_uid string,
                                            pl_name string,
                                            pl_desc string) bool {
    insert_success := false

    if jukeboxDB.db_connection != nil &&
       len(pl_uid) > 0 &&
       len(pl_name) > 0 {

        tx, err_tx := jukeboxDB.db_connection.Begin()
	if err_tx != nil {
            return false
        }

        sqlQuery := "INSERT INTO playlist VALUES (?,?,?)"
        stmt, err := tx.Prepare(sqlQuery)
        if err != nil {
            fmt.Printf("error: unable to prepare statement '%s'\n", sqlQuery)
            fmt.Printf("error: %v\n", err)
            return false
        }
        defer stmt.Close()

        stmt.Exec(pl_uid, pl_name, pl_desc)
        tx.Commit()
        insert_success = true
        //fmt.Println("error inserting playlist: " + e.args[0])
    }

    return insert_success
}

func (jukeboxDB *JukeboxDB) delete_playlist(pl_name string) bool {
    delete_success := false

    if jukeboxDB.db_connection != nil && len(pl_name) > 0 {
        sqlQuery := "DELETE FROM playlist " +
                    "WHERE playlist_name = ? "
        tx, err_tx := jukeboxDB.db_connection.Begin()
	if err_tx != nil {
            return false
        }
        stmt, err := tx.Prepare(sqlQuery)
        if err != nil {
            fmt.Printf("error: unable to prepare statement '%s'\n", sqlQuery)
            fmt.Printf("error: %v\n", err)
            return false
        }
        defer stmt.Close()
        stmt.Exec(pl_name)
        tx.Commit()
        delete_success = true
        //fmt.Println("error deleting playlist: " + e.args[0])
    }

    return delete_success
}

func (jukeboxDB *JukeboxDB) insert_song(song *SongMetadata) bool {
    insert_success := false

    if jukeboxDB.db_connection != nil && song != nil {
        sqlQuery := "INSERT INTO song VALUES (?,?,?,?,?,?,?,?,?,?,?,?,?,?)"
	tx, err_tx := jukeboxDB.db_connection.Begin()
	if err_tx != nil {
            return false
        }
        stmt, err := tx.Prepare(sqlQuery)
        if err != nil {
            fmt.Printf("error: unable to prepare statement '%s'\n", sqlQuery)
            fmt.Printf("error: %v\n", err)
            return false
        }
        defer stmt.Close()

        stmt.Exec(song.Fm.File_uid,
                  song.Fm.File_time,
                  song.Fm.Origin_file_size,
                  song.Fm.Stored_file_size,
                  song.Fm.Pad_char_count,
                  song.Artist_name,
                  "",
                  song.Song_name,
                  song.Fm.Md5_hash,
                  song.Fm.Compressed,
                  song.Fm.Encrypted,
                  song.Fm.Container_name,
                  song.Fm.Object_name,
                  song.Album_uid)
        tx.Commit()
        insert_success = true
        //fmt.Println("error inserting song: " + e.args[0])
    }

    return insert_success
}

func (jukeboxDB *JukeboxDB) update_song(song *SongMetadata) bool {
        update_success := false

        if jukeboxDB.db_connection != nil && song != nil && len(song.Fm.File_uid) > 0 {
            sqlQuery := `
                UPDATE song SET file_time=?,
                   origin_file_size=?,
                   stored_file_size=?,
                   pad_char_count=?,
                   artist_name=?,
                   artist_uid=?,
                   song_name=?,
                   md5_hash=?,
                   compressed=?,
                   encrypted=?,
                   container_name=?,
                   object_name=?,
                   album_uid=? WHERE song_uid = ?
            `
	    tx, err_tx := jukeboxDB.db_connection.Begin()
	    if err_tx != nil {
                return false
            }
            stmt, err := tx.Prepare(sqlQuery)
            if err != nil {
                fmt.Printf("error: unable to prepare statement '%s'\n", sqlQuery)
                fmt.Printf("error: %v\n", err)
                return false
            }

            defer stmt.Close()

            stmt.Exec(song.Fm.File_time,
                      song.Fm.Origin_file_size,
                      song.Fm.Stored_file_size,
                      song.Fm.Pad_char_count,
                      song.Artist_name,
                      "",
                      song.Song_name,
                      song.Fm.Md5_hash,
                      song.Fm.Compressed,
                      song.Fm.Encrypted,
                      song.Fm.Container_name,
                      song.Fm.Object_name,
                      song.Album_uid,
                      song.Fm.File_uid)
            tx.Commit()
            update_success = true
            //fmt.Println("error updating song: " + e.args[0])
        }

        return update_success
}

func (jukeboxDB *JukeboxDB) store_song_metadata(song *SongMetadata) bool {
    if song == nil {
       return false
    }
    db_song := jukeboxDB.retrieve_song(song.Fm.File_uid)
    if db_song != nil {
        if ! song.Equals(db_song) {
            return jukeboxDB.update_song(song)
        } else {
            return true  // no insert or update needed (already up-to-date)
        }
    } else {
        // song is not in the database, insert it
        return jukeboxDB.insert_song(song)
    }
}

func (jukeboxDB *JukeboxDB) sql_where_clause() string {
   var encryption int
   if jukeboxDB.use_encryption {
      encryption = 1
   } else {
      encryption = 0
   }

   var compression int
   if jukeboxDB.use_compression {
      compression = 1
   } else {
      compression = 0
   }

   where_clause := ""
   where_clause += " WHERE "
   where_clause += "encrypted = "
   where_clause += fmt.Sprintf("%d", encryption)
   where_clause += " AND "
   where_clause += "compressed = "
   where_clause += fmt.Sprintf("%d", compression)
   return where_clause
}

func (jukeboxDB *JukeboxDB) retrieve_songs(artist string,
                                           album string) []*SongMetadata {
    fmt.Printf("retrieve_songs entered\n")

    var songs []*SongMetadata
    if jukeboxDB.db_connection != nil {
        sqlQuery := `
            SELECT song_uid,
            file_time,
            origin_file_size,
            stored_file_size,
            pad_char_count,
            artist_name,
            artist_uid,
            song_name,
            md5_hash,
            compressed,
            encrypted,
            container_name,
            object_name,
            album_uid FROM song
        `

        sqlQuery += jukeboxDB.sql_where_clause()
        //if len(artist) > 0:
        //    sqlQuery += " AND artist_name='%s'" % artist
        if len(album) > 0 {
            encoded_artist := encode_value(artist)
            encoded_album := encode_value(album)
            added_clause := fmt.Sprintf(" AND object_name LIKE '%s--%s%%'",
                                        encoded_artist,
                                        encoded_album)
            sqlQuery += added_clause
        }

	//fmt.Printf("executing query: %s\n", sqlQuery)
	stmt, err_stmt := jukeboxDB.db_connection.Prepare(sqlQuery)
	if err_stmt != nil {
            fmt.Printf("error: unable to prepare statement '%s'\n", sqlQuery)
            fmt.Printf("error: %v\n", err_stmt)
            return nil
        }
	defer stmt.Close()

	rows, err := stmt.Query()
	if err != nil {
            fmt.Printf("error: unable to execute query of '%s'\n", sqlQuery)
	    fmt.Printf("error: %v\n", err)
            return nil
        }

        songs = jukeboxDB.songs_for_query_results(rows)
    }
    return songs
}

func (jukeboxDB* JukeboxDB) songs_for_artist(artist_name string) []*SongMetadata {
    songs := []*SongMetadata{}
    if jukeboxDB.db_connection != nil {
        sqlQuery := `
            SELECT song_uid,
            file_time,
            origin_file size,
            stored_file size,
            pad_char_count,
            artist_name,
            artist_uid,
            song_name,
            md5_hash,
            compressed,
            encrypted,
            container_name,
            object_name,
            album_uid FROM song
        `
        sqlQuery += jukeboxDB.sql_where_clause()
        sqlQuery += " AND artist = ?"
        stmt, err := jukeboxDB.db_connection.Prepare(sqlQuery)
        if err != nil {
           fmt.Printf("error: unable to prepare statement '%s'\n", sqlQuery)
           fmt.Printf("error: %v\n", err)
           return nil
        }
        defer stmt.Close()

        rows, err := stmt.Query(artist_name)
        if err != nil {
           return nil
        }
        songs = jukeboxDB.songs_for_query_results(rows)
    }
    return songs
}

func (jukeboxDB *JukeboxDB) show_listings() {
   if jukeboxDB.db_connection != nil {
      sqlQuery := "SELECT artist_name, song_name " +
                  "FROM song " +
                  "ORDER BY artist_name, song_name"
      stmt, err := jukeboxDB.db_connection.Prepare(sqlQuery)
      if err != nil {
          fmt.Printf("error: unable to prepare statement '%s'\n", sqlQuery)
	  fmt.Printf("error: %v\n", err)
          return
      }
      defer stmt.Close()

      rows, err := stmt.Query()
      if err != nil {
          return
      }

      for rows.Next() {
         var artist string
         var song string
         err := rows.Scan(&artist, &song)
         if err != nil {
         } else {
            fmt.Printf("%s, %s\n", artist, song)
         }
      }
   }
}

func (jukeboxDB *JukeboxDB) show_artists() {
   if jukeboxDB.db_connection != nil {
      sqlQuery := "SELECT DISTINCT artist_name " +
                  "FROM song " +
                  "ORDER BY artist_name"
      stmt, err := jukeboxDB.db_connection.Prepare(sqlQuery)
      if err != nil {
         fmt.Printf("error: unable to prepare statement '%s'\n", sqlQuery)
         fmt.Printf("error: %v\n", err)
         return
      }
      defer stmt.Close()
      rows, err := stmt.Query()
      if err != nil {
         return
      }
      for rows.Next() {
         var artist string
         err = rows.Scan(&artist)
         if err != nil {
         } else {
            fmt.Printf("%s\n", artist)
         }
      }
   }
}

func (jukeboxDB *JukeboxDB) show_genres() {
   if jukeboxDB.db_connection != nil {
      sqlQuery := "SELECT genre_name " +
                  "FROM genre " +
                  "ORDER BY genre_name"
      stmt, err := jukeboxDB.db_connection.Prepare(sqlQuery)
      if err != nil {
         fmt.Printf("error: unable to prepare statement '%s'\n", sqlQuery)
         fmt.Printf("error: %v\n", err)
         return
      }
      rows, err := stmt.Query()
      if err != nil {
         return
      }
      for rows.Next() {
         var genre_name string
         err = rows.Scan(&genre_name)
         if err != nil {
         } else {
            fmt.Printf("%s\n", genre_name)
         }
      }
   }
}

func (jukeboxDB *JukeboxDB) show_albums() {
   if jukeboxDB.db_connection != nil {
      sqlQuery := "SELECT album.album_name, artist.artist_name " +
                  "FROM album, artist " +
                  "WHERE album.artist_uid = artist.artist_uid " +
                  "ORDER BY album.album_name"
      stmt, err := jukeboxDB.db_connection.Prepare(sqlQuery)
      if err != nil {
         fmt.Printf("error: unable to prepare statement '%s'\n", sqlQuery)
         fmt.Printf("error: %v\n", err)
         return
      }
      rows, err := stmt.Query()
      if err != nil {
         return
      }

      for rows.Next() {
         var album_name string
         var artist_name string
         err = rows.Scan(&album_name, artist_name)
         if err != nil {
         } else {
            fmt.Printf("%s (%s)\n", album_name, artist_name)
         }
      }
   }
}

func (jukeboxDB *JukeboxDB) show_playlists() {
   if jukeboxDB.db_connection != nil {
      sqlQuery := "SELECT playlist_uid, playlist_name " +
                  "FROM playlist " +
                  "ORDER BY playlist_uid"
      stmt, err := jukeboxDB.db_connection.Prepare(sqlQuery)
      if err != nil {
         fmt.Printf("error: unable to prepare statement '%s'\n", sqlQuery)
         fmt.Printf("error: %v\n", err)
         return
      }
      rows, err := stmt.Query()
      if err != nil {
         return
      }
      for rows.Next() {
         var pl_uid string
         var pl_name string
         err = rows.Scan(&pl_uid, &pl_name)
         if err != nil {
         } else {
            fmt.Printf("%s - %s\n", pl_uid, pl_name)
         }
      }
   }
}

func (jukeboxDB *JukeboxDB) delete_song(song_uid string) bool {
   was_deleted := false
   if jukeboxDB.db_connection != nil {
      if len(song_uid) > 0 {
         sqlStatement := "DELETE FROM song WHERE song_uid = ?"
	 tx, err_tx := jukeboxDB.db_connection.Begin()
	 if err_tx != nil {
             return false
         }
         stmt, err := tx.Prepare(sqlStatement)
         if err != nil {
             fmt.Printf("error: unable to prepare statement '%s'\n", sqlStatement)
             fmt.Printf("error: %v\n", err)
             return false
         }
         defer stmt.Close()

         _, err = stmt.Exec(song_uid)
         if err != nil {
             tx.Rollback()
             fmt.Printf("error: unable to delete song '%s'\n", song_uid)
	     fmt.Printf("error: %v\n", err)
             return false
         }
	 tx.Commit()
         was_deleted = true
      }
   }

   return was_deleted
}
