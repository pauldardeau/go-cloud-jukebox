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
//
//	|         |       |           |       |
//	|         |       |           |       |----  file extension (e.g., 'mp3')
//	|         |       |           |
//	|         |       |           |---- name of the song (' ' replaced with '-')
//	|         |       |
//	|         |       |---- name of the album (' ' replaced with '-')
//	|         |
//	|         |---- double dashes to separate the artist name and song name
//	|
//	|---- artist name (' ' replaced with '-')
//
// For example, the MP3 version of the song 'Under My Thumb' from artist 'The
// Rolling Stones' from the album 'Aftermath' should be named:
//
//	The-Rolling-Stones--Aftermath--Under-My-Thumb.mp3
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
	"bytes"
	"compress/zlib"
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"os/signal"
	"path/filepath"
	"runtime"
	"strings"
	"syscall"
	"time"
)

const (
	downloadExtension   = ".download"
	albumContainer      = "albums"
	albumArtContainer   = "album-art"
	metadataContainer   = "music-metadata"
	playlistContainer   = "playlists"
	songContainerSuffix = "-artist-songs"
	albumArtImportDir   = "album-art-import"
	playlistImportDir   = "playlist-import"
	songImportDir       = "song-import"
	songPlayDir         = "song-play"
	defaultDbFileName   = "jukebox_db.sqlite3"
)

type AlbumTrack struct {
	Number int `json:"number"`
	Title string `json:"title"`
	Object string `json:"object"`
	Length string `json:"length"`
}

type Album struct {
	Artist string `json:"artist"`
	Album string `json:"album"`
	AlbumArt string `json:"album-art"`
	Year int `json:"year"`
	Genre []string `json:"genre"`
	AlbumType string `json:"type"`
	Wiki string `json:"wiki"`
	Tracks []AlbumTrack `json:"tracks"`
}

type PlaylistSong struct {
	Artist string `json:"artist"`
	Album string `json:"album"`
	Song string `json:"song"`
}

type Playlist struct {
	Name string `json:"name"`
	Tags string `json:"tags"`
	Songs []PlaylistSong `json:"songs"`
}

type Jukebox struct {
	jukeboxOptions          *JukeboxOptions
	storageSystem           *FSStorageSystem
	debugPrint              bool
	jukeboxDb               *JukeboxDB
	currentDir              string
	songImportDir           string
	playlistImportDir       string
	songPlayDir             string
	albumArtImportDir       string
	metadataDbFile          string
	songList                []*SongMetadata
	numberSongs             int
	songIndex               int
	audioPlayerExeFileName  string
	audioPlayerCommandArgs  string
	audioPlayerProcess      *os.Process
	songPlayLengthSeconds   int
	cumulativeDownloadBytes int64
	cumulativeDownloadTime  int
	exitRequested           bool
	isPaused                bool
	songSecondsOffset       int
}

func signalHandler(signalChannel chan os.Signal, jukebox *Jukebox) {
	for {
		s := <-signalChannel
		if jukebox != nil {
			if s == syscall.SIGUSR1 {
				jukebox.TogglePausePlay()
			} else if s == syscall.SIGUSR2 {
				jukebox.AdvanceToNextSong()
			} else if s == syscall.SIGINT {
				jukebox.PrepareForTermination()
			} else if s == syscall.SIGWINCH {
				jukebox.DisplayInfo()
			}
		}
	}
}

func (jukebox *Jukebox) installSignalHandlers() {
	signalChannel := make(chan os.Signal, 1)
	signal.Notify(signalChannel)
	go signalHandler(signalChannel, jukebox)
}

func NewJukebox(jbOptions *JukeboxOptions,
	storageSys *FSStorageSystem,
	debugPrint bool) *Jukebox {
	var jukebox Jukebox
	jukebox.jukeboxOptions = jbOptions
	jukebox.storageSystem = storageSys
	jukebox.debugPrint = debugPrint
	jukebox.jukeboxDb = nil
	cwd, err := os.Getwd()
	if err == nil {
		jukebox.currentDir = cwd
	}
	jukebox.songImportDir = PathJoin(jukebox.currentDir, songImportDir)
	jukebox.playlistImportDir = PathJoin(jukebox.currentDir, playlistImportDir)
	jukebox.songPlayDir = PathJoin(jukebox.currentDir, songPlayDir)
	jukebox.albumArtImportDir = PathJoin(jukebox.currentDir, albumArtImportDir)
	jukebox.metadataDbFile = defaultDbFileName
	jukebox.songList = []*SongMetadata{}
	jukebox.numberSongs = 0
	jukebox.songIndex = -1
	jukebox.audioPlayerExeFileName = ""
	jukebox.audioPlayerCommandArgs = ""
	jukebox.audioPlayerProcess = nil
	jukebox.songPlayLengthSeconds = 20
	jukebox.cumulativeDownloadBytes = 0
	jukebox.cumulativeDownloadTime = 0
	jukebox.exitRequested = false
	jukebox.isPaused = false
	jukebox.songSecondsOffset = 0

	if jukebox.jukeboxOptions != nil && jukebox.jukeboxOptions.DebugMode {
		jukebox.debugPrint = true
	}
	if jukebox.debugPrint {
		fmt.Printf("currentDir = '%s'\n", jukebox.currentDir)
		fmt.Printf("songImportDir = '%s'\n", jukebox.songImportDir)
		fmt.Printf("songPlayDir = '%s'\n", jukebox.songPlayDir)
	}
	return &jukebox
}

