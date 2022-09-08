// ******************************************************************************
// Cloud jukebox
// Copyright Paul Dardeau, SwampBits LLC, 2014
// BSD license -- see LICENSE file for details
//
// (1) create a directory for the jukebox (e.g., ~/jukebox)
//
// This cloud jukebox uses an abstract object storage system.
// (2) copy this source file to $JUKEBOX
// (3) create subdirectory for song imports (e.g., mkdir $JUKEBOX/song-import)
// (4) create subdirectory for song-play (e.g., mkdir $JUKEBOX/song-play)
//
// Song file naming convention:
//
// The-Artist-Name--Album-Name--The-Song-Name.ext
//       |         |       |           |       |
//       |         |       |           |       |----  file extension (e.g., 'mp3')
//       |         |       |           |
//       |         |       |           |---- name of the song (' ' replaced with '-')
//       |         |       |
//       |         |       |---- name of the album (' ' replaced with '-')
//       |         |
//       |         |---- double dashes to separate the artist name and song name
//       |
//       |---- artist name (' ' replaced with '-')
//
// For example, the MP3 version of the song 'Under My Thumb' from artist 'The
// Rolling Stones' from the album 'Aftermath' should be named:
//
//   The-Rolling-Stones--Aftermath--Under-My-Thumb.mp3
//
// first time use (or when new songs are added):
// (1) copy one or more song files to $JUKEBOX/song-import
// (2) import songs with command: 'python jukebox_main.py import-songs'
//
// show song listings:
// python jukebox_main.py list-songs
//
// play songs:
// python jukebox_main.py play
//
// ******************************************************************************
package jukebox

import (
   "fmt"
   "io/ioutil"
   "os"
   "path/filepath"
   "strings"
   "time"
)



var g_jukebox_instance *Jukebox


//def signal_handler(signum: int, frame):
//    if signum == signal.SIGUSR1 {
//        if g_jukebox_instance != nil {
//            g_jukebox_instance.toggle_pause_play()
//        }
//} else if signum == signal.SIGUSR2 {
//        if g_jukebox_instance != nil {
//            g_jukebox_instance.advance_to_next_song()
//}
//}



//def install_signal_handlers():
//    if os.name == 'posix':
//        signal.signal(signal.SIGUSR1, signal_handler)
//        signal.signal(signal.SIGUSR2, signal_handler)

type Jukebox struct {
   jukebox_options *JukeboxOptions
   storage_system StorageSystem
   debug_print bool
   jukebox_db *JukeboxDB
   current_dir string
   song_import_dir string
   playlist_import_dir string
   song_play_dir string
   album_art_import_dir string
   download_extension string
   metadata_db_file string
   metadata_container string
   playlist_container string
   album_art_container string
   song_list []*SongMetadata
   number_songs int
   song_index int
   //audio_player_command_args = []
   //audio_player_popen = nil
   song_play_length_seconds int
   cumulative_download_bytes int 
   cumulative_download_time int 
   exit_requested bool 
   is_paused bool
   song_start_time int
   song_seconds_offset int 
}

func NewJukebox(jb_options *JukeboxOptions,
                storage_sys StorageSystem,
                debug_print bool) *Jukebox {
   var jukebox Jukebox
   g_jukebox_instance = &jukebox
   jukebox.jukebox_options = jb_options
   jukebox.storage_system = storage_sys
   jukebox.debug_print = debug_print
   jukebox.jukebox_db = nil
   cwd, err := os.Getwd()
   if err == nil {
       jukebox.current_dir = cwd
   }
   jukebox.song_import_dir = Path_join(jukebox.current_dir, "song-import")
   jukebox.playlist_import_dir = Path_join(jukebox.current_dir, "playlist-import")
   jukebox.song_play_dir = Path_join(jukebox.current_dir, "song-play")
   jukebox.album_art_import_dir = Path_join(jukebox.current_dir, "album-art-import")
   jukebox.download_extension = ".download"
   jukebox.metadata_db_file = "jukebox_db.sqlite3"
   jukebox.metadata_container = "music-metadata"
   jukebox.playlist_container = "playlists"
   jukebox.album_art_container = "album-art"
   jukebox.song_list = []*SongMetadata{}
   jukebox.number_songs = 0
   jukebox.song_index = -1
   //jukebox.audio_player_command_args = []
   //jukebox.audio_player_popen = nil
   jukebox.song_play_length_seconds = 20
   jukebox.cumulative_download_bytes = 0
   jukebox.cumulative_download_time = 0
   jukebox.exit_requested = false
   jukebox.is_paused = false
   jukebox.song_start_time = 0
   jukebox.song_seconds_offset = 0

   if jukebox.jukebox_options != nil && jukebox.jukebox_options.Debug_mode {
      jukebox.debug_print = true
   }
   if jukebox.debug_print {
      fmt.Printf("current_dir = '%s'\n", jukebox.current_dir)
      fmt.Printf("song_import_dir = '%s'\n", jukebox.song_import_dir)
      fmt.Printf("song_play_dir = '%s'\n", jukebox.song_play_dir)
   }
   return &jukebox
}

