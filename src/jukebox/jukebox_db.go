package jukebox

import (
    "database/sql"
    "fmt"

    _ "github.com/mattn/go-sqlite3"
)

// https://pkg.go.dev/database/sql

type JukeboxDB struct {
   debugPrint bool
   useEncryption bool
   useCompression bool
   dbConnection *sql.DB
   metadataDbFilePath string
}


func NewJukeboxDB(metadataDbFilePath string,
                  useEncryption bool,
                  useCompression bool,
                  debugPrint bool) *JukeboxDB {
   var jukeboxDB JukeboxDB
   jukeboxDB.debugPrint = true //debugPrint
   jukeboxDB.useEncryption = useEncryption
   jukeboxDB.useCompression = useCompression
   jukeboxDB.dbConnection = nil
   if len(metadataDbFilePath) > 0 {
      jukeboxDB.metadataDbFilePath = metadataDbFilePath
   } else {
      jukeboxDB.metadataDbFilePath = "jukebox_db.sqlite3"
   }
   return &jukeboxDB
}

func (jukeboxDB *JukeboxDB) isOpen() bool {
   return jukeboxDB.dbConnection != nil
}

func (jukeboxDB *JukeboxDB) open() bool {
   jukeboxDB.close()
   openSuccess := false
   db, err := sql.Open("sqlite3", jukeboxDB.metadataDbFilePath)
   if err != nil {
      fmt.Printf("error: unable to open SQLite db: %v\n", err)
   } else {
      jukeboxDB.dbConnection = db
      if !jukeboxDB.haveTables() {
         openSuccess = jukeboxDB.createTables()
         if !openSuccess {
            fmt.Println("error: unable to create all tables")
         }
      } else {
         openSuccess = true
      }
   }
   return openSuccess
}

func (jukeboxDB *JukeboxDB) close() bool {
   didClose := false
   if jukeboxDB.dbConnection != nil {
      jukeboxDB.dbConnection.Close()
      jukeboxDB.dbConnection = nil
      didClose = true
   }
   return didClose
}

func (jukeboxDB *JukeboxDB) enter() bool {
    // look for stored metadata in the storage system
    if jukeboxDB.open() {
        if jukeboxDB.dbConnection != nil {
            if jukeboxDB.debugPrint {
                fmt.Println("have db connection")
            }
        }
    } else {
        fmt.Println("unable to connect to database")
        jukeboxDB.dbConnection = nil
    }

    return jukeboxDB.dbConnection != nil
}

func (jukeboxDB *JukeboxDB) exit() {
    if jukeboxDB.dbConnection != nil {
        jukeboxDB.dbConnection.Close()
        jukeboxDB.dbConnection = nil
    }
}

func (jukeboxDB *JukeboxDB) createTable(sqlStatement string) bool {
    if jukeboxDB.dbConnection != nil {
        stmt, err := jukeboxDB.dbConnection.Prepare(sqlStatement)
        if err != nil {
            fmt.Printf("prepare of statement failed: %s\n", sqlStatement)
            fmt.Printf("error: %v\n", err)
            return false
        }
        defer stmt.Close()

        _, errStmtExec := stmt.Exec()
        if errStmtExec != nil {
            fmt.Println("error: creation of table failed")
            fmt.Print(sqlStatement)
            fmt.Printf("error: %v\n", errStmtExec)
            return false
        } else {
            return true
        }
    } else {
        return false
    }
}