func (jukebox *Jukebox) Enter() bool {
	// look for stored metadata in the storage system
	if jukebox.storageSystem != nil &&
		jukebox.storageSystem.HasContainer(metadataContainer) &&
		!jukebox.jukeboxOptions.SuppressMetadataDownload {

		// metadata container exists, retrieve container listing
		metadataFileInContainer := false
		containerContents, err := jukebox.storageSystem.ListContainerContents(metadataContainer)
		if err == nil && len(containerContents) > 0 {
			for _, container := range containerContents {
				if container == jukebox.metadataDbFile {
					metadataFileInContainer = true
					break
				}
			}
		}

		// does our metadata DB file exist in the metadata container?
		if containerContents != nil && metadataFileInContainer {
			// download it
			metadataDbFilePath := jukebox.GetMetadataDbFilePath()
			downloadFile := metadataDbFilePath + downloadExtension
			if jukebox.storageSystem.GetObject(metadataContainer, jukebox.metadataDbFile, downloadFile) > 0 {
				// have an existing metadata DB file?
				if FileExists(metadataDbFilePath) {
					if jukebox.debugPrint {
						fmt.Println("deleting existing metadata DB file")
					}
					DeleteFile(metadataDbFilePath)
					// rename downloaded file
					if jukebox.debugPrint {
						fmt.Printf("renaming '%s' to '%s'\n", downloadFile, metadataDbFilePath)
					}
					RenameFile(downloadFile, metadataDbFilePath)
				} else {
					if jukebox.debugPrint {
						fmt.Println("error: unable to retrieve metadata DB file")
					}
				}
			} else {
				if jukebox.debugPrint {
					fmt.Println("no metadata DB file in metadata container")
				}
			}
		} else {
			if jukebox.debugPrint {
				fmt.Println("no metadata container in storage system")
			}
		}

		debugPrint := true
		jukebox.jukeboxDb = NewJukeboxDB(jukebox.GetMetadataDbFilePath(),
			jukebox.jukeboxOptions.UseEncryption,
			jukebox.jukeboxOptions.UseCompression,
			debugPrint)
		jukeboxDbSuccess := jukebox.jukeboxDb.enter()
		if !jukeboxDbSuccess {
			fmt.Println("unable to connect to database")
		}
		return jukeboxDbSuccess
	}

	return false
}

func (jukebox *Jukebox) Exit() {
	if jukebox.jukeboxDb != nil {
		jukebox.jukeboxDb.exit()
		jukebox.jukeboxDb = nil
	}
}

func (jukebox *Jukebox) TogglePausePlay() {
	jukebox.isPaused = !jukebox.isPaused
	if jukebox.isPaused {
		fmt.Println("paused")
		if jukebox.audioPlayerProcess != nil {
			// capture current song position (seconds into song)
			jukebox.audioPlayerProcess.Kill()
		}
	} else {
		fmt.Println("resuming play")
	}
}

func (jukebox *Jukebox) AdvanceToNextSong() {
	fmt.Println("advancing to next song")
	if jukebox.audioPlayerProcess != nil {
		jukebox.audioPlayerProcess.Kill()
	}
}

func (jukebox *Jukebox) PrepareForTermination() {
	fmt.Println("Ctrl-C detected, shutting down")

	// indicate that it's time to shutdown
	jukebox.exitRequested = true

	// terminate audio player if it's running
	if jukebox.audioPlayerProcess != nil {
		jukebox.audioPlayerProcess.Kill()
	}
}

func (jukebox *Jukebox) DisplayInfo() {
	if len(jukebox.songList) > 0 {
		maxIndex := len(jukebox.songList) - 1
		if jukebox.songIndex+3 <= maxIndex {
			fmt.Printf("----- songs on deck -----\n")
			firstSong := jukebox.songList[jukebox.songIndex+1]
			fmt.Printf("%s\n", firstSong.Fm.FileUid)
			secondSong := jukebox.songList[jukebox.songIndex+2]
			fmt.Printf("%s\n", secondSong.Fm.FileUid)
			thirdSong := jukebox.songList[jukebox.songIndex+3]
			fmt.Printf("%s\n", thirdSong.Fm.FileUid)
			fmt.Printf("-------------------------\n")
		}
	}
}

func (jukebox *Jukebox) GetMetadataDbFilePath() string {
	return PathJoin(jukebox.currentDir, jukebox.metadataDbFile)
}

func componentsFromFileName(fileName string) (string, string, string) {
	if len(fileName) == 0 {
		return "", "", ""
	}
	posExtension := strings.Index(fileName, ".")
	var baseFileName string
	if posExtension > -1 {
		baseFileName = fileName[0:posExtension]
	} else {
		baseFileName = fileName
	}
	components := strings.Split(baseFileName, "--")
	if len(components) == 3 {
		return UnencodeValue(components[0]),
			UnencodeValue(components[1]),
			UnencodeValue(components[2])
	} else {
		return "", "", ""
	}
}

func artistFromFileName(fileName string) string {
	if len(fileName) > 0 {
		artist, _, _ := componentsFromFileName(fileName)
		if len(artist) > 0 {
			return artist
		}
	}
	return ""
}

func albumFromFileName(fileName string) string {
	if len(fileName) > 0 {
		_, album, _ := componentsFromFileName(fileName)
		if len(album) > 0 {
			return album
		}
	}
	return ""
}

func songFromFileName(fileName string) string {
	if len(fileName) > 0 {
		_, _, song := componentsFromFileName(fileName)
		if len(song) > 0 {
			return song
		}
	}
	return ""
}

func (jukebox *Jukebox) storeSongMetadata(fsSong *SongMetadata) bool {
	dbSong := jukebox.jukeboxDb.retrieveSong(fsSong.Fm.FileUid)
	if dbSong != nil {
		if !fsSong.Equals(dbSong) {
			return jukebox.jukeboxDb.updateSong(fsSong)
		} else {
			return true // no insert or update needed (already up-to-date)
		}
	} else {
		// song is not in the database, insert it
		return jukebox.jukeboxDb.insertSong(fsSong)
	}
}

func (jukebox *Jukebox) storeSongPlaylist(fileName string, fileContents []byte) bool {
        var playlist Playlist
        err := json.Unmarshal([]byte(fileContents), &playlist)
	if err == nil {
		if len(playlist.Name) > 0 {
			return jukebox.jukeboxDb.insertPlaylist(fileName, playlist.Name, "")
		} else {
			fmt.Printf("error: playlist name is missing\n")
			return false
		}
	} else {
		fmt.Printf("error: unable to parse playlist json\n")
		fmt.Printf("error: %v\n", err)
		return false
	}
}

func (jukebox *Jukebox) containerSuffix() string {
	suffix := ""
	if jukebox.jukeboxOptions.UseEncryption &&
		jukebox.jukeboxOptions.UseCompression {
		suffix += "-ez"
	} else if jukebox.jukeboxOptions.UseEncryption {
		suffix += "-e"
	} else if jukebox.jukeboxOptions.UseCompression {
		suffix += "-z"
	}
	return suffix
}

func (jukebox *Jukebox) objectFileSuffix() string {
	suffix := ""
	if jukebox.jukeboxOptions.UseEncryption &&
		jukebox.jukeboxOptions.UseCompression {
		suffix = ".egz"
	} else if jukebox.jukeboxOptions.UseEncryption {
		suffix = ".e"
	} else if jukebox.jukeboxOptions.UseCompression {
		suffix = ".gz"
	}
	return suffix
}

