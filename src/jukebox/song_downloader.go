package jukebox


type SongDownloader struct {
   jukebox *Jukebox
   listSongs []*SongMetadata
}


func NewSongDownloader(jukebox *Jukebox,
                       listSongs []*SongMetadata) *SongDownloader {
    var sd SongDownloader;
    sd.jukebox = jukebox
    sd.listSongs = listSongs
    return &sd
}

func (sd *SongDownloader) run() {
    if sd.jukebox != nil && sd.listSongs != nil {
        sd.jukebox.batchDownloadStart()
        for _, song := range sd.listSongs {
            if sd.jukebox.exitRequested {
                break
            } else {
                sd.jukebox.downloadSong(song)
            }
        }
        sd.jukebox.batchDownloadComplete()
    }
}