func (jukebox *Jukebox) Enter() bool {
   // look for stored metadata in the storage system
   if jukebox.storage_system != nil &&
      jukebox.storage_system.Has_container(jukebox.metadata_container) &&
      !jukebox.jukebox_options.Suppress_metadata_download {

      // metadata container exists, retrieve container listing
      container_contents := jukebox.storage_system.List_container_contents(jukebox.metadata_container)

      metadata_file_in_container := false
      for _, container_item := range container_contents {
         if container_item == jukebox.metadata_db_file {
            metadata_file_in_container = true
	    break
         }
      }

      // does our metadata DB file exist in the metadata container?
      if container_contents != nil && metadata_file_in_container {
          // download it
	  metadata_db_file_path := jukebox.Get_metadata_db_file_path()
	  download_file := metadata_db_file_path + ".download"
          if jukebox.storage_system.Get_object(jukebox.metadata_container, jukebox.metadata_db_file, download_file) > 0 {
              // have an existing metadata DB file?
              if File_exists(metadata_db_file_path) {
                  if jukebox.debug_print {
                      fmt.Println("deleting existing metadata DB file")
                  }
                  Delete_file(metadata_db_file_path)
                  // rename downloaded file
                  if jukebox.debug_print {
                      fmt.Printf("renaming '%s' to '%s'\n", download_file, metadata_db_file_path)
                  }
                  os.Rename(download_file, metadata_db_file_path)
                } else {
                    if jukebox.debug_print {
                        fmt.Println("error: unable to retrieve metadata DB file")
                    }
                }
            } else {
                if jukebox.debug_print {
                    fmt.Println("no metadata DB file in metadata container")
                }
            }
        } else {
            if jukebox.debug_print {
                fmt.Println("no metadata container in storage system")
            }
        }

	debug_print := true
        jukebox.jukebox_db = NewJukeboxDB(jukebox.Get_metadata_db_file_path(),
	                                  debug_print,
                                          jukebox.jukebox_options.Use_encryption,
                                          jukebox.jukebox_options.Use_compression)
        return jukebox.jukebox_db.enter()
	/*
        if !jukebox.jukebox_db.open() {
            fmt.Println("unable to connect to database")
        }
	*/
    }

    return false
}

func (jukebox *Jukebox) Exit() {
   if jukebox.jukebox_db != nil {
       jukebox.jukebox_db.exit()
       jukebox.jukebox_db = nil
   }
}

func (jukebox *Jukebox) Toggle_pause_play() {
   jukebox.is_paused = ! jukebox.is_paused
   if jukebox.is_paused {
      fmt.Println("paused")
      //if jukebox.audio_player_popen != nil {
         // capture current song position (seconds into song)
      //   jukebox.audio_player_popen.terminate()
      //}
   } else {
      fmt.Println("resuming play")
   }
}

func (jukebox *Jukebox) Advance_to_next_song() {
   fmt.Println("advancing to next song")
   //if jukebox.audio_player_popen != nil {
   //   jukebox.audio_player_popen.terminate()
   //}
}

func (jukebox *Jukebox) Get_metadata_db_file_path() (string) {
   return Path_join(jukebox.current_dir, jukebox.metadata_db_file)
}

func unencode_value(encoded_value string) (string) {
   return strings.Replace(encoded_value, "-", " ", -1)
}

func encode_value(value string) (string) {
   return strings.Replace(value, " ", "-", -1)
}

func components_from_file_name(file_name string) (string,string,string) {
   if len(file_name) == 0 {
      return "", "", ""
   }
   pos_extension := strings.Index(file_name, ".")
   var base_file_name string
   if pos_extension > -1 {
      base_file_name = file_name[0:pos_extension]
   } else {
      base_file_name = file_name
   }
   components := strings.Split(base_file_name, "--")
   if len(components) == 3 {
      return unencode_value(components[0]),
             unencode_value(components[1]),
             unencode_value(components[2])
   } else {
      return "", "", ""
   }
}

func artist_from_file_name(file_name string) (*string) {
   if len(file_name) > 0 {
      artist, _, _ := components_from_file_name(file_name)
      if len(artist) > 0 {
         return &artist
      }
   }
   return nil
}

func album_from_file_name(file_name string) (*string) {
   if len(file_name) > 0 {
      _, album, _ := components_from_file_name(file_name)
      if len(album) > 0 {
         return &album
      }
   }
   return nil
}

func song_from_file_name(file_name string) (*string) {
   if len(file_name) > 0 {
      _, _, song := components_from_file_name(file_name)
      if len(song) > 0 {
         return &song
      }
   }
   return nil
}

func (jukebox *Jukebox) store_song_metadata(fs_song SongMetadata) (bool) {
   db_song := jukebox.jukebox_db.retrieve_song(fs_song.Fm.File_uid)
   if db_song != nil {
      p_fs_song := &fs_song
      if ! p_fs_song.Equals(db_song) {
         return jukebox.jukebox_db.update_song(fs_song)
      } else {
         return true  // no insert or update needed (already up-to-date)
      }
   } else {
      // song is not in the database, insert it
      return jukebox.jukebox_db.insert_song(fs_song)
   }
}

func (jukebox *Jukebox) store_song_playlist(file_name string, file_contents []byte) bool {
/*
        pl = json.loads(file_contents)
        if "name" in pl.keys() {
            pl_name = pl["name"]
            pl_uid = file_name
            return jukebox.jukebox_db.insert_playlist(pl_uid, pl_name)
        } else {
            return false
        }
*/
    return false
}

/*
func (jukebox *Jukebox) get_encryptor() {
    // key_block_size = 16  // AES-128
    // key_block_size = 24  // AES-192
    key_block_size = 32  // AES-256
    return AESBlockEncryption(key_block_size,
                              jukebox.jukebox_options.encryption_key,
                              jukebox.jukebox_options.encryption_iv)
}
*/

func (jukebox *Jukebox) container_suffix() (string) {
   suffix := ""
   if jukebox.jukebox_options.Use_encryption &&
      jukebox.jukebox_options.Use_compression {
      suffix += "-ez"
   } else if jukebox.jukebox_options.Use_encryption {
      suffix += "-e"
   } else if jukebox.jukebox_options.Use_compression {
      suffix += "-z"
   }
   return suffix
}

func (jukebox *Jukebox) object_file_suffix() (string) {
   suffix := ""
   if jukebox.jukebox_options.Use_encryption &&
      jukebox.jukebox_options.Use_compression {
      suffix = ".egz"
   } else if jukebox.jukebox_options.Use_encryption {
      suffix = ".e"
   } else if jukebox.jukebox_options.Use_compression {
      suffix = ".gz"
   }
   return suffix
}