func (jukebox *Jukebox) containerForSong(songUid string) string {
	if len(songUid) == 0 {
		return ""
	}
	containerSuffix := songContainerSuffix + jukebox.containerSuffix()

	artist := artistFromFileName(songUid)
	if len(artist) == 0 {
		return ""
	}

	var artistLetter string
	artistValue := artist

	if strings.HasPrefix(artistValue, "A ") {
		artistLetter = artistValue[2:3]
	} else if strings.HasPrefix(artistValue, "The ") {
		artistLetter = artistValue[4:5]
	} else {
		artistLetter = artistValue[0:1]
	}

	containerName := strings.ToLower(artistLetter) + containerSuffix
	return containerName
}

func (jukebox *Jukebox) ImportSongs() {
	if jukebox.jukeboxDb != nil && jukebox.jukeboxDb.isOpen() {
		dirListing, err := ListFilesInDirectory(jukebox.songImportDir)
		if err != nil {
			return
		}
		numEntries := float32(len(dirListing))
		progressbarChars := 0.0
		progressbarWidth := 40
		progressCharsPerIteration := float32(progressbarWidth) / numEntries
		progressbarChar := "#"
		barChars := 0

		if !jukebox.debugPrint {
			// setup progressbar
			fmt.Printf("[%s]", strings.Repeat(" ", progressbarWidth))
			fmt.Printf(strings.Repeat("\b", progressbarWidth+1)) // return to start of line, after '['
		}

		cumulativeUploadTime := 0
		cumulativeUploadBytes := 0
		fileImportCount := 0

		for _, listingEntry := range dirListing {
			fullPath := PathJoin(jukebox.songImportDir, listingEntry)
			// ignore it if it's not a file
			if FileExists(fullPath) {
				fileName := listingEntry
				_, extension := PathSplitExt(fullPath)
				if len(extension) > 0 {
					fileSize := GetFileSize(fullPath)
					artist := artistFromFileName(fileName)
					album := albumFromFileName(fileName)
					song := songFromFileName(fileName)
					if fileSize > 0 && len(artist) > 0 && len(album) > 0 && len(song) > 0 {
						objectName := fileName + jukebox.objectFileSuffix()
						fsSong := NewSongMetadata()
						fsSong.Fm = NewFileMetadata()
						fsSong.Fm.FileUid = objectName
						fsSong.AlbumUid = ""
						fsSong.Fm.OriginFileSize = fileSize
						mtime, errTime := PathGetMtime(fullPath)
						if errTime == nil {
							fsSong.Fm.FileTime = mtime.Format(time.RFC3339)
						}
						fsSong.ArtistName = artist
						fsSong.SongName = song
						md5Hash, errHash := Md5ForFile(fullPath)
						if errHash == nil {
							fsSong.Fm.Md5Hash = md5Hash
						}
						fsSong.Fm.Compressed = jukebox.jukeboxOptions.UseCompression
						fsSong.Fm.Encrypted = jukebox.jukeboxOptions.UseEncryption
						fsSong.Fm.ObjectName = objectName
						fsSong.Fm.PadCharCount = 0

						fsSong.Fm.ContainerName = jukebox.containerForSong(fileName)

						// read file contents
						fileRead := false

						fileContents, errFile := FileReadAllBytes(fullPath)
						if errFile == nil {
							fileRead = true
						} else {
							fmt.Printf("error: unable to read file %s\n", fullPath)
						}

						if fileRead && fileContents != nil {
							if len(fileContents) > 0 {
								// for general purposes, it might be useful or helpful to have
								// a minimum size for compressing
								if jukebox.jukeboxOptions.UseCompression {
									if jukebox.debugPrint {
										fmt.Println("compressing file")
									}

									level := zlib.BestCompression
									var buffer bytes.Buffer
									w, e := zlib.NewWriterLevel(&buffer, level)
									if e != nil {
										fmt.Printf("error: unable to create new zlib writer for level=%d\n", level)
										fmt.Printf("error: %v\n", e)
									} else {
										w.Write(fileContents)
										w.Close()
										fileContents = buffer.Bytes()
									}
								}

								if jukebox.jukeboxOptions.UseEncryption {
									if jukebox.debugPrint {
										fmt.Println("encrypting file")
									}

									// the length of the data to encrypt must be a multiple of 16
									numExtraChars := len(fileContents) % 16
									if numExtraChars > 0 {
										if jukebox.debugPrint {
											fmt.Println("padding file for encryption")
										}
										numPadChars := 16 - numExtraChars
										padding := []byte(strings.Repeat(" ", numPadChars))
										fileContents = append(fileContents, padding...)
										fsSong.Fm.PadCharCount = numPadChars
									}

									var encodedContents string
									var errEncrypt error
									encodedContents, errEncrypt = EncryptAES([]byte(jukebox.jukeboxOptions.EncryptionKey), fileContents)
									if errEncrypt != nil {
										fmt.Println("error: encryption failed")
										fmt.Printf("error: %v\n", errEncrypt)
									} else {
										fileContents = []byte(encodedContents)
									}
								}
							}

							// now that we have the data that will be stored, set the file size for
							// what's being stored
							fsSong.Fm.StoredFileSize = int64(len(fileContents))
							//startUploadTime := time.Now()

							// store song file to storage system
							if jukebox.storageSystem.PutObject(fsSong.Fm.ContainerName,
								fsSong.Fm.ObjectName,
								fileContents,
								nil) {
								//endUploadTime := time.Now()
								// endUploadTime - startUploadTime
								//uploadElapsedTime := endUploadTime.Add(-startUploadTime)
								//cumulativeUploadTime.Add(uploadElapsedTime)
								cumulativeUploadBytes += len(fileContents)

								// store song metadata in local database
								if !jukebox.storeSongMetadata(fsSong) {
									// we stored the song to the storage system, but were unable to store
									// the metadata in the local database. we need to delete the song
									// from the storage system since we won't have any way to access it
									// since we can't store the song metadata locally.
									fmt.Printf("unable to store metadata, deleting obj '%s'", fsSong.Fm.ObjectName)
									jukebox.storageSystem.DeleteObject(fsSong.Fm.ContainerName,
										fsSong.Fm.ObjectName)
								} else {
									fileImportCount += 1
								}
							} else {
								fmt.Printf("error: unable to upload '%s' to '%s'\n",
									fsSong.Fm.ObjectName,
									fsSong.Fm.ContainerName)
							}
						}
					}
				}

				if !jukebox.debugPrint {
					progressbarChars += float64(progressCharsPerIteration)
					if int(progressbarChars) > barChars {
						numNewChars := int(progressbarChars) - barChars
						if numNewChars > 0 {
							// update progress bar
							for j := 0; j < numNewChars; j++ {
								fmt.Print(progressbarChar)
							}
							barChars += numNewChars
						}
					}
				}
			}
		}

		if !jukebox.debugPrint {
			// if we haven't filled up the progress bar, fill it now
			if barChars < progressbarWidth {
				numNewChars := progressbarWidth - barChars
				for j := 0; j < numNewChars; j++ {
					fmt.Print(progressbarChar)
				}
			}
			fmt.Printf("\n")
		}

		if fileImportCount > 0 {
			jukebox.UploadMetadataDb()
		}

		fmt.Printf("%d song files imported\n", fileImportCount)

		if cumulativeUploadTime > 0 {
			cumulativeUploadKb := cumulativeUploadBytes / 1000.0
			fmt.Printf("average upload throughput = %d KB/sec\n",
				cumulativeUploadKb/cumulativeUploadTime)
		}
	}
}