func (jukeboxDB *JukeboxDB) createTables() bool {
    if jukeboxDB.dbConnection != nil {
        if jukeboxDB.debugPrint {
            fmt.Println("creating tables")
        }

        createGenreTable := "CREATE TABLE genre (" +
                            "genre_uid TEXT UNIQUE NOT NULL, " +
                            "genre_name TEXT UNIQUE NOT NULL, " +
                            "genre_description TEXT);"

        createArtistTable := "CREATE TABLE artist (" +
                             "artist_uid TEXT UNIQUE NOT NULL," +
                             "artist_name TEXT UNIQUE NOT NULL," +
                             "artist_description TEXT)"

        createAlbumTable := "CREATE TABLE album (" +
                            "album_uid TEXT UNIQUE NOT NULL," +
                            "album_name TEXT UNIQUE NOT NULL," +
                            "album_description TEXT," +
                            "artist_uid TEXT NOT NULL REFERENCES artist(artist_uid)," +
                            "genre_uid TEXT REFERENCES genre(genre_uid))"

        createSongTable := "CREATE TABLE song (" +
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

        createPlaylistTable := "CREATE TABLE playlist (" +
                               "playlist_uid TEXT UNIQUE NOT NULL," +
                               "playlist_name TEXT UNIQUE NOT NULL," +
                               "playlist_description TEXT)"

        createPlaylistSongTable := "CREATE TABLE playlist_song (" +
                                   "playlist_song_uid TEXT UNIQUE NOT NULL," +
                                   "playlist_uid TEXT NOT NULL REFERENCES playlist(playlist_uid)," +
                                   "song_uid TEXT NOT NULL REFERENCES song(song_uid))"

        return jukeboxDB.createTable(createGenreTable) &&
               jukeboxDB.createTable(createArtistTable) &&
               jukeboxDB.createTable(createAlbumTable) &&
               jukeboxDB.createTable(createSongTable) &&
               jukeboxDB.createTable(createPlaylistTable) &&
               jukeboxDB.createTable(createPlaylistSongTable)
    }

    return false
}

func (jukeboxDB *JukeboxDB) haveTables() bool {
   haveTablesInDb := false
   if jukeboxDB.dbConnection != nil {
      sqlQuery := "SELECT name " +
                  "FROM sqlite_master " +
                  "WHERE type='table' AND name='song'"
      stmt, err := jukeboxDB.dbConnection.Prepare(sqlQuery)
      if err != nil {
         fmt.Printf("error: unable to prepare sql: %s\n", sqlQuery)
         fmt.Printf("error: %v\n", err)
         return false
      }
      defer stmt.Close()

      var name string
      err = stmt.QueryRow().Scan(&name)
      if err == nil {
         haveTablesInDb = true
      }
   }

   return haveTablesInDb
}

func (jukeboxDB *JukeboxDB) getPlaylist(playlistName string) *string {
    var plObject string
    if len(playlistName) > 0 {
        sqlQuery := "SELECT playlist_uid FROM playlist WHERE playlist_name = ?"
        stmt, err := jukeboxDB.dbConnection.Prepare(sqlQuery)
        if err != nil {
           fmt.Printf("error: unable to prepare query string '%s'\n", sqlQuery)
	   fmt.Printf("error: %v\n", err)
	   return nil
        }
        defer stmt.Close()

        err = stmt.QueryRow(playlistName).Scan(&plObject)
        if err != nil {
           fmt.Printf("error: unable to execute query for statement '%s'\n", sqlQuery)
	   fmt.Printf("error: %v\n", err)
	   return nil
        }
    }
    return &plObject
}

func (jukeboxDB *JukeboxDB) songsForQueryResults(rows *sql.Rows) []*SongMetadata {
    resultSongs := []*SongMetadata{}

    for rows.Next() {
        var fileUid string
        var fileTime string
        var oFileSize int64
        var sFileSize int64
        var padCount int
        var artistName string
        var artistUid string
        var songName string
        var md5Hash string
        var compressed int
        var encrypted int
        var containerName string
        var objectName string
        var albumUid string

        err := rows.Scan(&fileUid, &fileTime, &oFileSize, &sFileSize,
                         &padCount, &artistName, &artistUid, &songName,
                         &md5Hash, &compressed, &encrypted, &containerName,
                         &objectName, &albumUid)

        if err != nil {
           fmt.Printf("error: scan of row values failed\n")
           fmt.Printf("error: %v\n", err)
           return resultSongs
        }

        song := NewSongMetadata()
        song.Fm = NewFileMetadata()
        song.Fm.FileUid = fileUid
        song.Fm.FileTime = fileTime
        song.Fm.OriginFileSize = oFileSize
        song.Fm.StoredFileSize = sFileSize
        song.Fm.PadCharCount = padCount
        song.ArtistName = artistName
        song.ArtistUid = artistUid
        song.SongName = songName
        song.Fm.Md5Hash = md5Hash
        song.Fm.Compressed = compressed == 1
        song.Fm.Encrypted = encrypted == 1
        song.Fm.ContainerName = containerName
        song.Fm.ObjectName = objectName
        song.AlbumUid = albumUid

        resultSongs = append(resultSongs, song)
    }

    return resultSongs
}