func (jukebox *Jukebox) container_for_song(song_uid string) *string {
   if len(song_uid) == 0 {
      return nil
   }
   container_suffix := "-artist-songs" + jukebox.container_suffix()

   artist := artist_from_file_name(song_uid)
   if artist == nil || len(*artist) == 0 {
      return nil
   }

   var artist_letter string
   artist_value := *artist

   if strings.HasPrefix(artist_value, "A ") {
      artist_letter = artist_value[2:3]
   } else if strings.HasPrefix(artist_value, "The ") {
      artist_letter = artist_value[4:5]
   } else {
      artist_letter = artist_value[0:1]
   }

   container_name := strings.ToLower(artist_letter) + container_suffix
   return &container_name
}

func (jukebox *Jukebox) Import_songs() {
   if jukebox.jukebox_db != nil && jukebox.jukebox_db.is_open() {
      /*
      dir_listing, err := ioutil.ReadDir(jukebox.song_import_dir)
      if err != nil {
         return
      }
      num_entries = float(len(dir_listing))
      progressbar_chars = 0.0
      progressbar_width = 40
      progress_chars_per_iteration = progressbar_width / num_entries
      progressbar_char = '#'
      bar_chars = 0

      if ! jukebox.debug_print {
         // setup progressbar
         sys.stdout.write("[%s]" % (" " * progressbar_width))
         sys.stdout.flush()
         sys.stdout.write("\b" * (progressbar_width + 1))  // return to start of line, after '['
      }

      //if jukebox.jukebox_options != nil && jukebox.jukebox_options.use_encryption {
      //   encryption = jukebox.get_encryptor()
      //} else {
      //   encryption = nil
      //}

      cumulative_upload_time = 0
      cumulative_upload_bytes = 0
      file_import_count = 0

      for _, listing_entry := range dir_listing {
         full_path = Path_join(jukebox.song_import_dir, listing_entry)
         // ignore it if it's not a file
         if os.path.isfile(full_path) {
            file_name = listing_entry
            extension = os.path.splitext(full_path)[1]
            if extension {
               file_size = Get_file_size(full_path)
               artist = jukebox.artist_from_file_name(file_name)
               album = jukebox.album_from_file_name(file_name)
               song = jukebox.song_from_file_name(file_name)
               if file_size > 0 && artist != nil && album != nil && song != nil {
                  object_name = file_name + jukebox.object_file_suffix()
                  fs_song = NewSongMetadata()
                  fs_song.Fm = NewFileMetadata()
                  fs_song.Fm.File_uid = object_name
                  fs_song.Album_uid = nil
                  fs_song.Fm.Origin_file_size = file_size
                  fs_song.Fm.File_time = datetime.datetime.fromtimestamp(os.path.getmtime(full_path))
                  fs_song.Artist_name = artist
                  fs_song.Song_name = song
                  fs_song.Fm.Md5_hash = Md5_for_file(full_path)
                  fs_song.Fm.Compressed = jukebox.jukebox_options.use_compression
                  fs_song.Fm.Encrypted = jukebox.jukebox_options.use_encryption
                  fs_song.Fm.Object_name = object_name
                  fs_song.Fm.Pad_char_count = 0

                  fs_song.Fm.Container_name = jukebox.container_for_song(file_name)

                  // read file contents
		  file_read := false
		  file_contents := nil

                  //try:
                  //   with open(full_path, 'rb') as content_file:
                  //      file_contents = content_file.read()
                  //   file_read = true
                  //except IOError:
                  //   fmt.Printf("error: unable to read file %s\n", full_path)

                  if file_read && file_contents != nil {
                     if file_contents {
                        // for general purposes, it might be useful or helpful to have
                        // a minimum size for compressing
			//TODO: add support for compression and encryption
                        if jukebox.jukebox_options.use_compression {
                           if jukebox.debug_print {
                              fmt.Println("compressing file")
                           }

                           file_bytes = bytes(file_contents, 'utf-8')
                           file_contents = zlib.compress(file_bytes, 9)

                           if jukebox.jukebox_options.use_encryption {
                              if jukebox.debug_print {
                                 fmt.Println("encrypting file")
                              }

                              // the length of the data to encrypt must be a multiple of 16
                              num_extra_chars = len(file_contents) % 16
                              if num_extra_chars > 0 {
                                 if jukebox.debug_print {
                                    fmt.Println("padding file for encryption")
                                 }
                                 num_pad_chars = 16 - num_extra_chars
                                 file_contents += "".ljust(num_pad_chars, ' ')
                                 fs_song.Fm.Pad_char_count = num_pad_chars
                              }

                              file_contents = encryption.encrypt(file_contents)
                           }

                           // now that we have the data that will be stored, set the file size for
                           // what's being stored
                           fs_song.Fm.Stored_file_size = len(file_contents)
                           start_upload_time = time.time()

                           // store song file to storage system
                           if jukebox.storage_system.Put_object(fs_song.Fm.Container_name,
                                                                fs_song.Fm.Object_name,
                                                                file_contents) {
                              end_upload_time = time.time()
                              upload_elapsed_time = end_upload_time - start_upload_time
                              cumulative_upload_time += upload_elapsed_time
                              cumulative_upload_bytes += len(file_contents)

                              // store song metadata in local database
                              if ! jukebox.store_song_metadata(fs_song) {
                                 // we stored the song to the storage system, but were unable to store
                                 // the metadata in the local database. we need to delete the song
                                 // from the storage system since we won't have any way to access it
                                 // since we can't store the song metadata locally.
                                 fmt.Printf("unable to store metadata, deleting obj '%s'", fs_song.Fm.Object_name)
                                 jukebox.storage_system.Delete_object(fs_song.Fm.Container_name,
                                                                      fs_song.Fm.Object_name)
                              } else {
                                 file_import_count += 1
                              }
                           } else {
                              fmt.Printf("error: unable to upload '%s' to '%s'\n",
                                         fs_song.Fm.Object_name,
                                         fs_song.Fm.Container_name)
                           }

                if ! jukebox.debug_print {
                    progressbar_chars += progress_chars_per_iteration
                    if int(progressbar_chars) > bar_chars {
                        num_new_chars = int(progressbar_chars) - bar_chars
                        if num_new_chars > 0 {
                            // update progress bar
                            for j in iter(range(num_new_chars)) {
                                sys.stdout.write(progressbar_char)
                            }
                            sys.stdout.flush()
                            bar_chars += num_new_chars
                        }
                    }
                }

            if ! jukebox.debug_print {
                // if we haven't filled up the progress bar, fill it now
                if bar_chars < progressbar_width {
                    num_new_chars = progressbar_width - bar_chars
                    for j in iter(range(num_new_chars)) {
                        sys.stdout.write(progressbar_char)
                    }
                    sys.stdout.flush()
                }
                sys.stdout.write("\n")
            }

            if file_import_count > 0 {
                jukebox.Upload_metadata_db()
            }

            fmt.Printf("%d song files imported\n", file_import_count)

            if cumulative_upload_time > 0 {
                cumulative_upload_kb := cumulative_upload_bytes / 1000.0
                fmt.Printf("average upload throughput = %d KB/sec\n",
                           cumulative_upload_kb / cumulative_upload_time)
            }
	    */
   }
}