func (jukebox *Jukebox) songPathInPlaylist(song *SongMetadata) string {
	return PathJoin(jukebox.songPlayDir, song.Fm.FileUid)
}

func (jukebox *Jukebox) checkFileIntegrity(song *SongMetadata) bool {
	fileIntegrityPassed := true

	if jukebox.jukeboxOptions != nil && jukebox.jukeboxOptions.CheckDataIntegrity {
		filePath := jukebox.songPathInPlaylist(song)
		if FileExists(filePath) {
			if jukebox.debugPrint {
				fmt.Printf("checking integrity for %s\n", song.Fm.FileUid)
			}

			if song.Fm != nil {
				playlistMd5, err := Md5ForFile(filePath)
				if err != nil {
					fmt.Printf("error: unable to calculate MD5 hash for file '%s'\n", filePath)
					fmt.Printf("error: %v\n", err)
					fileIntegrityPassed = false
				} else {
					if playlistMd5 == song.Fm.Md5Hash {
						if jukebox.debugPrint {
							fmt.Println("integrity check SUCCESS")
						}
						fileIntegrityPassed = true
					} else {
						fmt.Printf("file integrity check failed: %s\n", song.Fm.FileUid)
						fileIntegrityPassed = false
					}
				}
			}
		} else {
			// file doesn't exist
			fmt.Println("file doesn't exist")
			fileIntegrityPassed = false
		}
	} else {
		if jukebox.debugPrint {
			fmt.Println("file integrity bypassed, no jukebox options or check integrity not turned on")
		}
	}

	return fileIntegrityPassed
}

func (jukebox *Jukebox) batchDownloadStart() {
	jukebox.cumulativeDownloadBytes = 0
	jukebox.cumulativeDownloadTime = 0
}

func (jukebox *Jukebox) batchDownloadComplete() {
	if !jukebox.exitRequested {
		if jukebox.cumulativeDownloadTime > 0 {
			cumulativeDownloadKb := jukebox.cumulativeDownloadBytes / 1000.0
			fmt.Printf("average download throughput = %d KB/sec\n",
				cumulativeDownloadKb/int64(jukebox.cumulativeDownloadTime))
		}
		jukebox.cumulativeDownloadBytes = 0
		jukebox.cumulativeDownloadTime = 0
	}
}

func (jukebox *Jukebox) retrieveFile(fm *FileMetadata, dirPath string) int64 {
	var bytesRetrieved int64

	if jukebox.storageSystem != nil && fm != nil && len(dirPath) > 0 {
		localFilePath := PathJoin(dirPath, fm.FileUid)
		bytesRetrieved = jukebox.storageSystem.GetObject(fm.ContainerName, fm.ObjectName, localFilePath)
	}

	return bytesRetrieved
}

func (jukebox *Jukebox) downloadSong(song *SongMetadata) bool {
	if jukebox.exitRequested {
		return false
	}

	if song != nil {
		filePath := jukebox.songPathInPlaylist(song)
		//downloadStartTime := time.time()
		songBytesRetrieved := jukebox.retrieveFile(song.Fm, jukebox.songPlayDir)
		if jukebox.exitRequested {
			return false
		}

		if jukebox.debugPrint {
			fmt.Printf("bytes retrieved: %d\n", songBytesRetrieved)
		}

		if songBytesRetrieved > 0 {
			//downloadEndTime := time.time()
			//downloadElapsedTime := downloadEndTime - downloadStartTime
			//jukebox.cumulativeDownloadTime += downloadElapsedTime
			jukebox.cumulativeDownloadBytes += songBytesRetrieved

			// are we checking data integrity?
			// if so, verify that the storage system retrieved the same length that has been stored
			if jukebox.jukeboxOptions != nil && jukebox.jukeboxOptions.CheckDataIntegrity {
				if jukebox.debugPrint {
					fmt.Println("verifying data integrity")
				}

				if songBytesRetrieved != song.Fm.StoredFileSize {
					fmt.Printf("error: file size check failed for '%s'\n", filePath)
					return false
				}
			}

			// is it encrypted? if so, unencrypt it
			encrypted := song.Fm.Encrypted
			compressed := song.Fm.Compressed

			if encrypted || compressed {
				fileContents, errFile := FileReadAllBytes(filePath)
				if errFile != nil {
					fmt.Printf("error: unable to read file %s\n", filePath)
					fmt.Printf("error: %v\n", errFile)
					return false
				}

				if encrypted {
					encodedContents := string(fileContents[:])
					var errDecrypt error
					fileContents, errDecrypt = DecryptAES([]byte(jukebox.jukeboxOptions.EncryptionKey), encodedContents)
					if errDecrypt != nil {
						fmt.Printf("error: unable to decrypt file: %s\n", filePath)
						fmt.Printf("error: %v\n", errDecrypt)
						return false
					}
				}

				if compressed {
					//file_contents = zlib.decompress(file_contents)
				}

				// re-write out the uncompressed, unencrypted file contents
				if !FileWriteAllBytes(filePath, fileContents) {
					fmt.Printf("error: unable to write unencrypted/uncompressed file '%s'\n", filePath)
					return false
				}
			}

			if jukebox.checkFileIntegrity(song) {
				return true
			} else {
				// we retrieved the file, but it failed our integrity check
				// if file exists, remove it
				if FileExists(filePath) {
					DeleteFile(filePath)
				}
			}
		}
	}

	return false
}

