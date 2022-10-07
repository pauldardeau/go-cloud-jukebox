package main

import (
   "bufio"
   "fmt"
   "jukebox"
   "os"
   "strings"
)


func connect_storage_system(system_name string,
                            credentials map[string]string,
                            prefix string,
                            in_debug_mode bool,
                            is_update bool) *jukebox.FSStorageSystem {
   if system_name == "swift" {
      //return connect_swift_system(credentials, prefix, in_debug_mode, is_update)
   } else if system_name == "s3" {
      //return connect_s3_system(credentials, prefix, in_debug_mode, is_update)
   } else if system_name == "azure" {
      //return connect_azure_system(credentials, prefix, in_debug_mode, is_update)
   } else if system_name == "fs" {
      rootDir, exists := credentials["root_dir"]
      if exists && len(rootDir) > 0 {
         return jukebox.NewFSStorageSystem(rootDir, in_debug_mode)  
      }
   }
   return nil
}

func show_usage() {
   fmt.Println("Supported Commands:")
   fmt.Println("\tdelete-artist      - delete specified artist")
   fmt.Println("\tdelete-album       - delete specified album")
   fmt.Println("\tdelete-playlist    - delete specified playlist")
   fmt.Println("\tdelete-song        - delete specified song")
   fmt.Println("\thelp               - show this help message")
   fmt.Println("\timport-songs       - import all new songs from song-import subdirectory")
   fmt.Println("\timport-playlists   - import all new playlists from playlist-import subdirectory")
   fmt.Println("\timport-album-art   - import all album art from album-art-import subdirectory")
   fmt.Println("\tlist-songs         - show listing of all available songs")
   fmt.Println("\tlist-artists       - show listing of all available artists")
   fmt.Println("\tlist-containers    - show listing of all available storage containers")
   fmt.Println("\tlist-albums        - show listing of all available albums")
   fmt.Println("\tlist-genres        - show listing of all available genres")
   fmt.Println("\tlist-playlists     - show listing of all available playlists")
   fmt.Println("\tshow-playlist      - show songs in specified playlist")
   fmt.Println("\tplay               - start playing songs")
   fmt.Println("\tshuffle-play       - play songs randomly")
   fmt.Println("\tplay-playlist      - play specified playlist")
   fmt.Println("\tplay-album         - play specified album")
   fmt.Println("\tretrieve-catalog   - retrieve copy of music catalog")
   fmt.Println("\tupload-metadata-db - upload SQLite metadata")
   fmt.Println("\tinit-storage       - initialize storage system")
   fmt.Println("\tusage              - show this help message")
   fmt.Println("")
}

func initStorageSystem(storage_sys *jukebox.FSStorageSystem) {
   if jukebox.InitializeStorageSystem(storage_sys) {
      fmt.Println("storage system successfully initialized")
   } else {
      fmt.Println("error: unable to initialize storage system")
      os.Exit(1)
   }
}

