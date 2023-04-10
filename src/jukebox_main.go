package main

import (
	"bufio"
	"fmt"
	"jukebox"
	"os"
	"strings"
)

const (
	argDebug           = "debug"
	argFileCacheCount  = "file-cache-count"
	argIntegrityChecks = "integrity-checks"
	argEncrypt         = "encrypt"
	argKey             = "key"
	argKeyFile         = "keyfile"
	argStorage         = "storage"
	argArtist          = "artist"
	argPlaylist        = "playlist"
	argSong            = "song"
	argAlbum           = "album"
	argCommand         = "command"

	cmdDeleteAlbum      = "delete-album"
	cmdDeleteArtist     = "delete-artist"
	cmdDeletePlaylist   = "delete-playlist"
	cmdDeleteSong       = "delete-song"
	cmdExportAlbum      = "export-album"
	cmdExportPlaylist   = "export-playlist"
	cmdHelp             = "help"
	cmdImportAlbum      = "import-album"
	cmdImportAlbumArt   = "import-album-art"
	cmdImportPlaylists  = "import-playlists"
	cmdImportSongs      = "import-songs"
	cmdInitStorage      = "init-storage"
	cmdListAlbums       = "list-albums"
	cmdListArtists      = "list-artists"
	cmdListContainers   = "list-containers"
	cmdListGenres       = "list-genres"
	cmdListPlaylists    = "list-playlists"
	cmdListSongs        = "list-songs"
	cmdPlay             = "play"
	cmdPlayAlbum        = "play-album"
	cmdPlayPlaylist     = "play-playlist"
	cmdRetrieveCatalog  = "retrieve-catalog"
	cmdShowAlbum        = "show-album"
	cmdShowPlaylist     = "show-playlist"
	cmdShufflePlay      = "shuffle-play"
	cmdUploadMetadataDb = "upload-metadata-db"
	cmdUsage            = "usage"

	ssFs = "fs"
	ssS3 = "s3"
)

func connectS3StorageSystem(credentials map[string]string,
	prefix string,
	inDebugMode bool,
	isUpdate bool) jukebox.StorageSystem {

	awsAccessKey := ""
	awsSecretKey := ""
	updateAwsAccessKey := ""
	updateAwsSecretKey := ""

	if accessKey, ok := credentials["aws_access_key"]; ok {
		awsAccessKey = accessKey
	}
	if secretKey, ok := credentials["aws_secret_key"]; ok {
		awsSecretKey = secretKey
	}

	updateAccessKey, okAccessKey := credentials["update_aws_access_key"]
	updateSecretKey, okSecretKey := credentials["update_aws_secret_key"]

	if okAccessKey && okSecretKey {
		updateAwsAccessKey = updateAccessKey
		updateAwsSecretKey = updateSecretKey
	}

	if inDebugMode {
		fmt.Printf("aws_access_key=%s\n", awsAccessKey)
		fmt.Printf("aws_secret_key=%s\n", awsSecretKey)
		if len(updateAwsAccessKey) > 0 && len(updateAwsSecretKey) > 0 {
			fmt.Printf("update_aws_access_key=%s\n", updateAwsAccessKey)
			fmt.Printf("update_aws_secret_key=%s\n", updateAwsSecretKey)
		}
	}

	if len(awsAccessKey) == 0 || len(awsSecretKey) == 0 {
		fmt.Println("error: no s3 credentials given. please specify aws_access_key " +
			"and aws_secret_key in credentials file")
		return nil
	} else {
		var accessKey string
		var secretKey string

		if isUpdate {
			accessKey = updateAwsAccessKey
			secretKey = updateAwsSecretKey
		} else {
			accessKey = awsAccessKey
			secretKey = awsSecretKey
		}

		if inDebugMode {
			fmt.Println("Creating S3StorageSystem")
		}
		return jukebox.NewS3StorageSystem(accessKey, secretKey, prefix, inDebugMode)
	}
}

func connectStorageSystem(systemName string,
	credentials map[string]string,
	prefix string,
	inDebugMode bool,
	isUpdate bool) jukebox.StorageSystem {

	if systemName == ssS3 {
		return connectS3StorageSystem(credentials, prefix, inDebugMode, isUpdate)
	} else if systemName == ssFs {
		rootDir, exists := credentials["root_dir"]
		if exists && len(rootDir) > 0 {
			return jukebox.NewFSStorageSystem(rootDir, inDebugMode)
		}
	}
	return nil
}