func (jukeboxDB *JukeboxDB) retrieveSong(fileName string) *SongMetadata {
    if jukeboxDB.dbConnection != nil {
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
        stmt, err := jukeboxDB.dbConnection.Prepare(sqlQuery)
        if err != nil {
            fmt.Printf("error: unable to prepare statement '%s'\n", sqlQuery)
            fmt.Printf("error: %v\n", err)
            return nil
        }
        defer stmt.Close()

        rows, errRows := stmt.Query(fileName)
        if errRows != nil {
            return nil
        }
        songResults := jukeboxDB.songsForQueryResults(rows)
        if songResults != nil && len(songResults) > 0 {
            return songResults[0]
        }
    }
    return nil
}

func (jukeboxDB *JukeboxDB) insertPlaylist(plUid string,
                                           plName string,
                                           plDesc string) bool {
    insertSuccess := false

    if jukeboxDB.dbConnection != nil &&
       len(plUid) > 0 &&
       len(plName) > 0 {

        tx, errTx := jukeboxDB.dbConnection.Begin()
        if errTx != nil {
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

        stmt.Exec(plUid, plName, plDesc)
        tx.Commit()
        insertSuccess = true
    }

    return insertSuccess
}

func (jukeboxDB *JukeboxDB) deletePlaylist(plName string) bool {
    deleteSuccess := false

    if jukeboxDB.dbConnection != nil && len(plName) > 0 {
        sqlQuery := "DELETE FROM playlist " +
                    "WHERE playlist_name = ? "
        tx, errTx := jukeboxDB.dbConnection.Begin()
        if errTx != nil {
            return false
        }
        stmt, err := tx.Prepare(sqlQuery)
        if err != nil {
            fmt.Printf("error: unable to prepare statement '%s'\n", sqlQuery)
            fmt.Printf("error: %v\n", err)
            return false
        }
        defer stmt.Close()
        stmt.Exec(plName)
        tx.Commit()
        deleteSuccess = true
    }

    return deleteSuccess
}

func (jukeboxDB *JukeboxDB) insertSong(song *SongMetadata) bool {
    insertSuccess := false

    if jukeboxDB.dbConnection != nil && song != nil {
        sqlQuery := "INSERT INTO song VALUES (?,?,?,?,?,?,?,?,?,?,?,?,?,?)"
        tx, errTx := jukeboxDB.dbConnection.Begin()
        if errTx != nil {
            return false
        }
        stmt, err := tx.Prepare(sqlQuery)
        if err != nil {
            fmt.Printf("error: unable to prepare statement '%s'\n", sqlQuery)
            fmt.Printf("error: %v\n", err)
            return false
        }
        defer stmt.Close()

        stmt.Exec(song.Fm.FileUid,
                  song.Fm.FileTime,
                  song.Fm.OriginFileSize,
                  song.Fm.StoredFileSize,
                  song.Fm.PadCharCount,
                  song.ArtistName,
                  "",
                  song.SongName,
                  song.Fm.Md5Hash,
                  song.Fm.Compressed,
                  song.Fm.Encrypted,
                  song.Fm.ContainerName,
                  song.Fm.ObjectName,
                  song.AlbumUid)
        tx.Commit()
        insertSuccess = true
    }

    return insertSuccess
}

func (jukeboxDB *JukeboxDB) updateSong(song *SongMetadata) bool {
    updateSuccess := false

    if jukeboxDB.dbConnection != nil && song != nil && len(song.Fm.FileUid) > 0 {
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
        tx, errTx := jukeboxDB.dbConnection.Begin()
        if errTx != nil {
            return false
        }
        stmt, err := tx.Prepare(sqlQuery)
        if err != nil {
            fmt.Printf("error: unable to prepare statement '%s'\n", sqlQuery)
            fmt.Printf("error: %v\n", err)
            return false
        }

        defer stmt.Close()

        stmt.Exec(song.Fm.FileTime,
                  song.Fm.OriginFileSize,
                  song.Fm.StoredFileSize,
                  song.Fm.PadCharCount,
                  song.ArtistName,
                  "",
                  song.SongName,
                  song.Fm.Md5Hash,
                  song.Fm.Compressed,
                  song.Fm.Encrypted,
                  song.Fm.ContainerName,
                  song.Fm.ObjectName,
                  song.AlbumUid,
                  song.Fm.FileUid)
        tx.Commit()
        updateSuccess = true
    }

    return updateSuccess
}

func (jukeboxDB *JukeboxDB) storeSongMetadata(song *SongMetadata) bool {
    if song == nil {
       return false
    }
    dbSong := jukeboxDB.retrieveSong(song.Fm.FileUid)
    if dbSong != nil {
        if ! song.Equals(dbSong) {
            return jukeboxDB.updateSong(song)
        } else {
            return true  // no insert or update needed (already up-to-date)
        }
    } else {
        // song is not in the database, insert it
        return jukeboxDB.insertSong(song)
    }
}

func (jukeboxDB *JukeboxDB) sqlWhereClause() string {
   var encryption int
   if jukeboxDB.useEncryption {
      encryption = 1
   } else {
      encryption = 0
   }

   var compression int
   if jukeboxDB.useCompression {
      compression = 1
   } else {
      compression = 0
   }

   whereClause := ""
   whereClause += " WHERE "
   whereClause += "encrypted = "
   whereClause += fmt.Sprintf("%d", encryption)
   whereClause += " AND "
   whereClause += "compressed = "
   whereClause += fmt.Sprintf("%d", compression)
   return whereClause
}

func (jukeboxDB *JukeboxDB) retrieveSongs(artist string,
                                          album string) []*SongMetadata {
    var songs []*SongMetadata
    if jukeboxDB.dbConnection != nil {
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

        sqlQuery += jukeboxDB.sqlWhereClause()
        //if len(artist) > 0:
        //    sqlQuery += " AND artist_name='%s'" % artist
	if len(artist) > 0 {
            var addedClause string
            encodedArtist := EncodeValue(artist)
            if len(album) > 0 {
                encodedAlbum := EncodeValue(album)
                addedClause = fmt.Sprintf(" AND object_name LIKE '%s--%s%%'",
                                          encodedArtist,
                                          encodedAlbum)
            } else {
                addedClause = fmt.Sprintf(" AND object_name LIKE '%s--%%'",
                                          encodedArtist)
            }
            sqlQuery += addedClause
        }

        fmt.Printf("executing query: %s\n", sqlQuery)
        stmt, errStmt := jukeboxDB.dbConnection.Prepare(sqlQuery)
        if errStmt != nil {
            fmt.Printf("error: unable to prepare statement '%s'\n", sqlQuery)
            fmt.Printf("error: %v\n", errStmt)
            return nil
        }
        defer stmt.Close()

        rows, err := stmt.Query()
        if err != nil {
            fmt.Printf("error: unable to execute query of '%s'\n", sqlQuery)
            fmt.Printf("error: %v\n", err)
            return nil
        }

        songs = jukeboxDB.songsForQueryResults(rows)
    }
    return songs
}

func (jukeboxDB* JukeboxDB) songsForArtist(artistName string) []*SongMetadata {
    songs := []*SongMetadata{}
    if jukeboxDB.dbConnection != nil {
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
        sqlQuery += jukeboxDB.sqlWhereClause()
        sqlQuery += " AND artist = ?"
        stmt, err := jukeboxDB.dbConnection.Prepare(sqlQuery)
        if err != nil {
           fmt.Printf("error: unable to prepare statement '%s'\n", sqlQuery)
           fmt.Printf("error: %v\n", err)
           return nil
        }
        defer stmt.Close()

        rows, err := stmt.Query(artistName)
        if err != nil {
           fmt.Printf("error: query by artist name failed\n")
           fmt.Printf("error: %v\n", err)
           return nil
        }
        songs = jukeboxDB.songsForQueryResults(rows)
    }
    return songs
}

func (jukeboxDB *JukeboxDB) showListings() {
   if jukeboxDB.dbConnection != nil {
      sqlQuery := "SELECT artist_name, song_name " +
                  "FROM song " +
                  "ORDER BY artist_name, song_name"
      stmt, err := jukeboxDB.dbConnection.Prepare(sqlQuery)
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
            fmt.Printf("error: unable to scan values (artist, song)\n")
            fmt.Printf("error: %v\n", err)
            return
         } else {
            fmt.Printf("%s, %s\n", artist, song)
         }
      }
   }
}