func (jukebox *Jukebox) song_path_in_playlist(song *SongMetadata) string {
    return Path_join(jukebox.song_play_dir, song.Fm.File_uid)
}

func (jukebox *Jukebox) check_file_integrity(song *SongMetadata) bool {
    file_integrity_passed := true

    if jukebox.jukebox_options != nil && jukebox.jukebox_options.Check_data_integrity {
	file_path := jukebox.song_path_in_playlist(song)
        if File_exists(file_path) {
            if jukebox.debug_print {
                fmt.Printf("checking integrity for %s\n", song.Fm.File_uid)
            }

            if song.Fm != nil {
		playlist_md5, err := Md5_for_file(file_path)
		if err != nil {
                    // log error
		} else {
                    if playlist_md5 == song.Fm.Md5_hash {
                        if jukebox.debug_print {
                            fmt.Println("integrity check SUCCESS")
                        }
                        file_integrity_passed = true
                    } else {
                        fmt.Printf("file integrity check failed: %s\n", song.Fm.File_uid)
                        file_integrity_passed = false
                    }
                }
            }
        } else {
            // file doesn't exist
            fmt.Println("file doesn't exist")
            file_integrity_passed = false
        }
    } else {
        if jukebox.debug_print {
            fmt.Println("file integrity bypassed, no jukebox options or check integrity not turned on")
        }
    }

    return file_integrity_passed
}

func (jukebox *Jukebox) batch_download_start() {
   jukebox.cumulative_download_bytes = 0
   jukebox.cumulative_download_time = 0
}

func (jukebox *Jukebox) batch_download_complete() {
   if ! jukebox.exit_requested {
      if jukebox.cumulative_download_time > 0 {
         cumulative_download_kb := jukebox.cumulative_download_bytes / 1000.0
         fmt.Printf("average download throughput = %d KB/sec\n",
                    cumulative_download_kb / jukebox.cumulative_download_time)
      }
      jukebox.cumulative_download_bytes = 0
      jukebox.cumulative_download_time = 0
   }
}

func (jukebox *Jukebox) download_song(song *SongMetadata) (bool) {
   if jukebox.exit_requested {
      return false
   }

   if song != nil {
      file_path := jukebox.song_path_in_playlist(song)
      //download_start_time := time.time()
      song_bytes_retrieved := jukebox.storage_system.Retrieve_file(song.Fm, jukebox.song_play_dir)
      if jukebox.exit_requested {
         return false
      }

      if jukebox.debug_print {
         fmt.Println("bytes retrieved: %d\n", song_bytes_retrieved)
      }

      if song_bytes_retrieved > 0 {
         //download_end_time := time.time()
	 //download_elapsed_time := download_end_time - download_start_time
         //jukebox.cumulative_download_time += download_elapsed_time
         jukebox.cumulative_download_bytes += song_bytes_retrieved

         // are we checking data integrity?
         // if so, verify that the storage system retrieved the same length that has been stored
         if jukebox.jukebox_options != nil && jukebox.jukebox_options.Check_data_integrity {
            if jukebox.debug_print {
               fmt.Println("verifying data integrity")
            }

            if song_bytes_retrieved != song.Fm.Stored_file_size {
               fmt.Printf("error: data integrity check failed for '%s'\n", file_path)
               return false
            }
         }

	 //TODO: add support for encryption and compression
	 /*
         // is it encrypted? if so, unencrypt it
         encrypted = song.Fm.Encrypted
         compressed = song.Fm.Compressed

         if encrypted || compressed {
            //try {
               with open(file_path, 'rb') as content_file {
                  file_contents = content_file.read()
            //  }
            //except IOError {
            //   fmt.Printf("error: unable to read file %s\n", file_path)
            //   return false
            //}

            if encrypted {
               encryption = jukebox.get_encryptor()
               file_contents = encryption.decrypt(file_contents)
            }

            if compressed {
               file_contents = zlib.decompress(file_contents)
            }

            // re-write out the uncompressed, unencrypted file contents
            //try {
               with open(file_path, 'wb') as content_file:
                  content_file.write(file_contents)
	    //} except IOError {
            //   fmt.Printf("error: unable to write unencrypted/uncompressed file '%s'\n", file_path)
            //   return false
            //}
         }
	 */

         if jukebox.check_file_integrity(song) {
            return true
	 } else {
            // we retrieved the file, but it failed our integrity check
            // if file exists, remove it
            if File_exists(file_path) {
               Delete_file(file_path)
            }
         }
      }
   }

   return false
}

