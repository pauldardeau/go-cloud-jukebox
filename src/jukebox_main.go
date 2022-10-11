package main

import (
   "bufio"
   "fmt"
   "jukebox"
   "os"
   "strings"
)


func connectStorageSystem(systemName string,
                          credentials map[string]string,
                          prefix string,
                          inDebugMode bool,
                          isUpdate bool) *jukebox.FSStorageSystem {
   if systemName == "swift" {
      //return connectSwiftSystem(credentials, prefix, inDebugMode, isUpdate)
   } else if systemName == "s3" {
      //return connectS3System(credentials, prefix, inDebugMode, isUpdate)
   } else if systemName == "azure" {
      //return connectAzureSystem(credentials, prefix, inDebugMode, isUpdate)
   } else if systemName == "fs" {
      rootDir, exists := credentials["root_dir"]
      if exists && len(rootDir) > 0 {
         return jukebox.NewFSStorageSystem(rootDir, inDebugMode)
      }
   }
   return nil
}

func showUsage() {
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
   fmt.Println("\tshow-album         - show songs in a specified album")
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

func initStorageSystem(storageSys *jukebox.FSStorageSystem) bool {
   var success bool
   if jukebox.InitializeStorageSystem(storageSys) {
      fmt.Println("storage system successfully initialized")
      success = true
   } else {
      fmt.Println("error: unable to initialize storage system")
      success = false
   }
   return success
}

func main() {
   exitCode := 0
   debugMode := false
   storageType := "swift"
   artist := ""
   //shuffle := false
   playlist := ""
   song := ""
   album := ""

   optParser := jukebox.NewArgumentParser()
   optParser.AddOptionalBoolFlag("--debug", "run in debug mode")
   optParser.AddOptionalIntArgument("--file-cache-count", "number of songs to buffer in cache")
   optParser.AddOptionalBoolFlag("--integrity-checks", "check file integrity after download")
   optParser.AddOptionalBoolFlag("--compress", "use gzip compression")
   optParser.AddOptionalBoolFlag("--encrypt", "encrypt file contents")
   optParser.AddOptionalStringArgument("--key", "encryption key")
   optParser.AddOptionalStringArgument("--keyfile", "path to file containing encryption key")
   optParser.AddOptionalStringArgument("--storage", "storage system type (s3, swift, azure)")
   optParser.AddOptionalStringArgument("--artist", "limit operations to specified artist")
   optParser.AddOptionalStringArgument("--playlist", "limit operations to specified playlist")
   optParser.AddOptionalStringArgument("--song", "limit operations to specified song")
   optParser.AddOptionalStringArgument("--album", "limit operations to specified album")
   optParser.AddRequiredArgument("command", "command for jukebox")

   consoleArgs := os.Args[1:]

   ps := optParser.ParseArgs(consoleArgs)

   if ps == nil {
      fmt.Println("error: unable to obtain command-line arguments")
      os.Exit(1)
   }

   options := jukebox.NewJukeboxOptions()

   //fmt.Println("initial values for options:")
   //options.Show()

   if ps.Contains("debug") {
      debugMode = true
      options.DebugMode = true
   }

   if ps.Contains("file_cache_count") {
      //value := args["file_cache_count"]
      //if args.file_cache_count != nil && args.file_cache_count > 0 {
      //   if debug_mode {
      //      fmt.Printf("setting file cache count=%d", args.file_cache_count)
      //   }
      //   options.FileCacheCount = args.file_cache_count
      //}
   }

   if ps.Contains("integrity_checks") {
      if debugMode {
         fmt.Println("setting integrity checks on")
      }
      options.CheckDataIntegrity = true
   }

   if ps.Contains("compress") {
      if debugMode {
         fmt.Println("setting compression on")
      }
      options.UseCompression = true
   }

   if ps.Contains("encrypt") {
      if debugMode {
         fmt.Println("setting encryption on")
      }
      options.UseEncryption = true
   }

   if ps.Contains("key") {
      pv := ps.Get("key")
      keyValue := pv.GetStringValue()
      if debugMode {
         fmt.Printf("setting encryption key='%s'\n", keyValue)
      }
      options.EncryptionKey = keyValue
   }

   if ps.Contains("keyfile") {
      pvKeyFile := ps.Get("keyfile")
      keyFile := pvKeyFile.GetStringValue()
      if debugMode {
          fmt.Printf("reading encryption key file='%s'\n", keyFile)
      }

      encryptKey, errKey := jukebox.FileReadAllText(keyFile)
      if errKey != nil {
          fmt.Printf("error: unable to read key file '%s'\n", keyFile)
          os.Exit(1)
      }
      options.EncryptionKey = strings.TrimSpace(encryptKey)

      if len(options.EncryptionKey) == 0 {
          fmt.Printf("error: no key found in file '%s'\n", keyFile)
          os.Exit(1)
      }
   }

   if ps.Contains("storage") {
      pvStorage := ps.Get("storage")
      storageType = pvStorage.GetStringValue()

      supportedSystems := []string{"swift", "s3", "azure", "fs"}
      selectedSystemSupported := false
      for _, supportedSystem := range supportedSystems {
         if supportedSystem == storageType {
            selectedSystemSupported = true
            break
         }
      }

      if ! selectedSystemSupported {
         fmt.Printf("error: invalid storage type '%s'\n", storageType)
         //print("supported systems are: %s" % str(supportedSystems))
         os.Exit(1)
      } else {
         if debugMode {
            fmt.Printf("setting storage system to '%s'\n", storageType)
         }
      }
   }

   if ps.Contains("artist") {
      pvArtist := ps.Get("artist")
      artist = pvArtist.GetStringValue()
   }

   if ps.Contains("playlist") {
      pvPlaylist := ps.Get("playlist")
      playlist = pvPlaylist.GetStringValue()
   }

   if ps.Contains("song") {
      pvSong := ps.Get("song")
      song = pvSong.GetStringValue()
   }

   if ps.Contains("album") {
      pvAlbum := ps.Get("album")
      album = pvAlbum.GetStringValue()
   }

   if ps.Contains("command") {
      pvCommand := ps.Get("command")
      command := pvCommand.GetStringValue()

      if debugMode {
         fmt.Printf("using storage system type '%s'\n", storageType)
      }

      containerPrefix := "com.swampbits.jukebox."
      credsFile := storageType + "_creds.txt"
      var creds = make(map[string]string)
      credsFilePath := ""
      wd, errWd := os.Getwd()
      if errWd == nil {
         credsFilePath = jukebox.PathJoin(wd, credsFile)
      }

      if jukebox.FileExists(credsFilePath) {
         if debugMode {
            fmt.Printf("reading creds file '%s'\n", credsFilePath)
         }

         readFile, err := os.Open(credsFilePath)
         if err != nil {
            fmt.Println(err)
            os.Exit(1)
         }
         defer readFile.Close()

         fileScanner := bufio.NewScanner(readFile)
         fileScanner.Split(bufio.ScanLines)

         for fileScanner.Scan() {
            fileLine := strings.Trim(fileScanner.Text(), "\t \n")
            if len(fileLine) > 0 {
               lineTokens := strings.Split(fileLine, "=")
               if len(lineTokens) == 2 {
                  key := strings.Trim(lineTokens[0], " ")
                  value := strings.Trim(lineTokens[1], " ")
                  creds[key] = value
               }
            }
         }
      } else {
         fmt.Printf("no creds file (%s)\n", credsFilePath)
      }

      options.EncryptionIv = "sw4mpb1ts.juk3b0x"

      helpCmds := []string{"help", "usage"}
      nonHelpCmds := []string{"import-songs", "play", "shuffle-play", "list-songs",
                                "list-artists", "list-containers", "list-genres",
                                "list-albums", "retrieve-catalog", "import-playlists",
                                "list-playlists", "show-playlist", "play-playlist",
                                "delete-song", "delete-album", "delete-playlist",
                                "delete-artist", "upload-metadata-db",
                                "import-album-art", "play-album", "show-album"}
      updateCmds := []string{"import-songs", "import-playlists", "delete-song",
                              "delete-album", "delete-playlist", "delete-artist",
                              "upload-metadata-db", "import-album-art", "init-storage"}
      allCmds := []string{}
      for _, cmd := range helpCmds {
         allCmds = append(allCmds, cmd)
      }

      for _, cmd := range nonHelpCmds {
         allCmds = append(allCmds, cmd)
      }

      for _, cmd := range updateCmds {
         allCmds = append(allCmds, cmd)
      }

      commandInAllCmds := false
      commandInHelpCmds := false
      commandInUpdateCmds := false

      for _, cmd := range allCmds {
         if cmd == command {
            commandInAllCmds = true
            break
         }
      }

      for _, cmd := range helpCmds {
         if cmd == command {
            commandInHelpCmds = true
            break
         }
      }

      for _, cmd := range updateCmds {
         if cmd == command {
            commandInUpdateCmds = true
            break
         }
      }

      if ! commandInAllCmds {
          fmt.Printf("Unrecognized command '%s'\n", command)
          fmt.Println("")
          showUsage()
      } else {
          if commandInHelpCmds {
              showUsage()
          } else {
              if ! options.ValidateOptions() {
                  os.Exit(1)
              }

              if command == "upload-metadata-db" {
                  options.SuppressMetadataDownload = true
              } else {
                  options.SuppressMetadataDownload = false
              }

              isUpdate := false

              if commandInUpdateCmds {
                  isUpdate = true
              }

              storageSystem := connectStorageSystem(storageType,
                                                    creds,
                                                    containerPrefix,
                                                    debugMode,
                                                    isUpdate)
              if storageSystem != nil {
                  if storageSystem.Enter() {
                      defer storageSystem.Exit()
                      fmt.Println("storage system entered")

                      if command == "init-storage" {
                          if initStorageSystem(storageSystem) {
                             os.Exit(0)
                          } else {
                             os.Exit(1)
                          }
                      }

                      //fmt.Println("options given to jukebox:")
                      //options.Show()

                      jukebox := jukebox.NewJukebox(options, storageSystem, debugMode)
                      if jukebox.Enter() {
                          defer jukebox.Exit()
                          fmt.Println("jukebox entered")

                          if command == "import-songs" {
                              jukebox.ImportSongs()
                          } else if command == "import-playlists" {
                              jukebox.ImportPlaylists()
                          } else if command == "play" {
                              shuffle := false
                              jukebox.PlaySongs(shuffle, artist, album)
                          } else if command == "shuffle-play" {
                              shuffle := true
                              jukebox.PlaySongs(shuffle, artist, album)
                          } else if command == "list-songs" {
                              jukebox.ShowListings()
                          } else if command == "list-artists" {
                              jukebox.ShowArtists()
                          } else if command == "list-containers" {
                              jukebox.ShowListContainers()
                          } else if command == "list-genres" {
                              jukebox.ShowGenres()
                          } else if command == "list-albums" {
                              jukebox.ShowAlbums()
                          } else if command == "list-playlists" {
                              jukebox.ShowPlaylists()
                          } else if command == "show-album" {
                              if len(album) > 0 {
                                  jukebox.ShowAlbum(album)
                              } else {
                                  fmt.Println("error: album must be specified using --album option")
                                  exitCode = 1
                              }
                          } else if command == "show-playlist" {
                              if len(playlist) > 0 {
                                  jukebox.ShowPlaylist(playlist)
                              } else {
                                  fmt.Println("error: playlist must be specified using --playlist option")
                                  exitCode = 1
                              }
                          } else if command == "play-playlist" {
                              if len(playlist) > 0 {
                                  jukebox.PlayPlaylist(playlist)
                              } else {
                                  fmt.Println("error: playlist must be specified using --playlist option")
                                  exitCode = 1
                              }
                          } else if command == "play-album" {
                              if len(album) > 0 && len(artist) > 0 {
                                  jukebox.PlayAlbum(artist, album)
                              } else {
                                  fmt.Println("error: artist and album must be specified using --artist and --album options")
                              }
                          } else if command == "retrieve-catalog" {
                              //pass
                          } else if command == "delete-song" {
                              if len(song) > 0 {
                                  if jukebox.DeleteSong(song, false) {
                                      fmt.Println("song deleted")
                                  } else {
                                      fmt.Println("error: unable to delete song")
                                      exitCode = 1
                                  }
                              } else {
                                  fmt.Println("error: song must be specified using --song option")
                                  exitCode = 1
                              }
                          } else if command == "delete-artist" {
                              if len(artist) > 0 {
                                  if jukebox.DeleteArtist(artist) {
                                      fmt.Println("artist deleted")
                                  } else {
                                      fmt.Println("error: unable to delete artist")
                                      exitCode = 1
                                  }
                              } else {
                                  fmt.Println("error: artist must be specified using --artist option")
                                  exitCode = 1
                              }
                          } else if command == "delete-album" {
                              if len(album) > 0 {
                                  if jukebox.DeleteAlbum(album) {
                                      fmt.Println("album deleted")
                                  } else {
                                      fmt.Println("error: unable to delete album")
                                      exitCode = 1
                                  }
                              } else {
                                  fmt.Println("error: album must be specified using --album option")
                                  exitCode = 1
                              }
                          } else if command == "delete-playlist" {
                              if len(playlist) > 0 {
                                  if jukebox.DeletePlaylist(playlist) {
                                      fmt.Println("playlist deleted")
                                  } else {
                                      fmt.Println("error: unable to delete playlist")
                                      exitCode = 1
                                  }
                              } else {
                                  fmt.Println("error: playlist must be specified using --playlist option")
                                  exitCode = 1
                              }
                          } else if command == "upload-metadata-db" {
                              if jukebox.UploadMetadataDb() {
                                  fmt.Println("metadata db uploaded")
                              } else {
                                  fmt.Println("error: unable to upload metadata db")
                                  exitCode = 1
                              }
                          } else if command == "import-album-art" {
                              jukebox.ImportAlbumArt()
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
      showUsage()
   }

   os.Exit(exitCode)
}