func (jukeboxDB *JukeboxDB) showArtists() {
   if jukeboxDB.dbConnection != nil {
      sqlQuery := "SELECT DISTINCT artist_name " +
                  "FROM song " +
                  "ORDER BY artist_name"
      stmt, err := jukeboxDB.dbConnection.Prepare(sqlQuery)
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
            fmt.Printf("error: unable to scan row value (artist)\n")
            fmt.Printf("error: %v\n", err)
            return
         } else {
            fmt.Printf("%s\n", artist)
         }
      }
   }
}

func (jukeboxDB *JukeboxDB) showGenres() {
   if jukeboxDB.dbConnection != nil {
      sqlQuery := "SELECT genre_name " +
                  "FROM genre " +
                  "ORDER BY genre_name"
      stmt, err := jukeboxDB.dbConnection.Prepare(sqlQuery)
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
         var genreName string
         err = rows.Scan(&genreName)
         if err != nil {
            fmt.Printf("error: unable to Scan row value (genreName)\n")
            fmt.Printf("error: %v\n", err)
            return
         } else {
            fmt.Printf("%s\n", genreName)
         }
      }
   }
}

func (jukeboxDB *JukeboxDB) showAlbums() {
   if jukeboxDB.dbConnection != nil {
      sqlQuery := "SELECT album.album_name, artist.artist_name " +
                  "FROM album, artist " +
                  "WHERE album.artist_uid = artist.artist_uid " +
                  "ORDER BY album.album_name"
      stmt, err := jukeboxDB.dbConnection.Prepare(sqlQuery)
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
         var albumName string
         var artistName string
         err = rows.Scan(&albumName, artistName)
         if err != nil {
         } else {
            fmt.Printf("%s (%s)\n", albumName, artistName)
         }
      }
   }
}