func (jukebox *Jukebox) playSong(song *SongMetadata) {
	songFilePath := jukebox.songPathInPlaylist(song)

	if FileExists(songFilePath) {
		fmt.Printf("playing %s\n", song.Fm.FileUid)
		if len(jukebox.audioPlayerExeFileName) > 0 {
			var args []string
			if len(jukebox.audioPlayerCommandArgs) > 0 {
				vecAddlArgs := strings.Split(jukebox.audioPlayerCommandArgs, " ")
				for _, addlArg := range vecAddlArgs {
					args = append(args, addlArg)
				}
			}
			args = append(args, songFilePath)

			exitCode := -1
			startedAudioPlayer := false
			var cmd *exec.Cmd
			playerExe := jukebox.audioPlayerExeFileName

			numArgs := len(args)
			if numArgs == 1 {
				cmd = exec.Command(playerExe, args[0])
			} else if numArgs == 2 {
				cmd = exec.Command(playerExe, args[0], args[1])
			} else if numArgs == 3 {
				cmd = exec.Command(playerExe, args[0], args[1], args[2])
			} else if numArgs == 4 {
				cmd = exec.Command(playerExe, args[0], args[1], args[2], args[3])
			} else if numArgs == 5 {
				cmd = exec.Command(playerExe,
					args[0],
					args[1],
					args[2],
					args[3],
					args[4])
			} else if numArgs == 6 {
				cmd = exec.Command(playerExe,
					args[0],
					args[1],
					args[2],
					args[3],
					args[4],
					args[5])
			} else if numArgs == 7 {
				cmd = exec.Command(playerExe,
					args[0],
					args[1],
					args[2],
					args[3],
					args[4],
					args[5],
					args[6])
			} else if numArgs == 8 {
				cmd = exec.Command(playerExe,
					args[0],
					args[1],
					args[2],
					args[3],
					args[4],
					args[5],
					args[6],
					args[7])
			} else {
				fmt.Printf("error: too many arguments specified for audio player (max 8)\n")
				return
			}

			err := cmd.Start() // will not wait for process to run and exit
			if err == nil {
				startedAudioPlayer = true
				jukebox.audioPlayerProcess = cmd.Process
				errWait := cmd.Wait()
				if errWait != nil {
					if jukebox.debugPrint {
						fmt.Printf("error: unable to wait for audio player process\n")
						fmt.Printf("error: %v\n", errWait)
					}
				}
				jukebox.audioPlayerProcess = nil
			} else {
				fmt.Printf("error: unable to start audio player\n")
				fmt.Printf("error: %v\n", err)
				jukebox.audioPlayerExeFileName = ""
				jukebox.audioPlayerCommandArgs = ""
			}

			// if the audio player failed or is not present, just sleep
			// for the length of time that audio would be played
			if !startedAudioPlayer && exitCode != 0 {
				TimeSleepSeconds(jukebox.songPlayLengthSeconds)
			}
		} else {
			// we don't know about an audio player, so simulate a
			// song being played by sleeping
			TimeSleepSeconds(jukebox.songPlayLengthSeconds)
		}

		if !jukebox.isPaused {
			// delete the song file from the play list directory
			DeleteFile(songFilePath)
		}
	} else {
		fmt.Printf("song file doesn't exist: '%s'\n", songFilePath)

		f, err := os.OpenFile("404.txt",
			os.O_APPEND|os.O_CREATE|os.O_WRONLY,
			0644)
		if err != nil {
			fmt.Println("error: unable to open 404.txt to append song file")
			fmt.Println(err)
			return
		}
		defer f.Close()
		if _, err := f.WriteString(songFilePath + "\n"); err != nil {
			fmt.Println("error: unable to write to 404.txt")
			fmt.Println(err)
		}
	}
}

func (jukebox *Jukebox) downloadSongs() {
	// scan the play list directory to see if we need to download more songs
	dirListing, err := ListFilesInDirectory(jukebox.songPlayDir)
	if err != nil {
		// log error
		return
	}

	var dlSongs []*SongMetadata

	songFileCount := 0
	for _, fileName := range dirListing {
		extension := filepath.Ext(fileName)
		if len(extension) > 0 && extension != downloadExtension {
			songFileCount += 1
		}
	}

	fileCacheCount := jukebox.jukeboxOptions.FileCacheCount

	if songFileCount < fileCacheCount {
		// start looking at the next song in the list
		checkIndex := jukebox.songIndex + 1
		for j := 0; j < jukebox.numberSongs; j++ {
			if checkIndex >= jukebox.numberSongs {
				checkIndex = 0
			}
			if checkIndex != jukebox.songIndex {
				si := jukebox.songList[checkIndex]
				filePath := jukebox.songPathInPlaylist(si)
				if !FileExists(filePath) {
					dlSongs = append(dlSongs, si)
					if len(dlSongs) >= fileCacheCount {
						break
					}
				}
			}
			checkIndex += 1
		}
	}

	if len(dlSongs) > 0 {
		go downloadSongs(jukebox, dlSongs)
	}
}

func downloadSongs(jukebox *Jukebox, dlSongs []*SongMetadata) {
	downloader := NewSongDownloader(jukebox, dlSongs)
	downloader.run()
}

func (jukebox *Jukebox) PlaySongs(shuffle bool, artist string, album string) {
	songList := jukebox.jukeboxDb.retrieveSongs(artist, album)
	jukebox.playSongList(songList, shuffle)
}