func (jukebox *Jukebox) play_song(song_file_path string) {
   if File_exists(song_file_path) {
      fmt.Printf("playing %s\n", song_file_path)
      /*
      if jukebox.audio_player_command_args {
         cmd_args = jukebox.audio_player_command_args[:]
         cmd_args.append(song_file_path)
         exit_code = -1
         started_audio_player = false
         //try {
            audio_player_proc = Popen(cmd_args)
            if audio_player_proc != nil {
               started_audio_player = true
               jukebox.song_start_time = time.time()
               jukebox.audio_player_popen = audio_player_proc
               exit_code = audio_player_proc.wait()
               jukebox.audio_player_popen = nil
            }
         //} except OSError {
         //   // audio player not available
         //   jukebox.audio_player_command_args = []
         //   jukebox.audio_player_popen = nil
         //   exit_code = -1
         //}

         // if the audio player failed or is not present, just sleep
         // for the length of time that audio would be played
         if ! started_audio_player && exit_code != 0 {
            time.sleep(jukebox.song_play_length_seconds)
         }
      } else {
         // we don't know about an audio player, so simulate a
         // song being played by sleeping
         time.sleep(jukebox.song_play_length_seconds)
      }

      if ! jukebox.is_paused {
         // delete the song file from the play list directory
         Delete_file(song_file_path)
      }
      */
   } else {
      fmt.Printf("song file doesn't exist: '%s'\n", song_file_path)

      f, err := os.OpenFile("404.txt",
                            os.O_APPEND|os.O_CREATE|os.O_WRONLY,
                            0644)
      if err != nil {
          fmt.Println("error: unable to open 404.txt to append song file")
          fmt.Println(err)
	  return
      }
      defer f.Close()
      if _, err := f.WriteString(song_file_path + "\n"); err != nil {
          fmt.Println("error: unable to write to 404.txt")
          fmt.Println(err)
      }
   }
}

func (jukebox *Jukebox) download_songs() {
   // scan the play list directory to see if we need to download more songs
   dir_listing, err := ioutil.ReadDir(jukebox.song_play_dir)
   if err != nil {
      // log error
      return
   }
   song_file_count := 0
   for _, listing_entry := range dir_listing {
      if listing_entry.IsDir() {
          continue
      }
      full_path := Path_join(jukebox.song_play_dir, listing_entry.Name())
      extension := filepath.Ext(full_path)
      if len(extension) > 0 && extension != jukebox.download_extension {
          song_file_count += 1
      }
   }

   file_cache_count := jukebox.jukebox_options.File_cache_count

   if song_file_count < file_cache_count {
      //TODO: implement
      /*
      dl_songs := []SongMetadata{}
      // start looking at the next song in the list
      check_index := jukebox.song_index + 1
      for j in iter(range(jukebox.number_songs)) {
         if check_index >= jukebox.number_songs {
            check_index = 0
         }
         if check_index != jukebox.song_index {
            si = jukebox.song_list[check_index]
            file_path = jukebox.song_path_in_playlist(si)
            if ! File_exists(file_path) {
               dl_songs.append(si)
               if len(dl_songs) >= file_cache_count {
                  break
               }
            }
         }
         check_index += 1
      }
      */
   }

   /*
   if dl_songs {
      download_thread = NewSongDownloader(jukebox, dl_songs)
      download_thread.start()
   }
   */
}

func (jukebox *Jukebox) Play_songs(shuffle bool, artist string, album string) {
    song_list := jukebox.jukebox_db.retrieve_songs(artist, album)
    jukebox.play_song_list(song_list, shuffle)
}

func (jukebox *Jukebox) play_song_list(song_list []*SongMetadata, shuffle bool) {
    jukebox.song_list = song_list
    if jukebox.song_list != nil {
        jukebox.number_songs = len(jukebox.song_list)

        if jukebox.number_songs == 0 {
            fmt.Println("no songs in jukebox")
            return
        }

        // does play list directory exist?
        if ! File_exists(jukebox.song_play_dir) {
            if jukebox.debug_print {
                fmt.Println("song-play directory does not exist, creating it")
            }
            os.Mkdir(jukebox.song_play_dir, os.ModePerm)
        } else {
            // play list directory exists, delete any files in it
            if jukebox.debug_print {
                fmt.Println("deleting existing files in song-play directory")
            }
            dir_files, err_dir := ioutil.ReadDir(jukebox.song_play_dir)
	    if err_dir != nil {
                // log error
            } else {
	        for _, theFile := range dir_files {
                    if theFile.IsDir() {
                        continue
                    }
                    file_path := Path_join(jukebox.song_play_dir, theFile.Name())
                    //try:
                    Delete_file(file_path)
                    //except OSError:
                    //    pass
                }
            }
        }

        jukebox.song_index = 0
        //install_signal_handlers()

	//TODO: add media players
	/*
        if sys.platform == "darwin" {
            jukebox.audio_player_command_args = ["afplay"]
            // jukebox.audio_player_command_args.extend(["-t", str(jukebox.song_play_length_seconds)])
        } else if os.name == "posix" {
            jukebox.audio_player_command_args = ["mplayer", "-novideo", "-nolirc", "-really-quiet"]
            // jukebox.audio_player_command_args.extend(["-endpos", str(jukebox.song_play_length_seconds)])
        } else if sys.platform == "win32" {
            // we really need command-line support for /play and /close arguments. unfortunately,
            // this support used to be available in the built-in Windows Media Player, but is
            // no longer present.
            // jukebox.audio_player_command_args = ["C:\Program Files\Windows Media Player\wmplayer.exe"]
            jukebox.audio_player_command_args = ["C:\\Program Files\\MPC-HC\\mpc-hc64.exe",
                                              "/play", "/close", "/minimized"]
        } else {
            jukebox.audio_player_command_args = []
        }
	*/

        fmt.Println("downloading first song...")

        if shuffle {
            //TODO: add shuffling of song list
            //jukebox.song_list = random.sample(jukebox.song_list, len(jukebox.song_list))
        }

        if jukebox.download_song(jukebox.song_list[0]) {
            fmt.Println("first song downloaded. starting playing now.")
            //with open("jukebox.pid", "w") as f:
            //    f.write('%d\n' % os.getpid())

            for true {
                if ! jukebox.exit_requested {
                    if ! jukebox.is_paused {
                        jukebox.download_songs()
                        jukebox.play_song(jukebox.song_path_in_playlist(jukebox.song_list[jukebox.song_index]))
                    }
                    if ! jukebox.is_paused {
                        jukebox.song_index += 1
                        if jukebox.song_index >= jukebox.number_songs {
                            jukebox.song_index = 0
                        }
                    } else {
                        time.Sleep(1 * time.Second)
                    }
                }
            }
            Delete_file("jukebox.pid")
        } else {
            fmt.Println("error: unable to download songs")
            os.Exit(1)
        }
    }
}