func showUsage() {
	fmt.Println("Supported Commands:")
	fmt.Printf("\t%s      - delete specified artist\n", cmdDeleteArtist)
	fmt.Printf("\t%s       - delete specified album\n", cmdDeleteAlbum)
	fmt.Printf("\t%s    - delete specified playlist\n", cmdDeletePlaylist)
	fmt.Printf("\t%s        - delete specified song\n", cmdDeleteSong)
	fmt.Printf("\t%s               - show this help message\n", cmdHelp)
	fmt.Printf("\t%s       - import all new songs from song-import subdirectory\n", cmdImportSongs)
	fmt.Printf("\t%s   - import all new playlists from playlist-import subdirectory\n", cmdImportPlaylists)
	fmt.Printf("\t%s   - import all album art from album-art-import subdirectory\n", cmdImportAlbumArt)
	fmt.Printf("\t%s         - show listing of all available songs\n", cmdListSongs)
	fmt.Printf("\t%s       - show listing of all available artists\n", cmdListArtists)
	fmt.Printf("\t%s    - show listing of all available storage containers\n", cmdListContainers)
	fmt.Printf("\t%s        - show listing of all available albums\n", cmdListAlbums)
	fmt.Printf("\t%s        - show listing of all available genres\n", cmdListGenres)
	fmt.Printf("\t%s     - show listing of all available playlists\n", cmdListPlaylists)
	fmt.Printf("\t%s         - show songs in a specified album\n", cmdShowAlbum)
	fmt.Printf("\t%s      - show songs in specified playlist\n", cmdShowPlaylist)
	fmt.Printf("\t%s               - start playing songs\n", cmdPlay)
	fmt.Printf("\t%s       - play songs randomly\n", cmdShufflePlay)
	fmt.Printf("\t%s      - play specified playlist\n", cmdPlayPlaylist)
	fmt.Printf("\t%s         - play specified album\n", cmdPlayAlbum)
	fmt.Printf("\t%s   - retrieve copy of music catalog\n", cmdRetrieveCatalog)
	fmt.Printf("\t%s - upload SQLite metadata\n", cmdUploadMetadataDb)
	fmt.Printf("\t%s       - initialize storage system\n", cmdInitStorage)
	fmt.Printf("\t%s              - show this help message\n", cmdUsage)
	fmt.Println("")
}