func (jukebox *Jukebox) playSongList(songList []*SongMetadata, shuffle bool) {
	jukebox.songList = songList
	if jukebox.songList != nil {
		jukebox.numberSongs = len(jukebox.songList)

		if jukebox.numberSongs == 0 {
			fmt.Println("no songs in jukebox")
			return
		}

		// does play list directory exist?
		if !FileExists(jukebox.songPlayDir) {
			if jukebox.debugPrint {
				fmt.Println("song-play directory does not exist, creating it")
			}
			os.Mkdir(jukebox.songPlayDir, os.ModePerm)
		} else {
			// play list directory exists, delete any files in it
			if jukebox.debugPrint {
				fmt.Println("deleting existing files in song-play directory")
			}
			DeleteFilesInDirectory(jukebox.songPlayDir)
		}

		jukebox.songIndex = 0
		jukebox.installSignalHandlers()

		osId := runtime.GOOS
		if strings.HasPrefix(osId, "darwin") {
			jukebox.audioPlayerExeFileName = "afplay"
			jukebox.audioPlayerCommandArgs = ""
		} else if strings.HasPrefix(osId, "linux") ||
			strings.HasPrefix(osId, "freebsd") ||
			strings.HasPrefix(osId, "netbsd") ||
			strings.HasPrefix(osId, "openbsd") {

			jukebox.audioPlayerExeFileName = "/usr/bin/mplayer"
			jukebox.audioPlayerCommandArgs = "-novideo -nolirc -really-quiet"
		} else if strings.HasPrefix(osId, "windows") {
			// we really need command-line support for /play and /close arguments. unfortunately,
			// this support used to be available in the built-in Windows Media Player, but is
			// no longer present.
			jukebox.audioPlayerExeFileName = "C:\\Program Files\\MPC-HC\\mpc-hc64.exe"
			jukebox.audioPlayerCommandArgs = "/play /close /minimized"
		} else {
			fmt.Printf("error: %s is not a supported OS\n", osId)
			os.Exit(1)
		}

		fmt.Println("downloading first song...")

		if shuffle {
			//TODO: add shuffling of song list
			//jukebox.songList = random.sample(jukebox.songList, len(jukebox.songList))
		}

		if jukebox.downloadSong(jukebox.songList[0]) {
			fmt.Println("first song downloaded. starting playing now.")

			pidAsText := fmt.Sprintf("%d\n", os.Getpid())
			FileWriteAllText("jukebox.pid", pidAsText)

			for true {
				if !jukebox.exitRequested {
					if !jukebox.isPaused {
						jukebox.downloadSongs()
						jukebox.playSong(jukebox.songList[jukebox.songIndex])
					}
					if !jukebox.isPaused {
						jukebox.songIndex += 1
						if jukebox.songIndex >= jukebox.numberSongs {
							jukebox.songIndex = 0
						}
					} else {
						time.Sleep(1 * time.Second)
					}
				} else {
					break
				}
			}
			DeleteFile("jukebox.pid")
		} else {
			fmt.Println("error: unable to download songs")
			os.Exit(1)
		}
	}
}

func (jukebox *Jukebox) ShowListContainers() {
	if jukebox.storageSystem != nil {
		listContainers, err := jukebox.storageSystem.GetContainerNames()
		if err == nil {
			for _, containerName := range listContainers {
				fmt.Println(containerName)
			}
		} else {
			fmt.Println("error: unable to retrieve list of containers")
			fmt.Printf("error: %v\n", err)
		}
	}
}

func (jukebox *Jukebox) ShowListings() {
	if jukebox.jukeboxDb != nil {
		jukebox.jukeboxDb.showListings()
	}
}

func (jukebox *Jukebox) ShowArtists() {
	if jukebox.jukeboxDb != nil {
		jukebox.jukeboxDb.showArtists()
	}
}

func (jukebox *Jukebox) ShowGenres() {
	if jukebox.jukeboxDb != nil {
		jukebox.jukeboxDb.showGenres()
	}
}

func (jukebox *Jukebox) ShowAlbums() {
	if jukebox.jukeboxDb != nil {
		jukebox.jukeboxDb.showAlbums()
	}
}

func (jukebox *Jukebox) readFileContents(filePath string) (bool, []byte, int) {
	fileRead := false
	padChars := 0

	fileContents, errFile := FileReadAllBytes(filePath)
	if errFile != nil {
		fmt.Printf("error: unable to read file '%s'\n", filePath)
		fmt.Printf("error: %v\n", errFile)
		return false, nil, 0
	} else {
		fileRead = true
	}

	if fileRead && fileContents != nil {
		if len(fileContents) > 0 {
			// for general purposes, it might be useful or helpful to have
			// a minimum size for compressing
			if jukebox.jukeboxOptions.UseCompression {
				if jukebox.debugPrint {
					fmt.Println("compressing file")
				}

				level := zlib.BestCompression
				var buffer bytes.Buffer
				w, e := zlib.NewWriterLevel(&buffer, level)
				if e != nil {
					fmt.Printf("error: unable to create new zlib writer for level=%d\n", level)
					fmt.Printf("error: %v\n", e)
				} else {
					w.Write(fileContents)
					w.Close()
					fileContents = buffer.Bytes()
				}
			}

			if jukebox.jukeboxOptions.UseEncryption {
				if jukebox.debugPrint {
					fmt.Println("encrypting file")
				}

				// the length of the data to encrypt must be a multiple of 16
				numExtraChars := len(fileContents) % 16
				if numExtraChars > 0 {
					if jukebox.debugPrint {
						fmt.Println("padding file for encryption")
					}
					numPadChars := 16 - numExtraChars
					padding := []byte(strings.Repeat(" ", numPadChars))
					fileContents = append(fileContents, padding...)
					padChars = numPadChars
				}

				var encodedContents string
				var errEncrypt error
				encodedContents, errEncrypt = EncryptAES([]byte(jukebox.jukeboxOptions.EncryptionKey), fileContents)
				if errEncrypt != nil {
					fmt.Println("error: encryption failed")
					fmt.Printf("error: %v\n", errEncrypt)
				} else {
					fileContents = []byte(encodedContents)
				}
			}
		}
	}

	return fileRead, fileContents, padChars
}

func (jukebox *Jukebox) UploadMetadataDb() bool {
	metadataDbUpload := false
	haveMetadataContainer := false
	if !jukebox.storageSystem.HasContainer(metadataContainer) {
		haveMetadataContainer = jukebox.storageSystem.CreateContainer(metadataContainer)
	} else {
		haveMetadataContainer = true
	}

	if haveMetadataContainer {
		if jukebox.debugPrint {
			fmt.Println("uploading metadata db file to storage system")
		}

		jukebox.jukeboxDb.close()
		jukebox.jukeboxDb = nil

		metadataDbUpload := false
		dbFilePath := jukebox.GetMetadataDbFilePath()
		dbFileContents, errFile := FileReadAllBytes(dbFilePath)
		if errFile == nil {
			metadataDbUpload = jukebox.storageSystem.PutObject(metadataContainer,
				jukebox.metadataDbFile,
				dbFileContents,
				nil)
		} else {
			fmt.Printf("error: unable to read metadata db file\n")
			fmt.Printf("error: %v\n", errFile)
		}

		if jukebox.debugPrint {
			if metadataDbUpload {
				fmt.Println("metadata db file uploaded")
			} else {
				fmt.Println("unable to upload metadata db file")
			}
		}
	}

	return metadataDbUpload
}