func main() {
   debug_mode := false
   storage_type := "swift"
   artist := ""
   //shuffle := false
   playlist := ""
   song := ""
   album := ""

   opt_parser := jukebox.NewArgumentParser()
   opt_parser.Add_optional_bool_flag("--debug", "run in debug mode")
   opt_parser.Add_optional_int_argument("--file-cache-count", "number of songs to buffer in cache")
   opt_parser.Add_optional_bool_flag("--integrity-checks", "check file integrity after download")
   opt_parser.Add_optional_bool_flag("--compress", "use gzip compression")
   opt_parser.Add_optional_bool_flag("--encrypt", "encrypt file contents")
   opt_parser.Add_optional_string_argument("--key", "encryption key")
   opt_parser.Add_optional_string_argument("--keyfile", "path to file containing encryption key")
   opt_parser.Add_optional_string_argument("--storage", "storage system type (s3, swift, azure)")
   opt_parser.Add_optional_string_argument("--artist", "limit operations to specified artist")
   opt_parser.Add_optional_string_argument("--playlist", "limit operations to specified playlist")
   opt_parser.Add_optional_string_argument("--song", "limit operations to specified song")
   opt_parser.Add_optional_string_argument("--album", "limit operations to specified album")
   opt_parser.Add_required_argument("command", "command for jukebox")

   console_args := os.Args[1:]

   args := opt_parser.Parse_args(console_args)

   //args := make(map[string]string);

   if args == nil {
      fmt.Println("error: unable to obtain command-line arguments")
      os.Exit(1)
   }

   options := jukebox.NewJukeboxOptions()

   //fmt.Println("initial values for options:")
   //options.Show()

   _, debug_exists := args["debug"] 
   if debug_exists {
      debug_mode = true
      options.Debug_mode = true
   }

   _, file_cache_count_exists := args["file_cache_count"]
   if file_cache_count_exists {
      //value := args["file_cache_count"]
      //if args.file_cache_count != nil && args.file_cache_count > 0 {
      //   if debug_mode {
      //      fmt.Printf("setting file cache count=%d", args.file_cache_count)
      //   }
      //   options.File_cache_count = args.file_cache_count
      //}
   }

   _, integrity_checks_exists := args["integrity_checks"]
   if integrity_checks_exists {
      if debug_mode {
         fmt.Println("setting integrity checks on")
      }
      options.Check_data_integrity = true
   }

   _, compress_exists := args["compress"]
   if compress_exists {
      if debug_mode {
         fmt.Println("setting compression on")
      }
      options.Use_compression = true
   }

   _, encrypt_exists := args["encrypt"]
   if encrypt_exists {
      if debug_mode {
         fmt.Println("setting encryption on")
      }
      //options.Use_encryption = true
   }

   key_value, key_exists := args["key"]
   if key_exists {
      if debug_mode {
         fmt.Printf("setting encryption key='%s'\n", key_value)
      }
      options.Encryption_key = fmt.Sprintf("%v", key_value)
   }

   _, keyfile_exists := args["keyfile"]
   if keyfile_exists {
	   /*
        keyFile := pvKeyFile.GetStringValue()
        if debug_mode {
            fmt.Printf("reading encryption key file='%s'\n", keyFile)
        }

	encryptKey, errKey := jukebox.FileReadAllText(keyFile)
	if errKey != nil {
            fmt.Printf("error: unable to read key file '%s'\n", keyFile)
            os.Exit(1)
        }
        options.Encryption_key = strings.TrimSpace(encryptKey)

        if len(options.Encryption_key) == 0 {
            fmt.Printf("error: no key found in file '%s'\n", keyFile)
            os.Exit(1)
        }
	*/
   }

   storage_value, storage_exists := args["storage"]
   if storage_exists {
      supported_systems := []string{"swift", "s3", "azure", "fs"}
      selected_system_supported := false
      for _, supported_system := range supported_systems {
         if supported_system == storage_value {
            selected_system_supported = true
            break
         }
      }

      if ! selected_system_supported {
         fmt.Printf("error: invalid storage type '%s'\n", storage_value)
         //print("supported systems are: %s" % str(supported_systems))
         os.Exit(1)
      } else {
         if debug_mode {
            fmt.Printf("setting storage system to '%s'\n", storage_value)
         }
         storage_type = fmt.Sprintf("%v", storage_value)
      }
   }

   artist_value, artist_exists := args["artist"]
   if artist_exists {
      artist = fmt.Sprintf("%v", artist_value)
   }

   playlist_value, playlist_exists := args["playlist"]
   if playlist_exists {
      playlist = fmt.Sprintf("%v", playlist_value)
   }

   song_value, song_exists := args["song"]
   if song_exists {
      song = fmt.Sprintf("%v", song_value)
   }

   album_value, album_exists := args["album"]
   if album_exists {
      album = fmt.Sprintf("%v", album_value)
   }

   command_value, command_exists := args["command"]
   if command_exists {
      if debug_mode {
         fmt.Printf("using storage system type '%s'\n", storage_type)
      }

      container_prefix := "com.swampbits.jukebox."
      creds_file := storage_type + "_creds.txt"
      var creds = make(map[string]string)
      creds_file_path := ""
      wd, err_wd := os.Getwd()
      if err_wd == nil {
         creds_file_path = jukebox.PathJoin(wd, creds_file)
      }

      if jukebox.FileExists(creds_file_path) {
         if debug_mode {
            fmt.Printf("reading creds file '%s'\n", creds_file_path)
         }

         readFile, err := os.Open(creds_file_path)
         if err != nil {
            fmt.Println(err)
            os.Exit(1)
         }
         defer readFile.Close()

         fileScanner := bufio.NewScanner(readFile)
         fileScanner.Split(bufio.ScanLines)

         for fileScanner.Scan() {
            file_line := strings.Trim(fileScanner.Text(), "\t \n")
            if len(file_line) > 0 {
               line_tokens := strings.Split(file_line, "=")
               if len(line_tokens) == 2 {
                  key := strings.Trim(line_tokens[0], " ")
                  value := strings.Trim(line_tokens[1], " ")
                  creds[key] = value
               }
            }
         }
      } else {
         fmt.Printf("no creds file (%s)\n", creds_file_path)
      }

      options.Encryption_iv = "sw4mpb1ts.juk3b0x"

      command := fmt.Sprintf("%v", command_value)

      help_cmds := []string{"help", "usage"}
      non_help_cmds := []string{"import-songs", "play", "shuffle-play", "list-songs",
                                "list-artists", "list-containers", "list-genres",
                                "list-albums", "retrieve-catalog", "import-playlists",
                                "list-playlists", "show-playlist", "play-playlist",
                                "delete-song", "delete-album", "delete-playlist",
                                "delete-artist", "upload-metadata-db",
                                "import-album-art", "play-album"}
      update_cmds := []string{"import-songs", "import-playlists", "delete-song",
                              "delete-album", "delete-playlist", "delete-artist",
                              "upload-metadata-db", "import-album-art", "init-storage"}
      all_cmds := []string{}
      for _, cmd := range help_cmds {
         all_cmds = append(all_cmds, cmd)
      }

      for _, cmd := range non_help_cmds {
         all_cmds = append(all_cmds, cmd)
      } 

      for _, cmd := range update_cmds {
         all_cmds = append(all_cmds, cmd)
      }

      command_in_all_cmds := false
      command_in_help_cmds := false
      command_in_update_cmds := false

      for _, cmd := range all_cmds {
         if cmd == command {
            command_in_all_cmds = true
            break
         }
      }

      for _, cmd := range help_cmds {
         if cmd == command {
            command_in_help_cmds = true
            break
         }
      }

      for _, cmd := range update_cmds {
         if cmd == command {
            command_in_update_cmds = true
            break
         }
      }

      if ! command_in_all_cmds {
          fmt.Printf("Unrecognized command '%s'\n", command)
          fmt.Println("")
          show_usage()
      } else {
          if command_in_help_cmds {
              show_usage()
          } else {
              if ! options.Validate_options() {
                  os.Exit(1)
              }

              if command == "upload-metadata-db" {
                  options.Suppress_metadata_download = true
              } else {
                  options.Suppress_metadata_download = false
              }

              is_update := false

              if command_in_update_cmds {
                  is_update = true
              }

              storage_system := connect_storage_system(storage_type,
                                                       creds,
                                                       container_prefix,
                                                       debug_mode,
                                                       is_update)
              if storage_system != nil {
                  if storage_system.Enter() {
                      defer storage_system.Exit()
                      fmt.Println("storage system entered")

		      if command == "init-storage" {
                          initStorageSystem(storage_system)
			  os.Exit(0)
		      }

		      //fmt.Println("options given to jukebox:")
		      //options.Show()

                      jukebox := jukebox.NewJukebox(options, storage_system, debug_mode)
                      if jukebox.Enter() {
                          defer jukebox.Exit()
                          fmt.Println("jukebox entered")

                          if command == "import-songs" {
                              jukebox.Import_songs()
                          } else if command == "import-playlists" {
                              jukebox.Import_playlists()
                          } else if command == "play" {
                              shuffle := false
                              jukebox.Play_songs(shuffle, artist, album)
                          } else if command == "shuffle-play" {
                              shuffle := true
                              jukebox.Play_songs(shuffle, artist, album)
                          } else if command == "list-songs" {
                              jukebox.Show_listings()
                          } else if command == "list-artists" {
                              jukebox.Show_artists()
                          } else if command == "list-containers" {
                              jukebox.Show_list_containers()
                          } else if command == "list-genres" {
                              jukebox.Show_genres()
                          } else if command == "list-albums" {
                              jukebox.Show_albums()
                          } else if command == "list-playlists" {
                              jukebox.Show_playlists()
                          } else if command == "show-playlist" {
                              if len(playlist) > 0 {
                                  jukebox.Show_playlist(playlist)
                              } else {
                                  fmt.Println("error: playlist must be specified using --playlist option")
                                  os.Exit(1)
                              }
                          } else if command == "play-playlist" {
                              if len(playlist) > 0 {
                                  jukebox.Play_playlist(playlist)
                              } else {
                                  fmt.Println("error: playlist must be specified using --playlist option")
                                  os.Exit(1)
                              }
                          } else if command == "play-album" {
                              if len(album) > 0 && len(artist) > 0 {
                                  jukebox.Play_album(artist, album)
                              } else {
                                  fmt.Println("error: artist and album must be specified using --artist and --album options")
                              }
                          } else if command == "retrieve-catalog" {
                              //pass
                          } else if command == "delete-song" {
                              if len(song) > 0 {
                                  if jukebox.Delete_song(song, false) {
                                      fmt.Println("song deleted")
                                  } else {
                                      fmt.Println("error: unable to delete song")
                                      os.Exit(1)
                                  }
                              } else {
                                  fmt.Println("error: song must be specified using --song option")
                                  os.Exit(1)
                              }
                          } else if command == "delete-artist" {
                              if len(artist) > 0 {
                                  if jukebox.Delete_artist(artist) {
                                      fmt.Println("artist deleted")
                                  } else {
                                      fmt.Println("error: unable to delete artist")
                                      os.Exit(1)
                                  }
                              } else {
                                  fmt.Println("error: artist must be specified using --artist option")
                                  os.Exit(1)
                              }
                          } else if command == "delete-album" {
                              if len(album) > 0 {
                                  if jukebox.Delete_album(album) {
                                      fmt.Println("album deleted")
                                  } else {
                                      fmt.Println("error: unable to delete album")
                                      os.Exit(1)
                                  }
                              } else {
                                  fmt.Println("error: album must be specified using --album option")
                                  os.Exit(1)
                              }
                          } else if command == "delete-playlist" {
                              if len(playlist) > 0 {
                                  if jukebox.Delete_playlist(playlist) {
                                      fmt.Println("playlist deleted")
                                  } else {
                                      fmt.Println("error: unable to delete playlist")
                                      os.Exit(1)
                                  }
                              } else {
                                  fmt.Println("error: playlist must be specified using --playlist option")
                                  os.Exit(1)
                              }
                          } else if command == "upload-metadata-db" {
                              if jukebox.Upload_metadata_db() {
                                  fmt.Println("metadata db uploaded")
                              } else {
                                  fmt.Println("error: unable to upload metadata db")
                                  os.Exit(1)
                              }
                          } else if command == "import-album-art" {
                              jukebox.Import_album_art()
                          }
                      } else {
                          fmt.Println("unable to enter jukebox")
                      }
                  } else {
                      fmt.Println("unable to enter storage system")
                  }
              }
                //except requests.exceptions.ConnectionError:
                //    print("Error: unable to connect to storage system server")
                //    os.Exit(1)
            }
         }
   } else {
      fmt.Println("Error: no command given")
      show_usage()
   }
}