func initStorageSystem(storageSys jukebox.StorageSystem) bool {
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
	storageType := ssFs
	artist := ""
	shuffle := false
	playlist := ""
	song := ""
	album := ""

	optParser := jukebox.NewArgumentParser(debugMode)
	optParser.AddOptionalBoolFlag("--"+argDebug, "run in debug mode")
	optParser.AddOptionalIntArgument("--"+argFileCacheCount, "number of songs to buffer in cache")
	optParser.AddOptionalBoolFlag("--"+argIntegrityChecks, "check file integrity after download")
	optParser.AddOptionalBoolFlag("--"+argEncrypt, "encrypt file contents")
	optParser.AddOptionalStringArgument("--"+argKey, "encryption key")
	optParser.AddOptionalStringArgument("--"+argKeyFile, "path to file containing encryption key")
	optParser.AddOptionalStringArgument("--"+argStorage, "storage system type (s3, fs)")
	optParser.AddOptionalStringArgument("--"+argArtist, "limit operations to specified artist")
	optParser.AddOptionalStringArgument("--"+argPlaylist, "limit operations to specified playlist")
	optParser.AddOptionalStringArgument("--"+argSong, "limit operations to specified song")
	optParser.AddOptionalStringArgument("--"+argAlbum, "limit operations to specified album")
	optParser.AddRequiredArgument(argCommand, "command for jukebox")

	consoleArgs := os.Args[1:]

	ps := optParser.ParseArgs(consoleArgs)

	if ps == nil {
		fmt.Println("error: unable to obtain command-line arguments")
		os.Exit(1)
	}

	options := jukebox.NewJukeboxOptions()

	if ps.Contains(argDebug) {
		debugMode = true
		options.DebugMode = true
	}

	if ps.Contains(argFileCacheCount) {
		value := ps.Get(argFileCacheCount).GetIntValue()
		if debugMode {
			fmt.Printf("setting file cache count=%d\n", value)
		}
		options.FileCacheCount = value
	}

	if ps.Contains(argIntegrityChecks) {
		if debugMode {
			fmt.Println("setting integrity checks on")
		}
		options.CheckDataIntegrity = true
	}

	if ps.Contains(argEncrypt) {
		if debugMode {
			fmt.Println("setting encryption on")
		}
		options.UseEncryption = true
	}

	if ps.Contains(argKey) {
		keyValue := ps.Get(argKey).GetStringValue()
		if debugMode {
			fmt.Printf("setting encryption key='%s'\n", keyValue)
		}
		options.EncryptionKey = keyValue
	}

	if ps.Contains(argKeyFile) {
		keyFile := ps.Get(argKeyFile).GetStringValue()
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

	if ps.Contains(argStorage) {
		storageType = ps.Get(argStorage).GetStringValue()

		supportedSystems := []string{ssS3, ssFs}
		selectedSystemSupported := false
		for _, supportedSystem := range supportedSystems {
			if supportedSystem == storageType {
				selectedSystemSupported = true
				break
			}
		}

		if !selectedSystemSupported {
			fmt.Printf("error: invalid storage type '%s'\n", storageType)
			//TODO: print message indicating which storage systems are supported
			//print("supported systems are: %s" % str(supportedSystems))
			os.Exit(1)
		} else {
			if debugMode {
				fmt.Printf("setting storage system to '%s'\n", storageType)
			}
		}
	}

	if ps.Contains(argArtist) {
		artist = ps.Get(argArtist).GetStringValue()
	}

	if ps.Contains(argPlaylist) {
		playlist = ps.Get(argPlaylist).GetStringValue()
	}

	if ps.Contains(argSong) {
		song = ps.Get(argSong).GetStringValue()
	}

	if ps.Contains(argAlbum) {
		album = ps.Get(argAlbum).GetStringValue()
	}

	if ps.Contains(argCommand) {
		command := ps.Get(argCommand).GetStringValue()

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

			//TODO: convert the code below to use jukebox.FileReadAllText
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

		helpCmds := []string{cmdHelp, cmdUsage}
		nonHelpCmds := []string{cmdImportSongs, cmdPlay, cmdShufflePlay, cmdListSongs,
			cmdListArtists, cmdListContainers, cmdListGenres,
			cmdListAlbums, cmdRetrieveCatalog, cmdImportPlaylists,
			cmdListPlaylists, cmdShowPlaylist, cmdPlayPlaylist,
			cmdDeleteSong, cmdDeleteAlbum, cmdDeletePlaylist,
			cmdDeleteArtist, cmdUploadMetadataDb,
			cmdImportAlbumArt, cmdPlayAlbum, cmdShowAlbum}
		updateCmds := []string{cmdImportSongs, cmdImportPlaylists, cmdDeleteSong,
			cmdDeleteAlbum, cmdDeletePlaylist, cmdDeleteArtist,
			cmdUploadMetadataDb, cmdImportAlbumArt, cmdInitStorage}
		allCmds := make([]string, 0)
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

		if !commandInAllCmds {
			fmt.Printf("Unrecognized command '%s'\n", command)
			fmt.Println("")
			showUsage()
		} else {
			if commandInHelpCmds {
				showUsage()
			} else {
				if !options.ValidateOptions() {
					os.Exit(1)
				}

				if command == cmdUploadMetadataDb {
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

						if command == cmdInitStorage {
							if initStorageSystem(storageSystem) {
								os.Exit(0)
							} else {
								os.Exit(1)
							}
						}

						jukebox := jukebox.NewJukebox(options, storageSystem, debugMode)
						if jukebox.Enter() {
							defer jukebox.Exit()
							fmt.Println("jukebox entered")

							if command == cmdImportSongs {
								jukebox.ImportSongs()
							} else if command == cmdImportPlaylists {
								jukebox.ImportPlaylists()
							} else if command == cmdPlay {
								shuffle = false
								jukebox.PlaySongs(shuffle, artist, album)
							} else if command == cmdShufflePlay {
								shuffle = true
								jukebox.PlaySongs(shuffle, artist, album)
							} else if command == cmdListSongs {
								jukebox.ShowListings()
							} else if command == cmdListArtists {
								jukebox.ShowArtists()
							} else if command == cmdListContainers {
								jukebox.ShowListContainers()
							} else if command == cmdListGenres {
								jukebox.ShowGenres()
							} else if command == cmdListAlbums {
								jukebox.ShowAlbums()
							} else if command == cmdListPlaylists {
								jukebox.ShowPlaylists()
							} else if command == cmdShowAlbum {
								if len(album) > 0 {
									jukebox.ShowAlbum(album)
								} else {
									fmt.Printf("error: album must be specified using --%s option\n", argAlbum)
									exitCode = 1
								}
							} else if command == cmdShowPlaylist {
								if len(playlist) > 0 {
									jukebox.ShowPlaylist(playlist)
								} else {
									fmt.Printf("error: playlist must be specified using --%s option\n", argPlaylist)
									exitCode = 1
								}
							} else if command == cmdPlayPlaylist {
								if len(playlist) > 0 {
									jukebox.PlayPlaylist(playlist)
								} else {
									fmt.Printf("error: playlist must be specified using --%s option\n", argPlaylist)
									exitCode = 1
								}
							} else if command == cmdPlayAlbum {
								if len(album) > 0 && len(artist) > 0 {
									jukebox.PlayAlbum(artist, album)
								} else {
									fmt.Printf("error: artist and album must be specified using --%s and --%s options\n", argArtist, argAlbum)
								}
							} else if command == cmdRetrieveCatalog {
								//TODO: implement retrieve-catalog
								fmt.Printf("%s not yet implemented\n", cmdRetrieveCatalog)
							} else if command == cmdDeleteSong {
								if len(song) > 0 {
									if jukebox.DeleteSong(song, false) {
										fmt.Println("song deleted")
									} else {
										fmt.Println("error: unable to delete song")
										exitCode = 1
									}
								} else {
									fmt.Printf("error: song must be specified using --%s option\n", argSong)
									exitCode = 1
								}
							} else if command == cmdDeleteArtist {
								if len(artist) > 0 {
									if jukebox.DeleteArtist(artist) {
										fmt.Println("artist deleted")
									} else {
										fmt.Println("error: unable to delete artist")
										exitCode = 1
									}
								} else {
									fmt.Printf("error: artist must be specified using --%s option\n", argArtist)
									exitCode = 1
								}
							} else if command == cmdDeleteAlbum {
								if len(album) > 0 {
									if jukebox.DeleteAlbum(album) {
										fmt.Println("album deleted")
									} else {
										fmt.Println("error: unable to delete album")
										exitCode = 1
									}
								} else {
									fmt.Printf("error: album must be specified using --%s option\n", argAlbum)
									exitCode = 1
								}
							} else if command == cmdDeletePlaylist {
								if len(playlist) > 0 {
									if jukebox.DeletePlaylist(playlist) {
										fmt.Println("playlist deleted")
									} else {
										fmt.Println("error: unable to delete playlist")
										exitCode = 1
									}
								} else {
									fmt.Printf("error: playlist must be specified using --%s option\n", argPlaylist)
									exitCode = 1
								}
							} else if command == cmdUploadMetadataDb {
								if jukebox.UploadMetadataDb() {
									fmt.Println("metadata db uploaded")
								} else {
									fmt.Println("error: unable to upload metadata db")
									exitCode = 1
								}
							} else if command == cmdImportAlbumArt {
								jukebox.ImportAlbumArt()
							} else if command == cmdImportAlbum {
								//TODO: implement import album
								fmt.Printf("%s not yet implemented\n", cmdImportAlbum)
							} else if command == cmdExportAlbum {
								//TODO: implement export album
								fmt.Printf("%s not yet implemented\n", cmdExportAlbum)
							} else if command == cmdExportPlaylist {
								//TODO: implement export playlist
								fmt.Printf("%s not yet implemented\n", cmdExportPlaylist)
							}
						} else {
							fmt.Println("unable to enter jukebox")
						}
					} else {
						fmt.Println("unable to enter storage system")
					}
				}
			}
		}
	} else {
		fmt.Println("Error: no command given")
		showUsage()
	}

	os.Exit(exitCode)
}