func (jukebox *Jukebox) ImportPlaylists() {
	if jukebox.jukeboxDb != nil && jukebox.jukeboxDb.isOpen() {
		fileImportCount := 0
		dirListing, err := ListFilesInDirectory(jukebox.playlistImportDir)
		if err != nil {
			return
		}
		if len(dirListing) == 0 {
			fmt.Println("no playlists found")
			return
		}

		haveContainer := false
		if !jukebox.storageSystem.HasContainer(playlistContainer) {
			haveContainer = jukebox.storageSystem.CreateContainer(playlistContainer)
		} else {
			haveContainer = true
		}

		if !haveContainer {
			fmt.Println("error: unable to create container for playlists. unable to import")
			return
		}

		for _, fileName := range dirListing {
			fullPath := PathJoin(jukebox.playlistImportDir, fileName)
			objectName := fileName
			fileRead, fileContents, _ := jukebox.readFileContents(fullPath)
			if fileRead && fileContents != nil {
				if jukebox.storageSystem.PutObject(playlistContainer,
					objectName,
					fileContents,
					nil) {
					fmt.Println("put of playlist succeeded")
					if !jukebox.storeSongPlaylist(objectName, fileContents) {
						fmt.Println("storing of playlist to db failed")
						jukebox.storageSystem.DeleteObject(playlistContainer,
							objectName)
					} else {
						fmt.Println("storing of playlist succeeded")
						fileImportCount += 1
					}
				}
			}
		}

		if fileImportCount > 0 {
			fmt.Printf("%d playlists imported\n", fileImportCount)
			jukebox.UploadMetadataDb()
		} else {
			fmt.Println("no files imported")
		}
	}
}

func (jukebox *Jukebox) ShowPlaylists() {
	if jukebox.jukeboxDb != nil {
		jukebox.jukeboxDb.showPlaylists()
	}
}

func (jukebox *Jukebox) ShowAlbum(albumUid string) {
        album := jukebox.getAlbum(albumUid)
        if album != nil {
		fmt.Printf("%s %d (%s)\n", album.Album, album.Year, album.Artist)
		for _, track := range album.Tracks {
			fmt.Printf("%d %s (%s)\n", track.Number, track.Title, track.Length)
		}
	} else {
		fmt.Printf("error: unable to retrieve album '%s'\n", albumUid)
	}
}

func (jukebox *Jukebox) retrievePlaylist(playlist string) *Playlist {
        objectName := fmt.Sprintf("%s.json", EncodeValue(playlist))
        downloadFile := objectName
        if jukebox.storageSystem.GetObject(playlistContainer,
                objectName,
                downloadFile) > 0 {

                fileContents, err := FileReadAllText(downloadFile)
                if err != nil {
                        fmt.Printf("error: unable to read file %s\n", downloadFile)
                } else {
                        var playlist Playlist
                        err := json.Unmarshal([]byte(fileContents), &playlist)
                        if err != nil {
                                fmt.Printf("error: unable to parse json playlist\n")
			} else {
				return &playlist
			}
		}
	} else {
		fmt.Printf("error: unable to retrieve playlist from object storage\n")
	}
	return nil
}

func (jukebox *Jukebox) ShowPlaylist(playlistName string) {
	playlist := jukebox.retrievePlaylist(playlistName)
	if playlist != nil {
		for _, song := range playlist.Songs {
			fmt.Printf("%s - %s (%s)\n", song.Artist, song.Song, song.Album)
		}
	} else {
		fmt.Printf("error: unable to retrieve playlist '%s'\n", playlistName)
	}
}

func (jukebox *Jukebox) PlayPlaylist(playlistName string) {
	playlist := jukebox.retrievePlaylist(playlistName)
	if playlist != nil {
		songList := make([]*SongMetadata, 0)
		extList := make([]string, 0)
		extList = append(extList, ".flac")
		extList = append(extList, ".m4a")
		extList = append(extList, ".mp3")

		for _, song := range playlist.Songs {
			artistName := EncodeValue(song.Artist)
			albumName := EncodeValue(song.Album)
			songName := EncodeValue(song.Song)
			baseObjectName := fmt.Sprintf("%s--%s--%s", artistName, albumName, songName)
			for _, ext := range extList {
				objectName := baseObjectName + ext
				dbSong := jukebox.jukeboxDb.retrieveSong(objectName)
				if dbSong != nil {
					songList = append(songList, dbSong)
					break
				} else {
					fmt.Printf("No song file for %s\n", baseObjectName)
				}
			}
		}
		jukebox.playSongList(songList, false)
	} else {
		fmt.Printf("error: unable to retrieve playlist '%s'\n", playlistName)
	}
}

func (jukebox *Jukebox) getAlbum(albumUid string) *Album {
        downloadFile := albumUid
        if jukebox.storageSystem.GetObject(albumContainer,
                albumUid,
                downloadFile) > 0 {

                fileContents, err := FileReadAllText(downloadFile)
                if err != nil {
                        fmt.Printf("error: unable to read file %s\n", downloadFile)
                        fmt.Printf("error: %v\n", err)
                } else {
                        var album Album
                        err := json.Unmarshal([]byte(fileContents), &album)
                        if err != nil {
				fmt.Printf("error: unable to unmarshal json for album\n")
				fmt.Printf("error: %v\n", err)
			} else {
				return &album
			}
		}
	} else {
		fmt.Printf("error: unable to retrieve album json object\n")
	}
	return nil
}

func (jukebox *Jukebox) PlayAlbum(artist string, albumName string) {
	objectName := fmt.Sprintf("%s--%s.json", EncodeValue(artist), EncodeValue(albumName))
	album := jukebox.getAlbum(objectName)
	if album != nil {
		if len(album.Tracks) > 0 {
                     var extList []string
	             extList = append(extList, ".flac")
		     extList = append(extList, ".m4a")
		     extList = append(extList, ".mp3")
                     var songList []*SongMetadata

		     for _, albumTrack := range album.Tracks {
		         baseObjectName := albumTrack.Object
			 posDot := strings.Index(baseObjectName, ".")
		         if posDot > 0 {
		               baseObjectName = baseObjectName[0:posDot]
		         }
			 for _, ext := range extList {
		             objectName := baseObjectName + ext
			     dbSong := jukebox.jukeboxDb.retrieveSong(objectName)
		             if dbSong != nil {
		                  songList = append(songList, dbSong)
		                  break
		             } else {
		                  fmt.Printf("No song file for %s\n", baseObjectName)
			     }
			 }
		     }
		     if len(songList) > 0 {
                         jukebox.playSongList(songList, false)
	             }
	    }
    }
}