func (jukebox *Jukebox) Show_list_containers() {
   if jukebox.storage_system != nil {
      for _, container_name := range jukebox.storage_system.Get_container_names() {
         fmt.Println(container_name)
      }
   }
}

func (jukebox *Jukebox) Show_listings() {
   if jukebox.jukebox_db != nil {
      jukebox.jukebox_db.show_listings()
   }
}

func (jukebox *Jukebox) Show_artists() {
   if jukebox.jukebox_db != nil {
      jukebox.jukebox_db.show_artists()
   }
}

func (jukebox *Jukebox) Show_genres() {
   if jukebox.jukebox_db != nil {
      jukebox.jukebox_db.show_genres()
   }
}

func (jukebox *Jukebox) Show_albums() {
   if jukebox.jukebox_db != nil {
      jukebox.jukebox_db.show_albums()
   }
}

func (jukebox *Jukebox) read_file_contents(file_path string,
                                           allow_encryption bool) (bool, []byte, int) {
    file_read := false
    pad_chars := 0
    var file_contents []byte

    //try:
    //TODO: implement read_file_contents
    //    with open(file_path, 'r') as content_file:
    //        file_contents = content_file.read()
    //        file_read = true
    //except IOError:
    //    fmt.Printf("error: unable to read file %s\n", file_path)

    if file_read && file_contents != nil {
        if len(file_contents) > 0 {
            // for general purposes, it might be useful or helpful to have
            // a minimum size for compressing
	    //TODO: add support for compression
	    /*
            if jukebox.jukebox_options.use_compression {
                if jukebox.debug_print {
                    fmt.Println("compressing file")
                }

                file_bytes = bytes(file_contents, 'utf-8')
                file_contents = zlib.compress(file_bytes, 9)
            }
	    */

	    //TODO: add support for encryption
	    /*
            if allow_encryption && jukebox.jukebox_options.use_encryption {
                if jukebox.debug_print {
                    fmt.Println("encrypting file")
                }

                // the length of the data to encrypt must be a multiple of 16
                num_extra_chars = len(file_contents) % 16
                if num_extra_chars > 0 {
                    if jukebox.debug_print {
                        fmt.Println("padding file for encryption")
                    }
                    pad_chars = 16 - num_extra_chars
                    file_contents += "".ljust(pad_chars, ' ')
                }

                file_contents = encryption.encrypt(file_contents)
            }
	    */
        }
    }

    return file_read, file_contents, pad_chars
}

func (jukebox *Jukebox) Upload_metadata_db() bool {
    metadata_db_upload := false
    have_metadata_container := false
    if ! jukebox.storage_system.Has_container(jukebox.metadata_container) {
        have_metadata_container = jukebox.storage_system.Create_container(jukebox.metadata_container)
    } else {
        have_metadata_container = true
    }

    if have_metadata_container {
        if jukebox.debug_print {
            fmt.Println("uploading metadata db file to storage system")
        }

        jukebox.jukebox_db.close()
        jukebox.jukebox_db = nil

	//TODO: read metadata DB file and upload it
	/*
        db_file_contents = ''
        with open(jukebox.Get_metadata_db_file_path(), 'rb') as db_file:
            db_file_contents = db_file.read()

        metadata_db_upload = jukebox.storage_system.Put_object(jukebox.metadata_container,
                                                               jukebox.metadata_db_file,
                                                               db_file_contents)
        */

        if jukebox.debug_print {
            if metadata_db_upload {
                fmt.Println("metadata db file uploaded")
            } else {
                fmt.Println("unable to upload metadata db file")
            }
        }
    }

    return metadata_db_upload
}