func (jukeboxDB *JukeboxDB) showPlaylists() {
   if jukeboxDB.dbConnection != nil {
      sqlQuery := "SELECT playlist_uid, playlist_name " +
                  "FROM playlist " +
                  "ORDER BY playlist_uid"
      stmt, err := jukeboxDB.dbConnection.Prepare(sqlQuery)
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
         var plUid string
         var plName string
         err = rows.Scan(&plUid, &plName)
         if err != nil {
         } else {
            fmt.Printf("%s - %s\n", plUid, plName)
         }
      }
   }
}

func (jukeboxDB *JukeboxDB) deleteSong(songUid string) bool {
   wasDeleted := false
   if jukeboxDB.dbConnection != nil {
      if len(songUid) > 0 {
         sqlStatement := "DELETE FROM song WHERE song_uid = ?"
         tx, errTx := jukeboxDB.dbConnection.Begin()
         if errTx != nil {
             fmt.Printf("error: begin transaction failed\n")
             fmt.Printf("error: %v\n", errTx)
             return false
         }
         stmt, err := tx.Prepare(sqlStatement)
         if err != nil {
             tx.Rollback()
             fmt.Printf("error: unable to prepare statement '%s'\n", sqlStatement)
             fmt.Printf("error: %v\n", err)
             return false
         }
         defer stmt.Close()

         _, err = stmt.Exec(songUid)
         if err != nil {
             tx.Rollback()
             fmt.Printf("error: unable to delete song '%s'\n", songUid)
             fmt.Printf("error: %v\n", err)
             return false
         }
         tx.Commit()
         wasDeleted = true
      }
   }

   return wasDeleted
}