func (jukebox *Jukebox) DeleteSong(songUid string, uploadMetadata bool) bool {
	isDeleted := false
	if len(songUid) > 0 {
		dbDeleted := jukebox.jukeboxDb.deleteSong(songUid)
		container := jukebox.containerForSong(songUid)
		if len(container) > 0 {
			ssDeleted := jukebox.storageSystem.DeleteObject(container, songUid)
			if dbDeleted && uploadMetadata {
				jukebox.UploadMetadataDb()
			}
			isDeleted = dbDeleted || ssDeleted
		}
	}

	return isDeleted
}

func (jukebox *Jukebox) DeleteArtist(artist string) bool {
	isDeleted := false
	if len(artist) > 0 {
		songList := jukebox.jukeboxDb.retrieveSongs(artist, "")
		if songList != nil {
			if len(songList) == 0 {
				fmt.Println("no artist songs in jukebox")
			} else {
				for _, song := range songList {
					if !jukebox.DeleteSong(song.Fm.ObjectName, false) {
						fmt.Printf("error deleting song '%s'\n", song.Fm.ObjectName)
						return false
					}
				}
				jukebox.UploadMetadataDb()
				isDeleted = true
			}
		} else {
			fmt.Println("no songs in jukebox")
		}
	}

	return isDeleted
}

func (jukebox *Jukebox) DeleteAlbum(album string) bool {
	posDoubleDash := strings.Index(album, "--")
	if posDoubleDash > -1 {
		artist := album[0:posDoubleDash]
		albumName := album[posDoubleDash+2:]
		listAlbumSongs := jukebox.jukeboxDb.retrieveSongs(artist, albumName)
		if listAlbumSongs != nil && len(listAlbumSongs) > 0 {
			numSongsDeleted := 0
			for _, song := range listAlbumSongs {
				fmt.Printf("%s %s\n", song.Fm.ContainerName, song.Fm.ObjectName)
				// delete each song audio file
				if jukebox.storageSystem.DeleteObject(song.Fm.ContainerName,
					song.Fm.ObjectName) {
					numSongsDeleted += 1
					// delete song metadata
					jukebox.jukeboxDb.deleteSong(song.Fm.ObjectName)
				} else {
					fmt.Printf("error: unable to delete song %s\n", song.Fm.ObjectName)
				}
			}
			//TODO: delete song metadata if we got 404?
			if numSongsDeleted > 0 {
				// upload metadata db
				jukebox.UploadMetadataDb()
				return true
			}
		} else {
			fmt.Printf("no songs found for artist='%s' album name='%s'\n", artist, albumName)
		}
	} else {
		fmt.Println("specify album with 'the-artist--the-song-name' format")
	}

	return false
}

func (jukebox *Jukebox) DeletePlaylist(playlistName string) bool {
	isDeleted := false
	objectName := jukebox.jukeboxDb.getPlaylist(playlistName)
	if objectName != nil && len(*objectName) > 0 {
		objectNameValue := *objectName
		dbDeleted := jukebox.jukeboxDb.deletePlaylist(playlistName)
		if dbDeleted {
			fmt.Printf("container='%s', object='%s'\n", playlistContainer, objectNameValue)
			if jukebox.storageSystem.DeleteObject(playlistContainer, objectNameValue) {
				isDeleted = true
			} else {
				fmt.Println("error: object delete failed")
			}
		} else {
			fmt.Println("error: database delete failed")
			if isDeleted {
				jukebox.UploadMetadataDb()
			} else {
				fmt.Println("delete of playlist failed")
			}
		}
	} else {
		fmt.Println("invalid playlist name")
	}

	return isDeleted
}

func (jukebox *Jukebox) ImportAlbumArt() {
	if jukebox.jukeboxDb != nil && jukebox.jukeboxDb.isOpen() {
		fileImportCount := 0
		dirListing, err := ListFilesInDirectory(jukebox.albumArtImportDir)
		if err != nil {
			return
		} else {
			if len(dirListing) == 0 {
				fmt.Println("no album art found")
				return
			}
		}

		haveContainer := false

		if !jukebox.storageSystem.HasContainer(albumArtContainer) {
			haveContainer = jukebox.storageSystem.CreateContainer(albumArtContainer)
		} else {
			haveContainer = true
		}

		if !haveContainer {
			fmt.Println("error: unable to create container for album art. unable to import")
			return
		}

		for _, fileName := range dirListing {
			fullPath := PathJoin(jukebox.albumArtImportDir, fileName)
			objectName := fileName
			fileRead, fileContents, _ := jukebox.readFileContents(fullPath)
			if fileRead && fileContents != nil {
				if jukebox.storageSystem.PutObject(albumArtContainer,
					objectName,
					fileContents,
					nil) {
					fileImportCount += 1
				}
			}
		}

		if fileImportCount > 0 {
			fmt.Printf("%d album art files imported\n", fileImportCount)
		} else {
			fmt.Println("no files imported")
		}
	}
}

func InitializeStorageSystem(storageSys *FSStorageSystem) bool {
	// create the containers that will hold songs
	artistSongChars := "0123456789abcdefghijklmnopqrstuvwxyz"

	for _, ch := range artistSongChars {
		containerName := fmt.Sprintf("%c%s", ch, songContainerSuffix)
		if !storageSys.CreateContainer(containerName) {
			fmt.Printf("error: unable to create container '%s'\n", containerName)
			return false
		}
	}

	// create the other (non-song) containers
	containerNames := make([]string, 0)
	containerNames = append(containerNames, metadataContainer)
	containerNames = append(containerNames, albumArtContainer)
	containerNames = append(containerNames, albumContainer)
	containerNames = append(containerNames, playlistContainer)

	for _, containerName := range containerNames {
		if !storageSys.CreateContainer(containerName) {
			fmt.Printf("error: unable to create container '%s'\n", containerName)
			return false
		}
	}

	// delete metadata DB file if present
	metadataDbFile := "jukebox_db.sqlite3"
	if FileExists(metadataDbFile) {
		//if (debugPrint) {
		//   fmt.Printf("deleting existing metadata DB file\n");
		//}
		DeleteFile(metadataDbFile)
	}

	return true
}