func (jukebox *Jukebox) Import_playlists() {
   if jukebox.jukebox_db != nil && jukebox.jukebox_db.is_open() {
      file_import_count := 0
      dir_listing, err := ioutil.ReadDir(jukebox.playlist_import_dir)
      if err != nil {
         return
      }
      if len(dir_listing) == 0 {
         fmt.Println("no playlists found")
         return
      }

      have_container := false
      if ! jukebox.storage_system.Has_container(jukebox.playlist_container) {
         have_container = jukebox.storage_system.Create_container(jukebox.playlist_container)
      } else {
         have_container = true
      }

      if ! have_container {
         fmt.Println("error: unable to create container for playlists. unable to import")
         return
      }

      headers := make(map[string]string)

      for _, listing_entry := range dir_listing {
         if listing_entry.IsDir() {
            continue
         }

	 full_path := Path_join(jukebox.playlist_import_dir, listing_entry.Name())
         // ignore it if it's not a file
	 object_name := listing_entry.Name()
	 file_read, file_contents, _ := jukebox.read_file_contents(full_path, false)
         if file_read && file_contents != nil {
            if jukebox.storage_system.Put_object(jukebox.playlist_container,
                                                 object_name,
                                                 file_contents,
                                                 headers) {
               fmt.Println("put of playlist succeeded")
               if ! jukebox.store_song_playlist(object_name, file_contents) {
                  fmt.Println("storing of playlist to db failed")
                  jukebox.storage_system.Delete_object(jukebox.playlist_container,
                                                       object_name)
               } else {
                  fmt.Println("storing of playlist succeeded")
                  file_import_count += 1
               }
            }
         }
      }

      if file_import_count > 0 {
         fmt.Printf("%d playlists imported\n", file_import_count)
         // upload metadata DB file
         jukebox.Upload_metadata_db()
      } else {
         fmt.Println("no files imported")
      }
   }
}

func (jukebox *Jukebox) Show_playlists() {
   if jukebox.jukebox_db != nil {
      jukebox.jukebox_db.show_playlists()
   }
}

func (jukebox *Jukebox) Show_playlist(playlist string) {
   bucket_name := "cj-playlists"
   object_name := fmt.Sprintf("%s.json", encode_value(playlist))
   download_file := object_name
   if jukebox.storage_system.Get_object(bucket_name,
                                        object_name,
                                        download_file) > 0 {
      //file_read := false
      //try {
      //TODO: read playlist file
      //   with open(download_file, 'rb') as content_file {
      //      file_contents = content_file.read()
      //   }
      //   file_read = true
      //} except IOError {
      //   fmt.Printf("error: unable to read file %s\n", full_path)
      //   file_read = false
      //}

      /*
      if file_read {
         pl = json.loads(file_contents)
         if pl != nil {
            if "songs" in pl {
               song_list = []
               list_song_dicts = pl["songs"]
               for song_dict in list_song_dicts {
                  artist_name = song_dict["artist"]
                  if "'" in artist_name {
                     artist_name = strings.Replace(artist_name, "'", "", -1)
                  }
                  artist = Jukebox.encode_value(artist_name)
                  album_name = song_dict["album"]
                  if "'" in album_name {
                     album_name = strings.Replace(album_name, "'", "", -1)
                  }
                  album = Jukebox.encode_value(album_name)
                  song_name = song_dict["song"]
                  if "'" in song_name {
                     song_name = strings.Replace(song_name, "'", "", -1)
                  }
                  song = Jukebox.encode_value(song_name)
                  base_object_name = "%s--%s--%s" % (artist, album, song)
                  fmt.Println(base_object_name)
               }
            }
         }
      }
      */
   } else {
      fmt.Printf("error: unable to retrieve %s\n", object_name)
   }
}

func (jukebox *Jukebox) Play_playlist(playlist string) {
   bucket_name := "cj-playlists"
   object_name := fmt.Sprintf("%s.json", encode_value(playlist))
   download_file := object_name

   if jukebox.storage_system.Get_object(bucket_name,
                                        object_name,
                                        download_file) > 0 {
      //TODO: implement play_playlist
      /*
      with open(download_file, 'rb') as content_file:
         file_contents = content_file.read()
      file_read = true
      //fmt.Printf("error: unable to read file %s\n", full_path)
      //file_read = false

      if file_read {
          pl = json.loads(file_contents)
          if pl != nil {
              if "songs" in pl {
                  song_list = []
                  list_song_dicts = pl["songs"]
                  for song_dict in list_song_dicts {
                      artist_name = song_dict["artist"]
                      if "'" in artist_name {
                          artist_name = strings.Replace(artist_name, "'", "", -1)
                      }
                      artist = Jukebox.encode_value(artist_name)
                      album_name = song_dict["album"]
                      if "'" in album_name {
                          album_name = strings.Replace(album_name, "'", "", -1)
                      }
                      album = encode_value(album_name)
                      song_name = song_dict["song"]
                      if "'" in song_name {
                          song_name = strings.Replace(song_name, "'", "", -1)
                      }
                      song = encode_value(song_name)
                      base_object_name = "%s--%s--%s" % (artist, album, song)
                      ext_list = [".flac", ".m4a", ".mp3"]
                      for ext in ext_list {
                          object_name = base_object_name + ext
                          db_song = jukebox.jukebox_db.retrieve_song(object_name)
                          if db_song != nil {
                              song_list.append(db_song)
                              break
                          } else {
                              fmt.Printf("No song file for %s\n", base_object_name)
                          }
                      }
                      jukebox.play_song_list(song_list, false)
                  }
              }
          }
       }
       */
   } else {
      fmt.Printf("error: unable to retrieve %s\n", object_name)
   }
}

func (jukebox *Jukebox) Play_album(artist string, album string) {
   bucket_name := "cj-albums"
   object_name := fmt.Sprintf("%s--%s.json", encode_value(artist), encode_value(album))
   download_file := object_name
   if jukebox.storage_system.Get_object(bucket_name,
                                        object_name,
                                        download_file) > 0 {
      //TODO: implement play_album
      /*
      //try:
      with open(download_file, 'rb') as content_file:
         file_contents = content_file.read()
      file_read = true
      //except IOError:
      //    fmt.Printf("error: unable to read file %s\n", full_path)
      //    file_read = false

      if file_read {
          pl = json.loads(file_contents)
                if pl != nil {
                    if "tracks" in pl {
                        song_list = []
                        list_song_dicts = pl["tracks"]
                        for song_dict in list_song_dicts {
                            base_object_name = song_dict["object"]
                            pos_dot = base_object_name.find(".")
                            if pos_dot > 0 {
                                base_object_name = base_object_name[0:pos_dot]
                            }
                            ext_list = [".flac", ".m4a", ".mp3"]
                            for ext in ext_list:
                                object_name = base_object_name + ext
                                db_song = jukebox.jukebox_db.retrieve_song(object_name)
                                if db_song != nil {
                                    song_list.append(db_song)
                                    break
                                }
                            else:
                                fmt.Printf("No song file for %s\n", base_object_name)
                        }
                        jukebox.play_song_list(song_list, false)
      */
   } else {
      fmt.Printf("error: unable to retrieve %s\n", object_name)
   }
}

func (jukebox *Jukebox) Delete_song(song_uid string, upload_metadata bool) bool {
   is_deleted := false
   if len(song_uid) > 0 {
      db_deleted := jukebox.jukebox_db.delete_song(song_uid)
      container := jukebox.container_for_song(song_uid)
      if container != nil {
         container_value := *container
         ss_deleted := false
         if len(container_value) > 0 {
            ss_deleted = jukebox.storage_system.Delete_object(container_value, song_uid)
            if db_deleted && upload_metadata {
               jukebox.Upload_metadata_db()
            }
            is_deleted = db_deleted || ss_deleted
         }
      }
   }

   return is_deleted
}

func (jukebox *Jukebox) Delete_artist(artist string) bool {
   is_deleted := false
   if len(artist) > 0 {
      song_list := jukebox.jukebox_db.retrieve_songs(artist, "")
      if song_list != nil {
         if len(song_list) == 0 {
            fmt.Println("no songs in jukebox")
         } else {
            for _, song := range song_list {
               if ! jukebox.Delete_song(song.Fm.Object_name, false) {
                  fmt.Printf("error deleting song '%s'\n", song.Fm.Object_name)
                  return false
               }
            }
            jukebox.Upload_metadata_db()
            is_deleted = true
         }
      } else {
         fmt.Println("no songs in jukebox")
      }
   }

   return is_deleted
}

func (jukebox *Jukebox) Delete_album(album string) bool {
   pos_double_dash := strings.Index(album, "--")
   if pos_double_dash > -1 {
      artist := album[0:pos_double_dash]
      album_name := album[pos_double_dash+2:]
      list_album_songs := jukebox.jukebox_db.retrieve_songs(artist, album_name)
      if list_album_songs != nil && len(list_album_songs) > 0 {
	 num_songs_deleted := 0
	 for _, song := range list_album_songs {
            fmt.Printf("%s %s\n", song.Fm.Container_name, song.Fm.Object_name)
            // delete each song audio file
            if jukebox.storage_system.Delete_object(song.Fm.Container_name,
                                                    song.Fm.Object_name) {
               num_songs_deleted += 1
               // delete song metadata
               jukebox.jukebox_db.delete_song(song.Fm.Object_name)
            } else {
               fmt.Println("error: unable to delete song %s\n", song.Fm.Object_name)
            }
         }
         //TODO: delete song metadata if we got 404
         if num_songs_deleted > 0 {
            // upload metadata db
            jukebox.Upload_metadata_db()
            return true
         }
      } else {
         fmt.Printf("no songs found for artist='%s' album name='%s'\n", artist, album_name)
      }
   } else {
      fmt.Println("specify album with 'the-artist--the-song-name' format")
   }

   return false
}

func (jukebox *Jukebox) Delete_playlist(playlist_name string) (bool) {
   is_deleted := false
   object_name := jukebox.jukebox_db.get_playlist(playlist_name)
   if object_name != nil && len(*object_name) > 0 {
      object_name_value := *object_name
      db_deleted := jukebox.jukebox_db.delete_playlist(playlist_name)
      if db_deleted {
         fmt.Printf("container='%s', object='%s'\n", jukebox.playlist_container, object_name_value)
         if jukebox.storage_system.Delete_object(jukebox.playlist_container, object_name_value) {
            is_deleted = true
         } else {
            fmt.Println("error: object delete failed")
         }
      } else {
         fmt.Println("error: database delete failed")
         if is_deleted {
            jukebox.Upload_metadata_db()
         } else {
            fmt.Println("delete of playlist failed")
         }
      }
   } else {
      fmt.Println("invalid playlist name")
   }

   return is_deleted
}

func (jukebox *Jukebox) Import_album_art() {
   if jukebox.jukebox_db != nil && jukebox.jukebox_db.is_open() {
      file_import_count := 0
      dir_listing, err := ioutil.ReadDir(jukebox.album_art_import_dir)
      if err != nil {
         return
      } else {
         if len(dir_listing) == 0 {
            fmt.Println("no album art found")
            return
         }
      }

      have_container := false

      if ! jukebox.storage_system.Has_container(jukebox.album_art_container) {
         have_container = jukebox.storage_system.Create_container(jukebox.album_art_container)
      } else {
         have_container = true
      }

      if ! have_container {
         fmt.Println("error: unable to create container for album art. unable to import")
         return
      }

      headers := make(map[string]string)

      for _, listing_entry := range dir_listing {
         if listing_entry.IsDir() {
            continue
         }

         full_path := Path_join(jukebox.album_art_import_dir, listing_entry.Name())
	 object_name := listing_entry.Name()
	 file_read, file_contents, _ := jukebox.read_file_contents(full_path, false)
         if file_read && file_contents != nil {
            if jukebox.storage_system.Put_object(jukebox.album_art_container,
                                                 object_name,
                                                 file_contents,
                                                 headers) {
               file_import_count += 1
            }
         }
      }

      if file_import_count > 0 {
         fmt.Println("%d album art files imported", file_import_count)
      } else {
         fmt.Println("no files imported")
      }
   }
}
